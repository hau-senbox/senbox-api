package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	*usecase.GetOrganizationUseCase
	*usecase.CreateOrganizationUseCase
	*usecase.UserJoinOrganizationUseCase
	*usecase.GetUserFromTokenUseCase
}

func (receiver OrganizationController) GetAllOrganization(context *gin.Context) {
	organizations, err := receiver.GetOrganizationUseCase.GetAllOrganization()
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
			organizationResponse = append(organizationResponse, response.OrganizationResponse{
				ID:               organization.ID,
				OrganizationName: organization.OrganizationName,
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

	id, err := strconv.Atoi(organizationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	organization, err := receiver.GetOrganizationUseCase.GetOrganizationById(uint(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.OrganizationResponse{
			ID:               organization.ID,
			OrganizationName: organization.OrganizationName,
			Address:          organization.Address,
			Description:      organization.Description,
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

	req.UserId = user.ID.String()

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

	var userListResponse []response.UserListResponseData
	for _, user := range users {
		userListResponse = append(userListResponse, response.UserListResponseData{
			UserID:   user.ID.String(),
			Fullname: user.Username,
		})
	}

	context.JSON(200, response.SucceedResponse{
		Code: http.StatusOK,
		Data: userListResponse,
	})
}
