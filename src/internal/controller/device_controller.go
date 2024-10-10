package controller

import (
	"fmt"
	"net/http"
	"os"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"sen-global-api/internal/middleware"
	"sen-global-api/pkg/monitor"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceController struct {
	DBConn *gorm.DB
	*usecase.UpdateDeviceSheetUseCase
	*usecase.RegisterDeviceUseCase
	*usecase.GetDeviceByIdUseCase
	*usecase.GetDeviceListUseCase
	*usecase.UpdateDeviceUseCase
	*usecase.FindDeviceFromRequestCase
	*usecase.GetFormByIdUseCase
	*usecase.TakeNoteUseCase
	*usecase.SubmitFormUseCase
	*usecase.GetScreenButtonsByDeviceUseCase
	*usecase.GetTimeTableUseCase
	*usecase.GetSettingMessageUseCase
	*usecase.RefreshAccessTokenUseCase
	*usecase.UpdateDeviceInfoUseCase
	*usecase.GetDeviceStatusUseCase
	*usecase.SyncSubmissionUseCase
	*usecase.GetModeLUseCase
	*usecase.DiscoverUseCase
	*usecase.DeviceSignUpUseCases
	*usecase.GetBrandLogoCase
	*usecase.GetRecentSubmissionByFormIdUseCase
	*usecase.RegisterFcmDeviceUseCase
	*usecase.SendNotificationUseCase
	*usecase.ResetCodeCountingUseCase
}

// Init a device godoc
// @Summary      Register a new device
// @Description  Init a device
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param req body request.RegisterDeviceRequest true "Init Device Params"
// @Success      200  {object}  response.AuthorizedDeviceResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/init [post]
func (receiver *DeviceController) InitDeviceV1(c *gin.Context) {
	var req request.RegisterDeviceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			},
		)
		return
	}
	log.Info("Init Device: ", req)
	monitor.SendMessageViaTelegram("Init Device: ", req.DeviceUUID, req.Primary.Fullname, req.Secondary.Fullname, req.AppVersion)
	device, err := receiver.RegisterDeviceUseCase.RegisterDevice(req)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				},
			},
		)
		return
	}

	if err != nil {
		c.JSON(
			http.StatusUnauthorized, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusUnauthorized,
					Message: "Cannot authorize device",
				},
			},
		)
		return
	}

	c.JSON(http.StatusOK, device)

	return
}

// Deactivate Device godoc
// @Summary      Deactivate Device
// @Description  Deactivate Device
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param device_id path string true "Device Id"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/device/deactivate/{device_id} [put]
func (receiver *DeviceController) DeactivateDevice(context *gin.Context) {
	deviceId := context.Param("device_id")
	var req request.DeactivateDeviceRequest
	if err := context.ShouldBind(&req); err != nil {
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

	err := receiver.UpdateDeviceSheetUseCase.DeactivateDevice(deviceId, req)
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
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Deactivated",
		},
	})

	return
}

// Aactivate Device godoc
// @Summary      Activate Device
// @Description  Aactivate Device
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param device_id path string true "Device Id"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/device/activate/{device_id} [put]
func (receiver *DeviceController) ActivateDevice(context *gin.Context) {
	deviceId := context.Param("device_id")
	var req request.ReactivateDeviceRequest
	if err := context.ShouldBind(&req); err != nil {
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
	err := receiver.UpdateDeviceSheetUseCase.ActivateDevice(deviceId, req)
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
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Activated",
		},
	})

	return
}

// Update Device godoc
// @Summary      Update Device
// @Description  Update Device
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param device_id path string true "Device Id"
// @Param req body request.UpdateDeviceRequest true "Update Device Params"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/device/{device_id}/update/ [put]
func (receiver *DeviceController) UpdateDevice(context *gin.Context) {
	deviceId := context.Param("device_id")
	if deviceId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "Device Id is required",
				},
			},
		)
		return
	}
	var req request.UpdateDeviceRequest
	if err := context.BindJSON(&req); err != nil {
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
	device, err := receiver.UpdateDeviceUseCase.UpdateDevice(deviceId, req)
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

	context.JSON(http.StatusOK, response.UpdateDeviceResponse{
		Data: response.DeviceListResponseData{
			DeviceUUID:            device.DeviceId,
			DeviceName:            device.DeviceName,
			Attribute1:            device.PrimaryUserInfo,
			Attribute2:            device.SecondaryUserInfo,
			Attribute3:            device.TertiaryUserInfo,
			InputMode:             string(device.InputMode),
			Status:                value.GetDeviceStatusStringAtMode(device.Status),
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
		},
	})
}

