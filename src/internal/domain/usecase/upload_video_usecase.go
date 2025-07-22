package usecase

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/uploader"
	"strings"
	"time"
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
