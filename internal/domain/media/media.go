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

type ImageTask struct {
	ModelId uint32 `json:"modelId"`
	ImageId string `json:"imageId" validate:"required,len=24"`
}

type Status int32

const (
	StatusDone     Status = 0
	StatusFail     Status = 1
	StatusProgress Status = 2
)

type ImageTaskResult struct {
	Message string `json:"message" validate:"required"`
	Status  Status `json:"status" validate:"required"`
}

type Processor struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

type AppInfoTask struct {
	TotalCpuThreads     int         `json:"totalCpuThreads"`
	AvailableProcessors []Processor `json:"availableProcessors"`
}
