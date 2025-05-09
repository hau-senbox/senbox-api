package controller

import (
	"encoding/json"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeviceComponentValuesController struct {
	*usecase.GetDeviceComponentValuesUseCase
	*usecase.SaveDeviceComponentValuesUseCase
}

func (receiver *DeviceComponentValuesController) GetDeviceComponentValuesByOrganization(context *gin.Context) {
	organizationId := context.Param("organization_id")
	if organizationId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Organization Id is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(organizationId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Organization Id is invalid",
			},
		)
		return
	}

	setting, err := receiver.GetDeviceComponentValuesUseCase.GetDeviceComponentValuesByOrganization(request.GetDeviceComponentValuesByOrganizationRequest{ID: uint(id)})
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
		Data: response.DeviceComponentValuesResponse{
			ID:           int(setting.ID),
			Setting:      settingString,
			Organization: uint(*setting.OrganizationId),
		},
	})
}

func (receiver *DeviceComponentValuesController) GetDeviceComponentValuesByDevice(context *gin.Context) {
	organizationId := context.Param("organization_id")
	if organizationId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Organization Id is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(organizationId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Organization Id is invalid",
			},
		)
		return
	}

	setting, err := receiver.GetDeviceComponentValuesUseCase.GetDeviceComponentValuesByDevice(request.GetDeviceComponentValuesByDeviceRequest{ID: uint(id)})
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
