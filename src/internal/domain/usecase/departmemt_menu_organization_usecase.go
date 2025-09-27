package usecase

import (
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepartmentMenuOrganizationUseCase struct {
	DepartmentMenuOrganizationRepository *repository.DepartmentMenuOrganizationRepository
	ComponentRepo                        *repository.ComponentRepository
	DeviceRepository                     *repository.DeviceRepository
	OrganizationRepository               *repository.OrganizationRepository
	DepartmentGateway                    gateway.DepartmentGateway
	LanguageSettingRepo                  *repository.LanguageSettingRepository
}

func (uc *DepartmentMenuOrganizationUseCase) GetDepartmentMenuOrg4GW(ctx *gin.Context, departmentID, orgID string) (response.GetDepartmentMenuResponse, error) {
	// 1. Lấy danh sách menu của giáo viên trong org
	departmentMenusOrg, err := uc.DepartmentMenuOrganizationRepository.GetAllByDepartmentAndOrg(ctx, departmentID, orgID)
	if err != nil {
		return response.GetDepartmentMenuResponse{}, err
	}

	if len(departmentMenusOrg) == 0 {
		return response.GetDepartmentMenuResponse{}, nil
	}

	// 2. Chuẩn bị danh sách componentID + mapping order
	componentIDs := make([]uuid.UUID, 0, len(departmentMenusOrg))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range departmentMenusOrg {
		compID := cm.ComponentID
		componentIDs = append(componentIDs, compID)
		componentOrderMap[compID] = cm.Order
	}

	// 3. Lấy components theo ID
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
	if err != nil {
		return response.GetDepartmentMenuResponse{}, err
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
			langSetting, err := uc.LanguageSettingRepo.GetByID(comp.LanguageID)
			if err != nil {
				return response.GetDepartmentMenuResponse{}, err
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

	return response.GetDepartmentMenuResponse{
		Components: getMenus,
	}, nil
}

func (uc *DepartmentMenuOrganizationUseCase) GetDepartmentMenuOrg4App(ctx *gin.Context, req request.GetDepartmentMenuOrganizationRequest) ([]*response.GetDepartmentMenuOrganizationResponse, error) {

	// 1. Kiểm tra device có trong org hay không
	isExist, err := uc.DeviceRepository.CheckDeviceExistInOrganization(req.DeviceID, req.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, nil
	}

	// 2. Lấy danh sách department từ gateway
	departments, _ := uc.DepartmentGateway.GetDepartmentsByUser(ctx)

	deptOrgMenus := make([]*response.GetDepartmentMenuOrganizationResponse, 0, len(departments))

	for _, department := range departments {
		// 2.1 Lấy department menu trong org
		departmentMenusOrg, err := uc.DepartmentMenuOrganizationRepository.
			GetAllByDepartmentAndOrg(ctx, department.ID, req.OrganizationID)
		if err != nil {
			continue // skip nếu lỗi
		}

		// 2.2 Chuẩn bị componentIDs + mapping order
		componentIDs := make([]uuid.UUID, 0, len(departmentMenusOrg))
		componentOrderMap := make(map[uuid.UUID]int)

		for _, cm := range departmentMenusOrg {
			compID := cm.ComponentID
			componentIDs = append(componentIDs, compID)
			componentOrderMap[compID] = cm.Order
		}

		// 2.3 Lấy components theo ID
		appLanguage, _ := ctx.Get("app_language")
		components, err := uc.ComponentRepo.GetByIDsAndLanguage(componentIDs, appLanguage.(uint))
		if err != nil {
			return nil, err
		}

		// 2.4 Build component responses
		menus := make([]response.ComponentResponse, 0, len(components))
		for _, comp := range components {
			menus = append(menus, response.ComponentResponse{
				ID:    comp.ID.String(),
				Name:  comp.Name,
				Type:  comp.Type.String(),
				Key:   comp.Key,
				Value: helper.BuildSectionValueMenu(string(comp.Value), comp),
				Order: componentOrderMap[comp.ID],
			})
		}

		// 2.5 Build department menu response
		orgInfo, _ := uc.OrganizationRepository.GetByID(req.OrganizationID)

		departmentOrgMenus := &response.GetDepartmentMenuOrganizationResponse{
			Section:     department.Name + " Menu At " + orgInfo.OrganizationName,
			MenuIconKey: department.Icon,
			Components:  menus,
		}

		if len(menus) > 0 {
			deptOrgMenus = append(deptOrgMenus, departmentOrgMenus)
		}
	}

	return deptOrgMenus, nil
}
