package controller

import (
	"bufio"
	"net/http"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"

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
	*usecase.RoleOrgSignUpUseCase
	*usecase.ChildUseCase
	*usecase.StudentApplicationUseCase
	*usecase.TeacherApplicationUseCase
	*usecase.StaffApplicationUseCase
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

// func (receiver *UserEntityController) GetAllTeacherFormApplication(context *gin.Context) {
// 	applications, err := receiver.GetUserFormApplicationUseCase.GetAllTeacherFormApplication()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 			Error: err.Error(),
// 			Code:  http.StatusInternalServerError,
// 		})

// 		return
// 	}

// 	var applicationResponse []response.TeacherFormApplicationResponse
// 	if len(applications) > 0 {
// 		applicationResponse = make([]response.TeacherFormApplicationResponse, 0)
// 		for _, application := range applications {
// 			res := response.TeacherFormApplicationResponse{
// 				ID:         application.ID,
// 				Status:     application.Status.String(),
// 				ApprovedAt: "",
// 				CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
// 				UserID:     application.UserID.String(),
// 			}
// 			if application.ApprovedAt != defaultTime {
// 				res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
// 			}
// 			applicationResponse = append(applicationResponse, res)
// 		}
// 	}

// 	context.JSON(http.StatusOK, response.SucceedResponse{
// 		Code: http.StatusOK,
// 		Data: applicationResponse,
// 	})
// }

// func (receiver *UserEntityController) GetTeacherFormApplicationByID(context *gin.Context) {
// 	applicationID := context.Param("id")
// 	if applicationID == "" {
// 		context.JSON(
// 			http.StatusBadRequest, response.FailedResponse{
// 				Error: "id is required",
// 				Code:  http.StatusBadRequest,
// 			},
// 		)
// 		return
// 	}

// 	id, err := strconv.Atoi(applicationID)
// 	if err != nil {
// 		context.JSON(http.StatusBadRequest, response.FailedResponse{
// 			Error: "invalid id",
// 			Code:  http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	application, err := receiver.GetUserFormApplicationUseCase.GetTeacherFormApplicationByID(int64(id))
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 			Error: err.Error(),
// 			Code:  http.StatusInternalServerError,
// 		})

// 		return
// 	}

// 	res := response.TeacherFormApplicationResponse{
// 		ID:         application.ID,
// 		Status:     application.Status.String(),
// 		ApprovedAt: "",
// 		CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UserID:     application.UserID.String(),
// 	}
// 	if application.ApprovedAt != defaultTime {
// 		res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
// 	}

// 	context.JSON(http.StatusOK, response.SucceedResponse{
// 		Code: http.StatusOK,
// 		Data: res,
// 	})
// }

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

// func (receiver *UserEntityController) GetAllStaffFormApplication(context *gin.Context) {
// 	applications, err := receiver.GetUserFormApplicationUseCase.GetAllStaffFormApplication()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 			Error: err.Error(),
// 			Code:  http.StatusInternalServerError,
// 		})

// 		return
// 	}

// 	var applicationResponse []response.StaffFormApplicationResponse
// 	if len(applications) > 0 {
// 		applicationResponse = make([]response.StaffFormApplicationResponse, 0)
// 		for _, application := range applications {
// 			res := response.StaffFormApplicationResponse{
// 				ID:         application.ID,
// 				Status:     application.Status.String(),
// 				ApprovedAt: "",
// 				CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
// 				UserID:     application.UserID.String(),
// 			}
// 			if application.ApprovedAt != defaultTime {
// 				res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
// 			}
// 			applicationResponse = append(applicationResponse, res)
// 		}
// 	}

// 	context.JSON(http.StatusOK, response.SucceedResponse{
// 		Code: http.StatusOK,
// 		Data: applicationResponse,
// 	})
// }

// func (receiver *UserEntityController) GetStaffFormApplicationByID(context *gin.Context) {
// 	applicationID := context.Param("id")
// 	if applicationID == "" {
// 		context.JSON(
// 			http.StatusBadRequest, response.FailedResponse{
// 				Error: "id is required",
// 				Code:  http.StatusBadRequest,
// 			},
// 		)
// 		return
// 	}

// 	id, err := strconv.Atoi(applicationID)
// 	if err != nil {
// 		context.JSON(http.StatusBadRequest, response.FailedResponse{
// 			Error: "invalid id",
// 			Code:  http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	application, err := receiver.GetUserFormApplicationUseCase.GetStaffFormApplicationByID(int64(id))
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 			Error: err.Error(),
// 			Code:  http.StatusInternalServerError,
// 		})

// 		return
// 	}

// 	res := response.StaffFormApplicationResponse{
// 		ID:         application.ID,
// 		Status:     application.Status.String(),
// 		ApprovedAt: "",
// 		CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UserID:     application.UserID.String(),
// 	}
// 	if application.ApprovedAt != defaultTime {
// 		res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
// 	}

// 	context.JSON(http.StatusOK, response.SucceedResponse{
// 		Code: http.StatusOK,
// 		Data: res,
// 	})
// }

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

// func (receiver *UserEntityController) GetAllStudentFormApplication(context *gin.Context) {
// 	applications, err := receiver.GetUserFormApplicationUseCase.GetAllStudentFormApplication()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 			Error: err.Error(),
// 			Code:  http.StatusInternalServerError,
// 		})

// 		return
// 	}

// 	var applicationResponse []response.StudentFormApplicationResponse
// 	if len(applications) > 0 {
// 		applicationResponse = make([]response.StudentFormApplicationResponse, 0)
// 		for _, application := range applications {
// 			res := response.StudentFormApplicationResponse{
// 				ID:         application.ID,
// 				Status:     application.Status.String(),
// 				ApprovedAt: "",
// 				CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
// 				UserID:     application.UserID.String(),
// 			}
// 			if application.ApprovedAt != defaultTime {
// 				res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
// 			}
// 			applicationResponse = append(applicationResponse, res)
// 		}
// 	}

// 	context.JSON(http.StatusOK, response.SucceedResponse{
// 		Code: http.StatusOK,
// 		Data: applicationResponse,
// 	})
// }

// func (receiver *UserEntityController) GetStudentFormApplicationByID(context *gin.Context) {
// 	applicationID := context.Param("id")
// 	if applicationID == "" {
// 		context.JSON(
// 			http.StatusBadRequest, response.FailedResponse{
// 				Error: "id is required",
// 				Code:  http.StatusBadRequest,
// 			},
// 		)
// 		return
// 	}

// 	id, err := strconv.Atoi(applicationID)
// 	if err != nil {
// 		context.JSON(http.StatusBadRequest, response.FailedResponse{
// 			Error: "invalid id",
// 			Code:  http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	application, err := receiver.GetUserFormApplicationUseCase.GetStudentFormApplicationByID(int64(id))
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 			Error: err.Error(),
// 			Code:  http.StatusInternalServerError,
// 		})

// 		return
// 	}

// 	res := response.StudentFormApplicationResponse{
// 		ID:         application.ID,
// 		Status:     application.Status.String(),
// 		ApprovedAt: "",
// 		CreatedAt:  application.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UserID:     application.UserID.String(),
// 	}
// 	if application.ApprovedAt != defaultTime {
// 		res.ApprovedAt = application.ApprovedAt.Format("2006-01-02 15:04:05")
// 	}

// 	context.JSON(http.StatusOK, response.SucceedResponse{
// 		Code: http.StatusOK,
// 		Data: res,
// 	})
// }

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

func (receiver *UserEntityController) GetAllRoleOrgSignUp(context *gin.Context) {
	roles, err := receiver.RoleOrgSignUpUseCase.Get4App()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: roles,
	})
}

