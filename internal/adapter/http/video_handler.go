package http_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gomesmatheus/tc-hackaton/internal/core/port"
)

const maxUploadSize = 10 << 20 // 10MB

type VideoHandler struct {
	Service        port.VideoService
	UserRepository port.UserPort
}

func (h *VideoHandler) GenerateVideoFrames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	ownerID := r.URL.Query().Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Missing owner_id query parameter", http.StatusBadRequest)
		return
	}
	valid, err := h.UserRepository.ValidateToken(r.Header.Get("Authorization"), ownerID)
	if err != nil {
		http.Error(w, "Error validating token", http.StatusInternalServerError)
		fmt.Println("Validate token error:", err)
		return
	}
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		fmt.Println("Parse error:", err)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		fmt.Println("File error:", err)
		return
	}
	defer file.Close()

	h.Service.GenerateFrames(file, header, ownerID)
}

func (h *VideoHandler) GetZips(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	ownerID := r.URL.Query().Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Missing owner_id query parameter", http.StatusBadRequest)
		return
	}

	valid, err := h.UserRepository.ValidateToken(r.Header.Get("Authorization"), ownerID)
	if err != nil {
		http.Error(w, "Error validating token", http.StatusInternalServerError)
		fmt.Println("Validate token error:", err)
		return
	}
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	videos, err := h.Service.GetVideos(ownerID)

	if err != nil {
		http.Error(w, "Error retrieving videos", http.StatusInternalServerError)
		fmt.Println("Get videos error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(videos)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		fmt.Println("Error encoding response:", err)
		return
	}
}

func (h *VideoHandler) DownloadZip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	ownerID := r.URL.Query().Get("owner_id")
	if ownerID == "" {
		http.Error(w, "Missing owner_id query parameter", http.StatusBadRequest)
		return
	}

	valid, err := h.UserRepository.ValidateToken(r.Header.Get("Authorization"), ownerID)
	if err != nil {
		http.Error(w, "Error validating token", http.StatusInternalServerError)
		fmt.Println("Validate token error:", err)
		return
	}
	if !valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	videoID := r.URL.Query().Get("video_id")
	if videoID == "" {
		http.Error(w, "Missing video_id query parameter", http.StatusBadRequest)
		return
	}

	file, err := h.Service.DownloadZip(videoID, ownerID)
	if err != nil {
		http.Error(w, "Error downloading video", http.StatusInternalServerError)
		fmt.Println("Error downloading video:", err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", videoID))
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error writing file to response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
