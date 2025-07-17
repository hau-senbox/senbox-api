package controller

import (
	"net/http"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type UserTokenFCMController struct {
	CreateUserTokenFCMUseCase *usecase.CreateUserTokenFCMUseCase
	GetUserTokenFCMUseCase    *usecase.GetUserTokenFCMUseCase
}

type createFCMTokenRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

func (receiver *UserTokenFCMController) CreateFCMToken(ctx *gin.Context) {

	var req createFCMTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.CreateUserTokenFCMUseCase.CreateToken(req.UserID, req.DeviceID, req.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "token created successfully",
	})

}

func (receiver *UserTokenFCMController) GetAllFCMToken(ctx *gin.Context) {

	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	tokens, err := receiver.GetUserTokenFCMUseCase.GetAllToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "tokens retrieved successfully",
		Data:    tokens,
	})
}