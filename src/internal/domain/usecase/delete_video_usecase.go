package usecase

import (
	"context"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/pkg/uploader"
)

type DeleteVideoUseCase struct {
	uploader.UploadProvider
	*repository.VideoRepository
}

func (receiver *DeleteVideoUseCase) DeleteVideo(key string) error {
	// Fetch video metadata from DB
	videoData, err := receiver.VideoRepository.GetByKey(key)
	if err != nil {
		return fmt.Errorf("failed to get video by key: %w", err)
	}

	// Delete from S3
	err = receiver.UploadProvider.DeleteFileUploaded(context.Background(), videoData.Key)
	if err != nil {
		return err
	}

	// Delete from DB
	err = receiver.VideoRepository.DeleteVideo(videoData.ID)
	if err != nil {
		return err
	}

	return nil
}
