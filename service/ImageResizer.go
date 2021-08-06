package service

type ImageResizer interface {
	Resize(fileBytes []byte, size uint) ([]byte, error)
}
