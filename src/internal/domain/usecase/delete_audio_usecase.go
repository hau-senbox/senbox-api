package usecase

import (
	"context"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/pkg/uploader"
)

type DeleteAudioUseCase struct {
	uploader.UploadProvider
	*repository.AudioRepository
}

func (receiver *DeleteAudioUseCase) DeleteAudio(key string) error {
	// Fetch audio metadata from DB
	audioData, err := receiver.AudioRepository.GetByKey(key)
	if err != nil {
		return fmt.Errorf("failed to get audio by key: %w", err)
	}

	// Delete from S3
	err = receiver.UploadProvider.DeleteFileUploaded(context.Background(), audioData.Key)
	if err != nil {
		return err
	}

	// Delete from DB
	err = receiver.AudioRepository.DeleteAudio(audioData.ID)
	if err != nil {
		return err
	}

	return nil
}
