package controller

import (
	"bufio"
	"errors"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationController struct {
	*usecase.GetOrganizationUseCase
	*usecase.CreateOrganizationUseCase
	*usecase.UserJoinOrganizationUseCase
	*usecase.GetUserFromTokenUseCase
	*usecase.GetOrgFormApplicationUseCase
	*usecase.ApproveOrgFormApplicationUseCase
	*usecase.BlockOrgFormApplicationUseCase
	*usecase.CreateOrgFormApplicationUseCase
	*usecase.UploadOrgAvatarUseCase
	*usecase.OrganizationSettingUsecase
}

func (receiver OrganizationController) GetAllOrganization(context *gin.Context) {
	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	organizations, err := receiver.GetOrganizationUseCase.GetAllOrganization(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	var organizationResponse []response.OrganizationResponse
	if len(organizations) > 0 {
		organizationResponse = make([]response.OrganizationResponse, 0)
		for _, organization := range organizations {
			managers := make([]response.GetOrgManagerInfoResponse, 0)
			for _, userOrg := range organization.UserOrgs {
				managers = append(managers, response.GetOrgManagerInfoResponse{
					UserID:       userOrg.UserID.String(),
					UserNickName: userOrg.UserNickName,
					IsManager:    userOrg.IsManager,
				})
			}

			organizationResponse = append(organizationResponse, response.OrganizationResponse{
				ID:               organization.ID.String(),
				OrganizationName: organization.OrganizationName,
				Avatar:           organization.Avatar,
				AvatarURL:        organization.AvatarURL,
				Address:          organization.Address,
				Description:      organization.Description,
				Managers:         managers,
			})
		}
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: organizationResponse,
	})
}

func (receiver OrganizationController) Check4App(context *gin.Context) {
	deviceID := context.Param("device_id")
	if deviceID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "deviceID is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	organizationID := context.Param("organization_id")
	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "organizationID is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	isOK, err := receiver.CheckDeviceInOrg4App(deviceID, organizationID)

	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: err.Error(),
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: isOK,
	})
}

func (receiver OrganizationController) GetOrganizationByID(context *gin.Context) {
	organizationID := context.Param("id")
	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	organization, err := receiver.GetOrganizationUseCase.GetOrganizationByID(organizationID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	managers := make([]response.GetOrgManagerInfoResponse, 0)
	for _, userOrg := range organization.UserOrgs {
		managers = append(managers, response.GetOrgManagerInfoResponse{
			UserID:       userOrg.UserID.String(),
			UserNickName: userOrg.UserNickName,
			IsManager:    userOrg.IsManager,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.OrganizationResponse{
			ID:               organization.ID.String(),
			OrganizationName: organization.OrganizationName,
			Avatar:           organization.Avatar,
			AvatarURL:        organization.AvatarURL,
			Address:          organization.Address,
			Description:      organization.Description,
			Managers:         managers,
		},
	})
}

func (receiver OrganizationController) GetOrganizationByName(context *gin.Context) {
	organizationName := context.Query("organization_name")
	if organizationName == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "organization name is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	organization, err := receiver.GetOrganizationUseCase.GetByName(organizationName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	managers := make([]response.GetOrgManagerInfoResponse, 0)
	for _, userOrg := range organization.UserOrgs {
		managers = append(managers, response.GetOrgManagerInfoResponse{
			UserID:       userOrg.UserID.String(),
			UserNickName: userOrg.UserNickName,
			IsManager:    userOrg.IsManager,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.OrganizationResponse{
			ID:               organization.ID.String(),
			OrganizationName: organization.OrganizationName,
			Avatar:           organization.Avatar,
			AvatarURL:        organization.AvatarURL,
			Address:          organization.Address,
			Description:      organization.Description,
			Managers:         managers,
		},
	})
}

func (receiver OrganizationController) CreateOrganization(context *gin.Context) {
	var req request.CreateOrganizationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.CreateOrganizationUseCase.CreateOrganization(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Organization created successfully",
	})
}

func (receiver OrganizationController) UserJoinOrganization(context *gin.Context) {
	var req request.UserJoinOrganizationRequest
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

	req.UserID = user.ID.String()

	err = receiver.UserJoinOrganizationUseCase.UserJoinOrganization(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user join organization successfully",
	})
}

func (receiver OrganizationController) GetAllUserByOrganization(context *gin.Context) {
	organizationID := context.Param("id")
	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	users, err := receiver.GetOrganizationUseCase.GetAllUserByOrganization(organizationID)
	if err != nil {
		context.JSON(400, err.Error())
		return
	}

	var res []response.GetOrgManagerInfoResponse
	for _, user := range users {
		res = append(res, response.GetOrgManagerInfoResponse{
			UserID:       user.UserID.String(),
			UserNickName: user.UserNickName,
			IsManager:    user.IsManager,
		})
	}

	context.JSON(200, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

var defaultTime = time.Time{}

func (receiver OrganizationController) GetAllOrgFormApplication(context *gin.Context) {
	applications, err := receiver.GetOrgFormApplicationUseCase.GetAllOrgFormApplication()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	var applicationResponse []response.OrgFormApplicationResponse
	if len(applications) > 0 {
		applicationResponse = make([]response.OrgFormApplicationResponse, 0)
		for _, application := range applications {
			res := response.OrgFormApplicationResponse{
				ID:                 application.ID,
				OrganizationName:   application.OrganizationName,
				ApplicationContent: application.ApplicationContent,
				Status:             application.Status.String(),
				ApprovedAt:         "",
				CreatedAt:          application.CreatedAt.Format("2006-01-02 15:04:05"),
				UserID:             application.UserID.String(),
			}
			if application.ApprovedAt != defaultTime {
				res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
			}
			applicationResponse = append(applicationResponse, res)
		}
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: applicationResponse,
	})
}

func (receiver OrganizationController) GetOrgFormApplicationByID(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	application, err := receiver.GetOrgFormApplicationUseCase.GetOrgFormApplicationByID(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	res := response.OrgFormApplicationResponse{
		ID:                 application.ID,
		OrganizationName:   application.OrganizationName,
		ApplicationContent: application.ApplicationContent,
		Status:             application.Status.String(),
		ApprovedAt:         "",
		CreatedAt:          application.CreatedAt.Format("2006-01-02 15:04:05"),
		UserID:             application.UserID.String(),
	}
	if application.ApprovedAt != defaultTime {
		res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver OrganizationController) ApproveOrgFormApplication(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	err = receiver.ApproveOrgFormApplicationUseCase.ApproveOrgFormApplication(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application approved successfully",
	})
}

func (receiver OrganizationController) BlockOrgFormApplication(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	err = receiver.ApproveOrgFormApplicationUseCase.BlockOrgFormApplication(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application blocked successfully",
	})
}

func (receiver OrganizationController) CreateOrgFormApplication(context *gin.Context) {
	var req request.CreateOrgFormApplicationRequest
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

	req.UserID = user.ID.String()

	err = receiver.CreateOrgFormApplicationUseCase.CreateOrgFormApplication(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application created successfully",
	})
}

func (receiver OrganizationController) UploadAvatar(context *gin.Context) {
	fileHeader, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	orgID := context.PostForm("organization_id")
	if orgID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "organization id is required",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	defer file.Close()

	dataBytes := make([]byte, fileHeader.Size)
	if _, err := bufio.NewReader(file).Read(dataBytes); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	url, img, err := receiver.UploadOrgAvatarUseCase.UploadAvatar(orgID, dataBytes, fileHeader.Filename)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if url == nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Upload avatar failed",
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Upload avatar successfully",
		Data: response.ImageResponse{
			ImageName: img.ImageName,
			Key:       img.Key,
			Extension: img.Extension,
			Url:       *url,
			Width:     img.Width,
			Height:    img.Height,
		},
	})
}

