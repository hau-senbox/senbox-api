package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type MessageLanguageController struct {
	messageLanguageUseCase *usecase.MessageLanguageUseCase
}

func NewMessageLanguageController(messageLanguageUseCase *usecase.MessageLanguageUseCase) *MessageLanguageController {
	return &MessageLanguageController{
		messageLanguageUseCase: messageLanguageUseCase,
	}
}

func (ctrl *MessageLanguageController) UploadMessageLanguage(ctx *gin.Context) {
	var req request.UploadMessageLanguageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Error:   err.Error(),
				Message: "Invalid request data",
			},
		)
		return
	}

	err := ctrl.messageLanguageUseCase.UploadMessageLanguage(req)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Error:   err.Error(),
				Message: "Upload failed",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK, response.SucceedResponse{
			Code:    http.StatusOK,
			Message: "Upload success",
			Data:    nil,
		},
	)
}
