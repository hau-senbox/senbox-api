package usecase

import (
	"context"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/pkg/uploader"
)

type DeleteImageUseCase struct {
	uploader.UploadProvider
	*repository.ImageRepository
}

func (receiver *DeleteImageUseCase) DeleteImage(key string) error {
	// Fetch image metadata from DB
	imageData, err := receiver.ImageRepository.GetByKey(key)
	if err != nil {
		return fmt.Errorf("failed to get image by key: %w", err)
	}

	// Delete from S3
	err = receiver.UploadProvider.DeleteFileUploaded(context.Background(), imageData.Key)
	if err != nil {
		return err
	}

	// Delete from DB
	err = receiver.ImageRepository.DeleteImage(imageData.ID)
	if err != nil {
		return err
	}

	return nil
}
