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

func (receiver *DeviceComponentValuesController) GetDeviceComponentValuesByCompany(context *gin.Context) {
	companyId := context.Param("company_id")
	if companyId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "Company Id is required",
				},
			},
		)
		return
	}

	id, err := strconv.ParseUint(companyId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "Company Id is invalid",
				},
			},
		)
		return
	}

	setting, err := receiver.GetDeviceComponentValuesUseCase.GetDeviceComponentValuesByCompany(request.GetDeviceComponentValuesByCompanyRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	var settingString response.SettingResponse
	if err := json.Unmarshal(setting.Setting, &settingString); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.DeviceComponentValuesResponse{
			ID:      int(setting.ID),
			Setting: settingString,
			Company: uint(*setting.CompanyId),
		},
	})
}

func (receiver *DeviceComponentValuesController) GetDeviceComponentValuesByDevice(context *gin.Context) {
	companyId := context.Param("company_id")
	if companyId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "Company Id is required",
				},
			},
		)
		return
	}

	id, err := strconv.ParseUint(companyId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "Company Id is invalid",
				},
			},
		)
		return
	}

	setting, err := receiver.GetDeviceComponentValuesUseCase.GetDeviceComponentValuesByDevice(request.GetDeviceComponentValuesByDeviceRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	var settingString response.SettingResponse
	if err := json.Unmarshal(setting.Setting, &settingString); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
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

func (receiver *DeviceComponentValuesController) SaveDeviceComponentValuesByCompany(context *gin.Context) {
	var req request.SaveDeviceComponentValuesByCompanyRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := receiver.SaveDeviceComponentValuesUseCase.SaveByCompany(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Setting was saved successfully",
		},
	})
}

func (receiver *DeviceComponentValuesController) SaveDeviceComponentValuesByDevice(context *gin.Context) {
	var req request.SaveDeviceComponentValuesByDeviceRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	err := receiver.SaveDeviceComponentValuesUseCase.SaveByDevice(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Setting was saved successfully",
		},
	})
}
