package usecase

import (
	"context"
	"github.com/pkg/errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/uploader"
	"time"
)

type GetVideoUseCase struct {
	uploader.UploadProvider
	*repository.VideoRepository
}

func (receiver *GetVideoUseCase) GetAllByIDs(ids []int) ([]entity.SVideo, error) {
	return receiver.VideoRepository.GetAllByIDs(ids)
}

func (receiver *GetVideoUseCase) GetAllByName(imageName string) ([]entity.SVideo, error) {
	return receiver.VideoRepository.GetAllByName(imageName)
}

func (receiver *GetVideoUseCase) GetVideoByID(id uint64) (*entity.SVideo, error) {
	return receiver.VideoRepository.GetByID(id)
}

func (receiver *GetVideoUseCase) GetUrlByKey(key string, mode uploader.UploadMode) (*string, error) {
	video, err := receiver.VideoRepository.GetByKey(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case uploader.UploadPrivate:
		return receiver.GetFileUploaded(context.Background(), video.Key, nil)
	case uploader.UploadPublic:
		duration := time.Now().AddDate(10, 0, 0).Sub(time.Now())
		return receiver.GetFileUploaded(context.Background(), video.Key, &duration)
	default:
		return nil, errors.New("invalid upload mode")
	}
}
