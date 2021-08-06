package service

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hsaquib/ab-imagews/config"
	rLog "github.com/hsaquib/rest-log"
	"net/http"
	"strings"
)

type Uploader interface {
	Upload(filename string, fileBytes []byte) (string, error)
}

type S3FileUploader struct {
	Config *config.UploadConfig
	Log    rLog.Logger
}

func NewUploader(conf *config.UploadConfig, logger rLog.Logger) Uploader {
	return &S3FileUploader{
		Config: conf,
		Log:    logger,
	}
}

func (s *S3FileUploader) Upload(filename string, fileBytes []byte) (string, error) {
	mimeType := http.DetectContentType(fileBytes)
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	fmt.Println("new session")
	// Create an Uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	s.Log.Info("", "", "uploading file "+filename)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.Config.S3Bucket),
		Key:         aws.String(joinPath(s.Config.Folder, filename)),
		ContentType: aws.String(mimeType),
		Body:        bytes.NewBuffer(fileBytes),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	s.Log.Info("", "", "file uploaded: "+aws.StringValue(&result.Location))
	return s.Config.BaseUrl + filename, nil
}

func joinPath(joinWith string, joinPart string) string {

	path := ""
	if len(strings.TrimSpace(joinWith)) != 0 {
		path = joinWith + "/" + joinPart
	} else {
		path = joinPart
	}
	return path
}
