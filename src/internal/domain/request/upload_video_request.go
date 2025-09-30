package request

import "mime/multipart"

type UploadVideoRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Folder    string                `form:"folder" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	VideoName string                `form:"video_name"`
	Mode      string                `form:"mode" binding:"required"`
}