// child
func (receiver *UserEntityController) CreateChild(context *gin.Context) {
	var req request.CreateChildRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid input: " + err.Error(),
		})
		return
	}

	// Gọi UseCase để tạo child
	err := receiver.ChildUseCase.CreateChild(req, context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Child created successfully",
	})
}

func (receiver *UserEntityController) GetChildByID(context *gin.Context) {
	// Lấy childID từ URL param (giả sử endpoint là /children/:id)
	childID := context.Param("id")
	if childID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing child ID",
		})
		return
	}

	child, err := receiver.ChildUseCase.GetByID(childID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get child",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: child,
	})
}

func (ctl *UserEntityController) UpdateChild(c *gin.Context) {
	var req request.UpdateChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := ctl.ChildUseCase.UpdateChild(req, c); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update child",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: "Child updated successfully",
	})
}

func (receiver *UserEntityController) SearchUser4WebAdmin(c *gin.Context) {
	role := c.Request.URL.Query().Get("role")
	name := strings.ToLower(c.Request.URL.Query().Get("name"))

	callAll := role == "all"
	var (
		users    = make([]response.UserResponse, 0)
		children = make([]response.ChildrenResponse, 0)
		students = make([]response.StudentResponse, 0)
		teachers = make([]response.TeacherResponse, 0)
		staffs   = make([]response.StaffResponse, 0)
	)

	if callAll {
		users, _ := receiver.GetAllUsers4Search(c)
		children, _ := receiver.ChildUseCase.GetAll4Search(c)
		students, _ := receiver.StudentApplicationUseCase.GetAllStudents4Search(c)
		teachers, _ := receiver.TeacherApplicationUseCase.GetAllTeachers4Search(c)
		staffs, _ := receiver.StaffApplicationUseCase.GetAllStaff4Search(c)

		userResponse := make([]response.UserResponse, 0, len(users))
		for _, u := range users {
			userResponse = append(userResponse, response.UserResponse{
				ID:        u.ID.String(),
				Username:  u.Username,
				Nickname:  u.Nickname,
				Avatar:    u.Avatar,
				AvatarURL: u.AvatarURL,
			})
		}

		childrenResponse := make([]response.ChildrenResponse, 0, len(children))
		for _, c := range children {
			childrenResponse = append(childrenResponse, response.ChildrenResponse{
				ChildID:   c.ID.String(),
				ChildName: c.ChildName,
			})
		}

		studentResponse := make([]response.StudentResponse, 0, len(students))
		for _, s := range students {
			studentResponse = append(studentResponse, response.StudentResponse{
				StudentID:   s.StudentID,
				StudentName: s.StudentName,
			})
		}

		teachersResponse := make([]response.TeacherResponse, 0, len(teachers))
		for _, t := range teachers {
			teachersResponse = append(teachersResponse, response.TeacherResponse{
				TeacherID:   t.TeacherID,
				TeacherName: t.TeacherName,
			})
		}

		staffsResponse := make([]response.StaffResponse, 0, len(staffs))
		for _, s := range staffs {
			staffsResponse = append(staffsResponse, response.StaffResponse{
				StaffID:   s.StaffID,
				StaffName: s.StaffName,
			})
		}

		if name != "" {
			filteredChildren := make([]response.ChildrenResponse, 0)
			for _, child := range childrenResponse {
				if strings.Contains(strings.ToLower(child.ChildName), name) {
					filteredChildren = append(filteredChildren, child)
				}
			}
			childrenResponse = filteredChildren
		}

		if name != "" {
			filteredStudents := make([]response.StudentResponse, 0)
			for _, s := range students {
				if strings.Contains(strings.ToLower(s.StudentName), name) {
					filteredStudents = append(filteredStudents, s)
				}
			}
			studentResponse = filteredStudents
		}

		if name != "" {
			filteredTeachers := make([]response.TeacherResponse, 0)
			for _, t := range teachers {
				if strings.Contains(strings.ToLower(t.TeacherName), name) {
					filteredTeachers = append(filteredTeachers, t)
				}
			}
			teachersResponse = filteredTeachers
		}

		if name != "" {
			filteredUsers := make([]response.UserResponse, 0)
			for _, u := range userResponse {
				if strings.Contains(strings.ToLower(u.Nickname), name) {
					filteredUsers = append(filteredUsers, u)
				}
			}
			userResponse = filteredUsers
		}

		if name != "" {
			filteredStaffs := make([]response.StaffResponse, 0)
			for _, s := range staffsResponse {
				if strings.Contains(strings.ToLower(s.StaffName), name) {
					filteredStaffs = append(filteredStaffs, s)
				}
			}
			staffsResponse = filteredStaffs
		}

		c.JSON(http.StatusOK, response.SucceedResponse{
			Code: http.StatusOK,
			Data: response.SearchUserResponse{
				Users:     userResponse,
				Children:  childrenResponse,
				Students:  studentResponse,
				Teadchers: teachersResponse,
				Staffs:    staffsResponse,
			},
		})

		return
	}

	if !value.IsValidRoleSignUp(role) {
		c.JSON(http.StatusOK, response.SucceedResponse{
			Code: http.StatusOK,
			Data: response.SearchUserResponse{
				Users:     users,
				Children:  children,
				Students:  students,
				Teadchers: teachers,
				Staffs:    staffs,
			},
		})

		return
	}

	// Nếu role hợp lệ, ép kiểu và xử lý từng trường hợp
	roleSignUp := value.RoleSignUp(role)

	switch roleSignUp {
	case value.RoleOrganization:
		break
	case value.RoleChild:
		// get from childrepo
		childrenDataRepo, _ := receiver.ChildUseCase.GetAll4Search(c)
		for _, c := range childrenDataRepo {
			children = append(children, response.ChildrenResponse{
				ChildID:   c.ID.String(),
				ChildName: c.ChildName,
			})
		}
		if name != "" {
			filteredChildren := make([]response.ChildrenResponse, 0)
			for _, child := range children {
				if strings.Contains(strings.ToLower(child.ChildName), name) {
					filteredChildren = append(filteredChildren, child)
				}
			}
			children = filteredChildren
		}

	case value.RoleStudent:
		//getbfrom student repo
		studentsDataRepo, _ := receiver.StudentApplicationUseCase.GetAllStudents4Search(c)
		for _, s := range studentsDataRepo {
			students = append(students, response.StudentResponse{
				StudentID:   s.StudentID,
				StudentName: s.StudentName,
			})
		}

		if name != "" {
			filteredStudents := make([]response.StudentResponse, 0)
			for _, s := range students {
				if strings.Contains(strings.ToLower(s.StudentName), name) {
					filteredStudents = append(filteredStudents, s)
				}
			}
			students = filteredStudents
		}

	case value.RoleTeacher:
		// get from teacher repo
		teachersDataRepo, _ := receiver.TeacherApplicationUseCase.GetAllTeachers4Search(c)
		for _, t := range teachersDataRepo {
			teachers = append(teachers, response.TeacherResponse{
				TeacherID:   t.TeacherID,
				TeacherName: t.TeacherName,
			})
		}

		if name != "" {
			filteredTeachers := make([]response.TeacherResponse, 0)
			for _, t := range teachers {
				if strings.Contains(strings.ToLower(t.TeacherName), name) {
					filteredTeachers = append(filteredTeachers, t)
				}
			}
			teachers = filteredTeachers
		}

	case value.RoleStaff:
		// get from staff repo
		staffsDataRepo, _ := receiver.StaffApplicationUseCase.GetAllStaff4Search(c)
		for _, s := range staffsDataRepo {
			staffs = append(staffs, response.StaffResponse{
				StaffID:   s.StaffID,
				StaffName: s.StaffName,
			})
		}

		if name != "" {
			filteredStaffs := make([]response.StaffResponse, 0)
			for _, s := range staffs {
				if strings.Contains(strings.ToLower(s.StaffName), name) {
					filteredStaffs = append(filteredStaffs, s)
				}
			}
			staffs = filteredStaffs
		}
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.SearchUserResponse{
			Users:     users,
			Children:  children,
			Students:  students,
			Teadchers: teachers,
			Staffs:    staffs,
		},
	})
}

