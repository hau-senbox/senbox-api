package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type LogoutController struct {
	AuthorizeUseCase *usecase.AuthorizeUseCase
}

func (receiver LogoutController) UserLogout(c *gin.Context) {
	var req request.UserLogoutReqeust
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	err := receiver.AuthorizeUseCase.UserLogoutUsecase(c, req)

	if err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Logout successfully",
	})
}
