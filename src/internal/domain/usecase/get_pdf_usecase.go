package usecase

import (
	"context"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/uploader"
	"time"
)

type GetPdfByKeyUseCase struct {
	uploader.UploadProvider
	*repository.PdfRepository
}

func (u *GetPdfByKeyUseCase) GetPdfByKey(key string, mode uploader.UploadMode) (*response.PdfResponse, error) {

	pdf, err := u.PdfRepository.GetByKey(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case uploader.UploadPrivate:
		url, err := u.GetFileUploaded(context.Background(), pdf.Key, nil)
		if err != nil {
			return nil, err
		}

		return &response.PdfResponse{
			Url:            *url,
			PdfName:        pdf.PdfName,
			OrganizationID: pdf.OrganizationID,
			Key:            pdf.Key,
			Extension:      pdf.Extension,
		}, nil
	case uploader.UploadPublic:
		duration := time.Now().AddDate(10, 0, 0).Sub(time.Now())
		url, err := u.GetFileUploaded(context.Background(), pdf.Key, &duration)
		if err != nil {
			return nil, err
		}
		return &response.PdfResponse{
			Url:            *url,
			PdfName:        pdf.PdfName,
			OrganizationID: pdf.OrganizationID,
			Key:            pdf.Key,
			Extension:      pdf.Extension,
		}, nil
	default:
		return nil, errors.New("invalid upload mode")
	}
}

func (u *GetPdfByKeyUseCase) GetAllKeyByOrgID(orgID string) ([]*response.PdfResponse, error) {

	keyPdf, err := u.PdfRepository.GetAllKeyByOrgID(orgID)
	if err != nil {
		return nil, err
	}

	var pdfs []*response.PdfResponse
	for _, key := range keyPdf {
		pdf, err := u.GetPdfByKey(key, uploader.UploadPublic)
		if err != nil {
			return nil, err
		}
		pdfs = append(pdfs, pdf)
	}

	return pdfs, nil
}
