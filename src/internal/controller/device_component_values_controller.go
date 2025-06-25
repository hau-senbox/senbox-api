package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
)

type DeviceComponentValuesController struct {
	*usecase.GetDeviceComponentValuesUseCase
	*usecase.SaveDeviceComponentValuesUseCase
}

func (receiver *DeviceComponentValuesController) GetDeviceComponentValuesByOrganization(context *gin.Context) {
	organizationID := context.Param("organization_id")
	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Organization ID is required",
			},
		)
		return
	}

	setting, err := receiver.GetDeviceComponentValuesUseCase.GetDeviceComponentValuesByOrganization(request.GetDeviceComponentValuesByOrganizationRequest{ID: organizationID})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var settingString response.SettingResponse
	if err := json.Unmarshal(setting.Setting, &settingString); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.DeviceComponentValuesResponse{
			ID:           int(setting.ID),
			Setting:      settingString,
			Organization: setting.OrganizationID.String(),
		},
	})
}

func (receiver *DeviceComponentValuesController) GetDeviceComponentValuesByDevice(context *gin.Context) {
	organizationID := context.Param("organization_id")
	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Organization ID is required",
			},
		)
		return
	}

	setting, err := receiver.GetDeviceComponentValuesUseCase.GetDeviceComponentValuesByDevice(request.GetDeviceComponentValuesByDeviceRequest{ID: organizationID})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var settingString response.SettingResponse
	if err := json.Unmarshal(setting.Setting, &settingString); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.DeviceComponentValuesResponse{
			ID:      int(setting.ID),
			Setting: settingString,
		},
	})
}

func (receiver *DeviceComponentValuesController) SaveDeviceComponentValuesByOrganization(context *gin.Context) {
	var req request.SaveDeviceComponentValuesByOrganizationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.SaveDeviceComponentValuesUseCase.SaveByOrganization(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Setting was saved successfully",
	})
}

func (receiver *DeviceComponentValuesController) SaveDeviceComponentValuesByDevice(context *gin.Context) {
	var req request.SaveDeviceComponentValuesByDeviceRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	err := receiver.SaveDeviceComponentValuesUseCase.SaveByDevice(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Setting was saved successfully",
	})
}
