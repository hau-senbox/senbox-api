package usecase

import (
	"context"
	"github.com/pkg/errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/uploader"
	"time"
)

type GetImageUseCase struct {
	uploader.UploadProvider
	*repository.ImageRepository
}

func (receiver *GetImageUseCase) GetAllByIDs(ids []int) ([]entity.SImage, error) {
	return receiver.ImageRepository.GetAllByIDs(ids)
}

func (receiver *GetImageUseCase) GetAllByName(imageName string) ([]entity.SImage, error) {
	return receiver.ImageRepository.GetAllByName(imageName)
}

func (receiver *GetImageUseCase) GetImageByID(id uint64) (*entity.SImage, error) {
	return receiver.ImageRepository.GetByID(id)
}

func (receiver *GetImageUseCase) GetIcons() ([]entity.PublicImage, error) {
	return receiver.ImageRepository.GetIcons()
}

func (receiver *GetImageUseCase) GetIconByKey(key string) (*entity.PublicImage, error) {
	return receiver.ImageRepository.GetIconByKey(key)
}

func (receiver *GetImageUseCase) GetUrlByKey(key string, mode uploader.UploadMode) (*string, error) {
	img, err := receiver.ImageRepository.GetByKey(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case uploader.UploadPrivate:
		return receiver.GetFileUploaded(context.Background(), img.Key, nil)
	case uploader.UploadPublic:
		duration := time.Now().AddDate(10, 0, 0).Sub(time.Now())
		return receiver.GetFileUploaded(context.Background(), img.Key, &duration)
	default:
		return nil, errors.New("invalid upload mode")
	}
}
