package controller

import (
	"bufio"
	"github.com/samber/lo"
	"net/http"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserEntityController struct {
	*usecase.GetUserEntityUseCase
	*usecase.CreateUserEntityUseCase
	*usecase.CreateChildForParentUseCase
	*usecase.UpdateUserEntityUseCase
	*usecase.UpdateUserRoleUseCase
	*usecase.AuthorizeUseCase
	*usecase.UpdateUserOrgInfoUseCase
	*usecase.UpdateUserAuthorizeUseCase
	*usecase.DeleteUserAuthorizeUseCase
	*usecase.GetPreRegisterUseCase
	*usecase.CreatePreRegisterUseCase
	*usecase.GetUserFromTokenUseCase

	*usecase.GetUserFormApplicationUseCase
	*usecase.CreateUserFormApplicationUseCase
	*usecase.ApproveUserFormApplicationUseCase
	*usecase.BlockUserFormApplicationUseCase
	*usecase.UploadUserAvatarUseCase
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

	roleListResponse := make([]response.RoleListResponseData, 0)
	if len(userEntity.Roles) > 0 {
		roleListResponse = make([]response.RoleListResponseData, 0)
		for _, role := range userEntity.Roles {
			roleListResponse = append(roleListResponse, response.RoleListResponseData{
				ID:       role.ID,
				RoleName: role.Role.String(),
			})
		}
	}

	deviceListResponse := make([]string, 0)
	if len(userEntity.Devices) > 0 {
		deviceListResponse = make([]string, 0)
		for _, device := range userEntity.Devices {
			deviceListResponse = append(deviceListResponse, device.ID)
		}
	}

	organizations := make([]string, 0)
	if len(userEntity.Organizations) > 0 {
		organizations = lo.Map(userEntity.Organizations, func(item entity.SOrganization, index int) string {
			return item.ID.String()
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.UserEntityResponseV2{
			ID:           userEntity.ID.String(),
			Username:     userEntity.Username,
			Fullname:     userEntity.Fullname,
			Nickname:     userEntity.Nickname,
			Phone:        userEntity.Phone,
			Email:        userEntity.Email,
			Dob:          userEntity.Birthday.Format("2006-01-02"),
			QRLogin:      userEntity.QRLogin,
			Avatar:       userEntity.Avatar,
			AvatarURL:    userEntity.AvatarURL,
			IsBlocked:    userEntity.IsBlocked,
			BlockedAt:    userEntity.BlockedAt.Format("2006-01-02"),
			Organization: organizations,
			CreatedAt:    userEntity.CreatedAt.Format("2006-01-02"),
			Roles:        &roleListResponse,
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

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	isSuperAdmin := lo.ContainsBy(user.Roles, func(role entity.SRole) bool {
		return role.Role == entity.SuperAdmin
	})
	var users []entity.SUserEntity
	if isSuperAdmin {
		organizationID := context.Request.URL.Query().Get("organization_id")

		if organizationID != "" {
			users, err = receiver.GetAllByOrganization(organizationID)
		} else {
			users, err = receiver.GetAllUsers()
		}
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			})

			return
		}
	} else {
		organizationID := context.Request.URL.Query().Get("organization_id")
		if organizationID == "" {
			context.JSON(
				http.StatusBadRequest, response.FailedResponse{
					Code:  http.StatusBadRequest,
					Error: "organization is required",
				},
			)
			return
		}
		isOrganization := lo.ContainsBy(user.Organizations, func(org entity.SOrganization) bool {
			return org.ID.String() == organizationID
		})
		if !isOrganization {
			context.JSON(
				http.StatusUnauthorized, response.FailedResponse{
					Code:  http.StatusUnauthorized,
					Error: "access denied",
				},
			)
			return
		}

		users, err = receiver.GetAllByOrganization(organizationID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			})

			return
		}
	}

	userResponse := make([]response.UserEntityResponseData, 0)
	for _, user := range users {
		roles := make([]string, 0)
		for _, r := range user.Roles {
			if strings.ToLower(role) != "all" {
				if !strings.EqualFold(r.Role.String(), role) {
					continue
				}

				roles = append(roles, r.Role.String())
				break
			}
			roles = append(roles, r.Role.String())
		}

		userResponse = append(userResponse, response.UserEntityResponseData{
			ID:        user.ID.String(),
			Username:  user.Username,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			AvatarURL: user.AvatarURL,
			Roles:     roles,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: userResponse,
	})
}

func (receiver *UserEntityController) BlockUser(context *gin.Context) {
	userID := context.Param("id")
	if userID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	err := receiver.UpdateUserEntityUseCase.BlockUser(userID)
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

func (receiver *UserEntityController) GetUserEntityByID(context *gin.Context) {
	userID := context.Param("id")
	if userID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	userEntity, err := receiver.GetUserByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: "Your phone has been unable to locate your account. Try log in again or contact our customer support via our webpage.",
		})

		return
	}

	roleListResponse := make([]response.RoleListResponseData, 0)
	if len(userEntity.Roles) > 0 {
		roleListResponse = make([]response.RoleListResponseData, 0)
		for _, role := range userEntity.Roles {
			roleListResponse = append(roleListResponse, response.RoleListResponseData{
				ID:       role.ID,
				RoleName: role.Role.String(),
			})
		}
	}

	deviceListResponse := make([]string, 0)
	if len(userEntity.Devices) > 0 {
		deviceListResponse = make([]string, 0)
		for _, device := range userEntity.Devices {
			deviceListResponse = append(deviceListResponse, device.ID)
		}
	}

	organizations := make([]string, 0)
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
			QRLogin:      userEntity.QRLogin,
			Avatar:       userEntity.Avatar,
			AvatarURL:    userEntity.AvatarURL,
			IsBlocked:    userEntity.IsBlocked,
			BlockedAt:    userEntity.BlockedAt.Format("2006-01-02"),
			Organization: organizations,
			CreatedAt:    userEntity.CreatedAt.Format("2006-01-02"),
			Roles:        &roleListResponse,
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

	roleListResponse := make([]response.RoleListResponseData, 0)
	if len(userEntity.Roles) > 0 {
		roleListResponse = make([]response.RoleListResponseData, 0)
		for _, role := range userEntity.Roles {
			roleListResponse = append(roleListResponse, response.RoleListResponseData{
				ID:       role.ID,
				RoleName: role.Role.String(),
			})
		}
	}

	deviceListResponse := make([]string, 0)
	if len(userEntity.Devices) > 0 {
		deviceListResponse = make([]string, 0)
		for _, device := range userEntity.Devices {
			deviceListResponse = append(deviceListResponse, device.ID)
		}
	}

	organizations := make([]string, 0)
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
			QRLogin:      userEntity.QRLogin,
			Avatar:       userEntity.Avatar,
			AvatarURL:    userEntity.AvatarURL,
			IsBlocked:    userEntity.IsBlocked,
			BlockedAt:    userEntity.BlockedAt.Format("2006-01-02"),
			Organization: organizations,
			CreatedAt:    userEntity.CreatedAt.Format("2006-01-02"),
			Roles:        &roleListResponse,
			Devices:      &deviceListResponse,
		},
	})
}

