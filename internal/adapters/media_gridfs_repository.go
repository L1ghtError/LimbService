package adapters

import (
	"context"
	"io"
	"light-backend/internal/domain/media"
	"mime"
	"mime/multipart"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MediaGridFsRepository struct {
	bucket gridfs.Bucket
}

type ImageSchema = media.ImageSchema
type ImageMetadata = media.ImageMetadata

func NewMediaGridFsRepository(buc gridfs.Bucket) *MediaGridFsRepository {
	return &MediaGridFsRepository{bucket: buc}
}

func (r MediaGridFsRepository) UploadPicture(ctx context.Context, file multipart.FileHeader, meta ImageMetadata) (string, error) {
	// find a way to control data flow with context
	_ = ctx

	f, err := file.Open()
	if err != nil {
		return "", err
	}

	//TODO: consider move this to domain/service layer
	fileName := file.Filename
	if len(file.Filename) == 0 {
		uuid, _ := uuid.NewV7()
		fileName = uuid.String()

		fileType := meta.Header.Get("Content-Type")
		extensions, _ := mime.ExtensionsByType(fileType)
		fileExt := ""
		if len(extensions) > 0 {
			fileExt = extensions[0]
		}
		fileName += fileExt
	}

	uploadOpts := options.GridFSUpload().SetMetadata(meta)
	id, err := r.bucket.UploadFromStream(fileName, io.Reader(f), uploadOpts)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (r MediaGridFsRepository) DownloadPicture(ctx context.Context, imageId string) (io.Reader, ImageSchema, error) {

	id, err := primitive.ObjectIDFromHex(imageId)
	if err != nil {
		return nil, ImageSchema{}, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	cursor, err := r.bucket.FindContext(ctx, filter)
	if err != nil {
		return nil, ImageSchema{}, err
	}

	var file ImageSchema
	if cursor.RemainingBatchLength() == 0 {
		return nil, ImageSchema{}, err
	}
	cursor.Next(ctx)
	err = cursor.Decode(&file)
	if err != nil {
		return nil, ImageSchema{}, err
	}

	downloadStream, err := r.bucket.OpenDownloadStream(id)
	if err != nil {
		return nil, ImageSchema{}, err
	}

	return downloadStream, file, nil
}
