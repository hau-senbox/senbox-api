package usecase

import (
	"context"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentMenuOrganizationUseCase struct {
	StudentMenuOrganizationRepository *repository.StudentMenuOrganizationRepository
	ComponentRepo                     *repository.ComponentRepository
	StudentRepo                       *repository.StudentApplicationRepository
	DeviceRepository                  *repository.DeviceRepository
	UserImagesUsecase                 *UserImagesUsecase
	OrganizationRepository            *repository.OrganizationRepository
	LanguageSettingRepo               *repository.LanguageSettingRepository
}

func (uc *StudentMenuOrganizationUseCase) GetStudentMenuOrg4Admin(ctx context.Context, studentID, orgID string) ([]response.GetMenus4Web, error) {
	// 1. Lấy danh sách menu của giáo viên trong org
	studentMenus, err := uc.StudentMenuOrganizationRepository.GetAllByStudentAndOrg(ctx, studentID, orgID)
	if err != nil {
		return nil, err
	}

	if len(studentMenus) == 0 {
		return []response.GetMenus4Web{}, nil
	}

	// 2. Chuẩn bị danh sách componentID + mapping order
	componentIDs := make([]uuid.UUID, 0, len(studentMenus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range studentMenus {
		compID, err := uuid.Parse(cm.ComponentID)
		if err != nil {
			continue // skip nếu ComponentID không hợp lệ
		}
		componentIDs = append(componentIDs, compID)
		componentOrderMap[compID] = cm.Order
	}

	// 3. Lấy components theo ID
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
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
			langSetting, err := uc.LanguageSettingRepo.GetByID(comp.LanguageID)
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

func (uc *StudentMenuOrganizationUseCase) GetStudentMenuOrg4App(ctx *gin.Context, req request.GetStudentOrganizationMenuRequest) (*response.GetStudentOrganizationMenuResponse, error) {
	// kiem tra device dang co nam trong org hay khong
	isExist, err := uc.DeviceRepository.CheckDeviceExistInOrganization(req.DeviceID, req.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, nil
	}
	// get teacher by user ID
	student, err := uc.StudentRepo.GetByID(uuid.MustParse(req.StudentID))
	if err != nil {
		return nil, err
	}

	// Lấy danh sách menu của giáo viên trong org
	studentMenus, err := uc.StudentMenuOrganizationRepository.GetAllByStudentAndOrg(ctx, student.ID.String(), req.OrganizationID)
	if err != nil {
		return nil, err
	}

	if len(studentMenus) == 0 {
		return nil, nil
	}

	// 2. Chuẩn bị danh sách componentID + mapping order
	componentIDs := make([]uuid.UUID, 0, len(studentMenus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range studentMenus {
		compID, err := uuid.Parse(cm.ComponentID)
		if err != nil {
			continue // skip nếu ComponentID không hợp lệ
		}
		componentIDs = append(componentIDs, compID)
		componentOrderMap[compID] = cm.Order
	}

	// 3. Lấy components theo ID
	appLanguage, _ := ctx.Get("app_language")
	components, err := uc.ComponentRepo.GetByIDsAndLanguage(componentIDs, appLanguage.(uint))
	if err != nil {
		return nil, err
	}

	// 4. Build response
	menus := make([]response.ComponentResponse, 0, len(components))
	for _, comp := range components {
		menu := response.ComponentResponse{
			ID:    comp.ID.String(),
			Name:  comp.Name,
			Type:  comp.Type.String(),
			Key:   comp.Key,
			Value: helper.BuildSectionValueMenu(string(comp.Value), comp),
			Order: componentOrderMap[comp.ID],
		}
		menus = append(menus, menu)
	}

	// get menu icon key
	img, _ := uc.UserImagesUsecase.GetImg4Ownewr(student.ID.String(), value.OwnerRoleStudent)

	menuIconKey := ""
	if img != nil {
		menuIconKey = img.Key
	}

	// get organization
	orgInfo, _ := uc.OrganizationRepository.GetByID(req.OrganizationID)
	studentOrgMenus := &response.GetStudentOrganizationMenuResponse{
		Section:     "Student Menu At " + orgInfo.OrganizationName,
		MenuIconKey: menuIconKey,
		Components:  menus,
	}

	return studentOrgMenus, nil

}
