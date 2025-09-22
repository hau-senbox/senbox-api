package controller

import (
	"bufio"
	"net/http"
	"regexp"
	"sen-global-api/helper"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/uploader"
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
	*usecase.UserBlockSettingUsecase
	*usecase.ParentUseCase
	*usecase.StudentBlockSettingUsecase
	*usecase.GetUserOrganizationActiveUsecase
	*usecase.UploadImageUseCase
	*usecase.UserImagesUsecase
	*usecase.LanguagesConfigUsecase
	*usecase.UserSettingUseCase
	*usecase.OwnerAssignUseCase
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

	var orgAdminResp *response.OrganizationAdmin = nil

	if len(userEntity.Organizations) > 0 {
		// Lấy danh sách OrgID mà user là manager
		managedOrgIDs, err := userEntity.GetManagedOrganizationIDs(receiver.GetUserEntityUseCase.GetDB())
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: "failed to get managed organizations",
			})
			return
		}

		// So sánh với các org đã preload, map sang OrganizationAdmin nếu khớp
		for _, org := range userEntity.Organizations {
			if lo.Contains(managedOrgIDs, org.ID.String()) {
				orgAdminResp = &response.OrganizationAdmin{
					ID:               org.ID.String(),
					OrganizationName: org.OrganizationName,
					Avatar:           org.Avatar,
					AvatarURL:        org.AvatarURL,
					Address:          org.Address,
					Description:      org.Description,
					CreatedAt:        org.CreatedAt,
					UpdatedAt:        org.UpdatedAt,
				}
				break // chỉ lấy 1 org đầu tiên mà user là manager
			}
		}

		// Nếu không có org nào user quản lý, lấy org đầu tiên làm mặc định
		if orgAdminResp == nil && len(userEntity.Organizations) > 0 {
			org := userEntity.Organizations[0]
			orgAdminResp = &response.OrganizationAdmin{
				ID:               org.ID.String(),
				OrganizationName: org.OrganizationName,
				Avatar:           org.Avatar,
				AvatarURL:        org.AvatarURL,
				Address:          org.Address,
				Description:      org.Description,
				CreatedAt:        org.CreatedAt,
				UpdatedAt:        org.UpdatedAt,
			}
		}

	}

	// get is Deactive
	isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4User(userEntity.ID.String())

	// get avatars
	avatars, _ := receiver.UserImagesUsecase.GetAvt4Owner(userEntity.ID.String(), value.OwnerRoleUser)

	// get org name of student for parent
	studentOrgs, _ := receiver.StudentApplicationUseCase.GetStudentOrganizationsByUser(userEntity.ID.String())

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.UserEntityResponseV2{
			ID:                  userEntity.ID.String(),
			Username:            userEntity.Username,
			Fullname:            userEntity.Fullname,
			Nickname:            userEntity.Nickname,
			Phone:               userEntity.Phone,
			Email:               userEntity.Email,
			Dob:                 userEntity.Birthday.Format("2006-01-02"),
			QRLogin:             userEntity.QRLogin,
			Avatar:              userEntity.Avatar,
			AvatarURL:           userEntity.AvatarURL,
			IsBlocked:           userEntity.IsBlocked,
			BlockedAt:           userEntity.BlockedAt.Format("2006-01-02"),
			Organization:        organizations,
			CreatedAt:           userEntity.CreatedAt.Format("2006-01-02"),
			Roles:               &roleListResponse,
			Devices:             &deviceListResponse,
			OrganizationAdmin:   orgAdminResp,
			IsDeactive:          isDeactive,
			IsSuperAdmin:        userEntity.IsSuperAdmin(),
			Avatars:             avatars,
			StudentOrganization: studentOrgs,
			ReLoginWeb:          userEntity.ReLoginWeb,
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

	// get user organization active
	userOrgActive, err := receiver.GetUserOrganizationActiveUsecase.GetUserOrganizationActive(userEntity.ID.String())
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	// lấy toàn bộ Organizations trong userEntity làm default
	for _, org := range userEntity.Organizations {
		userOrgActive.DefaultOrganization = append(userOrgActive.DefaultOrganization, response.OrganizationActive{
			ID:               org.ID.String(),
			OrganizationName: org.OrganizationName,
			Avatar:           org.Avatar,
			AvatarURL:        org.AvatarURL,
			CreatedAt:        org.CreatedAt,
		})
	}

	// get avatars
	avatars, _ := receiver.UserImagesUsecase.GetAvt4Owner(userEntity.ID.String(), value.OwnerRoleUser)

	// get user setting
	settings, _ := receiver.UserSettingUseCase.GetByOwner(userEntity.ID.String())

	// get org name of student for parent
	studentOrgs, _ := receiver.StudentApplicationUseCase.GetStudentOrganizationsByUser(userEntity.ID.String())

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.UserEntityResponse{
			ID:                     userEntity.ID.String(),
			Username:               userEntity.Username,
			Fullname:               userEntity.Fullname,
			Nickname:               userEntity.Nickname,
			Phone:                  userEntity.Phone,
			Email:                  userEntity.Email,
			Dob:                    userEntity.Birthday.Format("2006-01-02"),
			QRLogin:                userEntity.QRLogin,
			Avatar:                 userEntity.Avatar,
			AvatarURL:              userEntity.AvatarURL,
			IsBlocked:              userEntity.IsBlocked,
			BlockedAt:              userEntity.BlockedAt.Format("2006-01-02"),
			Organization:           organizations,
			CreatedAt:              userEntity.CreatedAt.Format("2006-01-02"),
			Roles:                  &roleListResponse,
			Devices:                &deviceListResponse,
			CustomID:               userEntity.CustomID,
			UserOrganizationActive: *userOrgActive,
			Avatars:                avatars,
			Settings:               settings,
			StudentOrganization:    studentOrgs,
			CreatedIndex:           userEntity.CreatedIndex,
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
	}, value.ForRegister)

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
			Error: "Upload avatar fail",
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Upload avatar successfully",
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
	role := c.Query("role")
	name := strings.ToLower(c.Query("name"))
	statusParam := c.Query("status")

	var status value.SearchUserStatus = value.SearchUserStatusAll
	if value.IsValidSearchUserStatus(statusParam) {
		status = value.SearchUserStatus(statusParam)
	} else {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "status param is not valid",
		})
		return
	}

	isAll := role == "all"

	// Chuẩn bị cấu trúc response
	var (
		users    = make([]response.UserResponse, 0)
		children = make([]response.ChildrenResponse, 0)
		students = make([]response.StudentResponse, 0)
		teachers = make([]response.TeacherResponse, 0)
		staffs   = make([]response.StaffResponse, 0)
		parents  = make([]response.ParentResponse, 0)
	)

	if isAll {
		// Lấy tất cả
		rawUsers, _ := receiver.GetAllUsers4Search(c)
		rawChildren, _ := receiver.ChildUseCase.GetAll4Search(c)
		rawStudents, _ := receiver.StudentApplicationUseCase.GetAllStudents4Search(c)
		rawTeachers, _ := receiver.TeacherApplicationUseCase.GetAllTeachers4Search(c)
		rawStaffs, _ := receiver.StaffApplicationUseCase.GetAllStaff4Search(c)
		rawParents, _ := receiver.GetAllParents4Search(c)

		// Map sang response
		for _, u := range rawUsers {
			// get is_deactive
			isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4User(u.ID.String())
			users = append(users, response.UserResponse{
				ID:         u.ID.String(),
				Username:   u.Username,
				Nickname:   u.Nickname,
				IsDeactive: isDeactive,
			})
		}
		for _, c := range rawChildren {
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(c.ID.String(), value.OwnerRoleChild)
			children = append(children, response.ChildrenResponse{
				ChildID:      c.ID.String(),
				ChildName:    c.ChildName,
				CreatedIndex: c.CreatedIndex,
				Avatar:       avatar,
			})
		}
		for _, s := range rawStudents {
			isDeactive, _ := receiver.StudentBlockSettingUsecase.GetDeactive4Student(s.StudentID)
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(s.StudentID, value.OwnerRoleStudent)
			students = append(students, response.StudentResponse{
				StudentID:    s.StudentID,
				StudentName:  s.StudentName,
				IsDeactive:   isDeactive,
				CreatedIndex: s.CreatedIndex,
				Avatar:       avatar,
			})
		}
		for _, t := range rawTeachers {
			isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4Teacher(t.TeacherID)
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(t.TeacherID, value.OwnerRoleTeacher)
			teachers = append(teachers, response.TeacherResponse{
				TeacherID:    t.TeacherID,
				TeacherName:  t.TeacherName,
				IsDeactive:   isDeactive,
				CreatedIndex: t.CreatedIndex,
				Avatar:       avatar,
			})
		}
		for _, s := range rawStaffs {
			isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4Teacher(s.StaffID)
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(s.StaffID, value.OwnerRoleStaff)
			staffs = append(staffs, response.StaffResponse{
				StaffID:      s.StaffID,
				StaffName:    s.StaffName,
				IsDeactive:   isDeactive,
				CreatedIndex: s.CreatedIndex,
				Avatar:       avatar,
			})
		}
		for _, p := range rawParents {
			// get is_deactive
			isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4User(p.ID.String())
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(p.ID.String(), value.OwnerRoleParent)
			parents = append(parents, response.ParentResponse{
				ParentID:   p.ID.String(),
				ParentName: p.Nickname,
				IsDeactive: isDeactive,
				Avatar:     avatar,
			})
		}

		// filter by name
		users = helper.FilterUsersByName(users, name)
		children = helper.FilterChildrenByName(children, name)
		students = helper.FilterStudentByName(students, name)
		teachers = helper.FilterTeacherByName(teachers, name)
		staffs = helper.FilterStaffByName(staffs, name)
		parents = helper.FilterParentByName(parents, name)

		// filter by status
		users = helper.FilterUsersByStatus(users, status)
		teachers = helper.FilterTeachersByStatus(teachers, status)
		staffs = helper.FilterStaffsByStatus(staffs, status)
		parents = helper.FilterParentsByStatus(parents, status)

		// Trả kết quả
		c.JSON(http.StatusOK, response.SucceedResponse{
			Code: http.StatusOK,
			Data: response.SearchUserResponse{
				Users:    users,
				Children: children,
				Students: students,
				Teachers: teachers,
				Staffs:   staffs,
				Parents:  parents,
			},
		})
		return
	}

	// Nếu role không hợp lệ
	if !value.IsValidRoleSignUp(role) {
		c.JSON(http.StatusOK, response.SucceedResponse{
			Code: http.StatusOK,
			Data: response.SearchUserResponse{
				Users:    users,
				Children: children,
				Students: students,
				Teachers: teachers,
				Staffs:   staffs,
			},
		})
		return
	}

	// Xử lý từng role cụ thể
	switch value.RoleSignUp(role) {
	case value.RoleChild:
		rawChildren, _ := receiver.ChildUseCase.GetAll4Search(c)
		for _, c := range rawChildren {
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(c.ID.String(), value.OwnerRoleChild)
			children = append(children, response.ChildrenResponse{
				ChildID:      c.ID.String(),
				ChildName:    c.ChildName,
				CreatedIndex: c.CreatedIndex,
				Avatar:       avatar,
			})
		}
		children = helper.FilterChildrenByName(children, name)

	case value.RoleStudent:
		rawStudents, _ := receiver.StudentApplicationUseCase.GetAllStudents4Search(c)
		for _, s := range rawStudents {
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(s.StudentID, value.OwnerRoleStudent)
			students = append(students, response.StudentResponse{
				StudentID:    s.StudentID,
				StudentName:  s.StudentName,
				CreatedIndex: s.CreatedIndex,
				Avatar:       avatar,
			})
		}
		students = helper.FilterStudentByName(students, name)

	case value.RoleTeacher:
		rawTeachers, _ := receiver.TeacherApplicationUseCase.GetAllTeachers4Search(c)
		for _, t := range rawTeachers {
			isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4Teacher(t.TeacherID)
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(t.TeacherID, value.OwnerRoleTeacher)
			teachers = append(teachers, response.TeacherResponse{
				TeacherID:    t.TeacherID,
				TeacherName:  t.TeacherName,
				IsDeactive:   isDeactive,
				CreatedIndex: t.CreatedIndex,
				Avatar:       avatar,
			})
		}
		teachers = helper.FilterTeacherByName(teachers, name)
		teachers = helper.FilterTeachersByStatus(teachers, status)

	case value.RoleStaff:
		rawStaffs, _ := receiver.StaffApplicationUseCase.GetAllStaff4Search(c)
		for _, s := range rawStaffs {
			isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4Teacher(s.StaffID)
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(s.StaffID, value.OwnerRoleStaff)
			staffs = append(staffs, response.StaffResponse{
				StaffID:      s.StaffID,
				StaffName:    s.StaffName,
				IsDeactive:   isDeactive,
				CreatedIndex: s.CreatedIndex,
				Avatar:       avatar,
			})
		}
		staffs = helper.FilterStaffByName(staffs, name)
		staffs = helper.FilterStaffsByStatus(staffs, status)

	case value.User:
		rawUsers, _ := receiver.GetAllUsers4Search(c)
		for _, u := range rawUsers {
			// get is_deactive
			isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4User(u.ID.String())
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(u.ID.String(), value.OwnerRoleUser)
			users = append(users, response.UserResponse{
				ID:         u.ID.String(),
				Username:   u.Username,
				Nickname:   u.Nickname,
				IsDeactive: isDeactive,
				Avatar:     avatar,
			})
		}
		users = helper.FilterUsersByName(users, name)
		users = helper.FilterUsersByStatus(users, status)

	case value.Parent:
		rawParents, _ := receiver.GetAllParents4Search(c)
		for _, p := range rawParents {
			isDeactive, _ := receiver.GetDeactive4User(p.ID.String())
			avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(p.ID.String(), value.OwnerRoleParent)
			parents = append(parents, response.ParentResponse{
				ParentID:   p.ID.String(),
				ParentName: p.Nickname,
				IsDeactive: isDeactive,
				Avatar:     avatar,
			})
		}
		parents = helper.FilterParentByName(parents, name)
		parents = helper.FilterParentsByStatus(parents, status)

	case value.RoleOrganization:
		// Không xử lý gì
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.SearchUserResponse{
			Users:    users,
			Children: children,
			Students: students,
			Teachers: teachers,
			Staffs:   staffs,
			Parents:  parents,
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

func (receiver *UserEntityController) GetStaff4WebAdmin(context *gin.Context) {
	staffID := context.Param("id")
	if staffID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing staff ID",
		})
		return
	}

	staff, err := receiver.StaffApplicationUseCase.GetStaffByID(staffID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get staff",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: staff,
	})

}

