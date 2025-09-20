package model

import "encoding/binary"

type MUpscaleImage struct {
	ModelId uint32 `json:"modelid" validate:"required"`
	ImageId string `json:"imageid" validate:"required,len=24"`
}

// Host to Big endian raw
func (m *MUpscaleImage) Htoberaw() []byte {
	numBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(numBytes, uint32(m.ModelId))
	body := numBytes
	body = append(body, m.ImageId...)
	return body
}
