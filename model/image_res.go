package model

type UploadedImages struct {
	Original  string `json:"original,omitempty"`
	Medium    string `json:"medium,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}
