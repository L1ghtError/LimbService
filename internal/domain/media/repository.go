package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

type ErrNotFound struct {
	Id string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Picture %s, not found", e.Id)
}

type Repository interface {
	UploadPicture(ctx context.Context, file multipart.FileHeader, meta ImageMetadata) (string, error)

	DownloadPicture(ctx context.Context, imageId string) (io.Reader, ImageSchema, error)
}
