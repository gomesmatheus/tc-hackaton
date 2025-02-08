package usecase

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gomesmatheus/tc-hackaton/internal/core/entity"
	"github.com/gomesmatheus/tc-hackaton/internal/core/port"
)

const (
	allowedFileExtensions = ".mp4"
	allowedMimeTypes      = "video/mp4"
	secondsInterval       = 4
)

type VideoUseCase struct {
	Repository    port.VideoRepository
	ZipRepository port.ZipRepository
}

func NewVideoUseCase(repository port.VideoRepository, zipRepository port.ZipRepository) *VideoUseCase {
	return &VideoUseCase{
		Repository:    repository,
		ZipRepository: zipRepository,
	}
}

func (v *VideoUseCase) GenerateFrames(file multipart.File, header *multipart.FileHeader, ownerId string) error {
	videoFile, err := entity.NewVideoFile(file, header, ownerId)
	if err != nil {
		return err
	}

	v.Repository.Save(*videoFile)

	err = GenerateVideoFrames(videoFile.GetFileName(), secondsInterval)
	if err != nil {
		videoFile.ErrorProcessing()
		v.Repository.UpdateStatus(videoFile.Id, "error")
		fmt.Println("Error generating frames", err)
		return err
	}

	zipFilePath, err := ZipFrames(videoFile.Id, "frame_*.png")
	if err != nil {
		videoFile.ErrorProcessing()
		v.Repository.UpdateStatus(videoFile.Id, "error")
		fmt.Println("Error zipping frames", err)
		return err
	}

	zipFile, err := os.Open(zipFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer zipFile.Close()

	err = v.ZipRepository.UploadFile(videoFile.GetZipFileName(), zipFile)
	if err != nil {
		go videoFile.ErrorProcessing()
		go v.Repository.UpdateStatus(videoFile.Id, "error")
		go DeleteFrames("frame_*.png")
		go os.Remove(fmt.Sprintf("output_%s", videoFile.GetZipFileName()))
		fmt.Println("Error uploading zip file", err)
		return err
	}

	go v.Repository.UpdateStatus(videoFile.Id, "ready_to_download")
	go DeleteFrames("frame_*.png")
	go videoFile.Delete()
	go os.Remove(fmt.Sprintf("output_%s", videoFile.GetZipFileName()))

	return nil
}

func (v *VideoUseCase) GetVideos(ownerId string) ([]entity.VideoFileResponse, error) {
	videos, err := v.Repository.FindByOwnerId(ownerId)
	if err != nil {
		return nil, err
	}

	return GetVideosResponse(videos), nil
}

func (v *VideoUseCase) DownloadZip(videoId string, ownerId string) (io.Reader, error) {
	video, err := v.Repository.FindById(videoId)
	if err != nil {
		return nil, err
	}

	if video.OwnerId != ownerId {
		return nil, fmt.Errorf("Video not found")
	}

	if video.Status != "ready_to_download" {
		return nil, fmt.Errorf("Video not ready to download")
	}

	file, err := v.ZipRepository.DownloadFile(video.GetZipFileName())
	if err != nil {
		return nil, err
	}

	return file, nil
}

func GenerateVideoFrames(videoFilePath string, seconds int) error {
	cmd := exec.Command("ffmpeg", "-i", videoFilePath, "-vf", fmt.Sprintf("fps=1/%d", seconds), "frame_%04d.png")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return err
	}

	return nil
}

func ZipFrames(id string, filesPattern string) (string, error) {
	files, err := filepath.Glob(filesPattern)
	zipFilePath := fmt.Sprintf("output_%s.zip", id)
	if err != nil {
		fmt.Println("error getting files", err)
		return "", err
	}

	args := append([]string{zipFilePath}, files...)
	cmd := exec.Command("zip", args...)
	fmt.Println("zipping files", cmd)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("error zipping files", err, stderr.String())
		return "", err
	}

	return zipFilePath, nil
}

func DeleteFrames(filesPattern string) error {
	files, err := filepath.Glob(filesPattern)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("error deleting file", file, err)
			return err
		}
	}

	return nil
}

func GetVideosResponse(videos []entity.VideoFile) []entity.VideoFileResponse {
	response := make([]entity.VideoFileResponse, 0)
	for _, video := range videos {
		response = append(response, entity.VideoFileResponse{
			OwnerId: video.OwnerId,
			Id:      video.Id,
			Status:  video.Status,
		})
	}

	return response
}
