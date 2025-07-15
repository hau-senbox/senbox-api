package usecase

import (
	"encoding/json"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type GetMenuUseCase struct {
	MenuRepository          *repository.MenuRepository
	UserEntityRepository    *repository.UserEntityRepository
	OrganizationRepository  *repository.OrganizationRepository
	DeviceRepository        *repository.DeviceRepository
	RoleOrgSignUpRepository *repository.RoleOrgSignUpRepository
	FormRepository          *repository.FormRepository
	SubmissionRepository    *repository.SubmissionRepository
}

func (receiver *GetMenuUseCase) GetSuperAdminMenu() ([]menu.SuperAdminMenu, error) {
	return receiver.MenuRepository.GetSuperAdminMenu()
}

func (receiver *GetMenuUseCase) GetOrgMenu(orgID string) ([]menu.OrgMenu, error) {
	org, err := receiver.OrganizationRepository.GetByID(orgID)
	if err != nil {
		return nil, err
	}

	return receiver.MenuRepository.GetOrgMenu(org.ID.String())
}

func (receiver *GetMenuUseCase) GetStudentMenu(userID string) ([]menu.UserMenu, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	present := lo.ContainsBy(user.Roles, func(role entity.SRole) bool {
		return role.Role == entity.Student
	})
	if !present {
		return nil, errors.New("failed to get student menu")
	}

	return receiver.MenuRepository.GetUserMenu(user.ID.String())
}

func (receiver *GetMenuUseCase) GetTeacherMenu(userID string) ([]menu.UserMenu, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	present := lo.ContainsBy(user.Roles, func(role entity.SRole) bool {
		return role.Role == entity.Teacher
	})
	if !present {
		return nil, errors.New("failed to get teacher menu")
	}

	return receiver.MenuRepository.GetUserMenu(user.ID.String())
}

func (receiver *GetMenuUseCase) GetUserMenu(userID string) ([]menu.UserMenu, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	return receiver.MenuRepository.GetUserMenu(user.ID.String())
}

func (receiver *GetMenuUseCase) GetDeviceMenu(deviceID string) ([]menu.DeviceMenu, error) {
	device, err := receiver.DeviceRepository.GetDeviceByID(deviceID)
	if err != nil {
		return nil, err
	}

	return receiver.MenuRepository.GetDeviceMenu(device.ID)
}

func (receiver *GetMenuUseCase) GetDeviceMenuByOrg(organizationID string) ([]menu.DeviceMenu, error) {
	return receiver.MenuRepository.GetDeviceMenuByOrg(organizationID)
}

func (receiver *GetMenuUseCase) GetCommonMenu(ctx *gin.Context) response.GetCommonMenuResponse {
	componentsList := []response.ComponentResponse{
		buildComponent(uuid.NewString(), "My Account", "account_profile", "SENBOX.ORG/ACCOUNT-PROFILE"),
		buildComponent(uuid.NewString(), "Add Roles", "add_roles", "SENBOX.ORG/ADD-ROLES"),
	}

	// 1. Get role "Child"
	var childOrgCode string
	roleSignUp, err := receiver.RoleOrgSignUpRepository.GetByRoleName("Child")
	if err == nil && roleSignUp != nil && roleSignUp.OrgCode != "" {
		childOrgCode = roleSignUp.OrgCode
	}

	// 2. Get form by QRCode (childOrgCode)
	form, _ := receiver.FormRepository.GetFormByQRCode(childOrgCode)

	var formChildId uint64
	if form != nil {
		formChildId = form.ID
	}

	userID := ctx.GetString("user_id")
	submission, err := receiver.SubmissionRepository.GetByUserIdAndFormId(userID, formChildId)

	if err == nil && submission != nil {
		childComponent := buildComponent(uuid.NewString(), "Child Profile", "child_profile", "SENBOX.ORG/CHILD-PROFILE")
		componentsList = append(componentsList, childComponent)
	}

	return response.GetCommonMenuResponse{
		Component: componentsList,
	}
}

func buildComponent(id, name, key, formQR string) response.ComponentResponse {
	valueObject := map[string]interface{}{
		"visible": true,
		"icon":    "",
		"color":   "#86DEFF",
		"form_qr": formQR,
	}
	valueBytes, _ := json.Marshal(valueObject)

	return response.ComponentResponse{
		ID:    id,
		Name:  name,
		Type:  "button_form",
		Key:   key,
		Value: string(valueBytes),
	}
}
