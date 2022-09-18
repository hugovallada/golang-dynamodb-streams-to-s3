package models

type Image map[string]any

type Payload struct {
	SendToS3 bool
	Data     []byte
}