// Save device's Note godoc
// @Summary      Save device's Note
// @Description  Save device's Note
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param req body request.TakeNoteRequest true "Take Note Params"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/note/ [put]
func (receiver *DeviceController) TakeNote(context *gin.Context) {
	var takeNoteRequest request.TakeNoteRequest
	if err := context.ShouldBindJSON(&takeNoteRequest); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	err = receiver.TakeNoteUseCase.TakeNote(takeNoteRequest, device)
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
			Message: "Note taken",
		},
	})
}

// Submit a form godoc
// @Summary      Submit a form
// @Description  Submit a form
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param req body request.SubmitFormRequest true "Send Email Params"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/form/submit [post]
func (receiver *DeviceController) SubmitForm(context *gin.Context) {
	var req request.SubmitFormRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	form, err := receiver.GetFormByIdUseCase.GetFormByQRCode(req.QRCode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	reportMsgHeader := fmt.Sprintf("[FORM SUMBITTING]: %s", context.Request.Header)
	reportMsgBody := fmt.Sprintf("\n[FORM] ID: %d - \nNOTE:%s \nSUBMISSION TYPE: %s \n [SUBMITTED by DEVICE]: %s \n %s - \n %s \n %s\n", form.FormId, form.Note, form.SubmissionType, device.DeviceId, device.PrimaryUserInfo, device.SpreadsheetId, req)
	userInfo := fmt.Sprintf("\nSUBMITTED with [USER INFO] 1: %s \n [USER INFO] 2:%s \n", device.PrimaryUserInfo, device.SecondaryUserInfo)
	monitor.SendMessageViaTelegram(reportMsgHeader, reportMsgBody, userInfo)
	err = receiver.SubmitFormUseCase.AnswerForm(form.FormId, *device, req)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusNotAcceptable,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Succeed",
		},
	})

	defer func() {
		receiver.SyncSubmissionUseCase.Execute()
	}()
	return
}

type getLastSubmissionFromFormRequest struct {
	Code string `form:"code" bindding:"required"`
}

// Get last submission from a form godoc
// @Summary      Get last submission from a form
// @Description  Get last submission from a form
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param req body getLastSubmissionFromFormRequest true "Request"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/form/submission/last [get]
func (receiver *DeviceController) GetLastSubmissionByForm(context *gin.Context) {
	var req getLastSubmissionFromFormRequest
	if err := context.BindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	res, err := receiver.GetRecentSubmissionByFormIdUseCase.Execute(req.Code)
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
		Data: res,
	})
}

// Get Screen Buttons godoc
// @Summary      Get Screen Buttons
// @Description  Get Screen Buttons
// @Tags         Device
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.GetScreenButtonsResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/buttons/screen [get]
func (receiver *DeviceController) GetScreenButtons(context *gin.Context) {
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	buttons, err := receiver.GetScreenButtonsByDeviceUseCase.GetScreenButtons(*device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.GetScreenButtonsResponse{
		Data: response.GetScreenButtonsResponseData{
			Buttons: buttons,
		},
	})
}

// Get Top Buttons godoc
// @Summary      Get Top Buttons
// @Description  Get Top Buttons
// @Tags         Device
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.GetScreenButtonsResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/buttons/top [get]
func (receiver *DeviceController) GetTopButtons(context *gin.Context) {
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	buttons, err := receiver.GetScreenButtonsByDeviceUseCase.GetTopButtons(*device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.GetScreenButtonsResponse{
		Data: response.GetScreenButtonsResponseData{
			Buttons: buttons,
		},
	})
}

// Get Time Table godoc
// @Summary      Get Time Table
// @Description  Get Time Table
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Success      200  {object}  response.GetTimeTableResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/time-table [get]
func (receiver *DeviceController) GetTimeTable(context *gin.Context) {
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	times, err := receiver.GetTimeTableUseCase.Execute(*device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, times)
}

// Get Setting's Message godoc
// @Summary      Setting's Message
// @Description  Setting's Message
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Success      200  {object}  response.GetSettingMessageResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/messages [get]
func (receiver *DeviceController) GetSettingMessage(context *gin.Context) {
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	times, err := receiver.GetSettingMessageUseCase.Execute(*device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, times)
}

type refreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshAccessTokenResponseData struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshAccessTokenResponse struct {
	Data refreshAccessTokenResponseData `json:"data"`
}

