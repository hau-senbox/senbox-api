package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	*usecase.GetRoleUseCase
	*usecase.CreateRoleUseCase
	*usecase.UpdateRoleUseCase
	*usecase.DeleteRoleUseCase
}

func (receiver *RoleController) GetAllRole(context *gin.Context) {
	roles, err := receiver.GetRoleUseCase.GetAllRole()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var roleListResponse []response.RoleListResponseData
	for _, role := range roles {
		roleListResponse = append(roleListResponse, response.RoleListResponseData{
			ID:       role.ID,
			RoleName: role.Role.String(),
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: roleListResponse,
	})
}

func (receiver *RoleController) GetRoleByID(context *gin.Context) {
	roleID := context.Param("id")
	if roleID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Role ID is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Role ID is invalid",
			},
		)
		return
	}

	userRole, err := receiver.GetRoleUseCase.GetRoleByID(request.GetRoleByIDRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.RoleResponse{
			ID:       userRole.ID,
			RoleName: userRole.Role.String(),
		},
	})
}

func (receiver *RoleController) GetRoleByName(context *gin.Context) {
	roleName := context.Param("role")
	if roleName == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "role name is required",
			},
		)
		return
	}

	userRole, err := receiver.GetRoleUseCase.GetRoleByName(request.GetRoleByNameRequest{RoleName: roleName})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.RoleResponse{
			ID:       userRole.ID,
			RoleName: userRole.Role.String(),
		},
	})
}

func (receiver *RoleController) CreateRole(context *gin.Context) {
	var req request.CreateRoleRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.Create(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user role was create successfully",
	})
}

func (receiver *RoleController) UpdateRole(context *gin.Context) {
	var req request.UpdateRoleRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UpdateRoleUseCase.UpdateRole(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "role was update successfully",
	})
}

func (receiver *RoleController) DeleteRole(context *gin.Context) {
	roleID := context.Param("id")
	if roleID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Role ID is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(roleID, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "Role ID is invalid",
			},
		)
		return
	}

	err = receiver.DeleteRoleUseCase.DeleteRole(request.DeleteRoleRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "role was delete successfully",
	})
}
