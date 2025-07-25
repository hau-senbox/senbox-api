package controller

import (
	"net/http"
	"os"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tiendc/gofn"
	"gorm.io/gorm"
)

type DeviceController struct {
	DBConn *gorm.DB
	*usecase.UpdateDeviceSheetUseCase
	*usecase.RegisterDeviceUseCase
	*usecase.GetDeviceByIDUseCase
	*usecase.GetDeviceListUseCase
	*usecase.UpdateDeviceUseCase
	*usecase.FindDeviceFromRequestCase
	*usecase.GetFormByIDUseCase
	*usecase.TakeNoteUseCase
	*usecase.SubmitFormUseCase
	*usecase.RefreshAccessTokenUseCase
	*usecase.GetDeviceStatusUseCase
	*usecase.DiscoverUseCase
	*usecase.DeviceSignUpUseCases
	*usecase.GetRecentSubmissionByFormIDUseCase
	*usecase.RegisterFcmDeviceUseCase
	*usecase.SendNotificationUseCase
	*usecase.ResetCodeCountingUseCase
	*usecase.GetUserFromTokenUseCase
	*usecase.GetUserDeviceUseCase
	*usecase.OrgDeviceRegistrationUseCase
	*usecase.GetSubmissionByConditionUseCase
	*usecase.GetTotalNrSubmissionByConditionUseCase
	*usecase.GetUserEntityUseCase
	*usecase.GetSubmission4MemoriesFormUseCase
	*usecase.ChildUseCase
}

func (receiver *DeviceController) GetDeviceByID(c *gin.Context) {
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "device id is required",
			},
		)
		return
	}

	device, err := receiver.Get(deviceID)
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
		Code: http.StatusOK,
		Data: &response.DeviceResponseDataV2{
			ID:                device.ID,
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

func (receiver *DeviceController) GetAllDeviceByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	devices, err := receiver.GetDeviceListUseCase.GetDeviceListByUserID(userID)
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

	for _, device := range devices {
		deviceResponse = append(deviceResponse, response.DeviceResponseV2{
			ID:         device.ID,
			DeviceName: device.DeviceName,
		})
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: deviceResponse,
	})
}

func (receiver *DeviceController) GetAllDeviceByOrgID(c *gin.Context) {
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "organization id is required",
			},
		)
		return
	}

	user, err := receiver.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	present := lo.ContainsBy(user.Organizations, func(org entity.SOrganization) bool {
		return org.ID.String() == organizationID
	})
	if !present {
		c.JSON(http.StatusForbidden, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: "access denied",
		})
		return
	}

	userOrg, err := receiver.GetUserOrgInfo(user.ID.String(), organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	if !userOrg.IsManager {
		c.JSON(http.StatusForbidden, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: "access denied",
		})
		return
	}

	devices, err := receiver.GetDeviceListUseCase.GetDeviceListByOrgID(organizationID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			},
		)
		return
	}

	deviceResponse := make([]response.DeviceResponseV2, 0)
	for _, device := range devices {
		deviceResponse = append(deviceResponse, response.DeviceResponseV2{
			ID:         device.ID,
			DeviceName: device.DeviceName,
		})
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: deviceResponse,
	})
}

func (receiver *DeviceController) RegisterOrgDevice(c *gin.Context) {
	var req request.RegisteringDeviceForOrg
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	err := receiver.OrgDeviceRegistrationUseCase.RegisterOrgDevice(req)
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
		Code:    http.StatusOK,
		Message: "device was registered successfully",
	})
}

// InitDeviceV1 Init a device godoc
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

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	deviceID, err := receiver.RegisterDevice(user, req)
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
		Code: http.StatusOK,
		Data: deviceID,
	})
}

