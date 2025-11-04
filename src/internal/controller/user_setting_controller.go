package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type UserSettingController struct {
	UserSettingUsecase *usecase.UserSettingUseCase
}

func (c *UserSettingController) UploadUserSetting(ctx *gin.Context) {
	var req request.UploadUserSettingRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid request: " + err.Error(),
		})
		return
	}

	// Call usecase
	res, err := c.UserSettingUsecase.UploadUserSetting(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Upload failed",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Upload success",
		Data:    res,
	})
}

func (c *UserSettingController) UploadUserIsFirstLogin(ctx *gin.Context) {
	var req request.UploadUserIsFirstLoginRequest

	// Bind JSON request
	ctx.ShouldBindJSON(&req)

	// get user id from context
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "user_id is required",
		})
		return
	}

	req.UserID = userID

	err := c.UserSettingUsecase.UploadUserIsFirstLogin(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Error:   err.Error(),
			Message: "Upload failed",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Upload success",
		Data:    nil,
	})
}

func (c *UserSettingController) UploadUserWelcomeReminder(ctx *gin.Context) {
	var req request.UploadUserWelcomeReminderRequest

	// Bind JSON request
	ctx.ShouldBindJSON(&req)

	// get user id from context
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "user_id is required",
		})
		return
	}

	req.UserID = userID

	err := c.UserSettingUsecase.UploadUserWelcomeReminder(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Error:   err.Error(),
			Message: "Upload failed",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Upload success",
		Data:    nil,
	})
}
