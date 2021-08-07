package service

import (
	"github.com/google/uuid"
	"github.com/hsaquib/ab-imagews/config"
	req "github.com/hsaquib/ab-imagews/dto/request"
	"github.com/hsaquib/ab-imagews/dto/response"
	rLog "github.com/hsaquib/rest-log"
	"io/ioutil"
	"mime/multipart"
	"sync"
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

type BatchUploadStatus struct {
	UploadedImage *response.UploadedImages
	Err           error
}

func NewFileProcessor(conf *config.AppConfig, logger rLog.Logger) *FileProcessor {
	return &FileProcessor{
		Config:       conf,
		Log:          logger,
		Resizer:      NewVipsThumb(logger),
		FileUploader: NewUploader(&conf.Upload, logger),
	}
}

func (processor *FileProcessor) UploadImageVariants(file multipart.File, header *multipart.FileHeader) (*response.UploadedImages, error) {

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

	images := &response.UploadedImages{
		FileName:  header.Filename,
		Original:  orgStatus.Url,
		Medium:    mediumStatus.Url,
		Thumbnail: thumbStatus.Url,
	}
	close(resultOfOrigin)
	close(resultOfMedium)
	close(resultOfThumb)
	return images, nil
}

func (processor *FileProcessor) UploadVariantsOfImageList(fileList []req.FileInfo) ([]*response.UploadedImages, error) {

	var imageList []*response.UploadedImages

	wg := new(sync.WaitGroup)
	resultChan := make(chan BatchUploadStatus)
	for _, info := range fileList {
		wg.Add(1)
		fileBytes, err := ioutil.ReadAll(info.File)
		if err != nil {
			return nil, err
		}
		go processor.uploadVariants(info.FileName, fileBytes, resultChan, wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		//deal with the result in some way
		if result.Err != nil {
			return nil, result.Err
		}
		imageList = append(imageList, result.UploadedImage)
	}

	return imageList, nil
}

func (processor *FileProcessor) uploadVariants(filename string, fileBytes []byte, ch chan BatchUploadStatus, wg *sync.WaitGroup) {
	defer wg.Done()
	reqFileName := filename
	filename = uuid.New().String() + "-" + filename

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
		ch <- BatchUploadStatus{
			UploadedImage: nil,
			Err:           orgStatus.Err,
		}
		return
	}
	if mediumStatus.Err != nil {
		ch <- BatchUploadStatus{
			UploadedImage: nil,
			Err:           mediumStatus.Err,
		}
		return
	}
	if thumbStatus.Err != nil {
		ch <- BatchUploadStatus{
			UploadedImage: nil,
			Err:           thumbStatus.Err,
		}
		return
	}

	images := &response.UploadedImages{
		FileName:  reqFileName,
		Original:  orgStatus.Url,
		Medium:    mediumStatus.Url,
		Thumbnail: thumbStatus.Url,
	}
	ch <- BatchUploadStatus{
		UploadedImage: images,
		Err:           nil,
	}
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