func (receiver *UserEntityController) GetParent4WebAdmin(context *gin.Context) {
	parentID := context.Param("id")
	if parentID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing parent ID",
		})
		return
	}

	parent, err := receiver.ParentUseCase.GetParentByID(parentID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get parent",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: parent,
	})

}

func (receiver *UserEntityController) GetStudent4App(context *gin.Context) {

	studentID := context.Param("id")
	if studentID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "student id be required",
		})
		return
	}

	deviceID := context.Param("device_id")
	if deviceID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "device id be required",
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

	student, err := receiver.StudentApplicationUseCase.GetStudentByID4App(context, studentID, deviceID)
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

func (receiver *UserEntityController) AddCustomID2Student(context *gin.Context) {
	var req request.AddCustomId2StudentRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	err := receiver.StudentApplicationUseCase.AddCustomID(req)
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

func (receiver *UserEntityController) AddCustomID2User(context *gin.Context) {
	var req request.AddCustomID2UserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	err := receiver.UpdateUserEntityUseCase.UpdateCustomIDByUserID(req)
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

func (receiver *UserEntityController) UploadAvatarV2(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	ownerID := c.PostForm("owner_id")
	if ownerID == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "owner id is required",
		})
		return
	}

	ownerRole := c.PostForm("owner_role")
	if ownerRole == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "owner role is required",
		})
		return
	}

	fileName := c.PostForm("file_name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "file name is required",
		})
		return
	}
	re := regexp.MustCompile(`\s+`)
	fileName = re.ReplaceAllString(fileName, "_")

	indexStr := c.PostForm("index")
	if indexStr == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "index is required",
		})
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid index",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	defer file.Close()

	dataBytes := make([]byte, fileHeader.Size)
	if _, err := bufio.NewReader(file).Read(dataBytes); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// 1. Gọi usecase upload image
	_, img, err := receiver.UploadImageUseCase.UploadImagev2(
		dataBytes,
		"avatar",
		fileHeader.Filename,
		fileName,
		uploader.UploadPrivate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	// 2. Gọi usecase lưu vào user_images
	req := request.UploadAvatarRequest{
		OwnerID:   ownerID,
		OwnerRole: value.OwnerRole(ownerRole),
		ImageID:   img.ID,
		Index:     index,
	}

	if err := receiver.UserImagesUsecase.UploadAvt(req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	// 3. Trả về response
	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Uopload avatar successfully",
		Data:    nil,
	})
}

