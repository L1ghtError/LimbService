package media

import "net/textproto"

type ImageMetadata struct {
	UserId string               `json:"userId" bson:"userId"`
	Header textproto.MIMEHeader `json:"mimeHeader" bson:"mimeHeader"`
}

type ImageSchema struct {
	Filename string        `json:"filename" bson:"filename"`
	Length   int64         `json:"length" bson:"length"`
	Metadata ImageMetadata `json:"metadata" bson:"metadata"`
}