// Refresh Access Token godoc
// @Summary      Refresh Access Token
// @Description  Refresh Access Token
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param req body refreshAccessTokenRequest true "Refresh Token Params"
// @Success      200  {object}  refreshAccessTokenResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      401  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/refresh-token [post]
func (receiver *DeviceController) RefreshAccessToken(context *gin.Context) {
	log.Debug(context.Request.Body)
	var request refreshAccessTokenRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	accessToken, refreshToken, err := receiver.RefreshAccessTokenUseCase.Execute(request.RefreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, refreshAccessTokenResponse{
		Data: refreshAccessTokenResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

type updateDeviceInfoRequest struct {
	Version   *string `json:"app_version"`
	UserInfo3 *string `json:"user_info_3"`
}

// Update Device Info godoc
// @Summary      Update Device Info
// @Description  Update Device Info
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param req body updateDeviceInfoRequest true "Update Device Info Params"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      401  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/update-info [put]
func (receiver *DeviceController) UpdateDeviceInfo(context *gin.Context) {
	var request updateDeviceInfoRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	if request.UserInfo3 == nil && request.Version == nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	err = receiver.UpdateDeviceInfoUseCase.Execute(*device, request.Version, request.UserInfo3)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		monitor.SendMessageViaTelegram(fmt.Sprintf("Update Device Info Failed: %s", err.Error()))
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Update Device Info Succeed",
		},
	})
}

type getDeviceStatusResponse struct {
	Data response.GetDeviceStatusResponseData `json:"data"`
}

// Get Device Status godoc
// @Summary      Get Device Status
// @Description  Get Device Status
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Success      200  {object}  getDeviceStatusResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      401  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/status [get]
func (receiver *DeviceController) GetDeviceStatus(context *gin.Context) {
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		monitor.SendMessageViaTelegram(fmt.Sprintf("Get Device Status Failed: %s", err.Error()))
		return
	}

	status, err := receiver.GetDeviceStatusUseCase.Execute(*device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		monitor.SendMessageViaTelegram(fmt.Sprintf("Get Device Status Failed: %s", err.Error()))
		return
	}

	context.JSON(http.StatusOK, getDeviceStatusResponse{
		Data: status,
	})
}

type reserveRequest struct {
	DeviceId   string `json:"device_id" binding:"required"`
	AppVersion string `json:"app_version" binding:"required"`
}

type reserveResponse struct {
}

// Reserve Device godoc
// @Summary      Reserve Device
// @Description  Reserve Device
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param	   	X-API-KEY header string true "API Key"
// @Param req body reserveRequest true "Reserve Device Params"
// @Success      200  {object}  reserveResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      401  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Router       /v1/device/reserve [post]
func (receiver *DeviceController) Reserve(context *gin.Context) {
	//Extract API Key from header
	apiKey := context.GetHeader("X-API-KEY")
	if apiKey != os.Getenv("SENBOX_API_KEY") {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "API key is required",
			},
		})
		return
	}

	var request reserveRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	log.Debug(request)
	err := receiver.RegisterDeviceUseCase.Reserve(request.DeviceId, request.AppVersion)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[FAILED] Reserve Device: %s", err.Error()),
			fmt.Sprintf("Device ID: %s", request.DeviceId),
			fmt.Sprintf("App Version: %s", request.AppVersion),
		)
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Reserve Device Succeed",
		},
	})
}

type getBrandLogoResponse struct {
	ImageInBase64 string `json:"image_url" binding:"required"`
}

// Get Brand 	Logo
// @Summary 	Brand Logo
// @Description Brand Logo
// @Tag Device
// @Accept 		json
// @Produce 	json
// @Param Authorization header string true "Bearer {token}"
// @Success      200  {object}  getBrandLogoResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      401  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/brand-logo [get]
func (receiver *DeviceController) GetBrandLogo(context *gin.Context) {
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "invalid request",
			},
		})
		return
	}

	modeL, err := receiver.GetBrandLogoCase.Execute(*device)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "invalid request",
			},
		})
		return
	}

	context.JSON(http.StatusOK, getBrandLogoResponse{
		ImageInBase64: modeL,
	})
}

type getModeLResponse struct {
	Value string `json:"value" binding:"required"`
}

