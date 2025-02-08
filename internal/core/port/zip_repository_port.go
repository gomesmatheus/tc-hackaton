package port

import (
	"io"
)

type ZipRepository interface {
	UploadFile(id string, file io.Reader) error
	DownloadFile(id string) (io.Reader, error)
}
