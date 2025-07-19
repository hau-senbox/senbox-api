package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"strings"

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

	return response.GetCommonMenuResponse{
		Components: componentsList,
	}
}

func (receiver *GetMenuUseCase) GetCommonMenuByUser(ctx *gin.Context) response.GetCommonMenuByUserResponse {
	var componentMenus []response.ComponentCommonMenuByUser

	userID := ctx.GetString("user_id")
	children, err := receiver.ChildRepository.GetByParentID(userID)

	if err == nil && children != nil {
		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByRoleName("Child")
		var childOrg = ""
		if err == nil || roleOrg != nil {
			childOrg = roleOrg.OrgProfile
		}
		for _, child := range children {
			childComponent := buildComponent(
				uuid.NewString(),
				fmt.Sprintf("Child Profile: %s", child.ChildName),
				"child_profile",
				"icon/accident_and_injury_report_1745206766342940327.png",
				"button_form",
				childOrg,
			)

			componentMenus = append(componentMenus, response.ComponentCommonMenuByUser{
				ChildID:   child.ID.String(),
				Component: childComponent,
			})
		}
	}

	// check teacher menu
	if teacherComponent, _ := receiver.getProfileComponentByRole("Teacher", userID); teacherComponent != nil {
		componentMenus = append(componentMenus, *teacherComponent)
	}

	//checck student menu
	if teacherComponent, _ := receiver.getProfileComponentByRole("Student", userID); teacherComponent != nil {
		componentMenus = append(componentMenus, *teacherComponent)
	}

	//check staff menu
	if teacherComponent, _ := receiver.getProfileComponentByRole("Staff", userID); teacherComponent != nil {
		componentMenus = append(componentMenus, *teacherComponent)
	}

	//check org menu
	if teacherComponent, _ := receiver.getProfileComponentByRole("Sign up ORganise", userID); teacherComponent != nil {
		componentMenus = append(componentMenus, *teacherComponent)
	}

	return response.GetCommonMenuByUserResponse{
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
				Value: helper.BuildSectionValueMenu(string(c.Value), c),
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

func (receiver *GetMenuUseCase) getProfileComponentByRole(roleName, userID string) (*response.ComponentCommonMenuByUser, error) {
	// Lấy thông tin role theo tên
	roleOrg, err := receiver.RoleOrgSignUpRepository.GetByRoleName(roleName)
	if err != nil || roleOrg == nil {
		return nil, nil
	}

	// Lấy form theo OrgCode của role
	form, err := receiver.FormRepository.GetFormByQRCode(roleOrg.OrgCode)
	if err != nil || form == nil {
		return nil, nil
	}

	// Kiểm tra xem user đã nộp form hay chưa
	submission, err := receiver.SubmissionRepository.GetByUserIdAndFormId(userID, form.ID)
	if err != nil || submission == nil {
		return nil, nil
	}

	// Tạo component
	component := buildComponent(
		uuid.NewString(),
		fmt.Sprintf("%s Profile", roleName),
		fmt.Sprintf("%s_profile", strings.ToLower(roleName)),
		"icon/accident_and_injury_report_1745206766342940327.png",
		"button_form",
		roleOrg.OrgProfile,
	)

	return &response.ComponentCommonMenuByUser{
		Component: component,
	}, nil
}

func (receiver *GetMenuUseCase) GetSectionMenu4WebAdmin() ([]response.GetMenuSectionResponse, error) {
	var roleOrgChildId string
	roleOrgChild, _ := receiver.RoleOrgSignUpRepository.GetByRoleName("Child")
	if roleOrgChild != nil {
		roleOrgChildId = roleOrgChild.ID.String()
	}
	componentsList, err := receiver.ComponentRepository.GetBySectionID(roleOrgChildId)
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
				Value: helper.BuildSectionValueMenu(string(c.Value), c),
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
