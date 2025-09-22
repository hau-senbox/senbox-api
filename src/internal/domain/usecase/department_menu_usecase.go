package usecase

import (
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepartmentMenuUseCase struct {
	DepartmentMenuRepository *repository.DepartmentMenuRepository
	ComponentRepository      *repository.ComponentRepository
	LanguageSettingRepo      *repository.LanguageSettingRepository
}

func (uc *DepartmentMenuUseCase) GetDepartmentMenu4GW(departmentID string) (response.GetDepartmentMenuResponse, error) {

	// B1: Lấy các bản ghi teacher_menu
	departmentMenus, err := uc.DepartmentMenuRepository.GetByDepartmentID(departmentID)
	if err != nil {
		return response.GetDepartmentMenuResponse{}, err
	}

	// B2: Lấy componentID từ teacher_menu
	componentIDs := make([]uuid.UUID, 0, len(departmentMenus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, tm := range departmentMenus {
		componentIDs = append(componentIDs, tm.ComponentID)
		componentOrderMap[tm.ComponentID] = tm.Order
	}

	// B3: Lấy danh sách component tương ứng
	components, err := uc.ComponentRepository.GetByIDs(componentIDs)
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

func (uc *DepartmentMenuUseCase) GetDepartmentMenu4App(ctx *gin.Context, departmentID string) (response.GetDepartmentMenuResponse4App, error) {

	departmentMenus, err := uc.DepartmentMenuRepository.GetByDepartmentID(departmentID)
	if err != nil {
		return response.GetDepartmentMenuResponse4App{}, err
	}

	// B1: Lấy tất cả ComponentID từ departmentMenus
	componentIDs := make([]uuid.UUID, 0, len(departmentMenus))
	componentOrderMap := make(map[uuid.UUID]int) // lưu order theo ComponentID

	for _, cm := range departmentMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
	}

	// B2: Lấy danh sách Component theo IDs
	appLanguage, _ := ctx.Get("app_language")
	components, err := uc.ComponentRepository.GetByIDsAndLanguage(componentIDs, appLanguage.(uint))
	if err != nil {
		return response.GetDepartmentMenuResponse4App{}, err
	}

	// B3: Map sang ComponentChildResponse
	componentResponses := make([]response.ComponentResponse, 0, len(components))
	for _, comp := range components {
		visible, _ := helper.GetVisibleToValueComponent(comp.Value.String())
		if !visible {
			continue
		}
		componentResponses = append(componentResponses, response.ComponentResponse{
			ID:    comp.ID.String(),
			Name:  comp.Name,
			Type:  comp.Type.String(),
			Key:   comp.Key,
			Value: helper.BuildSectionValueMenu(string(comp.Value), comp),
			Order: componentOrderMap[comp.ID],
		})
	}

	return response.GetDepartmentMenuResponse4App{
		Components: componentResponses,
	}, nil
}
