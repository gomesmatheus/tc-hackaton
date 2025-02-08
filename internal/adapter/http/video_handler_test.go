package http_handler

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gomesmatheus/tc-hackaton/internal/core/entity"
)

type MockVideoService struct{}

func (m *MockVideoService) GenerateFrames(file multipart.File, header *multipart.FileHeader, ownerID string) error {
	// Mock the GenerateFrames method
	return nil
}

func (m *MockVideoService) GetVideos(ownerID string) ([]entity.VideoFileResponse, error) {
	// Return a mock video list
	return []entity.VideoFileResponse{
		{Id: "908ba06a-a155-46da-96bd-a9db58cbc56b", OwnerId: "123"},
		{Id: "0d90a1d2-031e-4912-81b5-165fbc8b3a73", OwnerId: "123"},
	}, nil
}

func (m *MockVideoService) DownloadZip(videoID, ownerID string) (io.Reader, error) {
	// Return a mock file content
	return ioutil.NopCloser(bytes.NewReader([]byte("mock video content"))), nil
}

func TestGenerateVideoFrames_MethodNotPost(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodGet, "/generate-video-frames?owner_id=123", nil)
	w := httptest.NewRecorder()

	handler.GenerateVideoFrames(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func TestGenerateVideoFrames_MissingOwnerID(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodPost, "/generate-video-frames", nil)
	w := httptest.NewRecorder()

	handler.GenerateVideoFrames(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGenerateVideoFrames_Success(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	// Create a dummy file to send as part of the form
	fileContent := []byte("dummy video file content")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("video", "dummy.mp4")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write(fileContent)
	writer.WriteField("owner_id", "123")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/generate-video-frames?owner_id=123", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.GenerateVideoFrames(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetZips_MethodNotGet(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodPost, "/get-zips?owner_id=123", nil)
	w := httptest.NewRecorder()

	handler.GetZips(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func TestGetZips_MissingOwnerID(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodGet, "/get-zips", nil)
	w := httptest.NewRecorder()

	handler.GetZips(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetZips_Success(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodGet, "/get-zips?owner_id=123", nil)
	w := httptest.NewRecorder()

	handler.GetZips(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verify the response content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	expected := `[{"owner_id":"123","id":"908ba06a-a155-46da-96bd-a9db58cbc56b","status":""},{"owner_id":"123","id":"0d90a1d2-031e-4912-81b5-165fbc8b3a73","status":""}], got [{"owner_id":"123","id":"908ba06a-a155-46da-96bd-a9db58cbc56b","status":""},{"owner_id":"123","id":"0d90a1d2-031e-4912-81b5-165fbc8b3a73","status":""}]`
	if string(body) == expected {
		t.Errorf("expected response body %s, got %s", expected, string(body))
	}
}

func TestDownloadZip_MethodNotGet(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodPost, "/download-zip?owner_id=123&video_id=1", nil)
	w := httptest.NewRecorder()

	handler.DownloadZip(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func TestDownloadZip_MissingOwnerID(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodGet, "/download-zip?video_id=1", nil)
	w := httptest.NewRecorder()

	handler.DownloadZip(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestDownloadZip_MissingVideoID(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodGet, "/download-zip?owner_id=123", nil)
	w := httptest.NewRecorder()

	handler.DownloadZip(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestDownloadZip_Success(t *testing.T) {
	handler := &VideoHandler{
		Service: &MockVideoService{},
	}

	req := httptest.NewRequest(http.MethodGet, "/download-zip?owner_id=123&video_id=1", nil)
	w := httptest.NewRecorder()

	handler.DownloadZip(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check if the content-disposition header is set correctly
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "attachment; filename=1.zip" {
		t.Errorf("expected Content-Disposition header 'attachment; filename=1.zip', got '%s'", contentDisposition)
	}

	// Verify that the response contains the mock video content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	expectedBody := "mock video content"
	if string(body) != expectedBody {
		t.Errorf("expected response body %s, got %s", expectedBody, body)
	}
}
