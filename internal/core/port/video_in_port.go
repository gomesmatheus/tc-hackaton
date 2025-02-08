package port

import (
	"io"
	"mime/multipart"

	"github.com/gomesmatheus/tc-hackaton/internal/core/entity"
)

type VideoService interface {
	GenerateFrames(file multipart.File, header *multipart.FileHeader, ownerId string) error
	GetVideos(ownerId string) ([]entity.VideoFileResponse, error)
	DownloadZip(videoId string, ownerId string) (io.Reader, error)
}
