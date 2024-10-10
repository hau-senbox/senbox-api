package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
)

type _modifyDeviceURI struct {
	DeviceID string `uri:"device_id" binding:"required"`
}

// Frontend Update Device godoc
// @Summary Update Device
// @Description Update Device
// @Tags Frontend
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param req body request.ModifyDeviceRequest true "Update Device Params"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 404 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/frontend/device/:device_id [post]
func UpdateDevice(context *gin.Context) {
	var modifyDeviceURI _modifyDeviceURI

	if err := context.BindUri(&modifyDeviceURI); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	var req request.ModifyDeviceRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := usecase.ModifyDevice(req, modifyDeviceURI.DeviceID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Device updated",
		},
	})
}

func SyncDevice(context *gin.Context) {
	var _modifyDeviceURI _modifyDeviceURI

	if err := context.BindUri(&_modifyDeviceURI); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := usecase.SyncDevice(_modifyDeviceURI.DeviceID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Device synced",
		},
	})
}

// Get Device List godoc
// @Summary      Get Device List
// @Description  Get Device List
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param page query int false "Page"
// @Param keyword query string false "Keyword"
// @Param limit query int false "Limit"
// @Success      200  {object}  response.DeviceListResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/device/list [get]
func GetDevicesList(context *gin.Context) {
	var req request.GetListDeviceRequest
	if err := context.BindQuery(&req); err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			},
		)
		return
	}
	devices, paging, err := usecase.GetDeviceList(req)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				},
			},
		)
		return
	}
	result := make([]response.DeviceListResponseData, 0)
	if devices != nil {
		type DeviceAtt struct {
			Att1      request.RegisterDeviceUser `json:"user_01"`
			Att2      request.RegisterDeviceUser `json:"user_02"`
			Att3      request.RegisterDeviceUser `json:"user_03"`
			InputMode string                     `json:"input_mode"`
		}
		type ScreenButton struct {
			ButtonType  value.ButtonType `json:"button_type"`
			ButtonTitle string           `json:"button_title"`
		}
		for _, device := range devices {
			result = append(result, response.DeviceListResponseData{
				DeviceUUID:            device.DeviceId,
				DeviceName:            device.DeviceName,
				Attribute1:            device.PrimaryUserInfo,
				Attribute2:            device.SecondaryUserInfo,
				Attribute3:            device.TertiaryUserInfo,
				InputMode:             string(device.InputMode),
				Status:                string(device.Status),
				ProfilePicture:        device.ProfilePictureUrl,
				SpreadsheetUrl:        "https://docs.google.com/spreadsheets/d/" + device.SpreadsheetId,
				Message:               device.Message,
				ButtonUrl:             device.ButtonUrl,
				ScreenButtonType:      device.ScreenButtonType,
				ScreenButtonValue:     device.ScreenButtonValue,
				AppVersion:            device.AppVersion,
				Note:                  device.Note,
				UpdatedAt:             device.UpdatedAt.Format("2006-01-02 15:04:05"),
				TeacherSpreadsheetUrl: "https://docs.google.com/spreadsheets/d/" + device.TeacherSpreadsheetId,
			})
		}
	}
	context.JSON(http.StatusOK, response.DeviceListResponse{
		Data:   result,
		Paging: *paging,
	})

	return
}
