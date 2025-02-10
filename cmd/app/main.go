package main

import (
	"fmt"
	"log"
	"net/http"

	http_handler "github.com/gomesmatheus/tc-hackaton/internal/adapter/http"
	"github.com/gomesmatheus/tc-hackaton/internal/adapter/repository"
	"github.com/gomesmatheus/tc-hackaton/internal/config"
	"github.com/gomesmatheus/tc-hackaton/internal/core/usecase"
)

func main() {
	db, err := config.NewPostgresDb("postgres://postgres:123@zip-db:5432/postgres")
	// db, err := config.NewPostgresDb("postgres://postgres:123@localhost:5432/postgres")
	if err != nil {
		log.Fatal("Error initializing database", err)
	}

	s3 := repository.NewS3Repository("fiap-hackaton")
	userRepository := repository.NewUserRepository()

	repository := repository.NewPostgresRepository(db)
	videoUseCase := usecase.NewVideoUseCase(repository, s3)
	videoHandler := http_handler.VideoHandler{
		Service:        videoUseCase,
		UserRepository: userRepository,
	}

	http.HandleFunc("/video", videoHandler.GenerateVideoFrames)
	http.HandleFunc("/zip/download", videoHandler.DownloadZip)
	http.HandleFunc("/zips", videoHandler.GetZips)
	fmt.Println("Poc hackaton is running!")

	log.Fatal(http.ListenAndServe(":3333", nil))
}