func (receiver *UserEntityController) GetUserOrgInfo(context *gin.Context) {
	userID := context.Param("user_id")
	organizationID := context.Param("organization_id")
	if userID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "organization id is required",
			},
		)
		return
	}

	user, err := receiver.GetUserEntityUseCase.GetUserOrgInfo(userID, organizationID)
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
	organizationID := context.Param("organization_id")
	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "organization id is required",
			},
		)
		return
	}

	users, err := receiver.GetUserEntityUseCase.GetAllOrgManagerInfo(organizationID)
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
			UserID:       user.UserID.String(),
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
	userID := context.Param("id")
	if userID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: "user id is required",
			},
		)
		return
	}

	rights, err := receiver.GetUserEntityUseCase.GetAllUserAuthorize(userID)
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
		FunctionClaimID int64  `json:"function_claim_id"`
		FunctionName    string `json:"function_name"`
		PermissionID    int64  `json:"permission_id"`
		PermissionName  string `json:"permission_name"`
	}
	type getAllUserAuthorizeResponse struct {
		UserID            string                      `json:"user_id"`
		Username          string                      `json:"username"`
		FunctionAuthorize []functionAuthorizeResponse `json:"function_authorize"`
	}

	var functionAuthorize []functionAuthorizeResponse
	for _, right := range rights {
		functionAuthorize = append(functionAuthorize, functionAuthorizeResponse{
			FunctionClaimID: right.FunctionClaimID,
			FunctionName:    right.FunctionClaim.FunctionName,
			PermissionID:    right.FunctionClaimPermissionID,
			PermissionName:  right.FunctionClaimPermission.PermissionName,
		})
	}

	res := &getAllUserAuthorizeResponse{
		UserID:            userID,
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

	req.Username = strings.ToLower(req.Username)

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

func (receiver *UserEntityController) CreateChildForParent(context *gin.Context) {
	var req request.CreateChildForParentRequest
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

	parent, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	req.Username = strings.ToLower(req.Username)

	err = receiver.CreateChildForParentUseCase.CreateChildForParent(parent.ID.String(), req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "child was created successfully",
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

// Teacher

func (receiver *UserEntityController) GetAllTeacherFormApplication(context *gin.Context) {
	applications, err := receiver.GetUserFormApplicationUseCase.GetAllTeacherFormApplication()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	var applicationResponse []response.TeacherFormApplicationResponse
	if len(applications) > 0 {
		applicationResponse = make([]response.TeacherFormApplicationResponse, 0)
		for _, application := range applications {
			res := response.TeacherFormApplicationResponse{
				ID:         application.ID,
				Status:     application.Status.String(),
				ApprovedAt: "",
				CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
				UserID:     application.UserID.String(),
			}
			if application.ApprovedAt != defaultTime {
				res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
			}
			applicationResponse = append(applicationResponse, res)
		}
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: applicationResponse,
	})
}

func (receiver *UserEntityController) GetTeacherFormApplicationByID(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	application, err := receiver.GetUserFormApplicationUseCase.GetTeacherFormApplicationByID(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	res := response.TeacherFormApplicationResponse{
		ID:         application.ID,
		Status:     application.Status.String(),
		ApprovedAt: "",
		CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
		UserID:     application.UserID.String(),
	}
	if application.ApprovedAt != defaultTime {
		res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *UserEntityController) ApproveTeacherFormApplication(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	err = receiver.ApproveUserFormApplicationUseCase.ApproveTeacherFormApplication(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application approved successfully",
	})
}

func (receiver *UserEntityController) BlockTeacherFormApplication(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	err = receiver.BlockUserFormApplicationUseCase.BlockTeacherFormApplication(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application blocked successfully",
	})
}

func (receiver *UserEntityController) CreateTeacherFormApplication(context *gin.Context) {
	var req request.CreateTeacherFormApplicationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	req.UserID = user.ID.String()

	err = receiver.CreateUserFormApplicationUseCase.CreateTeacherFormApplication(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application created successfully",
	})
}

// Staff

func (receiver *UserEntityController) GetAllStaffFormApplication(context *gin.Context) {
	applications, err := receiver.GetUserFormApplicationUseCase.GetAllStaffFormApplication()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	var applicationResponse []response.StaffFormApplicationResponse
	if len(applications) > 0 {
		applicationResponse = make([]response.StaffFormApplicationResponse, 0)
		for _, application := range applications {
			res := response.StaffFormApplicationResponse{
				ID:         application.ID,
				Status:     application.Status.String(),
				ApprovedAt: "",
				CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
				UserID:     application.UserID.String(),
			}
			if application.ApprovedAt != defaultTime {
				res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
			}
			applicationResponse = append(applicationResponse, res)
		}
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: applicationResponse,
	})
}

