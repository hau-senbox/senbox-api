package controller

import (
	"net/http"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"
	"strings"

	"github.com/samber/lo"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserEntityController struct {
	*usecase.GetUserEntityUseCase
	*usecase.CreateUserEntityUseCase
	*usecase.UpdateUserEntityUseCase
	*usecase.UpdateUserRoleUseCase
	*usecase.AuthorizeUseCase
	*usecase.UpdateUserOrgInfoUseCase
	*usecase.UpdateUserAuthorizeUseCase
	*usecase.DeleteUserAuthorizeUseCase
	*usecase.GetPreRegisterUseCase
	*usecase.CreatePreRegisterUseCase
	*usecase.GetUserFromTokenUseCase
}

func (receiver *UserEntityController) GetCurrentUser(context *gin.Context) {
	userEntity, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
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

	var organizations []string
	if len(userEntity.Organizations) > 0 {
		organizations = lo.Map(userEntity.Organizations, func(item entity.SOrganization, index int) string {
			return item.OrganizationName
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.UserEntityResponse{
			ID:           userEntity.ID.String(),
			Username:     userEntity.Username,
			Fullname:     userEntity.Fullname,
			Nickname:     userEntity.Nickname,
			Phone:        userEntity.Phone,
			Email:        userEntity.Email,
			Dob:          userEntity.Birthday.Format("2006-01-02"),
			IsBlocked:    userEntity.IsBlocked,
			BlockedAt:    userEntity.BlockedAt.Format("2006-01-02"),
			Organization: organizations,
			CreatedAt:    userEntity.CreatedAt.Format("2006-01-02"),
			Roles:        &roleListResponse,
			Guardians:    &guardianListResponse,
			Devices:      &deviceListResponse,
		},
	})
}

func (receiver *UserEntityController) GetAllUserEntity(context *gin.Context) {
	role := context.Request.URL.Query().Get("role")
	if role == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "role is required",
			},
		)
		return
	}

	users, err := receiver.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var userResponse []response.UserEntityResponseData
	for _, user := range users {
		var roles []string
		for _, r := range user.Roles {
			if strings.ToLower(role) != "all" {
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
				Nickname: user.Nickname,
				Roles:    roles,
			})
		}
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: userResponse,
	})
}

func (receiver *UserEntityController) GetChildrenOfGuardian(context *gin.Context) {
	userId := context.Param("id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	users, err := receiver.GetUserEntityUseCase.GetChildrenOfGuardian(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: *users,
	})
}
func (receiver *UserEntityController) BlockUser(context *gin.Context) {
	println("ENTERRRR")
	userId := context.Param("id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	err := receiver.UpdateUserEntityUseCase.BlockUser(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user was blocked successfully",
	})
}

func (receiver *UserEntityController) GetUserEntityById(context *gin.Context) {
	userId := context.Param("id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	userEntity, err := receiver.GetUserById(request.GetUserEntityByIdRequest{ID: userId})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
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

	var organizations []string
	if len(userEntity.Organizations) > 0 {
		organizations = lo.Map(userEntity.Organizations, func(item entity.SOrganization, index int) string {
			return item.OrganizationName
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.UserEntityResponse{
			ID:           userEntity.ID.String(),
			Username:     userEntity.Username,
			Fullname:     userEntity.Fullname,
			Nickname:     userEntity.Nickname,
			Phone:        userEntity.Phone,
			Email:        userEntity.Email,
			Dob:          userEntity.Birthday.Format("2006-01-02"),
			IsBlocked:    userEntity.IsBlocked,
			BlockedAt:    userEntity.BlockedAt.Format("2006-01-02"),
			Organization: organizations,
			CreatedAt:    userEntity.CreatedAt.Format("2006-01-02"),
			Roles:        &roleListResponse,
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
				Code:  http.StatusBadRequest,
				Error: "user name is required",
			},
		)
		return
	}

	userEntity, err := receiver.GetUserByUsername(request.GetUserEntityByUsernameRequest{Username: username})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
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

	var organizations []string
	if len(userEntity.Organizations) > 0 {
		organizations = lo.Map(userEntity.Organizations, func(item entity.SOrganization, index int) string {
			return item.OrganizationName
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.UserEntityResponse{
			ID:           userEntity.ID.String(),
			Username:     userEntity.Username,
			Fullname:     userEntity.Fullname,
			Nickname:     userEntity.Nickname,
			Phone:        userEntity.Phone,
			Email:        userEntity.Email,
			Dob:          userEntity.Birthday.Format("2006-01-02"),
			IsBlocked:    userEntity.IsBlocked,
			BlockedAt:    userEntity.BlockedAt.Format("2006-01-02"),
			Organization: organizations,
			CreatedAt:    userEntity.CreatedAt.Format("2006-01-02"),
			Roles:        &roleListResponse,
			Guardians:    &guardianListResponse,
			Devices:      &deviceListResponse,
		},
	})
}

func (receiver *UserEntityController) GetUserOrgInfo(context *gin.Context) {
	userId := context.Param("user_id")
	organizationId := context.Param("organization_id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	if organizationId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "organization id is required",
			},
		)
		return
	}

	orgID, err := strconv.Atoi(organizationId)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid organization id",
		})
		return
	}

	user, err := receiver.GetUserEntityUseCase.GetUserOrgInfo(userId, int64(orgID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: &response.GetUserOrgInfoResponse{
			UserNickName: user.UserNickName,
			IsManager:    user.IsManager,
		},
	})
}