// DeactivateDevice Deactivate Device godoc
// @Summary      Deactivate Device
// @Description  Deactivate Device
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param device_id path string true "Device ID"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/device/deactivate/{device_id} [put]
func (receiver *DeviceController) DeactivateDevice(context *gin.Context) {
	deviceID := context.Param("device_id")
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

	err := receiver.UpdateDeviceSheetUseCase.DeactivateDevice(deviceID, req)
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

// ActivateDevice Aactivate Device godoc
// @Summary      Activate Device
// @Description  Aactivate Device
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param device_id path string true "Device ID"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/device/activate/{device_id} [put]
func (receiver *DeviceController) ActivateDevice(context *gin.Context) {
	deviceID := context.Param("device_id")
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
	err := receiver.UpdateDeviceSheetUseCase.ActivateDevice(deviceID, req)
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

// UpdateDevice Update Device godoc
// @Summary      Update Device
// @Description  Update Device
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param device_id path string true "Device ID"
// @Param req body request.UpdateDeviceRequest true "Update Device Params"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/device/{device_id}/update/ [put]
func (receiver *DeviceController) UpdateDevice(context *gin.Context) {
	deviceID := context.Param("device_id")
	if deviceID == "" {
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
	device, err := receiver.UpdateDeviceUseCase.UpdateDevice(deviceID, req)
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
		Code: http.StatusOK,
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
	deviceID := context.Param("device_id")
	if deviceID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Device ID is required",
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
	device, err := receiver.UpdateDeviceUseCase.UpdateDeviceV2(deviceID, req)
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
		// err = usecase.SyncDevice(deviceID)
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

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
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

// TakeNote Save device's Note godoc
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

	userDevices, err := receiver.GetUserDeviceByID(takeNoteRequest.DeviceID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	var userIDs []string
	for _, userDevice := range *userDevices {
		userIDs = append(userIDs, userDevice.UserID.String())
	}
	isExist := gofn.ContainSlice(userIDs, []string{
		user.ID.String(),
	})
	if !isExist {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Code:  http.StatusUnauthorized,
			Error: "Unauthorized",
		})
		return
	}

	err = receiver.TakeNoteUseCase.TakeNote(takeNoteRequest, takeNoteRequest.DeviceID)
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

// SubmitForm Submit a form godoc
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

	// Lấy user từ token
	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusForbidden, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	// Nếu có childID → ưu tiên lấy parentID làm userID
	if req.ChildID != nil && *req.ChildID != "" {
		parentID, err := receiver.GetParentIDByChildID(*req.ChildID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: "Failed to get parent ID: " + err.Error(),
			})
			return
		}

		// Nếu user hiện tại KHÔNG phải parent và cũng KHÔNG phải super admin → cấm
		isSuperAdmin := lo.ContainsBy(user.Roles, func(role entity.SRole) bool {
			return role.Role == entity.SuperAdmin
		})

		if !isSuperAdmin && user.ID.String() != parentID {
			context.JSON(http.StatusUnauthorized, response.FailedResponse{
				Code:  http.StatusUnauthorized,
				Error: "Unauthorized: only parent or super admin can submit form for child",
			})
			return
		}

		req.UserID = parentID // Ưu tiên parentID
	} else {
		// Không có childID → lấy user từ token
		req.UserID = user.ID.String()
	}

	// Gửi câu trả lời form
	err = receiver.AnswerForm(form.ID, req)
	if err != nil {
		context.JSON(http.StatusNotAcceptable, response.FailedResponse{
			Code:  http.StatusNotAcceptable,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Succeed",
	})
}

type getLastSubmissionFromFormRequest struct {
	QRCode string `json:"qr_code" bindding:"required"`
}

// GetLastSubmissionByForm Get last submission from a form godoc
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

	res, err := receiver.GetRecentSubmissionByFormIDUseCase.Execute(req.QRCode, user.ID.String())
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
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

// RefreshAccessToken Refresh Access Token godoc
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
	var req refreshAccessTokenRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := receiver.RefreshAccessTokenUseCase.Execute(req.RefreshToken)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: refreshAccessTokenResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// GetDeviceStatus Get Device Status godoc
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
	deviceID := context.Param("device_id")
	if deviceID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Device ID is required",
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
	// 	return
	// }

	// userDevices, err := receiver.GetUserDeviceUseCase.GetUserDeviceByID(deviceID)
	// if err != nil {
	// 	context.JSON(http.StatusInternalServerError, response.FailedResponse{
	// 		Error: response.Cause{
	// 			Code:    http.StatusInternalServerError,
	// 			Message: err.Error(),
	// 		},
	// 	})
	// 	return
	// }
	// var userIDs []string
	// for _, userDevice := range *userDevices {
	// 	userIDs = append(userIDs, userDevice.UserID.String())
	// }
	// isExist := gofn.ContainSlice(userIDs, []string{
	// 	user.RoleID.String(),
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

	device, err := receiver.GetDeviceByIDUseCase.GetDeviceByID(deviceID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	status, err := receiver.GetDeviceStatusUseCase.Execute(*device)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: status,
	})
}

