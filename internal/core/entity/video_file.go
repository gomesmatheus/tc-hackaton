package entity

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	allowedFileExtensions = ".mp4"
	allowedMimeTypes      = "video/mp4"
)

type VideoFile struct {
	OwnerId string
	File    multipart.File
	Header  *multipart.FileHeader
	Id      string
	Status  string
}

type VideoFileResponse struct {
	OwnerId string `json:"owner_id"`
	Id      string `json:"id"`
	Status  string `json:"status"`
}

func NewVideoFile(file multipart.File, header *multipart.FileHeader, ownerId string) (*VideoFile, error) {
	if !isExtensionValid(header) {
		fmt.Println("Invalid file format. Only .mp4 files are allowed", http.StatusUnsupportedMediaType)
		return nil, fmt.Errorf("Invalid file format: %s", header.Filename)
	}

	if valid, mimeType := isMimeTypeValid(file); !valid {
		fmt.Println("Detected MIME type:", mimeType)
		return nil, fmt.Errorf("Invalid MIME type: %s", mimeType)
	}

	id := uuid.New().String()

	err := save(file, header, id)
	if err != nil {
		return nil, err
	}

	return &VideoFile{
		OwnerId: ownerId,
		File:    file,
		Header:  header,
		Id:      id,
		Status:  "processing",
	}, nil
}

func isExtensionValid(header *multipart.FileHeader) bool {
	return strings.HasSuffix(header.Filename, allowedFileExtensions)
}

func isMimeTypeValid(file multipart.File) (bool, string) {
	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return false, ""
	}
	file.Seek(0, io.SeekStart)

	mimeType := http.DetectContentType(buf)
	return strings.HasPrefix(mimeType, allowedMimeTypes), mimeType
}

func save(file multipart.File, header *multipart.FileHeader, id string) error {
	filename := fmt.Sprintf("./%s.mp4", id)
	dst, err := os.Create(filename)
	if err != nil {
		fmt.Println("File creation error:", err)
		return fmt.Errorf("Unable to create the file: %s", filename)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Copy error:", err)
		return fmt.Errorf("Unable to save the file: %s", filename)
	}
	fmt.Println("File uploaded:", header.Filename, filename)

	return nil
}

func (v *VideoFile) GetFileName() string {
	return fmt.Sprintf("%s.mp4", v.Id)
}

func (v *VideoFile) GetZipFileName() string {
	return fmt.Sprintf("%s.zip", v.Id)
}

func (v *VideoFile) Delete() error {
	err := os.Remove(v.GetFileName())
	if err != nil {
		fmt.Println("Error deleting file with id:", v.Id, err)
		return err
	}
	return nil
}

func (v *VideoFile) ErrorProcessing() {
	v.Status = "error_processing"
	v.Delete()
}
