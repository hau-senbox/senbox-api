package controller

import (
	"bufio"
	"fmt"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/randx"
	"sen-global-api/pkg/uploader"
	"strings"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	*usecase.GetImageUseCase
	*usecase.UploadImageUseCase
	*usecase.DeleteImageUseCase
	*usecase.UserImagesUsecase
}

type getUrlByKeyRequest struct {
	Key  string `json:"key" binding:"required"`
	Mode string `json:"mode" binding:"required"`
}

func (receiver *ImageController) GetUrlByKey(context *gin.Context) {
	var req getUrlByKeyRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	mode, err := uploader.UploadModeFromString(req.Mode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	url, err := receiver.GetImageUseCase.GetUrlByKey(req.Key, mode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "image was get successfully",
		Data:    *url,
	})
}

type iconResponse struct {
	ImageName string `json:"image_name"`
	Folder    string `json:"folder"`
	Key       string `json:"key"`
	URL       string `json:"url"`
	Extension string `json:"extension"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func (receiver *ImageController) GetIcons(context *gin.Context) {
	icons, err := receiver.GetImageUseCase.GetIcons()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	res := make([]iconResponse, 0)
	for _, icon := range icons {
		res = append(res, iconResponse{
			ImageName: icon.ImageName,
			Folder:    icon.Folder,
			Key:       icon.Key,
			URL:       icon.URL,
			Extension: icon.Extension,
			Width:     icon.Width,
			Height:    icon.Height,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "icons were get successfully",
		Data:    res,
	})
}

type getPublicImageByKeyRequest struct {
	Key string `json:"key" binding:"required"`
}

func (receiver *ImageController) GetIconByKey(context *gin.Context) {
	var req getPublicImageByKeyRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	icon, err := receiver.GetImageUseCase.GetIconByKey(req.Key)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "icons were get successfully",
		Data: iconResponse{
			ImageName: icon.ImageName,
			Folder:    icon.Folder,
			Key:       icon.Key,
			URL:       icon.URL,
			Extension: icon.Extension,
			Width:     icon.Width,
			Height:    icon.Height,
		},
	})
}

func (receiver *ImageController) CreateImage(context *gin.Context) {
	// userIDRaw, exists := context.Get("user_id")
	// if !exists {
	// 	context.JSON(http.StatusUnauthorized, response.FailedResponse{
	// 		Code:  http.StatusUnauthorized,
	// 		Error: "Unauthorized: user_id not found",
	// 	})
	// 	return
	// }

	// userID, ok := userIDRaw.(string)
	// if !ok {
	// 	context.JSON(http.StatusUnauthorized, response.FailedResponse{
	// 		Code:  http.StatusUnauthorized,
	// 		Error: "Unauthorized: user_id is not a valid string",
	// 	})
	// 	return
	// }

	fileHeader, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// folder := context.DefaultPostForm("folder", "img")
	// fileName := context.DefaultPostForm("file_name", randx.GenString(10))
	folder := strings.TrimSpace(context.DefaultPostForm("folder", "img"))
	fileName := strings.TrimSpace(context.DefaultPostForm("file_name", randx.GenString(10)))

	// Replace all whitespace sequences (tabs, spaces, multiple spaces, etc.) with "_"
	folder = strings.Join(strings.Fields(folder), "_")
	fileName = strings.Join(strings.Fields(fileName), "_")
	mode, err := uploader.UploadModeFromString(context.PostForm("mode"))
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	defer file.Close()

	dataBytes := make([]byte, fileHeader.Size)
	if _, err := bufio.NewReader(file).Read(dataBytes); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// Lấy topic_id và student_id nếu có
	topicIDRaw := context.PostForm("topic_id")
	var topicID *string
	if strings.TrimSpace(topicIDRaw) != "" {
		topicID = &topicIDRaw
	}

	// studentIDRaw := context.PostForm("student_id")
	// var studentID *string
	// if strings.TrimSpace(studentIDRaw) != "" {
	// 	studentID = &studentIDRaw
	// }

	// teacherIDRaw := context.PostForm("teacher_id")
	// var teacherID *string
	// if strings.TrimSpace(teacherIDRaw) != "" {
	// 	teacherID = &teacherIDRaw
	// }

	url, img, err := receiver.UploadImageUseCase.UploadImage(dataBytes, folder, fileHeader.Filename, fileName, mode, topicID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if url == nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "image was not created",
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "image was create successfully",
		Data: response.ImageResponse{
			ImageName: img.ImageName,
			Key:       img.Key,
			Extension: img.Extension,
			Url:       *url,
			Width:     img.Width,
			Height:    img.Height,
		},
	})
}

func (receiver *ImageController) DeleteImage(context *gin.Context) {
	var req request.DeleteImageRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "key is required",
		})
		return
	}

	err := receiver.DeleteImageUseCase.DeleteImage(req.Key)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "image deleted successfully",
	})
}

func (receiver *ImageController) GetUrlIsMain4Owner(context *gin.Context) {
	var req request.GetUrlIsMain4OwnerRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// check owner role valid
	if !value.OwnerRole(req.OwnerRole).IsValid() {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid owner role",
		})
		return
	}

	url, err := receiver.UserImagesUsecase.GetUrlIsMain4Owner(req.OwnerID, value.OwnerRole(req.OwnerRole))
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "image was get successfully",
		Data:    url,
	})
}

func (receiver *ImageController) CreateImages(context *gin.Context) {
	form, err := context.MultipartForm()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "failed to parse multipart form: " + err.Error(),
		})
		return
	}

	files := form.File["files"] // client phải gửi field name = "files"
	if len(files) == 0 {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "no files uploaded",
		})
		return
	}

	folder := strings.TrimSpace(context.DefaultPostForm("folder", "img"))
	fileName := strings.TrimSpace(context.DefaultPostForm("file_name", randx.GenString(10)))

	// Replace all whitespace sequences (tabs, spaces, multiple spaces, etc.) with "_"
	folder = strings.Join(strings.Fields(folder), "_")
	fileName = strings.Join(strings.Fields(fileName), "_")

	mode, err := uploader.UploadModeFromString(context.PostForm("mode"))
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// Chuẩn bị data cho UploadImages
	fileInputs := make([]struct {
		Data      []byte
		FileName  string
		ImageName string
	}, 0)

	for idx, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			context.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: fmt.Sprintf("failed to open file %s: %v", fileHeader.Filename, err),
			})
			return
		}
		defer file.Close()

		dataBytes := make([]byte, fileHeader.Size)
		if _, err := bufio.NewReader(file).Read(dataBytes); err != nil {
			context.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: fmt.Sprintf("failed to read file %s: %v", fileHeader.Filename, err),
			})
			return
		}

		imageName := fmt.Sprintf("%s_%d", fileName, idx+1)
		fileInputs = append(fileInputs, struct {
			Data      []byte
			FileName  string
			ImageName string
		}{
			Data:      dataBytes,
			FileName:  fileHeader.Filename,
			ImageName: imageName,
		})
	}

	// Gọi usecase
	result, err := receiver.UploadImageUseCase.UploadImages(fileInputs, folder, mode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("%d images uploaded successfully", len(result.Images)),
		Data:    result,
	})
}
