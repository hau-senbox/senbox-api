package request

import "mime/multipart"

type UploadAudioRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Folder    string                `form:"folder" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	AudioName string                `form:"audio_name"`
	Mode      string                `form:"mode" binding:"required"`
}
