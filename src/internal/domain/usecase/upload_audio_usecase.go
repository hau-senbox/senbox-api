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
