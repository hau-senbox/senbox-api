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
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type GetMenuUseCase struct {
	MenuRepository                     *repository.MenuRepository
	UserEntityRepository               *repository.UserEntityRepository
	OrganizationRepository             *repository.OrganizationRepository
	DeviceRepository                   *repository.DeviceRepository
	RoleOrgSignUpRepository            *repository.RoleOrgSignUpRepository
	FormRepository                     *repository.FormRepository
	SubmissionRepository               *repository.SubmissionRepository
	ComponentRepository                *repository.ComponentRepository
	ChildRepository                    *repository.ChildRepository
	StudentAppRepo                     *repository.StudentApplicationRepository
	ChildMenuUseCase                   *ChildMenuUseCase
	StudentMenuUseCase                 *StudentMenuUseCase
	GetUserEntityUseCase               *GetUserEntityUseCase
	TeacherMenuUseCase                 *TeacherMenuUseCase
	TeacherRepository                  *repository.TeacherApplicationRepository
	StaffMenuUsecase                   *StaffMenuUseCase
	StaffApplicationRepo               *repository.StaffApplicationRepository
	OrganizationMenuTemplateRepository *repository.OrganizationMenuTemplateRepository
	ParentMenuUsecase                  *ParentMenuUseCase
	UserImageUsecase                   *UserImagesUsecase
	DepartmentGateway                  gateway.DepartmentGateway
	DepartmentMenuUseCase              *DepartmentMenuUseCase
	SuperAdminEmergencyMenuRepo        *repository.SuperAdminEmergencyMenuRepository
	OrganizationEmergencyMenuRepo      *repository.OrganizationEmergencyMenuRepository
	LanguageSettingRepo                *repository.LanguageSettingRepository
	ParentRepo                         *repository.ParentRepository
}

func (receiver *GetMenuUseCase) GetSuperAdminMenu() ([]menu.SuperAdminMenu, error) {
	return receiver.MenuRepository.GetSuperAdminMenu()
}

