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

type UploadVideoUseCase struct {
	uploader.UploadProvider
	*repository.VideoRepository
}

var supportedVideoExts = []string{".mp4", ".mov", ".avi", ".wmv", ".flv", ".webm", ".mkv"}

func isVideo(ext string) bool {
	ext = strings.ToLower(ext)
	for _, v := range supportedVideoExts {
		if v == ext {
			return true
		}
	}
	return false
}

func (receiver *UploadVideoUseCase) UploadVideo(data []byte, folder, fileName, videoName string, mode uploader.UploadMode) (*string, *entity.SVideo, error) {
	fileExt := strings.ToLower(path.Ext(fileName))

	if !isVideo(fileExt) {
		return nil, nil, fmt.Errorf("file extension %s is not supported", fileExt)
	}

	// Generate the new filename
	timestamp := time.Now().UnixNano()
	finalFileName := fmt.Sprintf("%s_%d%s", videoName, timestamp, fileExt)

	if strings.TrimSpace(folder) == "" {
		folder = "videos"
	}

	// Upload the video
	uploadPath := fmt.Sprintf("%s/%s", folder, finalFileName)
	url, err := receiver.UploadProvider.SaveFileUploaded(context.Background(), data, uploadPath, mode)
	if err != nil {
		log.Errorf("error uploading file to S3: %v", err)
		return nil, nil, err
	}

	// Save video metadata
	video := entity.SVideo{
		VideoName: videoName,
		Folder:    folder,
		Key:       uploadPath,
		Extension: fileExt,
	}

	err = receiver.VideoRepository.CreateVideo(video)
	if err != nil {
		return nil, nil, err
	}

	return url, &video, nil
}

func (uc *UploadVideoUseCase) UploadVideoV2(data []byte, req request.UploadVideoRequest) (*response.VideoResponse, error) {
	// Chuẩn hóa folder, fileName, videoName
	req.Folder = helper.SanitizeName(req.Folder)
	req.FileName = helper.SanitizeName(req.FileName)
	req.VideoName = helper.SanitizeName(req.VideoName)

	// Lấy extension từ file upload
	fileExt := strings.ToLower(path.Ext(req.FileName))
	if fileExt == "" && req.File != nil {
		fileExt = strings.ToLower(path.Ext(req.File.Filename))
	}

	if !isVideo(fileExt) {
		return nil, fmt.Errorf("file extension %s is not supported", fileExt)
	}

	// Nếu folder trống → mặc định là "video"
	if req.Folder == "" {
		req.Folder = "video"
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

	// Save metadata
	video := entity.SVideo{
		VideoName: req.FileName,
		Folder:    req.Folder,
		Key:       uploadPath,
		Extension: fileExt,
	}
	if err := uc.VideoRepository.CreateVideo(video); err != nil {
		return nil, err
	}

	// Response
	res := &response.VideoResponse{
		VideoName: video.VideoName,
		Key:       video.Key,
		Extension: video.Extension,
		Url:       *url,
	}
	return res, nil
}
