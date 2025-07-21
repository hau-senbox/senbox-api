package usecase

import (
	"context"
	"fmt"
	"path"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/uploader"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type UploadPDFUseCase struct {
	uploader.UploadProvider
	*repository.PdfRepository
}

func (receiver *UploadPDFUseCase) UploadPDF(data []byte, folder, fileName, pdfName string, mode uploader.UploadMode, ogrID string) (*string, *entity.SPdf, error) {
	fileExt := strings.ToLower(path.Ext(fileName))

	if !isValidHandler(fileExt) {
		return nil, nil, fmt.Errorf("file extension %s is not supported", fileExt)
	}

	// Generate the new filename
	timestamp := time.Now().UnixNano()
	finalFileName := fmt.Sprintf("%s_%d%s", pdfName, timestamp, fileExt)

	if strings.TrimSpace(folder) == "" {
		folder = "pdf"
	}

	// Determine dimensions only for raster images
	// Upload the image
	uploadPath := fmt.Sprintf("%s/%s", folder, finalFileName)
	url, err := receiver.UploadProvider.SaveFileUploaded(context.Background(), data, uploadPath, mode)
	if err != nil {
		log.Errorf("error uploading file to S3: %v", err)
		return nil, nil, err
	}

	pdf := &entity.SPdf{
		PdfName:        pdfName,
		Folder:         folder,
		OrganizationID: ogrID,
		Key:            uploadPath,
		Extension:      fileExt,
	}

	err = receiver.PdfRepository.Save(pdf)
	if err != nil {
		return nil, nil, err
	}

	return url, pdf, nil

}

func isValidHandler(extName string) bool {
	extName = strings.ToLower(extName)
	if extName == ".pdf" {
		return true
	}
	return false
}
