package service

import (
	"io"
	"light-backend/model"
	"light-backend/mongoclient"
	"mime"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UploadPicture(c *fiber.Ctx, file *multipart.FileHeader, meta *model.ImageMetadata) (*primitive.ObjectID, error) {
	bucket, err := gridfs.NewBucket(mongoclient.DB)
	if err != nil {
		return nil, err
	}
	uploadOpts := options.GridFSUpload().SetMetadata(meta)
	f, err := file.Open()

	if err != nil {
		return nil, err
	}
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
	id, err := bucket.UploadFromStream(fileName, io.Reader(f), uploadOpts)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func DownloadPictureSt(c *fiber.Ctx, imageId *string) (io.Reader, *model.ImageSchema, error) {
	bucket, err := gridfs.NewBucket(mongoclient.DB)
	if err != nil {
		return nil, nil, err
	}

	id, err := primitive.ObjectIDFromHex(*imageId)
	if err != nil {
		return nil, nil, err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	cursor, err := bucket.Find(filter)
	if err != nil {
		return nil, nil, err
	}
	var file model.ImageSchema
	cursor.Next(c.Context())
	err = cursor.Decode(&file)
	if err != nil {
		return nil, nil, err
	}
	downloadStream, err := bucket.OpenDownloadStream(id)
	if err != nil {
		return nil, nil, err
	}
	return downloadStream, &file, nil
}