func (receiver *UserEntityController) GetStaffFormApplicationByID(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	application, err := receiver.GetUserFormApplicationUseCase.GetStaffFormApplicationByID(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	res := response.StaffFormApplicationResponse{
		ID:         application.ID,
		Status:     application.Status.String(),
		ApprovedAt: "",
		CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
		UserID:     application.UserID.String(),
	}
	if application.ApprovedAt != defaultTime {
		res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *UserEntityController) ApproveStaffFormApplication(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	err = receiver.ApproveUserFormApplicationUseCase.ApproveStaffFormApplication(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application approved successfully",
	})
}

func (receiver *UserEntityController) BlockStaffFormApplication(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	err = receiver.BlockUserFormApplicationUseCase.BlockStaffFormApplication(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application blocked successfully",
	})
}

func (receiver *UserEntityController) CreateStaffFormApplication(context *gin.Context) {
	var req request.CreateStaffFormApplicationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	req.UserID = user.ID.String()

	err = receiver.CreateUserFormApplicationUseCase.CreateStaffFormApplication(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application created successfully",
	})
}

// Student

func (receiver *UserEntityController) GetAllStudentFormApplication(context *gin.Context) {
	applications, err := receiver.GetUserFormApplicationUseCase.GetAllStudentFormApplication()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	var applicationResponse []response.StudentFormApplicationResponse
	if len(applications) > 0 {
		applicationResponse = make([]response.StudentFormApplicationResponse, 0)
		for _, application := range applications {
			res := response.StudentFormApplicationResponse{
				ID:         application.ID,
				Status:     application.Status.String(),
				ApprovedAt: "",
				CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
				UserID:     application.UserID.String(),
			}
			if application.ApprovedAt != defaultTime {
				res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
			}
			applicationResponse = append(applicationResponse, res)
		}
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: applicationResponse,
	})
}

