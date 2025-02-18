package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleClaimController struct {
	*usecase.GetRoleClaimUseCase
	*usecase.CreateRoleClaimUseCase
	*usecase.UpdateRoleClaimUseCase
	*usecase.DeleteRoleClaimUseCase
}

func (receiver *RoleClaimController) GetAllRoleClaim(context *gin.Context) {
	claims, err := receiver.GetRoleClaimUseCase.GetAllRoleClaim()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	var roleClaimListResponse []response.RoleClaimListResponseData
	for _, role := range claims {
		roleClaimListResponse = append(roleClaimListResponse, response.RoleClaimListResponseData{
			ID:        role.ID,
			ClaimName: role.ClaimName,
		})
	}

	context.JSON(http.StatusOK, response.RoleClaimListResponse{
		Data: roleClaimListResponse,
	})
}

func (receiver *RoleClaimController) GetAllRoleClaimByRole(context *gin.Context) {
	roleId := context.Param("role_id")
	if roleId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "role id is required",
				},
			},
		)
		return
	}

	id, err := strconv.ParseUint(roleId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "role id is invalid",
				},
			},
		)
		return
	}

	claims, err := receiver.GetRoleClaimUseCase.GetAllRoleClaimByRole(request.GetAllRoleClaimByRoleRequest{RoleId: uint(id)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	var roleClaimListResponse []response.RoleClaimListResponseData
	for _, role := range claims {
		roleClaimListResponse = append(roleClaimListResponse, response.RoleClaimListResponseData{
			ID:        role.ID,
			ClaimName: role.ClaimName,
		})
	}

	context.JSON(http.StatusOK, response.RoleClaimListResponse{
		Data: roleClaimListResponse,
	})
}

func (receiver *RoleClaimController) GetRoleClaimById(context *gin.Context) {
	claimId := context.Param("id")
	if claimId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "claim id is required",
				},
			},
		)
		return
	}

	id, err := strconv.ParseUint(claimId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "claim id is invalid",
				},
			},
		)
		return
	}

	roleClaim, err := receiver.GetRoleClaimUseCase.GetRoleClaimById(request.GetRoleClaimByIdRequest{ID: uint(id)})
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
		Data: response.RoleClaimResponse{
			ID:         roleClaim.ID,
			ClaimName:  roleClaim.ClaimName,
			ClaimValue: roleClaim.ClaimValue,
		},
	})
}

func (receiver *RoleClaimController) GetRoleClaimByName(context *gin.Context) {
	claimName := context.Param("claim_name")
	if claimName == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "claim name is required",
				},
			},
		)
		return
	}

	userRole, err := receiver.GetRoleClaimUseCase.GetRoleClaimByName(request.GetRoleClaimByNameRequest{ClaimName: claimName})
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
		Data: response.RoleClaimResponse{
			ID:         userRole.ID,
			ClaimName:  userRole.ClaimName,
			ClaimValue: userRole.ClaimValue,
		},
	})
}

func (receiver *RoleClaimController) CreateRoleClaim(context *gin.Context) {
	var req request.CreateRoleClaimRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := receiver.CreateRoleClaimUseCase.CreateRoleClaim(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "role claim was create successfully",
		},
	})
}

func (receiver *RoleClaimController) CreateRoleClaims(context *gin.Context) {
	var req request.CreateRoleClaimsRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := receiver.CreateRoleClaimUseCase.CreateRoleClaims(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "role claims was create successfully",
		},
	})
}

func (receiver *RoleClaimController) UpdateRoleClaim(context *gin.Context) {
	var req request.UpdateRoleClaimRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := receiver.UpdateRoleClaimUseCase.UpdateRoleClaim(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "role claim was update successfully",
		},
	})
}

func (receiver *RoleClaimController) DeleteRoleClaim(context *gin.Context) {
	claimId := context.Param("id")
	if claimId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "claim id is required",
				},
			},
		)
		return
	}

	id, err := strconv.ParseUint(claimId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "claim id is invalid",
				},
			},
		)
		return
	}

	err = receiver.DeleteRoleClaimUseCase.DeleteRoleClaim(request.DeleteRoleClaimRequest{ID: uint(id)})
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "role claim was delete successfully",
		},
	})
}
