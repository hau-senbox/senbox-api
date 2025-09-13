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
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			},
		)
		return
	}

	req.Method = ctx.Request.Method
	req.Endpoint = ctx.FullPath()

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
			Code: http.StatusOK,
			Data: "Accounts log created successfully",
		},
	)
}
