package controller

import (
	"fmt"
	"net/http"
	"os"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"sen-global-api/internal/middleware"
	"sen-global-api/pkg/monitor"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tiendc/gofn"
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
	*usecase.GetDevicesByUserIdUseCase
	*usecase.GetUserFromTokenUseCase
	*usecase.GetUserDeviceUseCase
}

func (receiver *DeviceController) GetDeviceById(c *gin.Context) {
	deviceId := c.Param("device_id")
	if deviceId == "" {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "device id is required",
				},
			},
		)
		return
	}

	device, err := receiver.GetDeviceByIdUseCase.Get(deviceId)
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

	c.JSON(http.StatusOK, response.SucceedResponse{
		Data: &response.DeviceResponseDataV2{
			Id:                device.ID,
			DeviceName:        device.DeviceName,
			InputMode:         string(device.InputMode),
			Status:            string(device.Status),
			DeactivateMessage: device.DeactivateMessage,
			ButtonUrl:         device.ButtonUrl,
			AppVersion:        device.AppVersion,
			Note:              device.Note,
			CreatedAt:         device.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:         device.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func (receiver *DeviceController) GetAllDeviceByUserId(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "user id is required",
				},
			},
		)
		return
	}

	devices, err := receiver.GetDevicesByUserIdUseCase.GetDevicesByUserId(userId)
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

	var deviceResponse []response.DeviceResponseV2

	for _, device := range *devices {
		deviceResponse = append(deviceResponse, response.DeviceResponseV2{
			ID:         device.ID,
			DeviceName: device.DeviceName,
		})
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Data: deviceResponse,
	})
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
func (receiver *DeviceController) InitDeviceV1(context *gin.Context) {
	var req request.RegisterDeviceRequest
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

	monitor.SendMessageViaTelegram("Init Device: ", req.DeviceUUID, req.AppVersion)

	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	deviceId, err := receiver.RegisterDeviceUseCase.RegisterDevice(user, req)
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
		Data: deviceId,
	})
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
					Message: "device id is required",
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
		Data: response.DeviceResponseData{
			DeviceUUID:        device.ID,
			DeviceName:        device.DeviceName,
			InputMode:         string(device.InputMode),
			Status:            value.GetDeviceStatusStringAtMode(device.Status),
			DeactivateMessage: device.DeactivateMessage,
			ButtonUrl:         device.ButtonUrl,
			ScreenButtonType:  device.ScreenButtonType,
			AppVersion:        device.AppVersion,
			Note:              device.Note,
			UpdatedAt:         device.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func (receiver *DeviceController) UpdateDeviceV2(context *gin.Context) {
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
	var req request.UpdateDeviceRequestV2
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
	device, err := receiver.UpdateDeviceUseCase.UpdateDeviceV2(deviceId, req)
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

	go func() {
		// err = usecase.SyncDevice(deviceId)
		// if err != nil {
		// 	context.JSON(http.StatusInternalServerError, response.FailedResponse{
		// 		Error: response.Cause{
		// 			Code:    http.StatusInternalServerError,
		// 			Message: err.Error(),
		// 		},
		// 	})
		// 	return
		// }
	}()

	context.JSON(http.StatusOK, response.UpdateDeviceResponse{
		Data: response.DeviceResponseData{
			DeviceUUID:        device.ID,
			DeviceName:        device.DeviceName,
			InputMode:         string(device.InputMode),
			Status:            value.GetDeviceStatusStringAtMode(device.Status),
			DeactivateMessage: device.DeactivateMessage,
			ButtonUrl:         device.ButtonUrl,
			ScreenButtonType:  device.ScreenButtonType,
			AppVersion:        device.AppVersion,
			Note:              device.Note,
			UpdatedAt:         device.UpdatedAt.Format("2006-01-02 15:04:05"),
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
	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	userDevices, err := receiver.GetUserDeviceUseCase.GetUserDeviceById(takeNoteRequest.DeviceId)
	var userIds []string
	for _, userDevice := range *userDevices {
		userIds = append(userIds, userDevice.UserId.String())
	}
	isExist := gofn.ContainSlice(userIds, []string{
		user.ID.String(),
	})
	if !isExist {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			},
		})
		return
	}

	err = receiver.TakeNoteUseCase.TakeNote(takeNoteRequest, takeNoteRequest.DeviceId)
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

	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	userDevices, err := receiver.GetUserDeviceUseCase.GetUserDeviceById(req.DeviceId)
	var userIds []string
	for _, userDevice := range *userDevices {
		userIds = append(userIds, userDevice.UserId.String())
	}
	isExist := gofn.ContainSlice(userIds, []string{
		user.ID.String(),
	})
	if !isExist {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized",
			},
		})
		return
	}

	ud, found := gofn.Find(*userDevices, func(ud entity.SUserDevices) bool {
		return ud.UserId.String() == user.ID.String()
	})
	if !found {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			},
		})
		return
	}

	reportMsgHeader := fmt.Sprintf("[FORM SUMBITTING]: %s", context.Request.Header)
	userInfo := fmt.Sprintf("\nSUBMITTED with [USER INFO] 1: %s \n [USER INFO] 2:%s \n", user.Fullname, user.Company.CompanyName)
	reportMsgBody := fmt.Sprintf("\n[FORM] ID: %d - \nNOTE:%s \n [SUBMITTED by DEVICE]: %s \n %s - \n %v\n", form.ID, form.Note, ud.DeviceId, user.Fullname, req)
	monitor.SendMessageViaTelegram(reportMsgHeader, reportMsgBody, userInfo)
	err = receiver.SubmitFormUseCase.AnswerForm(form.ID, ud, req)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusNotAcceptable,
				Message: err.Error(),
			},
		})
		return
	}

	// defer receiver.SyncSubmissionUseCase.Execute()

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Succeed",
		},
	})
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
// @Router       /v1/buttons/screen/:device_id [get]
func (receiver *DeviceController) GetScreenButtons(context *gin.Context) {
	deviceId := context.Param("device_id")
	if deviceId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "device id is required",
				},
			},
		)
		return
	}

	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	buttons, err := receiver.GetScreenButtonsByDeviceUseCase.GetScreenButtons(deviceId, user)
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
	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	buttons, err := receiver.GetScreenButtonsByDeviceUseCase.GetTopButtons(*user)
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
	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	times, err := receiver.GetTimeTableUseCase.Execute(*user)
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
	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	times, err := receiver.GetSettingMessageUseCase.Execute(*user)
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

