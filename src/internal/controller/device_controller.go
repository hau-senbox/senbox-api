package controller

import (
	"fmt"
	"net/http"
	"os"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
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
	*usecase.RefreshAccessTokenUseCase
	*usecase.GetDeviceStatusUseCase
	*usecase.DiscoverUseCase
	*usecase.DeviceSignUpUseCases
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
				Code:  http.StatusBadRequest,
				Error: "device id is required",
			},
		)
		return
	}

	device, err := receiver.Get(deviceId)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
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
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	devices, err := receiver.GetDevicesByUserId(userId)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
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
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	monitor.SendMessageViaTelegram("Init Device: ", req.DeviceUUID, req.AppVersion)

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	deviceId, err := receiver.RegisterDevice(user, req)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
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
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	err := receiver.UpdateDeviceSheetUseCase.DeactivateDevice(deviceId, req)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			},
		)
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Deactivated",
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
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	err := receiver.UpdateDeviceSheetUseCase.ActivateDevice(deviceId, req)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			},
		)
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Activated",
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
				Code:  http.StatusBadRequest,
				Error: "device id is required",
			},
		)
		return
	}
	var req request.UpdateDeviceRequest
	if err := context.BindJSON(&req); err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	device, err := receiver.UpdateDeviceUseCase.UpdateDevice(deviceId, req)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
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
				Code:  http.StatusBadRequest,
				Error: "Device Id is required",
			},
		)
		return
	}
	var req request.UpdateDeviceRequestV2
	if err := context.BindJSON(&req); err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}
	device, err := receiver.UpdateDeviceUseCase.UpdateDeviceV2(deviceId, req)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
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
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	userDevices, err := receiver.GetUserDeviceById(takeNoteRequest.DeviceId)
	var userIds []string
	for _, userDevice := range *userDevices {
		userIds = append(userIds, userDevice.UserId.String())
	}
	isExist := gofn.ContainSlice(userIds, []string{
		user.ID.String(),
	})
	if !isExist {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Code:  http.StatusUnauthorized,
			Error: err.Error(),
		})
		return
	}

	err = receiver.TakeNoteUseCase.TakeNote(takeNoteRequest, takeNoteRequest.DeviceId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Note taken",
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
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	form, err := receiver.GetFormByQRCode(req.QRCode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	req.UserId = user.ID.String()

	// reportMsgHeader := fmt.Sprintf("[FORM SUMBITTING]: %s", context.Request.Header)
	// userInfo := fmt.Sprintf("\nSUBMITTED with [USER INFO] 1: %s \n [USER INFO] 2:%s \n", user.Fullname, user.Organization.OrganizationName)
	// reportMsgBody := fmt.Sprintf("\n[FORM] RoleId: %d - \nNOTE:%s \n [SUBMITTED by USER]: %s \n %s - \n %v\n", form.RoleId, form.Note, user.Username, user.Fullname, req)
	// monitor.SendMessageViaTelegram(reportMsgHeader, reportMsgBody, userInfo)
	err = receiver.AnswerForm(form.ID, req)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, response.FailedResponse{
			Code:  http.StatusNotAcceptable,
			Error: err.Error(),
		})
		return
	}

	// defer receiver.SyncSubmissionUseCase.Execute()

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Succeed",
	})
}

type getLastSubmissionFromFormRequest struct {
	QRCode string `json:"qr_code" bindding:"required"`
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
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	res, err := receiver.GetRecentSubmissionByFormIdUseCase.Execute(req.QRCode, user.ID.String())
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: res,
	})
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
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := receiver.RefreshAccessTokenUseCase.Execute(request.RefreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
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
				Code:  http.StatusBadRequest,
				Error: "Device Id is required",
			},
		)
		return
	}

	// user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	// if err != nil {
	// 	context.JSON(http.StatusInternalServerError, response.FailedResponse{
	// 		Error: response.Cause{
	// 			Code:    http.StatusForbidden,
	// 			Message: err.Error(),
	// 		},
	// 	})
	// 	monitor.SendMessageViaTelegram(fmt.Sprintf("Get Device Status Failed: %s", err.Error()))
	// 	return
	// }

	// userDevices, err := receiver.GetUserDeviceUseCase.GetUserDeviceById(deviceId)
	// if err != nil {
	// 	context.JSON(http.StatusInternalServerError, response.FailedResponse{
	// 		Error: response.Cause{
	// 			Code:    http.StatusInternalServerError,
	// 			Message: err.Error(),
	// 		},
	// 	})
	// 	return
	// }
	// var userIds []string
	// for _, userDevice := range *userDevices {
	// 	userIds = append(userIds, userDevice.UserId.String())
	// }
	// isExist := gofn.ContainSlice(userIds, []string{
	// 	user.RoleId.String(),
	// })
	// if !isExist {
	// 	context.JSON(http.StatusUnauthorized, response.FailedResponse{
	// 		Error: response.Cause{
	// 			Code:    http.StatusUnauthorized,
	// 			Message: "unauthorized",
	// 		},
	// 	})
	// 	return
	// }

	device, err := receiver.GetDeviceByIdUseCase.GetDeviceById(deviceId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		monitor.SendMessageViaTelegram(fmt.Sprintf("Get Device Status Failed: %s", err.Error()))
		return
	}

	status, err := receiver.GetDeviceStatusUseCase.Execute(*device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
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
			Code:  http.StatusBadRequest,
			Error: "API key is required",
		})
		return
	}

	var request reserveRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid request",
		})
		return
	}

	log.Debug(request)
	err := receiver.RegisterDeviceUseCase.Reserve(request.DeviceId, request.AppVersion)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[FAILED] Reserve Device: %s", err.Error()),
			fmt.Sprintf("Device RoleId: %s", request.DeviceId),
			fmt.Sprintf("App Version: %s", request.AppVersion),
		)
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Reserve Device Succeed",
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
			Code:  http.StatusBadRequest,
			Error: "invalid request",
		})
		return
	}

	res, err := receiver.DiscoverUseCase.Execute(rq.DeviceId)
	if err != nil {
		//device is not in api mode(aka not in `VIA Spreadsheet`)
		if err.Error() == "not_api_input_mode" {
			context.JSON(461, response.FailedResponse{
				Code:  461,
				Error: "device is not in api mode(aka not in `VIA Spreadsheet`)",
			})
			return
		} else {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: "invalid request",
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
	setting, err := receiver.GetSignUpSetting()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
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
	r := receiver.GetSignUpFormQuestions()
	if r == nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: "Could not get sign up form",
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: r.Data,
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
	r := receiver.GetSignUpPreset2Setting()
	if r == nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: "Could not get sign up form",
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
	r := receiver.GetSignUpPreset1Setting()
	if r == nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: "Could not get sign up form",
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
			Code:  http.StatusBadRequest,
			Error: "invalid request",
		})
		return
	}

	err := receiver.RegisterFcmDeviceUseCase.Execute(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Succeed",
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
			Code:  http.StatusBadRequest,
			Error: "invalid request",
		})
		return
	}

	err := receiver.SendNotificationUseCase.Execute(req)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Succeed",
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
			Code:  http.StatusBadRequest,
			Error: "invalid request",
		})
		return
	}

	err := receiver.ResetCodeCountingUseCase.Execute(rq)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Succeed",
	})
}