func (receiver *GetMenuUseCase) GetSuperAdminMenu4Web() (*response.GetSuperAdminMenuResponse4Web, error) {
	menus, err := receiver.MenuRepository.GetSuperAdminMenu()
	if err != nil {
		return nil, err
	}

	// Gom ComponentID
	componentIDs := make([]uuid.UUID, 0, len(menus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentDirectionMap := make(map[uuid.UUID]menu.Direction)

	for _, m := range menus {
		componentIDs = append(componentIDs, m.ComponentID)
		componentOrderMap[m.ComponentID] = m.Order
		componentDirectionMap[m.ComponentID] = m.Direction
	}

	// Lấy tất cả components
	components, err := receiver.ComponentRepository.GetByIDs(componentIDs)
	if err != nil {
		return nil, err
	}

	// Gom theo direction + language
	topByLang := make(map[uint][]response.ComponentResponse)
	bottomByLang := make(map[uint][]response.ComponentResponse)
	langMap := make(map[uint]entity.LanguageSetting)

	for _, comp := range components {
		menuResp := response.ComponentResponse{
			ID:         comp.ID.String(),
			Name:       comp.Name,
			Type:       comp.Type.String(),
			Key:        comp.Key,
			Value:      string(comp.Value),
			Order:      componentOrderMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		// cache language
		if _, ok := langMap[comp.LanguageID]; !ok {
			langSetting, err := receiver.LanguageSettingRepo.GetByID(comp.LanguageID)
			if err != nil {
				return nil, err
			}
			if langSetting != nil {
				langMap[comp.LanguageID] = *langSetting
			}
		}

		if componentDirectionMap[comp.ID] == menu.Top {
			topByLang[comp.LanguageID] = append(topByLang[comp.LanguageID], menuResp)
		} else {
			bottomByLang[comp.LanguageID] = append(bottomByLang[comp.LanguageID], menuResp)
		}
	}

	// Build response
	resp := &response.GetSuperAdminMenuResponse4Web{
		Top:    make([]response.GetMenus4Web, 0, len(topByLang)),
		Bottom: make([]response.GetMenus4Web, 0, len(bottomByLang)),
	}

	for langID, comps := range topByLang {
		resp.Top = append(resp.Top, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	for langID, comps := range bottomByLang {
		resp.Bottom = append(resp.Bottom, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	return resp, nil
}

func (receiver *GetMenuUseCase) GetSuperAdminMenu4App(ctx *gin.Context) ([]menu.SuperAdminMenu, error) {
	appLanguage, _ := ctx.Get("app_language")
	return receiver.MenuRepository.GetSuperAdminMenuByLanguage(appLanguage.(uint))
}

func (receiver *GetMenuUseCase) GetOrgMenu(orgID string) ([]menu.OrgMenu, error) {
	org, err := receiver.OrganizationRepository.GetByID(orgID)
	if err != nil {
		return nil, err
	}

	return receiver.MenuRepository.GetOrgMenu(org.ID.String())
}

func (receiver *GetMenuUseCase) GetOrgMenu4Web(orgID string) (*response.GetOrganizationAdminMenuResponse4Web, error) {
	menus, err := receiver.MenuRepository.GetOrgMenu(orgID)
	if err != nil {
		return nil, err
	}

	// Gom ComponentID
	componentIDs := make([]uuid.UUID, 0, len(menus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentDirectionMap := make(map[uuid.UUID]menu.Direction)

	for _, m := range menus {
		componentIDs = append(componentIDs, m.ComponentID)
		componentOrderMap[m.ComponentID] = m.Order
		componentDirectionMap[m.ComponentID] = m.Direction
	}

	// Lấy tất cả components
	components, err := receiver.ComponentRepository.GetByIDs(componentIDs)
	if err != nil {
		return nil, err
	}

	// Gom theo direction + language
	topByLang := make(map[uint][]response.ComponentResponse)
	bottomByLang := make(map[uint][]response.ComponentResponse)
	langMap := make(map[uint]entity.LanguageSetting)

	for _, comp := range components {
		menuResp := response.ComponentResponse{
			ID:         comp.ID.String(),
			Name:       comp.Name,
			Type:       comp.Type.String(),
			Key:        comp.Key,
			Value:      string(comp.Value),
			Order:      componentOrderMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		// cache language
		if _, ok := langMap[comp.LanguageID]; !ok {
			langSetting, err := receiver.LanguageSettingRepo.GetByID(comp.LanguageID)
			if err != nil {
				return nil, err
			}
			if langSetting != nil {
				langMap[comp.LanguageID] = *langSetting
			}
		}

		if componentDirectionMap[comp.ID] == menu.Top {
			topByLang[comp.LanguageID] = append(topByLang[comp.LanguageID], menuResp)
		} else {
			bottomByLang[comp.LanguageID] = append(bottomByLang[comp.LanguageID], menuResp)
		}
	}

	// Build response
	resp := &response.GetOrganizationAdminMenuResponse4Web{
		Top:    make([]response.GetMenus4Web, 0, len(topByLang)),
		Bottom: make([]response.GetMenus4Web, 0, len(bottomByLang)),
	}

	for langID, comps := range topByLang {
		resp.Top = append(resp.Top, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	for langID, comps := range bottomByLang {
		resp.Bottom = append(resp.Bottom, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	return resp, nil
}

func (receiver *GetMenuUseCase) GetOrgMenu4App(ctx *gin.Context, orgID string) ([]menu.OrgMenu, error) {
	org, err := receiver.OrganizationRepository.GetByID(orgID)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")
	return receiver.MenuRepository.GetOrgMenuByLanguage(org.ID.String(), appLanguage.(uint))
}

func (receiver *GetMenuUseCase) GetStudentMenu4App(ctx *gin.Context, studentID string) (*response.GetStudentMenuResponse, error) {
	studentMenu, err := receiver.StudentMenuUseCase.GetByStudentID(ctx, studentID, true)

	if err != nil {
		return nil, fmt.Errorf("failed to get student menu: %w", err)
	}

	if len(studentMenu.Components) == 0 {
		return nil, errors.New("student menu not found")
	}

	// get menu icon key
	img, _ := receiver.UserImageUsecase.GetImg4Ownewr(studentID, value.OwnerRoleStudent)

	menuIconKey := ""
	if img != nil {
		menuIconKey = img.Key
	}

	studentMenu.MenuIconKey = menuIconKey
	return studentMenu, nil
}

func (receiver *GetMenuUseCase) GetTeacherMenu4App(ctx *gin.Context, userID string) (*response.GetTeacherMenuResponse, error) {

	teacher, _ := receiver.TeacherRepository.GetByUserID(userID)

	teacherMenu, err := receiver.TeacherMenuUseCase.GetByTeacherID(ctx, teacher.ID.String(), true)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher menu: %w", err)
	}

	if len(teacherMenu.Components) == 0 {
		return nil, errors.New("teacher menu not found")
	}

	// get menu icon key
	img, _ := receiver.UserImageUsecase.GetImg4Ownewr(teacher.ID.String(), value.OwnerRoleTeacher)

	menuIconKey := ""
	if img != nil {
		menuIconKey = img.Key
	}
	teacherMenu.MenuIconKey = menuIconKey

	return &teacherMenu, nil
}

func (receiver *GetMenuUseCase) GetUserMenu(userID string) ([]menu.UserMenu, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	return receiver.MenuRepository.GetUserMenu(user.ID.String())
}

func (receiver *GetMenuUseCase) GetUserMenu4Web(userID string) ([]response.GetMenus4Web, error) {
	// Lấy user
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	// Lấy danh sách menu gốc (UserMenu)
	userMenus, err := receiver.MenuRepository.GetUserMenu(user.ID.String())
	if err != nil {
		return nil, err
	}

	// Chuẩn bị dữ liệu
	componentIDs := make([]uuid.UUID, 0, len(userMenus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, um := range userMenus {
		componentIDs = append(componentIDs, um.ComponentID)
		componentOrderMap[um.ComponentID] = um.Order
	}

	// Lấy components theo IDs
	components, err := receiver.ComponentRepository.GetByIDs(componentIDs)
	if err != nil {
		return nil, err
	}

	// Gom theo language_id
	menusByLang := make(map[uint][]response.ComponentResponse)
	langMap := make(map[uint]entity.LanguageSetting)

	for _, comp := range components {
		menu := response.ComponentResponse{
			ID:         comp.ID.String(),
			Name:       comp.Name,
			Type:       comp.Type.String(),
			Key:        comp.Key,
			Value:      string(comp.Value),
			Order:      componentOrderMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		menusByLang[comp.LanguageID] = append(menusByLang[comp.LanguageID], menu)

		// Cache language
		if _, ok := langMap[comp.LanguageID]; !ok {
			langSetting, err := receiver.LanguageSettingRepo.GetByID(comp.LanguageID)
			if err != nil {
				return nil, err
			}
			if langSetting != nil {
				langMap[comp.LanguageID] = *langSetting
			}
		}
	}

	// Build []GetMenus4Web
	getMenus := make([]response.GetMenus4Web, 0, len(menusByLang))
	for langID, comps := range menusByLang {
		getMenus = append(getMenus, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	return getMenus, nil
}

func (receiver *GetMenuUseCase) GetUserMenu4App(ctx *gin.Context, userID string) ([]menu.UserMenu, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")
	return receiver.MenuRepository.GetUserMenuByLanguage(user.ID.String(), appLanguage.(uint))
}

func (receiver *GetMenuUseCase) GetDeviceMenu(deviceID string) ([]menu.DeviceMenu, error) {
	device, err := receiver.DeviceRepository.GetDeviceByID(deviceID)
	if err != nil {
		return nil, err
	}

	return receiver.MenuRepository.GetDeviceMenu(device.ID)
}

func (receiver *GetMenuUseCase) GetDeviceMenu4App(ctx *gin.Context, deviceID string) ([]menu.DeviceMenu, error) {
	device, err := receiver.DeviceRepository.GetDeviceByID(deviceID)
	if err != nil {
		return nil, err
	}

	appLanguage, _ := ctx.Get("app_language")
	return receiver.MenuRepository.GetDeviceMenuByLanguage(device.ID, appLanguage.(uint))
}

func (receiver *GetMenuUseCase) GetDeviceMenuByOrg(organizationID string) ([]menu.DeviceMenu, error) {
	return receiver.MenuRepository.GetDeviceMenuByOrg(organizationID)
}

func (receiver *GetMenuUseCase) GetDeviceMenuByOrg4Web(organizationID string) ([]response.GetMenus4Web, error) {
	menus, err := receiver.MenuRepository.GetDeviceMenuByOrg(organizationID)
	if err != nil {
		return nil, err
	}

	// Tạo danh sách componentID để lấy Component
	componentIDs := make([]uuid.UUID, 0, len(menus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range menus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
	}

	// Lấy components theo ID
	components, err := receiver.ComponentRepository.GetByIDs(componentIDs)
	if err != nil {
		return nil, err
	}

	// Gom components theo language_id
	menusByLang := make(map[uint][]response.ComponentResponse)
	langMap := make(map[uint]entity.LanguageSetting)

	for _, comp := range components {
		menu := response.ComponentResponse{
			ID:         comp.ID.String(),
			Name:       comp.Name,
			Type:       comp.Type.String(),
			Key:        comp.Key,
			Value:      string(comp.Value),
			Order:      componentOrderMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		menusByLang[comp.LanguageID] = append(menusByLang[comp.LanguageID], menu)

		// nếu chưa có languageID trong cache -> query DB
		if _, ok := langMap[comp.LanguageID]; !ok {
			langSetting, err := receiver.LanguageSettingRepo.GetByID(comp.LanguageID)
			if err != nil {
				return nil, err
			}
			if langSetting != nil {
				langMap[comp.LanguageID] = *langSetting
			}
		}
	}

	// Build []GetMenus4Web
	getMenus := make([]response.GetMenus4Web, 0, len(menusByLang))
	for langID, comps := range menusByLang {
		getMenus = append(getMenus, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	return getMenus, nil
}

func (receiver *GetMenuUseCase) GetDeviceMenuByOrg4App(ctx *gin.Context, organizationID string) ([]menu.DeviceMenu, error) {
	appLanguage, _ := ctx.Get("app_language")
	return receiver.MenuRepository.GetDeviceMenuByOrgByLanguage(organizationID, appLanguage.(uint))
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
	children, _ := receiver.ChildRepository.GetByParentID(userID)
	students, _ := receiver.StudentAppRepo.GetByUserIDApproved(userID)
	teachers, _ := receiver.TeacherRepository.GetByUserIDApproved(userID)

	// Child btn
	if children != nil {
		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleChild))
		var childOrg = ""
		if err == nil || roleOrg != nil {
			childOrg = roleOrg.OrgProfile
		}
		for _, child := range children {
			childComponent := buildComponent(
				uuid.NewString(),
				"Child Profile",
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

	// student btn
	if students != nil {
		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleStudent))
		var studentOrg = ""
		if err == nil && roleOrg != nil {
			studentOrg = roleOrg.OrgProfile
		}
		for range students {
			studentComponent := buildComponent(
				uuid.NewString(),
				"Student Profile",
				"student_profile",
				"icon/accident_and_injury_report_1745206766342940327.png",
				"button_form",
				studentOrg,
			)

			componentMenus = append(componentMenus, response.ComponentCommonMenuByUser{
				Component: studentComponent,
			})
		}
	}

	// teacher btn
	if teachers != nil {
		roleOrg, err := receiver.RoleOrgSignUpRepository.GetByRoleName(string(value.RoleTeacher))
		var teacherOrg = ""
		if err == nil && roleOrg != nil {
			teacherOrg = roleOrg.OrgProfile
		}
		for range teachers {
			teacherComponent := buildComponent(
				uuid.NewString(),
				"Teacher Profile",
				"teacher_profile",
				"icon/accident_and_injury_report_1745206766342940327.png",
				"button_form",
				teacherOrg,
			)

			componentMenus = append(componentMenus, response.ComponentCommonMenuByUser{
				Component: teacherComponent,
			})
		}
	}

	// check teacher menu
	// if teacherComponent, _ := receiver.getProfileComponentByRole("Teacher", userID); teacherComponent != nil {
	// 	componentMenus = append(componentMenus, *teacherComponent)
	// }

	//checck student menu
	// if studentComponent, _ := receiver.getProfileComponentByRole("Student", userID); studentComponent != nil {
	// 	componentMenus = append(componentMenus, *studentComponent)
	// }

	//check staff menu
	// if teacherComponent, _ := receiver.getProfileComponentByRole("Staff", userID); teacherComponent != nil {
	// 	componentMenus = append(componentMenus, *teacherComponent)
	// }

	//check org menu
	// if teacherComponent, _ := receiver.getProfileComponentByRole("Sign up ORganise", userID); teacherComponent != nil {
	// 	componentMenus = append(componentMenus, *teacherComponent)
	// }

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

func (receiver *GetMenuUseCase) GetSectionMenu(context *gin.Context) ([]response.GetMenuSectionResponse, error) {
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

func (receiver *GetMenuUseCase) GetSectionMenu4WebAdmin(ctx *gin.Context) ([]response.GetMenuSectionResponse, error) {
	user, err := receiver.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	var roleNames []string
	if user.IsSuperAdmin() {
		// SuperAdmin: lấy child
		roleNames = []string{string(value.RoleChild), string(value.Parent)}

	} else {
		// Không phải SuperAdmin: chỉ lấy role Student và Teacher
		roleNames = []string{string(value.RoleStudent), string(value.RoleTeacher)}
	}

	roleIDs := make([]string, 0)
	for _, roleName := range roleNames {
		role, err := receiver.RoleOrgSignUpRepository.GetByRoleName(roleName)
		if err != nil {
			return nil, err
		}
		if role != nil {
			roleIDs = append(roleIDs, role.ID.String())
		}
	}

	// Lấy tất cả components theo RoleID (SectionID)
	var allComponents []components.Component
	for _, roleID := range roleIDs {
		comps, err := receiver.ComponentRepository.GetBySectionID(roleID)
		if err != nil {
			return nil, err
		}
		allComponents = append(allComponents, comps...)
	}

	// Neu khong la super admin, loc theo organization menu template
	if !user.IsSuperAdmin() {
		// Nếu không phải SuperAdmin → lấy orgIDs mà user đang quản lý
		orgIDs, err := user.GetManagedOrganizationIDs(receiver.StudentAppRepo.GetDB())
		if err != nil {
			return nil, err
		}
		orgMenuTemplates, err := receiver.OrganizationMenuTemplateRepository.GetOrganizationMenuTemplatesByOrgID(orgIDs[0])
		if err != nil {
			return nil, err
		}
		if len(orgMenuTemplates) > 0 {
			// Lọc components theo OrganizationMenuTemplate
			allComponents = lo.Filter(allComponents, func(c components.Component, _ int) bool {
				return lo.ContainsBy(orgMenuTemplates, func(template entity.OrganizationMenuTemplate) bool {
					return template.SectionID == c.SectionID && template.ComponentID == c.ID.String()
				})
			})
		} else {
			allComponents = []components.Component{}
		}
	}

	// Gom nhóm theo SectionID
	grouped := make(map[string][]components.Component)
	for _, c := range allComponents {
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

func (receiver *GetMenuUseCase) GetSectionMenu4App(context *gin.Context) ([]response.GetMenuSectionResponse, error) {
	userID := context.GetString("user_id")
	var result []response.GetMenuSectionResponse
	// Lay danh sach child, students, teachers, staffs by userId
	children, _ := receiver.ChildRepository.GetByParentID(userID)
	// students, _ := receiver.StudentAppRepo.GetByUserIDApproved(userID)
	teachers, _ := receiver.TeacherRepository.GetByUserIDApproved(userID)
	staffs, _ := receiver.StaffApplicationRepo.GetByUserIDApproved(userID)

	// Get menu
	for _, child := range children {
		// get 4 app
		childMenu, err := receiver.ChildMenuUseCase.GetByChildID(context, child.ID.String(), true)
		// get menu icon key
		img, _ := receiver.UserImageUsecase.GetImg4Ownewr(child.ID.String(), value.OwnerRoleChild)
		menuIconKey := ""
		if img != nil {
			menuIconKey = img.Key
		}
		if err == nil && len(childMenu.Components) > 0 {
			result = append(result, response.GetMenuSectionResponse{
				SectionName: child.ChildName,
				MenuIconKey: menuIconKey,
				Components:  childMenu.Components,
			})
		}
	}

	// for _, student := range students {
	// 	studentMenu, _ := receiver.StudentMenuUseCase.GetByStudentID(student.ID.String())
	// 	result = append(result, response.GetMenuSectionResponse{
	// 		SectionName: studentMenu.StudentName,
	// 		SectionID:   studentMenu.StudentID,
	// 		Components:  studentMenu.Components,
	// 	})
	// }

	for _, teacher := range teachers {
		// get 4 app
		teacherMenu, err := receiver.TeacherMenuUseCase.GetByTeacherID(context, teacher.ID.String(), true)
		// get menu icon key
		img, _ := receiver.UserImageUsecase.GetImg4Ownewr(teacher.ID.String(), value.OwnerRoleTeacher)
		menuIconKey := ""
		if img != nil {
			menuIconKey = img.Key
		}
		if err == nil && len(teacherMenu.Components) > 0 {
			result = append(result, response.GetMenuSectionResponse{
				SectionName: teacherMenu.TeacherName,
				SectionID:   teacherMenu.TeacherID,
				MenuIconKey: menuIconKey,
				Components:  teacherMenu.Components,
			})
		}
	}

	for _, staff := range staffs {
		// get 4 app
		staffMenu, err := receiver.StaffMenuUsecase.GetByStaffID(context, staff.ID.String(), true)
		// get menu icon key
		img, _ := receiver.UserImageUsecase.GetImg4Ownewr(staff.ID.String(), value.OwnerRoleStaff)
		menuIconKey := ""
		if img != nil {
			menuIconKey = img.Key
		}
		if err == nil && len(staffMenu.Components) > 0 {
			result = append(result, response.GetMenuSectionResponse{
				SectionName: staffMenu.StaffName,
				SectionID:   staffMenu.StaffID,
				MenuIconKey: menuIconKey,
				Components:  staffMenu.Components,
			})
		}
	}

	// Get Parent Menu
	parent, _ := receiver.ParentRepo.GetByUserID(context, userID)
	parentMenu, err := receiver.ParentMenuUsecase.GetByParentID(context, parent.ID.String(), userID)
	img, _ := receiver.UserImageUsecase.GetImg4Ownewr(parent.ID.String(), value.OwnerRoleParent)
	menuIconKey := ""
	if img != nil {
		menuIconKey = img.Key
	}
	if err == nil && parentMenu.ParentID != "" && len(parentMenu.Components) > 0 {
		result = append(result, response.GetMenuSectionResponse{
			SectionName: "Parent Menu",
			MenuIconKey: menuIconKey,
			Components:  parentMenu.Components,
		})
	}

	// get department menu
	departments, _ := receiver.DepartmentGateway.GetDepartmentsByUser(context)
	for _, department := range departments {
		departmentMenu, err := receiver.DepartmentMenuUseCase.GetDepartmentMenu4App(context, department.ID)
		if err == nil && len(departmentMenu.Components) > 0 {
			result = append(result, response.GetMenuSectionResponse{
				SectionName: department.Name,
				SectionID:   department.ID,
				MenuIconKey: department.Icon,
				Components:  departmentMenu.Components,
			})
		}
	}

	return result, nil
}

func (receiver *GetMenuUseCase) GetEmergencyMenu4WebAdmin(ctx *gin.Context) ([]response.GetMenus4Web, error) {
	user, err := receiver.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	var emMenus []struct {
		ComponentID uuid.UUID
		Order       int
	}

	if user.IsSuperAdmin() {
		menus, err := receiver.SuperAdminEmergencyMenuRepo.GetAll()
		if err != nil {
			return nil, err
		}
		for _, m := range menus {
			emMenus = append(emMenus, struct {
				ComponentID uuid.UUID
				Order       int
			}{
				ComponentID: m.ComponentID,
				Order:       m.Order,
			})
		}
	} else {
		if len(user.Organizations) == 0 {
			return nil, errors.New("user does not belong to any organization")
		}

		orgIDsManaged, err := user.GetManagedOrganizationIDs(receiver.UserEntityRepository.GetDB())
		if err != nil {
			return nil, err
		}
		if len(orgIDsManaged) == 0 {
			return nil, errors.New("user does not manage any organization")
		}

		menus, err := receiver.OrganizationEmergencyMenuRepo.GetByOrganizationID(orgIDsManaged[0])
		if err != nil {
			return nil, err
		}
		for _, m := range menus {
			emMenus = append(emMenus, struct {
				ComponentID uuid.UUID
				Order       int
			}{
				ComponentID: m.ComponentID,
				Order:       m.Order,
			})
		}
	}

	// === Tái sử dụng logic gom components theo lang ===

	// Tạo danh sách componentID
	componentIDs := make([]uuid.UUID, 0, len(emMenus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range emMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
	}

	// Lấy tất cả components theo IDs
	components, err := receiver.ComponentRepository.GetByIDs(componentIDs)
	if err != nil {
		return nil, err
	}

	// Gom components theo language_id
	menusByLang := make(map[uint][]response.ComponentResponse)
	langMap := make(map[uint]entity.LanguageSetting)

	for _, comp := range components {
		menu := response.ComponentResponse{
			ID:         comp.ID.String(),
			Name:       comp.Name,
			Type:       comp.Type.String(),
			Key:        comp.Key,
			Value:      string(comp.Value),
			Order:      componentOrderMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		menusByLang[comp.LanguageID] = append(menusByLang[comp.LanguageID], menu)

		// Nếu chưa có language cache thì lấy từ DB
		if _, ok := langMap[comp.LanguageID]; !ok {
			langSetting, err := receiver.LanguageSettingRepo.GetByID(comp.LanguageID)
			if err != nil {
				return nil, err
			}
			if langSetting != nil {
				langMap[comp.LanguageID] = *langSetting
			}
		}
	}

	// Build []GetMenus4Web
	getMenus := make([]response.GetMenus4Web, 0, len(menusByLang))
	for langID, comps := range menusByLang {
		getMenus = append(getMenus, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	return getMenus, nil
}

func (receiver *GetMenuUseCase) GetEmergencyMenu4App(ctx *gin.Context, organizationID string) ([]response.GetMenuSectionResponse, error) {
	appLanguage, _ := ctx.Get("app_language")
	menus, err := receiver.OrganizationEmergencyMenuRepo.GetByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	var result []response.GetMenuSectionResponse
	for _, menu := range menus {
		comp, _ := receiver.ComponentRepository.GetByIDAndLanguage(menu.ComponentID.String(), appLanguage.(uint))

		components := []response.ComponentResponse{}
		if comp != nil {
			components = append(components, response.ComponentResponse{
				ID:    menu.ComponentID.String(),
				Name:  comp.Name,
				Key:   comp.Key,
				Type:  comp.Type.String(),
				Value: comp.Value.String(),
				Order: menu.Order,
			})
		}
		result = append(result, response.GetMenuSectionResponse{
			SectionName: "Emergency Menu",
			MenuIconKey: "",
			Components:  components,
		})
	}

	return result, nil
}
