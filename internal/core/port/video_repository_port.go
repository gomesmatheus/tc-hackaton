package port

import "github.com/gomesmatheus/tc-hackaton/internal/core/entity"

type VideoRepository interface {
	Save(video entity.VideoFile) error
	FindById(id string) (*entity.VideoFile, error)
	UpdateStatus(id string, status string) error
	FindByOwnerId(ownerId string) ([]entity.VideoFile, error)
}
