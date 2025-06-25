package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FunctionClaimController struct {
	*usecase.GetFunctionClaimUseCase
	*usecase.CreateFunctionClaimUseCase
	*usecase.UpdateFunctionClaimUseCase
	*usecase.DeleteFunctionClaimUseCase
}

func (receiver *FunctionClaimController) GetAllFunctionClaim(context *gin.Context) {
	claims, err := receiver.GetFunctionClaimUseCase.GetAllFunctionClaim()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var functionClaimListResponse []response.FunctionClaimListResponseData
	for _, function := range claims {
		var functionClaimPermissionListResponse []response.FunctionClaimPermissionListResponseData
		for _, permission := range function.ClaimPermissions {
			functionClaimPermissionListResponse = append(functionClaimPermissionListResponse, response.FunctionClaimPermissionListResponseData{
				ID:             permission.ID,
				PermissionName: permission.PermissionName,
			})
		}

		functionClaimListResponse = append(functionClaimListResponse, response.FunctionClaimListResponseData{
			ID:           function.ID,
			FunctionName: function.FunctionName,
			Permissions:  functionClaimPermissionListResponse,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "function claims was get successfully",
		Data:    functionClaimListResponse,
	})
}

func (receiver *FunctionClaimController) GetFunctionClaimByID(context *gin.Context) {
	claimID := context.Param("id")
	if claimID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "function claim id is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(claimID, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "function claim id is invalid",
			},
		)
		return
	}

	functionClaim, err := receiver.GetFunctionClaimUseCase.GetFunctionClaimByID(request.GetFunctionClaimByIDRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var functionClaimPermissionListResponse []response.FunctionClaimPermissionListResponseData
	for _, permission := range functionClaim.ClaimPermissions {
		functionClaimPermissionListResponse = append(functionClaimPermissionListResponse, response.FunctionClaimPermissionListResponseData{
			ID:             permission.ID,
			PermissionName: permission.PermissionName,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.FunctionClaimResponse{
			ID:           functionClaim.ID,
			FunctionName: functionClaim.FunctionName,
			Permissions:  functionClaimPermissionListResponse,
		},
	})
}

func (receiver *FunctionClaimController) GetFunctionClaimByName(context *gin.Context) {
	functionName := context.Param("function_name")
	if functionName == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "claim name is required",
			},
		)
		return
	}

	functionClaim, err := receiver.GetFunctionClaimUseCase.GetFunctionClaimByName(request.GetFunctionClaimByNameRequest{FunctionName: functionName})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var functionClaimPermissionListResponse []response.FunctionClaimPermissionListResponseData
	for _, permission := range functionClaim.ClaimPermissions {
		functionClaimPermissionListResponse = append(functionClaimPermissionListResponse, response.FunctionClaimPermissionListResponseData{
			ID:             permission.ID,
			PermissionName: permission.PermissionName,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.FunctionClaimResponse{
			ID:           functionClaim.ID,
			FunctionName: functionClaim.FunctionName,
			Permissions:  functionClaimPermissionListResponse,
		},
	})
}

func (receiver *FunctionClaimController) CreateFunctionClaim(context *gin.Context) {
	var req request.CreateFunctionClaimRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.CreateFunctionClaimUseCase.CreateFunctionClaim(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "function claim was create successfully",
	})
}

func (receiver *FunctionClaimController) CreateFunctionClaims(context *gin.Context) {
	var req request.CreateFunctionClaimsRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.CreateFunctionClaimUseCase.CreateFunctionClaims(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "function claims was create successfully",
	})
}

func (receiver *FunctionClaimController) UpdateFunctionClaim(context *gin.Context) {
	var req request.UpdateFunctionClaimRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UpdateFunctionClaimUseCase.UpdateFunctionClaim(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "function claim was update successfully",
	})
}

func (receiver *FunctionClaimController) DeleteFunctionClaim(context *gin.Context) {
	claimID := context.Param("id")
	if claimID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "function claim id is required",
			},
		)
		return
	}

	id, err := strconv.ParseUint(claimID, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "function claim id is invalid",
			},
		)
		return
	}

	err = receiver.DeleteFunctionClaimUseCase.DeleteFunctionClaim(request.DeleteFunctionClaimRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "function claim was delete successfully",
	})
}
