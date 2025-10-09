package usecase

import (
	"context"
	"fmt"
	"path"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/uploader"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type UploadAudioUseCase struct {
	uploader.UploadProvider
	*repository.AudioRepository
}

var supportedAudioExts = []string{".mp3", ".wav", ".aac", ".ogg", ".flac"}

func isAudio(extName string) bool {
	extName = strings.ToLower(extName)
	for _, ext := range supportedAudioExts {
		if ext == extName {
			return true
		}
	}
	return false
}

func (receiver *UploadAudioUseCase) UploadAudio(data []byte, folder, fileName, videoName string, mode uploader.UploadMode) (*string, *entity.SAudio, error) {
	fileExt := strings.ToLower(path.Ext(fileName))

	if !isAudio(fileExt) {
		return nil, nil, fmt.Errorf("file extension %s is not supported", fileExt)
	}

	// Generate the new filename
	timestamp := time.Now().UnixNano()
	finalFileName := fmt.Sprintf("%s_%d%s", videoName, timestamp, fileExt)

	if strings.TrimSpace(folder) == "" {
		folder = "audio"
	}

	// Upload the video
	uploadPath := fmt.Sprintf("%s/%s", folder, finalFileName)
	url, err := receiver.UploadProvider.SaveFileUploaded(context.Background(), data, uploadPath, mode)
	if err != nil {
		log.Errorf("error uploading file to S3: %v", err)
		return nil, nil, err
	}

	// Save video metadata
	audio := entity.SAudio{
		AudioName: videoName,
		Folder:    folder,
		Key:       uploadPath,
		Extension: fileExt,
	}

	err = receiver.AudioRepository.CreateAudio(audio)
	if err != nil {
		return nil, nil, err
	}

	return url, &audio, nil
}

func (uc *UploadAudioUseCase) UploadAudiov2(data []byte, req request.UploadAudioRequest) (*response.AudioResponse, error) {
	// Chuẩn hóa folder, fileName, audioName
	req.Folder = helper.SanitizeName(req.Folder)
	req.FileName = helper.SanitizeName(req.FileName)
	req.AudioName = helper.SanitizeName(req.AudioName)

	// Lấy extension từ file
	fileExt := strings.ToLower(path.Ext(req.FileName))
	if fileExt == "" && req.File != nil {
		fileExt = strings.ToLower(path.Ext(req.File.Filename))
	}

	if !isAudio(fileExt) {
		return nil, fmt.Errorf("file extension %s is not supported", fileExt)
	}

	// Nếu folder trống → mặc định là "audio"
	if req.Folder == "" {
		req.Folder = "audio"
	}

	// Sinh finalFileName
	timestamp := time.Now().UnixNano()
	finalFileName := fmt.Sprintf("%s_%d%s", req.FileName, timestamp, fileExt)

	// Convert mode string → enum
	mode, err := uploader.UploadModeFromString(req.Mode)
	if err != nil {
		return nil, err
	}

	// Upload
	uploadPath := fmt.Sprintf("%s/%s", req.Folder, finalFileName)
	url, err := uc.UploadProvider.SaveFileUploaded(context.Background(), data, uploadPath, mode)
	if err != nil {
		return nil, err
	}

	// Lưu metadata
	audio := entity.SAudio{
		AudioName: req.AudioName,
		Folder:    req.Folder,
		Key:       uploadPath,
		Extension: fileExt,
	}
	if err := uc.AudioRepository.CreateAudio(audio); err != nil {
		return nil, err
	}

	// Response
	res := &response.AudioResponse{
		AudioName: audio.AudioName,
		Key:       audio.Key,
		Extension: audio.Extension,
		Url:       *url,
	}

	return res, nil
}
