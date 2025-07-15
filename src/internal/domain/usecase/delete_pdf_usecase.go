package usecase

import (
	"context"
	"sen-global-api/internal/data/repository"
	"sen-global-api/pkg/uploader"
)

type DeletePDFUseCase struct {
	uploader.UploadProvider
	*repository.PdfRepository
}

func (receiver *DeletePDFUseCase) DeletePDF(key string) error {
	
	pdf, err := receiver.PdfRepository.GetByKey(key)
	if err != nil {
		return err
	}

	// Delete from S3
	err = receiver.UploadProvider.DeleteFileUploaded(context.Background(), pdf.Key)
	if err != nil {
		return err
	}

	// Delete from DB
	err = receiver.PdfRepository.DeletePDF(pdf.Key)
	if err != nil {
		return err
	}

	return nil
}