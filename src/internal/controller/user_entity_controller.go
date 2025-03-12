package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserEntityController struct {
	*usecase.GetUserEntityUseCase
	*usecase.CreateUserEntityUseCase
	*usecase.UpdateUserEntityUseCase
	*usecase.AuthorizeUseCase
}

func (receiver *UserEntityController) GetAllUserEntity(context *gin.Context) {
	role := context.Request.URL.Query().Get("role")
	if role == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "role is required",
				},
			},
		)
		return
	}

	users, err := receiver.GetUserEntityUseCase.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	var userResponse []response.UserEntityResponseData
	for _, user := range users {
		var roles []string
		for _, r := range user.Roles {
			if role != "" && strings.ToLower(role) != "all" {
				if !strings.EqualFold(r.RoleName, role) {
					continue
				}

				roles = append(roles, r.RoleName)
				break
			}
			roles = append(roles, r.RoleName)
		}

		if len(roles) > 0 {
			userResponse = append(userResponse, response.UserEntityResponseData{
				ID:       user.ID.String(),
				Username: user.Username,
				Roles:    roles,
			})
		}
	}

	context.JSON(http.StatusOK, response.UserEntityDataResponse{
		Data: userResponse,
	})
}

func (receiver *UserEntityController) GetChildrenOfGuardian(context *gin.Context) {
	userId := context.Param("id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "user id is required",
				},
			},
		)
		return
	}

	users, err := receiver.GetUserEntityUseCase.GetChildrenOfGuardian(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	context.JSON(http.StatusOK, response.UserEntityDataResponse{
		Data: *users,
	})
}

func (receiver *UserEntityController) GetUserEntityById(context *gin.Context) {
	userId := context.Param("id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "user id is required",
				},
			},
		)
		return
	}

	userEntity, err := receiver.GetUserById(request.GetUserEntityByIdRequest{ID: userId})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	var roleListResponse []response.RoleListResponseData
	if len(userEntity.Roles) > 0 {
		roleListResponse = make([]response.RoleListResponseData, 0)
		for _, role := range userEntity.Roles {
			roleListResponse = append(roleListResponse, response.RoleListResponseData{
				ID:       role.ID,
				RoleName: role.RoleName,
			})
		}
	}

	var rolePoliciesListResponse []response.RolePolicyListResponseData
	if len(userEntity.RolePolicies) > 0 {
		rolePoliciesListResponse = make([]response.RolePolicyListResponseData, 0)
		for _, rolePolicy := range userEntity.RolePolicies {
			rolePoliciesListResponse = append(rolePoliciesListResponse, response.RolePolicyListResponseData{
				ID:         rolePolicy.ID,
				PolicyName: rolePolicy.PolicyName,
			})
		}
	}

	var guardianListResponse []response.UserEntityResponseData
	if len(userEntity.Guardians) > 0 {
		guardianListResponse = make([]response.UserEntityResponseData, 0)
		for _, guardian := range userEntity.Guardians {
			guardianListResponse = append(guardianListResponse, response.UserEntityResponseData{
				ID:       guardian.ID.String(),
				Username: guardian.Username,
			})
		}
	}

	var deviceListResponse []string
	if len(userEntity.Devices) > 0 {
		deviceListResponse = make([]string, 0)
		for _, device := range userEntity.Devices {
			deviceListResponse = append(deviceListResponse, device.ID)
		}
	}

	var userConfig *response.UserConfigResponse
	if userEntity.UserConfig != nil {
		userConfig = &response.UserConfigResponse{
			TopButtonConfig:      userEntity.UserConfig.TopButtonConfig,
			StudentOutputSheetId: userEntity.UserConfig.StudentOutputSheetId,
			TeacherOutputSheetId: userEntity.UserConfig.TeacherOutputSheetId,
		}
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.UserEntityResponse{
			ID:           userEntity.ID.String(),
			Username:     userEntity.Username,
			Fullname:     userEntity.Fullname,
			Phone:        userEntity.Phone,
			Email:        userEntity.Email,
			Dob:          userEntity.Birthday.Format("2006-01-02"),
			Company:      userEntity.Company.CompanyName,
			CreatedAt:    userEntity.CreatedAt.Format("2006-01-02"),
			UserConfig:   userConfig,
			Roles:        &roleListResponse,
			RolePolicies: &rolePoliciesListResponse,
			Guardians:    &guardianListResponse,
			Devices:      &deviceListResponse,
		},
	})
}

