package handler

import (
	"fmt"
	req "github.com/hsaquib/ab-imagews/dto/request"
	"github.com/hsaquib/ab-imagews/dto/response"
	"github.com/hsaquib/ab-imagews/service"
	"github.com/hsaquib/ab-imagews/utils"
	rLog "github.com/hsaquib/rest-log"
	"mime/multipart"
	"net/http"
)

type apiHandler struct {
	ServiceProvider *service.Provider
	Log             rLog.Logger
}

func NewApiHandler(srvProvider *service.Provider, log rLog.Logger) *apiHandler {
	return &apiHandler{
		ServiceProvider: srvProvider,
		Log:             log,
	}
}

// UploadImageWithVariantsByAdmin godoc
// @Summary upload image
// @Description Upload an image with three variants: Original, Medium(500X500) & ThumbNail(100X100)
// @Tags Image
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  file formData file true "file is mandatory"
// @Success 200 {object} response.ImageVariantUploadSuccessResponse
// @Failure 400 {object} response.EmptyErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyErrorRes "API sever or db unreachable."
// @Router /api/v1/private/admin/image/upload [post]
func (handler *apiHandler) UploadImageWithVariantsByAdmin(writer http.ResponseWriter, request *http.Request) {

	handler.uploadImageVariants(writer, request)
}

// UploadImageWithVariantsByMerchant godoc
// @Summary upload image
// @Description Upload an image with three variants: Original, Medium(500X500) & ThumbNail(100X100)
// @Tags Image
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  file formData file true "Some fields are mandatory"
// @Success 200 {object} response.ImageVariantUploadSuccessResponse
// @Failure 400 {object} response.EmptyErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyErrorRes "API sever or db unreachable."
// @Router /api/v1/private/merchant/image/upload [post]
func (handler *apiHandler) UploadImageWithVariantsByMerchant(writer http.ResponseWriter, request *http.Request) {

	handler.uploadImageVariants(writer, request)
}

func (handler *apiHandler) uploadImageVariants(writer http.ResponseWriter, request *http.Request) {
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

func (handler *apiHandler) UploadImageListWithVariantsByAdmin(writer http.ResponseWriter, request *http.Request) {
	handler.uploadImageVariants(writer, request)
}

func (handler *apiHandler) UploadImageListWithVariantsByMerchant(writer http.ResponseWriter, request *http.Request) {
	handler.uploadImageVariants(writer, request)
}

func (handler *apiHandler) uploadImageListWithVariants(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}
	var fileList []req.FileInfo
	files := request.MultipartForm.File["files"]
	for i, header := range files {
		file, err := files[i].Open()
		if err != nil {
			return
		}
		defer file.Close()
		fileInfo := req.FileInfo{
			File:     file,
			FileName: header.Filename,
		}
		fileList = append(fileList, fileInfo)
	}
	list, err := handler.ServiceProvider.FileProcessor.UploadVariantsOfImageList(fileList)
	if err != nil {
		utils.HandleListError(writer, err)
	}
	meta := &response.ListMeta{
		Page:  1,
		Pages: 1,
		Limit: int64(len(fileList)),
		Count: int64(len(fileList)),
	}
	utils.ServeJSONList(writer, http.StatusOK, "Successful", list, meta, true)
}
