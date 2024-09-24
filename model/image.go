package model

type ImageSchema struct {
	Filename string        `json:"filename" bson:"filename"`
	Length   int64         `json:"length" bson:"length"`
	Metadata ImageMetadata `json:"metadata" bson:"metadata"`
}