func (receiver *UserEntityController) UpdateIsMain(context *gin.Context) {
	var req request.UpdateIsMainAvatar
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := receiver.UserImagesUsecase.UpdateIsMain(req); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update main avatar",
			Error:   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Update main avatar successfully",
	})

}

func (receiver *UserEntityController) DeleteUserAvatar(context *gin.Context) {
	ownerID := context.Query("owner_id")
	ownerRole := context.Query("owner_role")
	indexStr := context.Query("index")

	if ownerID == "" || ownerRole == "" || indexStr == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "owner_id, owner_role and index are required",
		})
		return
	}

	// validate owner role
	role := value.OwnerRole(ownerRole)
	if !role.IsValid() {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid owner role",
		})
		return
	}
	// Convert index về int
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "index must be a number",
			Error:   err.Error(),
		})
		return
	}

	req := request.DeleteUserAvatarRequest{
		OwnerID:   ownerID,
		OwnerRole: role,
		Index:     index,
	}

	if err := receiver.UserImagesUsecase.DeleteUserAvatar(req); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete avatar",
			Error:   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Delete avatar successfully",
	})
}

func (receiver *UserEntityController) GetStudent4Gateway(context *gin.Context) {
	studentID := context.Param("student_id")
	if studentID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing student ID",
		})
		return
	}

	student, err := receiver.StudentApplicationUseCase.GetStudent4Gateway(studentID)
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

