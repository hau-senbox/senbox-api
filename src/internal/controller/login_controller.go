package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
)

type LoginController struct {
	DBConn *gorm.DB
	usecase.AuthorizeUseCase
}

func (receiver LoginController) Login(c *gin.Context) {
	var req request.LoginInputReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			},
		)
		return
	}

	data, err := receiver.AuthorizeUseCase.LoginInputDao(req)

	if err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			},
		)
		return
	}

	c.JSON(http.StatusOK, data)
	return
}

// Login godoc
// @Summary      Retrieve a token
// @Description  login using loginId and password
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param req body request.LoginInputReq true "Login Params"
// @Success      200  {object}  response.LoginResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/login [post]
func (receiver LoginController) LoginV1(c *gin.Context) {
	var req request.LoginInputReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			},
		)
		return
	}

	data, err := receiver.AuthorizeUseCase.LoginInputDao(req)

	if err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			},
		)
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		Data: data,
	})
	return
}
