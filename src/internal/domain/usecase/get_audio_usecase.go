package usecase

import (
	"context"
	"github.com/pkg/errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/uploader"
	"time"
)

type GetAudioUseCase struct {
	uploader.UploadProvider
	*repository.AudioRepository
}

func (receiver *GetAudioUseCase) GetAllByIDs(ids []int) ([]entity.SAudio, error) {
	return receiver.AudioRepository.GetAllByIDs(ids)
}

func (receiver *GetAudioUseCase) GetAllByName(imageName string) ([]entity.SAudio, error) {
	return receiver.AudioRepository.GetAllByName(imageName)
}

func (receiver *GetAudioUseCase) GetAudioByID(id uint64) (*entity.SAudio, error) {
	return receiver.AudioRepository.GetByID(id)
}

func (receiver *GetAudioUseCase) GetUrlByKey(key string, mode uploader.UploadMode) (*string, error) {
	audio, err := receiver.AudioRepository.GetByKey(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case uploader.UploadPrivate:
		return receiver.GetFileUploaded(context.Background(), audio.Key, nil)
	case uploader.UploadPublic:
		duration := time.Now().AddDate(10, 0, 0).Sub(time.Now())
		return receiver.GetFileUploaded(context.Background(), audio.Key, &duration)
	default:
		return nil, errors.New("invalid upload mode")
	}
}