func (receiver *UserEntityController) GetTeacher4Gateway(context *gin.Context) {
	teacherID := context.Param("teacher_id")
	if teacherID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing teacher ID",
		})
		return
	}

	teacher, err := receiver.TeacherApplicationUseCase.GetTeacher4Gateway(teacherID)
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

func (receiver *UserEntityController) GetStaff4Gateway(context *gin.Context) {
	staffID := context.Param("staff_id")
	if staffID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing staff ID",
		})
		return
	}

	staff, err := receiver.StaffApplicationUseCase.GetStaff4Gateway(staffID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get staff",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: staff,
	})
}

func (receiver *UserEntityController) GetUser4Gateway(context *gin.Context) {
	userID := context.Param("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	userEntity, err := receiver.GetUserByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: "Failed to get user: " + err.Error(),
		})

		return
	}

	// get avatars
	avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(userEntity.ID.String(), value.OwnerRoleUser)

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: response.GetUser4Gateway{
			UserID:   userEntity.ID.String(),
			UserName: userEntity.Nickname,
			Avatar:   avatar,
		},
	})
}

func (receiver *UserEntityController) GetStudentLanguageConfig(context *gin.Context) {
	studentID := context.Param("id")
	if studentID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing student ID",
		})
		return
	}

	studentLangConfig, err := receiver.LanguagesConfigUsecase.GetLanguagesConfigByOwner(context, studentID, value.OwnerRoleLangStudent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get student study language config",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: studentLangConfig,
	})
}

