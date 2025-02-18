package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RolePolicyController struct {
	*usecase.GetRolePolicyUseCase
	*usecase.CreateRolePolicyUseCase
	*usecase.UpdateRolePolicyUseCase
	*usecase.DeleteRolePolicyUseCase
}

func (receiver *RolePolicyController) GetAllRolePolicy(context *gin.Context) {
	policies, err := receiver.GetRolePolicyUseCase.GetAllRolePolicy()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	var policyListResponse []response.RolePolicyListResponseData
	for _, policy := range policies {
		policyListResponse = append(policyListResponse, response.RolePolicyListResponseData{
			ID:         policy.ID,
			PolicyName: policy.PolicyName,
		})
	}

	context.JSON(http.StatusOK, response.RolePolicyListResponse{
		Data: policyListResponse,
	})
}

func (receiver *RolePolicyController) GetRolePolicyById(context *gin.Context) {
	policyId := context.Param("id")
	if policyId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "policy id is required",
				},
			},
		)
		return
	}

	id, err := strconv.ParseUint(policyId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "policy id is invalid",
				},
			},
		)
		return
	}

	userRolePolicy, err := receiver.GetRolePolicyUseCase.GetRolePolicyById(request.GetRolePolicyByIdRequest{ID: uint(id)})
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
		Data: response.RolePolicyResponse{
			ID:          userRolePolicy.ID,
			PolicyName:  userRolePolicy.PolicyName,
			Description: userRolePolicy.Description,
		},
	})
}

func (receiver *RolePolicyController) GetRolePolicyByName(context *gin.Context) {
	policyName := context.Param("policy_name")
	if policyName == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "policy name is required",
				},
			},
		)
		return
	}

	userRolePolicy, err := receiver.GetRolePolicyUseCase.GetRolePolicyByName(request.GetRolePolicyByNameRequest{PolicyName: policyName})
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
		Data: response.RolePolicyResponse{
			ID:          userRolePolicy.ID,
			PolicyName:  userRolePolicy.PolicyName,
			Description: userRolePolicy.Description,
		},
	})
}

func (receiver *RolePolicyController) CreateRolePolicy(context *gin.Context) {
	var req request.CreateRolePolicyRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := receiver.CreateRolePolicyUseCase.Create(req)
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
			Message: "role policy was create successfully",
		},
	})
}

func (receiver *RolePolicyController) UpdateRolePolicy(context *gin.Context) {
	var req request.UpdateRolePolicyRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := receiver.UpdateRolePolicyUseCase.UpdateRolePolicy(req)
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
			Message: "role policy was update successfully",
		},
	})
}

func (receiver *RolePolicyController) DeleteRolePolicy(context *gin.Context) {
	policyId := context.Param("id")
	if policyId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "policy id is required",
				},
			},
		)
		return
	}

	id, err := strconv.ParseUint(policyId, 10, 32)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "policy id is invalid",
				},
			},
		)
		return
	}

	err = receiver.DeleteRolePolicyUseCase.DeleteRolePolicy(request.DeleteRolePolicyRequest{ID: uint(id)})
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
			Message: "role policy was delete successfully",
		},
	})
}