type reserveRequest struct {
	DeviceID   string `json:"device_id" binding:"required"`
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

	var req reserveRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid req",
		})
		return
	}

	err := receiver.RegisterDeviceUseCase.Reserve(req.DeviceID, req.AppVersion)

	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Reserve Device Succeed",
	})
}

type discoverRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

type discoverResponse struct {
	DeviceID string `json:"device_id"`
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

	res, err := receiver.DiscoverUseCase.Execute(rq.DeviceID)
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
		Code: http.StatusOK,
		Data: discoverResponse{
			DeviceID: res.DeviceID,
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

// GetSignUp Get Sign Up godoc
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

// GetSignUpForm Get Registration Form godoc
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
		Code: http.StatusOK,
		Data: r.Data,
	})
}

type getPresetResponse struct {
	Value *string `json:"value"`
}

// GetPreset2 Get Preset Form godoc
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
		Code: http.StatusOK,
		Data: getPresetResponse{
			Value: r,
		},
	})
}

// GetPreset1 Get Preset Form godoc
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
		Code: http.StatusOK,
		Data: getPresetResponse{
			Value: r,
		},
	})
}

// RegisterFCM Get Register FCM Device godoc
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

// SenNotification Send Notification Logo
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

// ResetCodeCounting Reset Code Counting
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

// GetSubmissionByCondition Get Submission By Condition
// @Summary      Get Submission By Condition
// @Description  Get Submission By Condition
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param  req body request.ResetCodeCountingRequest true "Reset Code Counting Request"
// @Success 200 {object} response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/form/get-submission-by-condition [post]
func (receiver *DeviceController) GetSubmissionByCondition(context *gin.Context) {
	var req request.GetSubmissionByConditionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// Parse atr_value_string
	attr := helper.ParseAtrValueStringToStruct(req.AtrValueString)

	// Ưu tiên lấy userID từ parentID nếu có childID
	if req.ChildID != nil && *req.ChildID != "" {
		// Lấy user từ token
		user, err := receiver.GetUserFromToken(context)
		if err != nil {
			context.JSON(http.StatusForbidden, response.FailedResponse{
				Code:  http.StatusForbidden,
				Error: "Failed to get user from token: " + err.Error(),
			})
			return
		}

		// Kiểm tra nếu là SuperAdmin
		isSuperAdmin := lo.ContainsBy(user.Roles, func(role entity.SRole) bool {
			return role.Role == entity.SuperAdmin
		})

		// Lấy ParentID từ ChildID
		parentID, err := receiver.GetParentIDByChildID(*req.ChildID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: "Failed to get parent ID: " + err.Error(),
			})
			return
		}

		// Nếu không phải super admin và cũng không phải chính parent thì từ chối
		if !isSuperAdmin && user.ID != uuid.MustParse(parentID) {
			context.JSON(http.StatusUnauthorized, response.FailedResponse{
				Code:  http.StatusUnauthorized,
				Error: "Unauthorized: only super admin or child's parent can access",
			})
			return
		}

		// Gán parentID làm userID thực thi
		attr.UserID = parentID
	}

	// Nếu sau bước trên vẫn chưa có UserID thì lấy từ context
	if attr.UserID == "" {
		userIDRaw, exists := context.Get("user_id")
		if !exists {
			context.JSON(http.StatusUnauthorized, response.FailedResponse{
				Code:  http.StatusUnauthorized,
				Error: "Unauthorized: user_id not found in context",
			})
			return
		}

		userID, ok := userIDRaw.(string)
		if !ok {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: "Invalid user_id type in context",
			})
			return
		}

		attr.UserID = userID
	}

	// Gọi usecase trả về list
	res, err := receiver.GetSubmissionByConditionUseCase.Execute(usecase.GetSubmissionByConditionInput{
		UserID:       attr.UserID,
		Key:          attr.Key,
		DB:           attr.DB,
		TimeSort:     attr.TimeSort,
		Quantity:     attr.Quantity,
		DateDuration: attr.DateDuration,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})

}

