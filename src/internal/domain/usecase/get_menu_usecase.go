package usecase

import (
	"encoding/json"
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"strings"

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

func (receiver *GetMenuUseCase) GetStudentMenu(studentID string) (response.GetStudentMenuResponse, error) {
	studentMenu, err := receiver.StudentMenuUseCase.GetByStudentID(studentID)
	if err != nil {
		return response.GetStudentMenuResponse{}, fmt.Errorf("failed to get student menu: %w", err)
	}

	return studentMenu, nil
}

func (receiver *GetMenuUseCase) GetTeacherMenu(userID string) (response.GetTeacherMenuResponse, error) {

	teacher, _ := receiver.TeacherRepository.GetByUserID(userID)

	teacherMenu, err := receiver.TeacherMenuUseCase.GetByTeacherID(teacher.ID.String())
	if err != nil {
		return response.GetTeacherMenuResponse{}, fmt.Errorf("failed to get teacher menu: %w", err)
	}

	return teacherMenu, nil
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

func (receiver *GetMenuUseCase) GetSectionMenu4WebAdmin(ctx *gin.Context) ([]response.GetMenuSectionResponse, error) {
	user, err := receiver.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	var roleNames []string
	if user.IsSuperAdmin() {
		// SuperAdmin: lấy child
		roleNames = []string{string(value.RoleChild)}

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
	// Lay danh sach child, students, teacher by userId
	children, _ := receiver.ChildRepository.GetByParentID(userID)
	// students, _ := receiver.StudentAppRepo.GetByUserIDApproved(userID)
	teachers, _ := receiver.TeacherRepository.GetByUserIDApproved(userID)
	staffs, _ := receiver.StaffApplicationRepo.GetByUserIDApproved(userID)
	// neu co child lay menu cua child
	for _, child := range children {
		childMenu, _ := receiver.ChildMenuUseCase.GetByChildID(child.ID.String())
		result = append(result, response.GetMenuSectionResponse{
			SectionName: child.ChildName,
			Components:  childMenu.Components,
		})
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
		teacherMenu, _ := receiver.TeacherMenuUseCase.GetByTeacherID(teacher.ID.String())
		result = append(result, response.GetMenuSectionResponse{
			SectionName: teacherMenu.TeacherName,
			SectionID:   teacherMenu.TeacherID,
			Components:  teacherMenu.Components,
		})
	}

	for _, staff := range staffs {
		staffMenu, _ := receiver.StaffMenuUsecase.GetByStaffID(staff.ID.String())
		result = append(result, response.GetMenuSectionResponse{
			SectionName: staffMenu.StaffName,
			SectionID:   staffMenu.StaffID,
			Components:  staffMenu.Components,
		})
	}

	return result, nil
}