func (receiver OrganizationController) UploadOrgSetting(c *gin.Context) {
	var req request.UploadOrgSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	err := receiver.OrganizationSettingUsecase.UploadOrgSetting(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to upload organization setting",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Upload success",
		Data:    nil,
	})
}

func (receiver OrganizationController) GetOrgSetting(c *gin.Context) {
	deviceID := c.Param("device_id")

	if deviceID == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing device ID",
		})
		return
	}

	orgSetting, err := receiver.OrganizationSettingUsecase.GetOrgSetting(deviceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, response.SucceedResponse{
				Code:    http.StatusOK,
				Message: "Not Found",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get organization setting",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    orgSetting,
	})
}

// Org Setting News
func (receiver *OrganizationController) UploadOrgSettingNewsDevice(c *gin.Context) {
	var req request.UploadOrgSettingDeviceNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	orgID := c.Param("organization_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing organization ID",
		})
		return
	}

	req.OrganizationID = orgID

	if err := receiver.OrganizationSettingUsecase.UploadOrgSettingNewsDevice(req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to upload organization device news setting",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Device news updated successfully",
	})
}

func (receiver *OrganizationController) UploadOrgSettingNewsPortal(c *gin.Context) {
	var req request.UploadOrgSettingPortalNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	orgID := c.Param("organization_id")

	if orgID == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing organization ID",
		})
		return
	}

	req.OrganizationID = orgID

	if err := receiver.OrganizationSettingUsecase.UploadOrgSettingNewsPortal(req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to upload organization portal news setting",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Portal news updated successfully",
	})
}

func (receiver *OrganizationController) GetOrgSettingNews(c *gin.Context) {
	orgID := c.Param("organization_id")

	data, err := receiver.OrganizationSettingUsecase.GetOrgSettingNews(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, response.SucceedResponse{
				Code:    http.StatusOK,
				Message: "Record not found",
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get organization device news setting",
			Error:   err.Error(),
		})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, response.FailedResponse{
			Code:    http.StatusNotFound,
			Message: "Device news not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}
