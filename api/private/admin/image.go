package admin

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hsaquib/ab-imagews/utils"
	"mime/multipart"
	"net/http"
)

type FileInfo struct {
	MimeType string
	FileName string
	FileSize int64
}

func (ar *adminRouter) imageRouter() http.Handler {
	r := chi.NewRouter()
	reqHandler := NewApiHandler(ar.ServiceProvider, ar.Log)

	r.Post("/upload", reqHandler.UploadImage)
	return r
}

func (handler *apiHandler) UploadImage(writer http.ResponseWriter, request *http.Request) {
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

	images, err := handler.ServiceProvider.FileProcessor.UploadImageVariants(file, fileHeader)
	if err != nil {
		utils.HandleObjectError(writer, err)
		return
	}

	err = utils.ServeJSONObject(writer, http.StatusOK, "Successful", images, nil, true)
	if err != nil {
		fmt.Println(err.Error())
	}
}