func (receiver *UserEntityController) GetUserEntityByName(context *gin.Context) {
	username := context.Param("username")
	if username == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "user name is required",
				},
			},
		)
		return
	}

	userEntity, err := receiver.GetUserByUsername(request.GetUserEntityByUsernameRequest{Username: username})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	var roleListResponse []response.RoleListResponseData
	if len(userEntity.Roles) > 0 {
		roleListResponse = make([]response.RoleListResponseData, 0)
		for _, role := range userEntity.Roles {
			roleListResponse = append(roleListResponse, response.RoleListResponseData{
				ID:       role.ID,
				RoleName: role.RoleName,
			})
		}
	}

	var rolePoliciesListResponse []response.RolePolicyListResponseData
	if len(userEntity.RolePolicies) > 0 {
		rolePoliciesListResponse = make([]response.RolePolicyListResponseData, 0)
		for _, rolePolicy := range userEntity.RolePolicies {
			rolePoliciesListResponse = append(rolePoliciesListResponse, response.RolePolicyListResponseData{
				ID:         rolePolicy.ID,
				PolicyName: rolePolicy.PolicyName,
			})
		}
	}

	var guardianListResponse []response.UserEntityResponseData
	if len(userEntity.Guardians) > 0 {
		guardianListResponse = make([]response.UserEntityResponseData, 0)
		for _, guardian := range userEntity.Guardians {
			guardianListResponse = append(guardianListResponse, response.UserEntityResponseData{
				ID:       guardian.ID.String(),
				Username: guardian.Username,
			})
		}
	}

	var deviceListResponse []string
	if len(userEntity.Devices) > 0 {
		deviceListResponse = make([]string, 0)
		for _, device := range userEntity.Devices {
			deviceListResponse = append(deviceListResponse, device.ID)
		}
	}

	var userConfig *response.UserConfigResponse
	if userEntity.UserConfig != nil {
		userConfig = &response.UserConfigResponse{
			TopButtonConfig:      userEntity.UserConfig.TopButtonConfig,
			StudentOutputSheetId: userEntity.UserConfig.StudentOutputSheetId,
			TeacherOutputSheetId: userEntity.UserConfig.TeacherOutputSheetId,
		}
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.UserEntityResponse{
			ID:           userEntity.ID.String(),
			Username:     userEntity.Username,
			Fullname:     userEntity.Fullname,
			Phone:        userEntity.Phone,
			Email:        userEntity.Email,
			Company:      userEntity.Company.CompanyName,
			UserConfig:   userConfig,
			Roles:        &roleListResponse,
			RolePolicies: &rolePoliciesListResponse,
			Guardians:    &guardianListResponse,
			Devices:      &deviceListResponse,
		},
	})
}

func (receiver *UserEntityController) CreateUserEntity(context *gin.Context) {
	var req request.CreateUserEntityRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// Validate username
	if err := req.IsUsernameValid(); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// Validate the user's age
	if err := req.IsOver18(); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := receiver.CreateUserEntityUseCase.CreateUserEntity(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	data, err := receiver.AuthorizeUseCase.UserLoginUsecase(request.UserLoginFromDeviceReqest{
		Username:   req.Username,
		Password:   req.Password,
		DeviceUUID: req.DeviceUUID,
	})

	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			},
		)
		return
	}

	context.JSON(http.StatusOK, response.LoginResponse{
		Data: data,
	})
}

func (receiver *UserEntityController) UpdateUserEntity(context *gin.Context) {
	var req request.UpdateUserEntityRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	// // Validate the user's email
	// if ok := req.ValidateEmail(); !ok {
	// 	context.JSON(http.StatusBadRequest, response.FailedResponse{
	// 		Error: response.Cause{
	// 			Code:    http.StatusBadRequest,
	// 			Message: "invalid email",
	// 		},
	// 	})
	// 	return
	// }

	// Validate the user's phone number
	if ok := req.ValidatePhone(); !ok {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid phone number",
			},
		})
		return
	}

	err := receiver.UpdateUserEntityUseCase.UpdateUserEntity(req)
	if err != nil {
		log.Error(err)

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
			Message: "user was update successfully",
		},
	})
}