func (receiver *UserEntityController) GetStudentFormApplicationByID(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	application, err := receiver.GetUserFormApplicationUseCase.GetStudentFormApplicationByID(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	res := response.StudentFormApplicationResponse{
		ID:         application.ID,
		Status:     application.Status.String(),
		ApprovedAt: "",
		CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
		UserID:     application.UserID.String(),
	}
	if application.ApprovedAt != defaultTime {
		res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *UserEntityController) ApproveStudentFormApplication(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	err = receiver.ApproveUserFormApplicationUseCase.ApproveStudentFormApplication(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application approved successfully",
	})
}

func (receiver *UserEntityController) BlockStudentFormApplication(context *gin.Context) {
	applicationID := context.Param("id")
	if applicationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	id, err := strconv.Atoi(applicationID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: "invalid id",
			Code:  http.StatusBadRequest,
		})
		return
	}

	err = receiver.BlockUserFormApplicationUseCase.BlockStudentFormApplication(int64(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application blocked successfully",
	})
}

func (receiver *UserEntityController) CreateStudentFormApplication(context *gin.Context) {
	var req request.CreateStudentFormApplicationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	req.UserID = user.ID.String()

	err = receiver.CreateUserFormApplicationUseCase.CreateStudentFormApplication(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Application created successfully",
	})
}

func (receiver *UserEntityController) UploadAvatar(context *gin.Context) {
	fileHeader, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	userID := context.PostForm("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "user id is required",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	defer file.Close()

	dataBytes := make([]byte, fileHeader.Size)
	if _, err := bufio.NewReader(file).Read(dataBytes); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	url, img, err := receiver.UploadUserAvatarUseCase.UploadAvatar(userID, dataBytes, fileHeader.Filename)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if url == nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "avatar was not created",
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "avatar was create successfully",
		Data: response.ImageResponse{
			ImageName: img.ImageName,
			Key:       img.Key,
			Extension: img.Extension,
			Url:       *url,
			Width:     img.Width,
			Height:    img.Height,
		},
	})
}