func (receiver *UserEntityController) GetTeacherByUser4Gateway(context *gin.Context) {
	userID := context.Param("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	teacher, err := receiver.TeacherApplicationUseCase.GetTeacherByUser4Gateway(userID)
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

func (receiver *UserEntityController) GetTeachersByUser4Gateway(context *gin.Context) {
	userID := context.Param("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	teacher, err := receiver.TeacherApplicationUseCase.GetTeachersByUser4Gateway(userID)
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

func (receiver *UserEntityController) GetStaffByUser4Gateway(context *gin.Context) {
	userID := context.Param("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	staff, err := receiver.StaffApplicationUseCase.GetStaffByUser4Gateway(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get staff",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: staff,
	})
}

func (receiver *UserEntityController) GetStaffsByUser4Gateway(context *gin.Context) {
	userID := context.Param("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	staff, err := receiver.StaffApplicationUseCase.GetStaffsByUser4Gateway(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get staff",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: staff,
	})
}

func (receiver *UserEntityController) GetListOwner2Assign(context *gin.Context) {
	organizationID := context.Param("organization_id")
	if organizationID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing organization ID",
		})
		return
	}

	listOwner, err := receiver.OwnerAssignUseCase.GetListOwner2Assign(context, organizationID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get list owner",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: listOwner,
	})
}

func (receiver *UserEntityController) GetTeacherByOrgAndUser4Gateway(context *gin.Context) {
	userID := context.Param("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	organizationID := context.Param("organization_id")
	if organizationID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	teacher, err := receiver.TeacherApplicationUseCase.GetTeacherByOrgAndUser4Gateway(userID, organizationID)
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

func (receiver *UserEntityController) GetStaffByOrgAndUser4Gateway(context *gin.Context) {
	userID := context.Param("user_id")
	if userID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	organizationID := context.Param("organization_id")
	if organizationID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Missing user ID",
		})
		return
	}

	staff, err := receiver.StaffApplicationUseCase.GetStaffByOrgAndUser4Gateway(userID, organizationID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get staff",
			Error:   err.Error(),
		})
		return
	}

	// Thành công
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: staff,
	})
}

func (receiver *UserEntityController) SetReLogin(context *gin.Context) {

	err := receiver.UpdateUserEntityUseCase.SetReLogin()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
	})
}
