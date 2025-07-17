package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
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
	ComponentRepository     *repository.ComponentRepository
	ChildRepository         *repository.ChildRepository
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
	// Tạo component dạng ComponentResponse
	componentsList := []response.ComponentResponse{
		buildComponent(
			uuid.NewString(),
			"My Account Profiles",
			"my_account_profile",
			"icon/accident_and_injury_report_1745206766342940327.png",
			"button_form",
			"SENBOX.ORG/MY-ACCOUNT-PROFILES",
		),
	}

	// Gói vào ComponentCommonMenu (vì định nghĩa như vậy)
	componentCommon := response.ComponentCommonMenu{
		ChildID:    "", // Hoặc bạn có thể gán userID nếu cần
		Components: componentsList,
	}

	return response.GetCommonMenuResponse{
		Components: []response.ComponentCommonMenu{componentCommon},
	}
}

func (receiver *GetMenuUseCase) GetCommonMenuByUser(ctx *gin.Context) response.GetCommonMenuResponse {
	var componentMenus []response.ComponentCommonMenu

	userID := ctx.GetString("user_id")
	children, err := receiver.ChildRepository.GetByParentID(userID)
	if err == nil && children != nil {
		for _, child := range children {
			childComponent := buildComponent(
				uuid.NewString(),
				fmt.Sprintf("Child Profile: %s", child.ChildName),
				"child_profile",
				"icon/accident_and_injury_report_1745206766342940327.png",
				"button_form",
				"SENBOX.ORG/CHILD-PROFILE",
			)

			componentMenus = append(componentMenus, response.ComponentCommonMenu{
				ChildID:    child.ID.String(),
				Components: []response.ComponentResponse{childComponent},
			})
		}
	}

	return response.GetCommonMenuResponse{
		Components: componentMenus,
	}
}

func buildComponent(id, name, key, icon, typeName, formQR string) response.ComponentResponse {
	valueObject := map[string]interface{}{
		"id":   id,
		"name": name,
		"type": typeName,
		"key":  "",
		"value": map[string]interface{}{
			"visible": true,
			"icon":    icon,
			"color":   "#86DEFF",
			"form_qr": formQR,
		},
		"visible": true,
		"icon":    icon,
		"color":   "#86DEFF",
		"form_qr": formQR,
	}

	valueBytes, _ := json.Marshal(valueObject)

	return response.ComponentResponse{
		ID:    id,
		Name:  name,
		Type:  typeName,
		Key:   key,
		Value: string(valueBytes),
	}
}

func (receiver *GetMenuUseCase) GetSectionMenu() ([]response.GetMenuSectionResponse, error) {
	componentsList, err := receiver.ComponentRepository.GetAllByKey("section-menu")
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]components.Component)
	for _, c := range componentsList {
		grouped[c.SectionID] = append(grouped[c.SectionID], c)
	}

	var result []response.GetMenuSectionResponse
	for sectionID, comps := range grouped {
		var componentResponses []response.ComponentResponse
		for i, c := range comps {
			componentResponses = append(componentResponses, response.ComponentResponse{
				ID:    c.ID.String(),
				Name:  c.Name,
				Type:  string(c.Type),
				Key:   c.Key,
				Value: string(c.Value),
				Order: i,
			})
		}

		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByID(sectionID)
		if err != nil {
			return nil, err
		}
		sectionName := ""
		if roleOrg != nil {
			sectionName = roleOrg.RoleName
		}

		result = append(result, response.GetMenuSectionResponse{
			SectionID:   sectionID,
			SectionName: sectionName,
			Components:  componentResponses,
		})
	}

	return result, nil
}
