package repository

import (
	"context"
	"fmt"

	"github.com/gomesmatheus/tc-hackaton/internal/core/entity"
	"github.com/gomesmatheus/tc-hackaton/internal/core/port"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) port.VideoRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(video entity.VideoFile) error {
	r.db.QueryRow(context.Background(), "INSERT INTO videos (id, owner_id, status) VALUES ($1, $2, $3)", video.Id, video.OwnerId, video.Status)

	return nil
}

func (r *PostgresRepository) FindById(id string) (*entity.VideoFile, error) {
	video := entity.VideoFile{}
	row := r.db.QueryRow(context.Background(), "SELECT id, owner_id, status FROM videos WHERE id = $1", id)
	err := row.Scan(&video.Id, &video.OwnerId, &video.Status)
	if err != nil {
		fmt.Println("Error scanning video", err)
		return nil, err
	}

	return &video, nil
}

func (r *PostgresRepository) FindByOwnerId(ownerId string) ([]entity.VideoFile, error) {
	videos := []entity.VideoFile{}
	rows, err := r.db.Query(context.Background(), "SELECT id, owner_id, status FROM videos WHERE owner_id = $1", ownerId)
	if err != nil {
		fmt.Println("Error querying videos", err)
		return nil, err
	}

	for rows.Next() {
		video := entity.VideoFile{}
		err = rows.Scan(&video.Id, &video.OwnerId, &video.Status)
		if err != nil {
			fmt.Println("Error scanning video", err)
			return nil, err
		}

		videos = append(videos, video)
	}

	return videos, nil
}

func (r *PostgresRepository) UpdateStatus(id string, status string) error {
	_, err := r.db.Exec(context.Background(), "UPDATE videos SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		fmt.Println("Error updating video status", err)
	}

	return err
}
