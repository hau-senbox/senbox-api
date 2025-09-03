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