// Get Status L godoc
// @Summary      Get Status L
// @Description  Get Status L
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Success      200  {object}  getModeLResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      401  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/mode-l [get]
func (receiver *DeviceController) GetModeL(context *gin.Context) {
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "invalid request",
			},
		})
		return
	}

	modeL, err := receiver.GetModeLUseCase.Execute(*device)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "invalid request",
			},
		})
		return
	}

	context.JSON(http.StatusOK, getModeLResponse{
		Value: modeL,
	})
}

type discoverRequest struct {
	DeviceId string `json:"device_id" binding:"required"`
}

type discoverResponse struct {
	DeviceId  string `json:"device_id"`
	UserInfo1 string `json:"user_info_1"`
	UserInfo2 string `json:"user_info_2"`
	UserInfo3 string `json:"user_info_3"`
}

// Discover godoc
// @Summary      Discover
// @Description  Discover
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param req body discoverRequest true "Discover Params"
// @Success      200  {object}  discoverResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      401  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      461  {object}  response.FailedResponse
// @Router       /v1/device/discover [post]
func (receiver *DeviceController) Discover(context *gin.Context) {
	var rq discoverRequest
	if err := context.ShouldBindJSON(&rq); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	res, err := receiver.DiscoverUseCase.Execute(rq.DeviceId)
	if err != nil {
		//device is not in api mode(aka not in `VIA Spreadsheet`)
		if err.Error() == "not_api_input_mode" {
			context.JSON(461, response.FailedResponse{
				Error: response.Cause{
					Code:    461,
					Message: "device is not in api mode(aka not in `VIA Spreadsheet`)",
				},
			})
			return
		} else {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusInternalServerError,
					Message: "invalid request",
				},
			})
			return
		}
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: discoverResponse{
			DeviceId:  res.DeviceId,
			UserInfo1: res.UserInfo1,
			UserInfo2: res.UserInfo2,
			UserInfo3: res.UserInfo3,
		},
	})
}

type getSignUpResponseTextButton struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type getSignUpResponse struct {
	Button1 getSignUpResponseTextButton `json:"button_1" binding:"required"`
	Button2 getSignUpResponseTextButton `json:"button_2" binding:"required"`
	Button3 getSignUpResponseTextButton `json:"button_3" binding:"required"`
	Button4 getSignUpResponseTextButton `json:"button_4" binding:"required"`
}

// Get Sign Up godoc
// @Summary      Get Sign Up
// @Description  Get Sign Up
// @Tags         Device
// @Accept       json
// @Produce      json
// @Response     200  {object}  getSignUpResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/sign-up [get]
func (receiver *DeviceController) GetSignUp(context *gin.Context) {
	setting, err := receiver.DeviceSignUpUseCases.GetSignUpSetting()
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
		Data: getSignUpResponse{
			Button1: getSignUpResponseTextButton{
				Name:  setting.Button1.Name,
				Value: setting.Button1.Value,
			},
			Button2: getSignUpResponseTextButton{
				Name:  setting.Button2.Name,
				Value: setting.Button2.Value,
			},
			Button3: getSignUpResponseTextButton{
				Name:  setting.Button3.Name,
				Value: setting.Button3.Value,
			},
			Button4: getSignUpResponseTextButton{
				Name:  setting.Button4.Name,
				Value: setting.Button4.Value,
			},
		},
	})
}

type getSignUpFormRequest struct {
	Code string `form:"code"`
}

// Get Registration Form godoc
// @Summary      Get Registration Form
// @Description  Get Registration Form
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param req query getSignUpFormRequest true "Get Registration Form"
// @Success 200 {object} response.QuestionListResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/sign-up/form [get]
func (receiver *DeviceController) GetSignUpForm(context *gin.Context) {
	var rq getSignUpFormRequest
	if err := context.ShouldBindQuery(&rq); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}
	if rq.Code == "" {
		r := receiver.DeviceSignUpUseCases.GetSignUpFormQuestions()
		if r == nil {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusInternalServerError,
					Message: "Could not get sign up form",
				},
			})
			return
		}

		context.JSON(http.StatusOK, response.SucceedResponse{
			Data: r,
		})

		return
	}

	r := receiver.DeviceSignUpUseCases.GetSignUpFormQuestionsByFormNote(rq.Code)
	if r == nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "Could not get sign up form",
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: r,
	})
}

type signUpRequestUri struct {
	FormId int `uri:"form_id" binding:"required"`
}

type signUpRequestBody struct {
	Answers  []request.Answer `json:"answers" binding:"required"`
	OpenedAt time.Time        `json:"opened_at,default=now()"`
}

type SubmitFormRequest struct {
}

