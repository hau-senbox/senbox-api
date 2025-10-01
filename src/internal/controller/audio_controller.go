package controller

import (
	"bufio"
	"io"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/pkg/randx"
	"sen-global-api/pkg/uploader"

	"github.com/gin-gonic/gin"
)

type AudioController struct {
	*usecase.GetAudioUseCase
	*usecase.UploadAudioUseCase
	*usecase.DeleteAudioUseCase
}

func (receiver *AudioController) GetUrlByKey(context *gin.Context) {
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

	url, err := receiver.GetAudioUseCase.GetUrlByKey(req.Key, mode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "audio was get successfully",
		Data:    *url,
	})
}

func (receiver *AudioController) CreateAudio(context *gin.Context) {
	fileHeader, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	folder := context.DefaultPostForm("folder", "audios")
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

	url, audio, err := receiver.UploadAudioUseCase.UploadAudio(dataBytes, folder, fileHeader.Filename, fileName, mode)
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
			Error: "audio upload failed",
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "audio was create successfully",
		Data: response.AudioResponse{
			AudioName: audio.AudioName,
			Key:       audio.Key,
			Extension: audio.Extension,
			Url:       *url,
		},
	})
}

func (receiver *AudioController) DeleteAudio(context *gin.Context) {
	var req request.DeleteAudioRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "key is required",
		})
		return
	}

	err := receiver.DeleteAudioUseCase.DeleteAudio(req.Key)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "audio deleted successfully",
	})
}

func (ac *AudioController) UploadAudio4GW(c *gin.Context) {
	var req request.UploadAudioRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// mở file
	file, err := req.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	res, err := ac.UploadAudioUseCase.UploadAudiov2(data, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "audio uploaded successfully",
		Data:    res,
	})
}

func (receiver *AudioController) DeleteAudio4GW(context *gin.Context) {
	key := context.Param("key")
	if key == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "key is required",
		})
		return
	}
	var req request.DeleteAudioRequest
	req.Key = key

	err := receiver.DeleteAudioUseCase.DeleteAudio(req.Key)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "audio deleted successfully",
	})
}
