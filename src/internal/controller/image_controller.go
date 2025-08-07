package controller

import (
	"bufio"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/pkg/randx"
	"sen-global-api/pkg/uploader"
	"strings"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	*usecase.GetImageUseCase
	*usecase.UploadImageUseCase
	*usecase.DeleteImageUseCase
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
	userIDRaw, exists := context.Get("user_id")
	if !exists {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Code:  http.StatusUnauthorized,
			Error: "Unauthorized: user_id not found",
		})
		return
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Code:  http.StatusUnauthorized,
			Error: "Unauthorized: user_id is not a valid string",
		})
		return
	}

	fileHeader, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	folder := context.DefaultPostForm("folder", "img")
	fileName := context.DefaultPostForm("file_name", randx.GenString(10))
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

	studentIDRaw := context.PostForm("student_id")
	var studentID *string
	if strings.TrimSpace(studentIDRaw) != "" {
		studentID = &studentIDRaw
	}

	teacherIDRaw := context.PostForm("teacher_id")
	var teacherID *string
	if strings.TrimSpace(teacherIDRaw) != "" {
		teacherID = &teacherIDRaw
	}

	url, img, err := receiver.UploadImageUseCase.UploadImage(dataBytes, folder, fileHeader.Filename, fileName, mode, topicID, &userID, studentID, teacherID)
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
