package usecase

import (
	"context"
	"fmt"
	"path"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/uploader"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type UploadPDFUseCase struct {
	uploader.UploadProvider
	*repository.PdfRepository
}

func (receiver *UploadPDFUseCase) UploadPDF(data []byte, folder, fileName, pdfName string, mode uploader.UploadMode) (*string, *entity.SPdf, error) {
	fileExt := strings.ToLower(path.Ext(fileName))

	if !isValidPDF(fileExt) {
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
		PdfName: pdfName,
		Folder:  folder,
		//OrganizationID: ogrID,
		Key:       uploadPath,
		Extension: fileExt,
	}

	err = receiver.PdfRepository.Save(pdf)
	if err != nil {
		return nil, nil, err
	}

	return url, pdf, nil

}

func (uc *UploadPDFUseCase) UploadPDFv2(data []byte, req request.UploadPdfRequest) (*response.PdfResponse, error) {
	// Chuẩn hóa dữ liệu đầu vào
	req.Folder = helper.SanitizeName(req.Folder)
	req.FileName = helper.SanitizeName(req.FileName)

	// Lấy extension
	fileExt := strings.ToLower(path.Ext(req.FileName))
	if fileExt == "" && req.File != nil {
		fileExt = strings.ToLower(path.Ext(req.File.Filename))
	}

	// Kiểm tra định dạng hợp lệ
	if !isValidPDF(fileExt) {
		return nil, fmt.Errorf("file extension %s is not supported (only .pdf allowed)", fileExt)
	}

	// Nếu folder trống → mặc định là "pdf"
	if req.Folder == "" {
		req.Folder = "pdf"
	}

	// Sinh tên file mới
	timestamp := time.Now().UnixNano()
	finalFileName := fmt.Sprintf("%s_%d%s", req.FileName, timestamp, fileExt)

	// Chuyển mode string → enum
	mode, err := uploader.UploadModeFromString(req.Mode)
	if err != nil {
		return nil, fmt.Errorf("invalid upload mode: %w", err)
	}

	// Upload file
	uploadPath := fmt.Sprintf("%s/%s", req.Folder, finalFileName)
	url, err := uc.UploadProvider.SaveFileUploaded(context.Background(), data, uploadPath, mode)
	if err != nil {
		log.Errorf("error uploading file to S3: %v", err)
		return nil, err
	}

	// Lưu thông tin PDF vào database
	pdf := &entity.SPdf{
		PdfName:   req.FileName,
		Folder:    req.Folder,
		Key:       uploadPath,
		Extension: fileExt,
	}

	if err := uc.PdfRepository.Save(pdf); err != nil {
		log.Errorf("error saving pdf metadata: %v", err)
		return nil, err
	}

	// Chuẩn bị response
	res := &response.PdfResponse{
		PdfName:   pdf.PdfName,
		Key:       pdf.Key,
		Url:       *url,
		Extension: pdf.Extension,
	}

	return res, nil
}

func isValidPDF(ext string) bool {
	return strings.ToLower(ext) == ".pdf"
}
