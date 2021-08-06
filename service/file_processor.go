package service

import (
	"github.com/google/uuid"
	"github.com/hsaquib/ab-imagews/config"
	"github.com/hsaquib/ab-imagews/model"
	rLog "github.com/hsaquib/rest-log"
	"io/ioutil"
	"mime/multipart"
)

type FileProcessor struct {
	Config       *config.AppConfig
	Log          rLog.Logger
	Resizer      ImageResizer
	FileUploader Uploader
}

type UploadStatus struct {
	Url string
	Err error
}

func NewFileProcessor(conf *config.AppConfig, logger rLog.Logger) *FileProcessor {
	return &FileProcessor{
		Config:       conf,
		Log:          logger,
		Resizer:      NewVipsThumb(logger),
		FileUploader: NewUploader(&conf.Upload, logger),
	}
}

func (processor *FileProcessor) UploadImageVariants(file multipart.File, header *multipart.FileHeader) (*model.UploadedImages, error) {

	filename := uuid.New().String() + "-" + header.Filename
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	resultOfOrigin := make(chan UploadStatus)
	resultOfMedium := make(chan UploadStatus)
	resultOfThumb := make(chan UploadStatus)

	go processor.uploadOriginal(filename, fileBytes, resultOfOrigin)
	go processor.uploadMedium(filename, fileBytes, resultOfMedium)
	go processor.uploadThumb(filename, fileBytes, resultOfThumb)

	orgStatus := <-resultOfOrigin
	mediumStatus := <-resultOfMedium
	thumbStatus := <-resultOfThumb

	if orgStatus.Err != nil {
		return nil, orgStatus.Err
	}
	if mediumStatus.Err != nil {
		return nil, mediumStatus.Err
	}
	if thumbStatus.Err != nil {
		return nil, thumbStatus.Err
	}

	images := &model.UploadedImages{
		Original:  orgStatus.Url,
		Medium:    mediumStatus.Url,
		Thumbnail: thumbStatus.Url,
	}
	return images, nil
}

func (processor *FileProcessor) uploadOriginal(filename string, fileBytes []byte, result chan UploadStatus) {

	url, err := processor.FileUploader.Upload(filename, fileBytes)
	if err != nil {
		result <- UploadStatus{
			Url: "",
			Err: err,
		}
		return
	}
	result <- UploadStatus{
		Url: url,
		Err: nil,
	}
}

func (processor *FileProcessor) uploadMedium(filename string, fileBytes []byte, result chan UploadStatus) {

	filename = "m/" + filename
	resized, err := processor.Resizer.Resize(fileBytes, 500)
	if err != nil {
		result <- UploadStatus{
			Url: "",
			Err: err,
		}
	}
	url, err := processor.FileUploader.Upload(filename, resized)
	if err != nil {
		result <- UploadStatus{
			Url: "",
			Err: err,
		}
		return
	}
	result <- UploadStatus{
		Url: url,
		Err: nil,
	}
}

func (processor *FileProcessor) uploadThumb(filename string, fileBytes []byte, result chan UploadStatus) {

	filename = "s/" + filename
	resized, err := processor.Resizer.Resize(fileBytes, 100)
	if err != nil {
		result <- UploadStatus{
			Url: "",
			Err: err,
		}
	}
	url, err := processor.FileUploader.Upload(filename, resized)
	if err != nil {
		result <- UploadStatus{
			Url: "",
			Err: err,
		}
		return
	}
	result <- UploadStatus{
		Url: url,
		Err: nil,
	}
}