func (receiver *UserEntityController) GetChild4WebAdmin(context *gin.Context) {
	childID := context.Param("id")
	if childID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing child ID",
		})
		return
	}

	child, err := receiver.ChildUseCase.GetByID4WebAdmin(childID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get child",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: child,
	})

}

func (receiver *UserEntityController) GetStudent4WebAdmin(context *gin.Context) {
	studentID := context.Param("id")
	if studentID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing child ID",
		})
		return
	}

	student, err := receiver.StudentApplicationUseCase.GetStudentByID(studentID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get student",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: student,
	})

}

func (receiver *UserEntityController) GetTeacher4WebAdmin(context *gin.Context) {
	teacherID := context.Param("id")
	if teacherID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing teacher ID",
		})
		return
	}

	teacher, err := receiver.TeacherApplicationUseCase.GetTeacherByID(teacherID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get teacher",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: teacher,
	})

}

func (receiver *UserEntityController) GetStudent4App(context *gin.Context) {

	studentID := context.Param("id")
	if studentID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing child ID",
		})
		return
	}

	if _, err := uuid.Parse(studentID); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid student ID format (must be UUID)",
		})
		return
	}

	student, err := receiver.StudentApplicationUseCase.GetStudentByID4App(context, studentID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: student,
	})

}

func (receiver *UserEntityController) UpdateStudent4App(context *gin.Context) {
	var req request.UpdateStudentRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := req.Validate(); err != nil {
		// validate UUID lỗi
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	err := receiver.StudentApplicationUseCase.UpdateStudentName(req)
	if err != nil {
		context.JSON(http.StatusOK, response.SucceedResponse{
			Code: http.StatusOK,
			Data: err.Error(),
		})
		return
	}
	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: "Updated",
	})

}
