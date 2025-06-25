package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FunctionClaimPermissionController struct {
	*usecase.GetFunctionClaimPermissionUseCase
	*usecase.CreateFunctionClaimPermissionUseCase
	*usecase.UpdateFunctionClaimPermissionUseCase
	*usecase.DeleteFunctionClaimPermissionUseCase
}

func (receiver *FunctionClaimPermissionController) GetAllFunctionClaimPermission(context *gin.Context) {
	functionClaimID := context.Param("function_claim_id")
	if functionClaimID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "function claim id is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(functionClaimID, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "function claim id is invalid",
			},
		)
		return
	}

	permissions, err := receiver.GetFunctionClaimPermissionUseCase.GetAllFunctionClaimPermission(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var permissionListResponse []response.FunctionClaimPermissionListResponseData
	for _, permission := range permissions {
		permissionListResponse = append(permissionListResponse, response.FunctionClaimPermissionListResponseData{
			ID:             permission.ID,
			PermissionName: permission.PermissionName,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: permissionListResponse,
	})
}

func (receiver *FunctionClaimPermissionController) GetFunctionClaimPermissionByID(context *gin.Context) {
	permissionID := context.Param("id")
	if permissionID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "permission id is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "permission id is invalid",
			},
		)
		return
	}

	functionClaimPermission, err := receiver.GetFunctionClaimPermissionUseCase.GetFunctionClaimPermissionByID(request.GetFunctionClaimPermissionByIDRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.FunctionClaimPermissionResponse{
			ID:             functionClaimPermission.ID,
			PermissionName: functionClaimPermission.PermissionName,
		},
	})
}

func (receiver *FunctionClaimPermissionController) GetFunctionClaimPermissionByName(context *gin.Context) {
	permissionName := context.Param("permission_name")
	if permissionName == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "permission name is required",
			},
		)
		return
	}

	functionClaimPermission, err := receiver.GetFunctionClaimPermissionUseCase.GetFunctionClaimPermissionByName(request.GetFunctionClaimPermissionByNameRequest{PermissionName: permissionName})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.FunctionClaimPermissionResponse{
			ID:             functionClaimPermission.ID,
			PermissionName: functionClaimPermission.PermissionName,
		},
	})
}

func (receiver *FunctionClaimPermissionController) CreateFunctionClaimPermission(context *gin.Context) {
	var req request.CreateFunctionClaimPermissionRequest
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
		Message: "function claim permission was create successfully",
	})
}

func (receiver *FunctionClaimPermissionController) UpdateRoleClaimPermission(context *gin.Context) {
	var req request.UpdateFunctionClaimPermissionRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UpdateFunctionClaimPermissionUseCase.UpdateFunctionClaimPermission(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "function claim permission was update successfully",
	})
}

func (receiver *FunctionClaimPermissionController) DeleteRoleClaimPermission(context *gin.Context) {
	permissionID := context.Param("id")
	if permissionID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "permission id is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "permission id is invalid",
			},
		)
		return
	}

	err = receiver.DeleteFunctionClaimPermissionUseCase.DeleteFunctionClaimPermission(request.DeleteFunctionClaimPermissionRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "function claim permission was delete successfully",
	})
}
