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

func (ctrl *MessageLanguageController) UploadMessageLanguages(ctx *gin.Context) {
	var req request.UploadMessageLanguagesRequest
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

	err := ctrl.messageLanguageUseCase.UploadMessageLanguages(req)
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

func (ctrl *MessageLanguageController) GetMessageLanguages4GW(ctx *gin.Context) {
	typeStr := ctx.Query("type")
	typeID := ctx.Query("type_id")

	if typeStr == "" || typeID == "" {
		ctx.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Error:   "type and type_id are required",
				Message: "Invalid query parameters",
			},
		)
		return
	}

	result, err := ctrl.messageLanguageUseCase.GetMessageLanguages4GW(typeStr, typeID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Error:   err.Error(),
				Message: "Failed to get message languages",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK, response.SucceedResponse{
			Code:    http.StatusOK,
			Message: "Get message languages success",
			Data:    result,
		},
	)
}

func (ctrl *MessageLanguageController) GetMessageLanguage4GW(ctx *gin.Context) {
	typeStr := ctx.Query("type")
	typeID := ctx.Query("type_id")
	languageID := ctx.GetUint("language_id")

	if typeStr == "" || typeID == "" || languageID == 0 {
		ctx.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Error:   "type and type_id are required",
				Message: "Invalid query parameters",
			},
		)
		return
	}

	result, err := ctrl.messageLanguageUseCase.GetMessageLanguage4GW(typeStr, typeID, languageID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Error:   err.Error(),
				Message: "Failed to get message languages",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK, response.SucceedResponse{
			Code:    http.StatusOK,
			Message: "Get message languages success",
			Data:    result,
		},
	)
}
