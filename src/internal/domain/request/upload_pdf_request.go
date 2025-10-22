package request

import "mime/multipart"

type UploadPdfRequest struct {
	File     *multipart.FileHeader `form:"file" binding:"required"`
	Folder   string                `form:"folder" binding:"required"`
	FileName string                `form:"file_name" binding:"required"`
	Mode     string                `form:"mode" binding:"required"`
}
