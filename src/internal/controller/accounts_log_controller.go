package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type AccountsLogController struct {
	AccountsLogUseCase *usecase.AccountsLogUseCase
}

func (c *AccountsLogController) CreateAccountsLog(ctx *gin.Context) {
	var req request.CreateAccountsLogRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Error:   err.Error(),
				Message: "Invalid request body",
			},
		)
		return
	}

	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Error:   "user_id not found",
				Message: "Failed to create accounts log",
			},
		)
		return
	}
	req.UserID = userID.(string)

	err := c.AccountsLogUseCase.CreateAccountsLog(req)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Error:   err.Error(),
				Message: "Failed to create accounts log",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK, response.SucceedResponse{
			Code:    http.StatusOK,
			Message: "Accounts log created successfully",
			Data:    nil,
		},
	)
}