// GetTotalNrSubmissionByCondition Get Total Number Submission By Condition
// @Summary      Get Total Submission Number
// @Description  Get Total Submission By Condition
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        req body request.GetSubmissionByConditionRequest true "Get Submission Total Request"
// @Success      200 {object} response.SucceedResponse{data=response.GetSubmissionTotalNrResponse}
// @Failure      400 {object} response.FailedResponse
// @Failure      500 {object} response.FailedResponse
// @Router       /v1/form/get-submission-total [post]
func (receiver *DeviceController) GetTotalNrSubmissionByCondition(context *gin.Context) {
	var req request.GetSubmissionByConditionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// Parse atr_value_string
	attr := helper.ParseAtrValueStringToStruct(req.AtrValueString)

	// check question key, question db NR
	if attr.Key != nil && !strings.Contains(*attr.Key, "NR") {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request condition: the question key must contain NR!",
		})
		return
	}

	if attr.DB != nil && !strings.Contains(*attr.DB, "NR") {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request condition: the question db must contain NR!",
		})
		return
	}

	//check duration format

	// Gọi use case trả về tổng
	res, err := receiver.GetTotalNrSubmissionByConditionUseCase.Execute(usecase.GetTotalNrSubmissionByConditionInput{
		UserID:       attr.UserID,
		Key:          attr.Key,
		DB:           attr.DB,
		TimeSort:     attr.TimeSort,
		DateDuration: attr.DateDuration,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

// GetSubmissionChildProfile get the latest submission for a form and user
// @Summary      Get Child Submission Profile
// @Description  Get the latest submission profile for a given form and user
// @Tags         Form
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer {token}"
// @Param        id path int true "Form ID"
// @Success      200 {object} response.SucceedResponse{data=[]response.SubmissionDataItem}
// @Failure      400 {object} response.FailedResponse
// @Failure      500 {object} response.FailedResponse
// @Router       /v1/form/get-submission-child-profile/{id} [get]
func (receiver *DeviceController) GetSubmission4Memories(c *gin.Context) {

	var req request.GetSubmission4MemmoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	form, err := receiver.GetFormByIDUseCase.GetFormByQRCode(req.QrCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.FailedResponse{
			Code:  http.StatusUnauthorized,
			Error: "Unauthorized: user_id not found",
		})
		return
	}

	// Gọi use case hoặc repository
	res, err := receiver.GetSubmission4MemoriesFormUseCase.Execute(repository.GetSubmission4MemoriesFormParam{
		FormID: form.ID,
		UserId: userID.(string),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

// func (receiver *DeviceController) GetSubmission4Memories(c *gin.Context) {
// 	var req request.GetSubmission4MemmoriesRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, response.FailedResponse{
// 			Code:  http.StatusBadRequest,
// 			Error: err.Error(),
// 		})
// 		return
// 	}

// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, response.FailedResponse{
// 			Code:  http.StatusUnauthorized,
// 			Error: "Unauthorized: user_id not found",
// 		})
// 		return
// 	}
// }
