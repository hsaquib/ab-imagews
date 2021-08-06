package request

import "mime/multipart"

type FileInfo struct {
	File     multipart.File
	FileName string
}