func (receiver *UserEntityController) GetAllOrgManagerInfo(context *gin.Context) {
	organizationId := context.Param("organization_id")
	if organizationId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "organization id is required",
			},
		)
		return
	}

	users, err := receiver.GetUserEntityUseCase.GetAllOrgManagerInfo(organizationId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	var res []response.GetOrgManagerInfoResponse
	for _, user := range *users {
		res = append(res, response.GetOrgManagerInfoResponse{
			UserId:       user.UserId.String(),
			UserNickName: user.UserNickName,
			IsManager:    user.IsManager,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *UserEntityController) GetAllUserAuthorize(context *gin.Context) {
	userId := context.Param("id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	rights, err := receiver.GetUserEntityUseCase.GetAllUserAuthorize(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	if len(rights) == 0 {
		context.JSON(http.StatusOK, response.SucceedResponse{
			Code: http.StatusOK,
			Data: nil,
		})
		return
	}

	type functionAuthorizeResponse struct {
		FunctionClaimId int64  `json:"function_claim_id"`
		FunctionName    string `json:"function_name"`
		PermissionId    int64  `json:"permission_id"`
		PermissionName  string `json:"permission_name"`
	}
	type getAllUserAuthorizeResponse struct {
		UserId            string                      `json:"user_id"`
		Username          string                      `json:"username"`
		FunctionAuthorize []functionAuthorizeResponse `json:"function_authorize"`
	}

	var functionAuthorize []functionAuthorizeResponse
	for _, right := range rights {
		functionAuthorize = append(functionAuthorize, functionAuthorizeResponse{
			FunctionClaimId: right.FunctionClaimId,
			FunctionName:    right.FunctionClaim.FunctionName,
			PermissionId:    right.FunctionClaimPermissionId,
			PermissionName:  right.FunctionClaimPermission.PermissionName,
		})
	}

	res := &getAllUserAuthorizeResponse{
		UserId:            userId,
		Username:          rights[0].User.Username,
		FunctionAuthorize: functionAuthorize,
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *UserEntityController) UpdateUserAuthorize(context *gin.Context) {
	var req request.UpdateUserAuthorizeRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UpdateUserAuthorizeUseCase.UpdateUserAuthorize(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user authorize was updated successfully",
	})
}

func (receiver *UserEntityController) DeleteUserAuthorize(context *gin.Context) {
	var req request.DeleteUserAuthorizeRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.DeleteUserAuthorizeUseCase.DeleteUserAuthorize(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user authorize was deleted successfully",
	})
}

func (receiver *UserEntityController) CreateUserEntity(context *gin.Context) {
	var req request.CreateUserEntityRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// Validate username
	if err := req.IsUsernameValid(); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// Validate the user's age
	if err := req.IsOver18(); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.CreateUserEntityUseCase.CreateUserEntity(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	data, err := receiver.UserLoginUsecase(request.UserLoginFromDeviceReqest{
		Username:   req.Username,
		Password:   req.Password,
		DeviceUUID: req.DeviceUUID,
	})

	if err != nil {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			},
		)
		return
	}

	context.JSON(http.StatusOK, response.LoginResponse{
		Data: *data,
	})
}

func (receiver *UserEntityController) UpdateUserEntity(context *gin.Context) {
	var req request.UpdateUserEntityRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
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
			Code:  http.StatusBadRequest,
			Error: "invalid phone number",
		})
		return
	}

	err := receiver.UpdateUserEntityUseCase.UpdateUserEntity(req)
	if err != nil {
		log.Error(err)

		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user was update successfully",
	})
}

func (receiver *UserEntityController) UpdateUserRole(context *gin.Context) {
	var req request.UpdateUserRoleRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UpdateUserRoleUseCase.UpdateUserRole(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user role was updated successfully",
	})
}

func (receiver *UserEntityController) UpdateUserOrgInfo(context *gin.Context) {
	var req request.UpdateUserOrgInfoRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UpdateUserOrgInfoUseCase.UpdateUserOrgInfo(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user org info was updated successfully",
	})
}

type registerResponse struct {
	Email string `json:"email"`
}

func (receiver *UserEntityController) GetAllPreRegisterUser(context *gin.Context) {
	registers, err := receiver.GetPreRegisterUseCase.GetAllPreRegisterUser()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	var res []registerResponse
	for _, register := range registers {
		res = append(res, registerResponse{
			Email: register.Email,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

type createPreRegisterResponse struct {
	Email string `json:"email"`
}

func (receiver *UserEntityController) CreatePreRegister(context *gin.Context) {
	var req createPreRegisterResponse
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.CreatePreRegisterUseCase.CreatePreRegister(req.Email)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "register was created successfully",
	})
}
