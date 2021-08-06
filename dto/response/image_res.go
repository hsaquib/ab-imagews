package response

type UploadedImages struct {
	FileName  string `json:"filename,omitempty" example:"xyz.jpg"`
	Original  string `json:"original,omitempty" example:"https://....jpg"`
	Medium    string `json:"medium,omitempty" example:"https://..../m/..jpg"`
	Thumbnail string `json:"thumbnail,omitempty" example:"https://..../s/..jpg"`
}

type ImageVariantUploadSuccessResponse struct {
	Status    string         `json:"status,omitempty"`
	Message   string         `json:"message,omitempty"`
	Success   bool           `json:"success"`
	Data      UploadedImages `json:"data"`
	Timestamp string         `json:"timestamp,omitempty"`
}
