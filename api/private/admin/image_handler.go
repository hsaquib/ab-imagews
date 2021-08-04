package admin

import (
	"fmt"
	"github.com/hsaquib/ab-imagews/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type FileInfo struct {
	MimeType string
	FileName string
	FileSize int64
}

func UploadImage(writer http.ResponseWriter, request *http.Request) {

	// ParseMultipartForm parses a request body as multipart/form-data
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		utils.HandleObjectError(writer, err)
		return
	}

	file, fileHeader, err := request.FormFile("file") // Retrieve the file from form data

	if err != nil {
		utils.HandleObjectError(writer, err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			utils.HandleObjectError(writer, err)
			return
		}
	}(file) // Close the file when we finish

	// This is path which we want to store the file
	f, err := os.OpenFile("/Users/hasibussaquib/IdeaProjects/cp-"+fileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		utils.HandleObjectError(writer, err)
		return
	}

	// Copy the file to the destination path
	written, err := io.Copy(f, file)
	if err != nil {
		utils.HandleObjectError(writer, err)
		return
	}
	fmt.Println(written)

	info := FileInfo{
		MimeType: fmt.Sprintf("%s", fileHeader.Header),
		FileName: fileHeader.Filename,
		FileSize: fileHeader.Size,
	}

	utils.ServeJSONObject(writer, http.StatusOK, "Successful", info, nil, true)
	return
}