func (receiver *DeviceController) GetSettingMessageV2(context *gin.Context) {
	deviceId := context.Param("device_id")
	if deviceId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "device id is required",
				},
			},
		)
		return
	}

	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	userDevices, err := receiver.GetUserDeviceUseCase.GetUserDeviceById(deviceId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	var userIds []string
	for _, userDevice := range *userDevices {
		userIds = append(userIds, userDevice.UserId.String())
	}
	isExist := gofn.ContainSlice(userIds, []string{
		user.ID.String(),
	})
	if !isExist {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized",
			},
		})
		return
	}

	times, err := receiver.GetSettingMessageUseCase.Execute(*user)
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
	DeviceId string  `json:"device_id" binding:"required"`
	Version  *string `json:"app_version"`
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
	if request.Version == nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	err = receiver.UpdateDeviceInfoUseCase.Execute(*user, request.DeviceId, request.Version)
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
// @Router       /v1/device/status/:device_id [get]
func (receiver *DeviceController) GetDeviceStatus(context *gin.Context) {
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

	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
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

	userDevices, err := receiver.GetUserDeviceUseCase.GetUserDeviceById(deviceId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	var userIds []string
	for _, userDevice := range *userDevices {
		userIds = append(userIds, userDevice.UserId.String())
	}
	isExist := gofn.ContainSlice(userIds, []string{
		user.ID.String(),
	})
	if !isExist {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized",
			},
		})
		return
	}

	device, err := receiver.GetDeviceByIdUseCase.GetDeviceById(deviceId)
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
// @Router       /v1/user/brand-logo [get]
func (receiver *DeviceController) GetBrandLogo(context *gin.Context) {
	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "invalid request",
			},
		})
		return
	}

	modeL, err := receiver.GetBrandLogoCase.Execute(*user)

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
	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "invalid request",
			},
		})
		return
	}

	modeL, err := receiver.GetModeLUseCase.Execute(*user)

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
	DeviceId string `json:"device_id"`
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
			DeviceId: res.DeviceId,
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
	Button5 getSignUpResponseTextButton `json:"button_5" binding:"required"`
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
			Button5: getSignUpResponseTextButton{
				Name:  setting.Button5.Name,
				Value: setting.Button5.Value,
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
			Data: r.Data,
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
		Data: r.Data,
	})
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
	Value *string `json:"value"`
}

// Get Preset Form godoc
// @Summary      Get Preset 2 Form
// @Description  Get Preset 2 Form
// @Tags         Device
// @Accept       json
// @Produce      json
// @Success 200 {object} getPresetResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/sign-up/pre-set-2 [get]
func (receiver *DeviceController) GetPreset2(context *gin.Context) {
	r := receiver.DeviceSignUpUseCases.GetSignUpPreset2Setting()
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
			Value: r,
		},
	})
}

// Get Preset Form godoc
// @Summary      Get Preset 1 Form
// @Description  Get Preset 1 Form
// @Tags         Device
// @Accept       json
// @Produce      json
// @Success 200 {object} getPresetResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/device/sign-up/pre-set-1 [get]
func (receiver *DeviceController) GetPreset1(context *gin.Context) {
	r := receiver.DeviceSignUpUseCases.GetSignUpPreset1Setting()
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
			Value: r,
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

	err := receiver.SendNotificationUseCase.Execute(req)

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
