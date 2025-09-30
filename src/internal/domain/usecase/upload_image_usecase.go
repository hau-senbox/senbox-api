package usecase

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/uploader"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type UploadImageUseCase struct {
	uploader.UploadProvider
	*repository.ImageRepository
}

func getImageDimensions(data []byte) (int, int, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		log.Println("Error decoding image: ", err)
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}

var (
	supportedImageExts       = []string{".jpg", ".jpeg", ".png", ".ico", ".svg", ".bmp", ".gif"}
	supportedRasterImageExts = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".bmp":  true,
		".gif":  true,
	}
)

func isImage(extName string) bool {
	extName = strings.ToLower(extName)
	for _, ext := range supportedImageExts {
		if ext == extName {
			return true
		}
	}
	return false
}

func (receiver *UploadImageUseCase) UploadImage(
	data []byte,
	folder, fileName, imageName string,
	mode uploader.UploadMode,
	topicID *string,
) (*string, *entity.SImage, error) {

	fileExt := strings.ToLower(path.Ext(fileName))

	if !isImage(fileExt) {
		return nil, nil, fmt.Errorf("file extension %s is not supported", fileExt)
	}

	// Generate the new filename
	timestamp := time.Now().UnixNano()
	finalFileName := fmt.Sprintf("%s_%d%s", imageName, timestamp, fileExt)

	if strings.TrimSpace(folder) == "" {
		folder = "img"
	}

	// Determine dimensions only for raster images
	var width, height int
	var err error
	if supportedRasterImageExts[fileExt] {
		width, height, err = getImageDimensions(data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get image dimensions: %w", err)
		}
	}

	// Upload the image
	uploadPath := fmt.Sprintf("%s/%s", folder, finalFileName)
	url, err := receiver.UploadProvider.SaveFileUploaded(context.Background(), data, uploadPath, mode)
	if err != nil {
		log.Errorf("error uploading file to S3: %v", err)
		return nil, nil, err
	}

	// Create SImage entity
	img := entity.SImage{
		ImageName: imageName,
		Folder:    folder,
		Key:       uploadPath,
		Width:     width,
		Height:    height,
		Extension: fileExt,
	}

	// add topic id
	if topicID != nil {
		img.TopicID = *topicID
	} else {
		img.TopicID = ""
	}

	// Save image metadata
	err = receiver.ImageRepository.CreateImage(&img)
	if err != nil {
		return nil, nil, err
	}

	// Optionally create public image
	if mode == uploader.UploadPublic {
		err = receiver.ImageRepository.CreatePublicImage(entity.PublicImage{
			ImageName: imageName,
			Folder:    folder,
			Key:       uploadPath,
			URL:       *url,
			Width:     width,
			Height:    height,
			Extension: fileExt,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	return url, &img, nil
}

func (receiver *UploadImageUseCase) UploadImagev2(
	data []byte,
	folder, fileName, imageName string,
	mode uploader.UploadMode,
) (*string, *entity.SImage, error) {

	fileExt := strings.ToLower(path.Ext(fileName))

	if !isImage(fileExt) {
		return nil, nil, fmt.Errorf("file extension %s is not supported", fileExt)
	}

	// Generate the new filename
	timestamp := time.Now().UnixNano()
	finalFileName := fmt.Sprintf("%s_%d%s", imageName, timestamp, fileExt)

	if strings.TrimSpace(folder) == "" {
		folder = "img"
	}

	// Determine dimensions only for raster images
	var width, height int
	var err error
	if supportedRasterImageExts[fileExt] {
		width, height, err = getImageDimensions(data)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get image dimensions: %w", err)
		}
	}

	// Upload the image
	uploadPath := fmt.Sprintf("%s/%s", folder, finalFileName)
	url, err := receiver.UploadProvider.SaveFileUploaded(context.Background(), data, uploadPath, mode)
	if err != nil {
		log.Errorf("error uploading file to S3: %v", err)
		return nil, nil, err
	}

	// Create SImage entity
	img := entity.SImage{
		ImageName: imageName,
		Folder:    folder,
		Key:       uploadPath,
		Width:     width,
		Height:    height,
		Extension: fileExt,
	}

	// Save image metadata
	err = receiver.ImageRepository.CreateImage(&img)
	if err != nil {
		return nil, nil, err
	}

	// Optionally create public image
	if mode == uploader.UploadPublic {
		err = receiver.ImageRepository.CreatePublicImage(entity.PublicImage{
			ImageName: imageName,
			Folder:    folder,
			Key:       uploadPath,
			URL:       *url,
			Width:     width,
			Height:    height,
			Extension: fileExt,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	return url, &img, nil
}

func (receiver *UploadImageUseCase) UploadImages(
	files []struct {
		Data      []byte
		FileName  string
		ImageName string
	},
	folder string,
	mode uploader.UploadMode,
) (*response.UploadImagesResponse, error) {
	// Khởi tạo response rỗng
	results := &response.UploadImagesResponse{
		Images: make([]response.ImageResponse, 0),
	}

	for _, f := range files {
		url, img, err := receiver.UploadImagev2(
			f.Data, folder, f.FileName, f.ImageName, mode,
		)
		if err != nil {
			return nil, err
		}

		results.Images = append(results.Images, response.ImageResponse{
			ImageName: img.ImageName,
			Key:       img.Key,
			Extension: img.Extension,
			Url:       *url,
			Width:     img.Width,
			Height:    img.Height,
		})
	}

	return results, nil
}
