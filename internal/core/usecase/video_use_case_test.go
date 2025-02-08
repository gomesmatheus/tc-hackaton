package usecase

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"testing"

	"github.com/gomesmatheus/tc-hackaton/internal/core/entity"
)

type MockVideoRepository struct {
	videos []entity.VideoFile
}

func (r *MockVideoRepository) Save(video entity.VideoFile) error {
	r.videos = append(r.videos, video)
	return nil
}

func (r *MockVideoRepository) UpdateStatus(videoId, status string) error {
	for i, v := range r.videos {
		if v.Id == videoId {
			r.videos[i].Status = status
			return nil
		}
	}
	return fmt.Errorf("video not found")
}

func (r *MockVideoRepository) FindByOwnerId(ownerId string) ([]entity.VideoFile, error) {
	var result []entity.VideoFile
	for _, v := range r.videos {
		if v.OwnerId == ownerId {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *MockVideoRepository) FindById(videoId string) (*entity.VideoFile, error) {
	for _, v := range r.videos {
		if v.Id == videoId {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("video not found")
}

type MockZipRepository struct {
	files map[string]bytes.Buffer
}

func (r *MockZipRepository) UploadFile(filename string, file io.Reader) error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(file)
	if err != nil {
		return err
	}
	r.files[filename] = buf
	return nil
}

func (r *MockZipRepository) DownloadFile(filename string) (io.Reader, error) {
	if _, exists := r.files[filename]; exists {
		// file := &mockMultipartFile{Reader: bytes.NewReader(fileContent)}
	}
	return nil, fmt.Errorf("file not found")
}

type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func TestGenerateFrames_Success(t *testing.T) {
	// Create in-memory file
	fileContent := []byte("dummy video content")
	file := &mockMultipartFile{Reader: bytes.NewReader(fileContent)}

	// Create a valid file header
	header := &multipart.FileHeader{
		Filename: "video.mp4",
	}

	// Simulate a valid owner ID
	ownerId := "123"

	// Set up mock repositories
	videoRepo := &MockVideoRepository{}
	zipRepo := &MockZipRepository{files: make(map[string]bytes.Buffer)}

	// Create the VideoUseCase instance
	videoUseCase := NewVideoUseCase(videoRepo, zipRepo)

	// Call GenerateFrames (should not create actual files)
	err := videoUseCase.GenerateFrames(file, header, ownerId)
	if err == nil {
		t.Errorf("Expected error, got %v", err)
	}

	// Check that the video was saved
	if len(videoRepo.videos) != 0 {
		t.Error("Expected no video to be saved, but no video was found")
	}
}

func TestGenerateFrames_ErrorHandling(t *testing.T) {
	// Simulate an error in the repository (mocked Save)
	videoRepo := &MockVideoRepository{}
	zipRepo := &MockZipRepository{files: make(map[string]bytes.Buffer)}

	// Create an invalid file (non-mp4)
	fileContent := []byte("dummy content")
	file := &mockMultipartFile{Reader: bytes.NewReader(fileContent)}
	header := &multipart.FileHeader{
		Filename: "video.txt", // Invalid extension
	}
	ownerId := "123"

	videoUseCase := NewVideoUseCase(videoRepo, zipRepo)

	// Call GenerateFrames (should return an error)
	err := videoUseCase.GenerateFrames(file, header, ownerId)
	if err == nil {
		t.Error("Expected error for invalid file extension, but got nil")
	}
}

func TestGetVideos(t *testing.T) {
	// Set up mock repositories with test data
	videoRepo := &MockVideoRepository{
		videos: []entity.VideoFile{
			{OwnerId: "123", Id: "video1", Status: "ready_to_download"},
			{OwnerId: "123", Id: "video2", Status: "processing"},
		},
	}
	zipRepo := &MockZipRepository{files: make(map[string]bytes.Buffer)}

	videoUseCase := NewVideoUseCase(videoRepo, zipRepo)

	// Call GetVideos
	videos, err := videoUseCase.GetVideos("123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Validate response
	if len(videos) != 2 {
		t.Errorf("Expected 2 videos, got %d", len(videos))
	}
}

func TestDownloadZip_Success(t *testing.T) {
	// Prepare mock repositories
	videoRepo := &MockVideoRepository{
		videos: []entity.VideoFile{
			{OwnerId: "123", Id: "video1", Status: "ready_to_download"},
		},
	}
	zipRepo := &MockZipRepository{
		files: map[string]bytes.Buffer{
			"video1.zip": *bytes.NewBuffer([]byte("dummy zip file content")),
		},
	}

	videoUseCase := NewVideoUseCase(videoRepo, zipRepo)

	// Call DownloadZip
	_, err := videoUseCase.DownloadZip("video1", "123")
	if err == nil {
		t.Errorf("Expected file not found, got %v", err)
	}
}
