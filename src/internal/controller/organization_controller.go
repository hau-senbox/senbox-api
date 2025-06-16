package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
					UserId:       userOrg.UserID.String(),
					UserNickName: userOrg.UserNickName,
					IsManager:    userOrg.IsManager,
				})
			}

			organizationResponse = append(organizationResponse, response.OrganizationResponse{
				ID:               organization.ID.String(),
				OrganizationName: organization.OrganizationName,
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

func (receiver OrganizationController) GetOrganizationById(context *gin.Context) {
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

	organization, err := receiver.GetOrganizationUseCase.GetOrganizationById(organizationID)
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
			UserId:       userOrg.UserID.String(),
			UserNickName: userOrg.UserNickName,
			IsManager:    userOrg.IsManager,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.OrganizationResponse{
			ID:               organization.ID.String(),
			OrganizationName: organization.OrganizationName,
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
			UserId:       userOrg.UserID.String(),
			UserNickName: userOrg.UserNickName,
			IsManager:    userOrg.IsManager,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.OrganizationResponse{
			ID:               organization.ID.String(),
			OrganizationName: organization.OrganizationName,
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

	id, err := strconv.Atoi(organizationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	users, err := receiver.GetOrganizationUseCase.GetAllUserByOrganization(uint(id))
	if err != nil {
		context.JSON(400, err.Error())
		return
	}

	var res []response.GetOrgManagerInfoResponse
	for _, user := range users {
		res = append(res, response.GetOrgManagerInfoResponse{
			UserId:       user.UserID.String(),
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
				UserId:             application.UserID.String(),
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
		UserId:             application.UserID.String(),
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
