package model

import (
	"net/textproto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Image metadata for mongodb
type ImageMetadata struct {
	UserId primitive.ObjectID   `json:"userId" bson:"userId"`
	Header textproto.MIMEHeader `json:"mimeHeader" bson:"mimeHeader"`
}
