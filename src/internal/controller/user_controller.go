package controller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
)

type UserController struct {
	GetUsersUseCase       usecase.GetUsersUseCase
	CreateUserUseCase     usecase.CreateUserUseCase
	ChangePasswordUseCase usecase.ChangePasswordUseCase
}

// Get user list
// @Summary Get user list
// @Description Get user list
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} response.UserListResponse
// @Failure 400 {object} response.FailedResponse
// @Router /v1/admin/users [get]
func (receiver UserController) GetUserList(context *gin.Context) {
	users, err := receiver.GetUsersUseCase.GetUsers(value.User)
	if err != nil {
		context.JSON(400, err.Error())
		return
	}
	var userListResponse []response.UserListResponseData
	for _, user := range users {
		userListResponse = append(userListResponse, response.UserListResponseData{
			UserID: user.UserId,
			Name:   user.Fullname,
		})
	}
	context.JSON(200, response.UserListResponse{
		Data: userListResponse,
	})
}

// Create user godoc
// @Summary Create user
// @Description Create user
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param body body request.CreateUserRequest true "Create user request"
// @Success 200 {object} response.CreateUserResponse
// @Failure 400 {object} response.FailedResponse
// @Router /v1/admin/user/create [post]
func (receiver UserController) CreateUser(context *gin.Context) {
	var rq request.CreateUserRequest
	err := context.ShouldBindJSON(&rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.SucceedResponse{
			Data: err.Error(),
		})
		return
	}
	user, err := receiver.CreateUserUseCase.CreateUser(rq)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.SucceedResponse{
			Data: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: gin.H{
			"userId":    user.UserId,
			"username":  user.Username,
			"fullname":  user.Fullname,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
		},
	})
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type changePasswordResponse struct {
	Success bool `json:"success"`
}

// Change password godoc
// @Summary Change password
// @Description Change password
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param body body changePasswordRequest true "Change password request"
// @Params body body changePasswordRequest true "Change password request"
// @Success 200 {object} changePasswordResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 422 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/me/new-password [put]
func (receiver UserController) ChangePassword(context *gin.Context) {
	var rq changePasswordRequest
	err := context.ShouldBindJSON(&rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	userId := context.GetString("user_id")

	err = receiver.ChangePasswordUseCase.VerifyNewPassword(rq.NewPassword)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusUnprocessableEntity,
				Message: err.Error(),
			},
		})
		return
	}

	user, err := receiver.ChangePasswordUseCase.ValidateCurrentPassword(userId, rq.CurrentPassword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusUnprocessableEntity,
				Message: "Current password is not valid",
			},
		})
		return
	}

	err = receiver.ChangePasswordUseCase.ChangePassword(user, rq.NewPassword)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: "Failed to change password, please try again later.",
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: changePasswordResponse{
			Success: true,
		},
	})
}
