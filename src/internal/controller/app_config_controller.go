package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type AppConfigController struct {
	AppConfigUsecase        *usecase.AppConfigUseCase
	UserBlockSettingUsecase *usecase.UserBlockSettingUsecase
}

// Get all configs
func (ctrl *AppConfigController) GetAll(c *gin.Context) {
	configs, err := ctrl.AppConfigUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: configs,
	})
}

func (ctrl *AppConfigController) Upload(c *gin.Context) {
	var req request.UploadAppConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if err := ctrl.AppConfigUsecase.Upload(req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: "Uploaded successfully",
	})
}

// On is need to update
func (ctrl *AppConfigController) OnIsNeedToUpdate(c *gin.Context) {
	if err := ctrl.UserBlockSettingUsecase.OnIsNeedToUpdate(); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: "On is need to update successfully",
	})
}