// Sign Up godoc
// @Summary      Sign Up
// @Description  Sign Up
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param req body request.SubmitFormRequest true "Sign Up Params"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Router       /v1/device/sign-up/form [post]
func (receiver *DeviceController) SubmitSignUpForm(context *gin.Context) {
	var rq request.SubmitFormRequest
	if err := context.ShouldBindJSON(&rq); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	form, err := receiver.DeviceSignUpUseCases.FindSignUpForm()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "Sign up form not found",
			},
		})
		return
	}

	if rq.DeviceId == "" {
		//Legacy app does not send device id
		err = receiver.SubmitFormUseCase.SubmitSignUpForm(form, rq)
	} else {
		err = receiver.SubmitFormUseCase.SubmitSignUpMemoryForm(form, rq)
	}

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
			Message: "Succeed",
		},
	})
}

type getPresetResponse struct {
	Value12 *string `json:"value_12"`
}

// Get Preset Form godoc
// @Summary      Get Preset Form
// @Description  Get Preset Form
// @Tags         Device
// @Accept       json
// @Produce      json
// @Success 200 {object} getPresetResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/sign-up/pre-set [get]
func (receiver *DeviceController) GetPreset(context *gin.Context) {
	r := receiver.DeviceSignUpUseCases.GetSigGnUpPresetSetting()
	if r == nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "Could not get sign up form",
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: getPresetResponse{
			Value12: r,
		},
	})
}

// Get Register FCM Device godoc
// @Summary      Register FCM Device
// @Description  Register FCM Device
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param  req body request.RegisterFCMRequest true "Register FCM Params"
// @Success 200 {object} response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/messaging/fcm/register [post]
func (receiver *DeviceController) RegisterFCM(context *gin.Context) {
	var req request.RegisterFCMRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	err := receiver.RegisterFcmDeviceUseCase.Execute(req)
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
			Message: "Succeed",
		},
	})
}

// Send Notification	Logo
// @Summary 	Send Notification
// @Description Send Notification
// @Tag Device
// @Accept 		json
// @Produce 	json
// @Param Authorization header string true "Bearer {token}"
// @Param  req body request.SendNotificationRequest true "Send Notification Request"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      401  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/buttons/notification [post]
func (receiver *DeviceController) SenNotification(context *gin.Context) {
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	var req request.SendNotificationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	err = receiver.SendNotificationUseCase.Execute(req, *device)

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
			Message: "Succeed",
		},
	})
}

type GetDeviceSignUpResponse struct {
	FormNote string `json:"form_note" binding:"required"`
}

// Get Device's Sign Up Form godoc
// @Summary      Get  Device's Sign Up Form
// @Description  Get  Device's Sign Up Form
// @Tags         Device
// @Accept       json
// @Produce      json
// @Success 200 {object} GetDeviceSignUpResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/form/device/sign-up [get]
func (receiver *DeviceController) GetDeviceSignUp(context *gin.Context) {
	deviceId := context.GetString(middleware.ContextKeyDeviceId)

	r := receiver.DeviceSignUpUseCases.GetSignUpFormQuestionsByDevice(deviceId)
	if r == nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "Could not get sign up form for this device",
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: GetDeviceSignUpResponse{
			FormNote: r.Data.FormName,
		},
	})
}

// Update Device's Sign Up Form godoc
// @Summary      Update Device's Sign Up Form
// @Description  Update Device's Sign Up Form
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param req body request.SubmitFormRequest true "Update Device's Sign Up Form Params"
// @Success 200 {object} response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/form/device/sign-up [put]
func (receiver *DeviceController) UpdateDeviceSignUp(context *gin.Context) {
	var rq request.SubmitFormRequest
	if err := context.ShouldBindJSON(&rq); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	deviceId := context.GetString(middleware.ContextKeyDeviceId)
	rq.DeviceId = deviceId

	err := receiver.SubmitFormUseCase.UpdateSignUpMemoryForm(rq)

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
			Message: "Succeed",
		},
	})
}

// Reset Code Counting
// @Summary      Reset Code Counting
// @Description  Reset Code Counting
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param  req body request.ResetCodeCountingRequest true "Reset Code Counting Request"
// @Success 200 {object} response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/code-counting/reset [put]
func (receiver *DeviceController) ResetCodeCounting(context *gin.Context) {
	var rq request.ResetCodeCountingRequest
	if err := context.ShouldBindJSON(&rq); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	err := receiver.ResetCodeCountingUseCase.Execute(rq)

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
			Message: "Succeed",
		},
	})
}
