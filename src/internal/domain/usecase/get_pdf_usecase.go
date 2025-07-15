package usecase

import (
	"context"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/pkg/uploader"
	"time"
)

type GetPdfByKeyUseCase struct {
	uploader.UploadProvider
	*repository.PdfRepository
}

func (u *GetPdfByKeyUseCase) GetPdfByKey(key string, mode uploader.UploadMode) (*string, error) {

	pdf, err := u.PdfRepository.GetByKey(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case uploader.UploadPrivate:
		return u.GetFileUploaded(context.Background(), pdf.Key, nil)
	case uploader.UploadPublic:
		duration := time.Now().AddDate(10, 0, 0).Sub(time.Now())
		return u.GetFileUploaded(context.Background(), pdf.Key, &duration)
	default:
		return nil, errors.New("invalid upload mode")
	}
}

func (u *GetPdfByKeyUseCase) GetAllKeyByOrgID(orgID int64) ([]string, error) {
	return u.PdfRepository.GetAllKeyByOrgID(orgID)
}
