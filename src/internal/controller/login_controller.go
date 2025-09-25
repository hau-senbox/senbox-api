package controller

import (
	"errors"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type LoginController struct {
	DBConn *gorm.DB
	usecase.AuthorizeUseCase
}

func (receiver LoginController) Login(c *gin.Context) {
	var req request.UserLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	data, err := receiver.LoginInputDao(req)

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
		Code: http.StatusOK,
		Data: data,
	})
}

// Login godoc
// @Summary      Retrieve a token
// @Description  login using username and password
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param req body request.LoginInputReq true "Login Params"
// @Success      200  {object}  response.LoginResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/login [post]
func (receiver LoginController) UserLogin(c *gin.Context) {
	var req request.UserLoginFromDeviceReqest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	data, err := receiver.UserLoginUsecase(req, value.ForScan)

	if err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Error:   err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: *data,
	})
}

func (receiver LoginController) RefreshToken(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "missing or invalid Authorization header",
		})
		return
	}
	tokenString := strings.Split(authorizationHeader, " ")[1]

	// validate token
	_, err := receiver.SessionRepository.ValidateToken(tokenString)
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// extract user_id
	userID, err := receiver.SessionRepository.ExtractUserIDIgnoreExp(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cannot extract user id"})
		return
	}

	// generate new token
	data, err := receiver.AuthorizeUseCase.RefreshToken(*userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success response
	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: *data,
	})
}
