package usecase

import (
	"context"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParentUseCase struct {
	DBConn                 *gorm.DB
	UserRepo               *repository.UserEntityRepository
	ParentMenuRepo         *repository.ParentMenuRepository
	ComponentRepo          *repository.ComponentRepository
	LanguagesConfigUsecase *LanguagesConfigUsecase
	UserImagesUsecase      *UserImagesUsecase
	LanguageSettingRepo    *repository.LanguageSettingRepository
	ParentRepo             *repository.ParentRepository
	ParentChildsRepo       *repository.ParentChildsRepository
	StudentRepo            *repository.StudentApplicationRepository
}

func NewParentUseCase(
	userRepo *repository.UserEntityRepository,
	parentMenuRepo *repository.ParentMenuRepository,
	componentRepo *repository.ComponentRepository,
	languagesConfigUsecase *LanguagesConfigUsecase,
	userImagesUsecase *UserImagesUsecase,
	languageSettingRepo *repository.LanguageSettingRepository,
	parentRepo *repository.ParentRepository,
) *ParentUseCase {
	return &ParentUseCase{
		UserRepo:               userRepo,
		ParentMenuRepo:         parentMenuRepo,
		ComponentRepo:          componentRepo,
		LanguagesConfigUsecase: languagesConfigUsecase,
		UserImagesUsecase:      userImagesUsecase,
		LanguageSettingRepo:    languageSettingRepo,
		ParentRepo:             parentRepo,
	}
}

func (uc *ParentUseCase) GetParentByID(parentID string) (*response.ParentResponseBase, error) {
	parent, err := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: parentID})
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent not found")
	}

	// Lấy danh sách ParentMenus
	parentMenus, err := uc.ParentMenuRepo.GetByParentID(parentID)
	if err != nil {
		return nil, err
	}

	// Tạo danh sách componentID để lấy Component
	componentIDs := make([]uuid.UUID, 0, len(parentMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, cm := range parentMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
		componentIsShowMap[cm.ComponentID] = cm.IsShow
	}

	// Lấy tất cả components theo danh sách ID
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
			IsShow:     componentIsShowMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		menusByLang[comp.LanguageID] = append(menusByLang[comp.LanguageID], menu)

		// nếu chưa có language ILanguageIDtrong cache -> query DB
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

	// get languages config
	languageConfig, _ := uc.LanguagesConfigUsecase.GetLanguagesConfigByOwnerNoCtx(parentID, value.OwnerRoleLangParent)

	// get avts
	avatars, _ := uc.UserImagesUsecase.GetAvt4Owner(parentID, value.OwnerRoleParent)

	return &response.ParentResponseBase{
		ParentID:       parentID,
		ParentName:     parent.Nickname,
		Avatar:         "",
		AvatarURL:      "",
		Menus:          getMenus,
		CustomID:       parent.CustomID,
		LanguageConfig: languageConfig,
		Avatars:        avatars,
		CreatedIndex:   parent.CreatedIndex,
	}, nil
}

func (uc *ParentUseCase) GetParentByID4Web(ctx context.Context, parentID string) (*response.ParentResponseBase, error) {
	parent, err := uc.ParentRepo.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent not found")
	}

	// Lấy danh sách ParentMenus
	parentMenus, err := uc.ParentMenuRepo.GetByParentID(parentID)
	if err != nil {
		return nil, err
	}

	// Tạo danh sách componentID để lấy Component
	componentIDs := make([]uuid.UUID, 0, len(parentMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, cm := range parentMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
		componentIsShowMap[cm.ComponentID] = cm.IsShow
	}

	// Lấy tất cả components theo danh sách ID
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
			IsShow:     componentIsShowMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		menusByLang[comp.LanguageID] = append(menusByLang[comp.LanguageID], menu)

		// nếu chưa có language ILanguageIDtrong cache -> query DB
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

	// get languages config
	languageConfig, _ := uc.LanguagesConfigUsecase.GetLanguagesConfigByOwnerNoCtx(parentID, value.OwnerRoleLangParent)

	// get avts
	avatars, _ := uc.UserImagesUsecase.GetAvt4Owner(parentID, value.OwnerRoleParent)

	// get user
	user, _ := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: parent.UserID})

	return &response.ParentResponseBase{
		ParentID:       parentID,
		ParentName:     user.Nickname,
		Avatar:         "",
		AvatarURL:      "",
		Menus:          getMenus,
		LanguageConfig: languageConfig,
		Avatars:        avatars,
		CustomID:       "",
		CreatedIndex:   parent.CreatedIndex,
	}, nil
}

func (uc *ParentUseCase) GetParentByUser4Web(ctx context.Context, userID string) (*response.ParentResponseBase, error) {
	parent, err := uc.ParentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent not found")
	}

	// Lấy danh sách ParentMenus
	parentMenus, err := uc.ParentMenuRepo.GetByParentID(parent.ID.String())
	if err != nil {
		return nil, err
	}

	// Tạo danh sách componentID để lấy Component
	componentIDs := make([]uuid.UUID, 0, len(parentMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, cm := range parentMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
		componentIsShowMap[cm.ComponentID] = cm.IsShow
	}

	// Lấy tất cả components theo danh sách ID
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
			IsShow:     componentIsShowMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		menusByLang[comp.LanguageID] = append(menusByLang[comp.LanguageID], menu)

		// nếu chưa có language ILanguageIDtrong cache -> query DB
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

	// get languages config
	languageConfig, _ := uc.LanguagesConfigUsecase.GetLanguagesConfigByOwnerNoCtx(parent.ID.String(), value.OwnerRoleLangParent)

	// get avts
	avatars, _ := uc.UserImagesUsecase.GetAvt4Owner(parent.ID.String(), value.OwnerRoleParent)

	// get user
	user, _ := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: parent.UserID})

	return &response.ParentResponseBase{
		ParentID:       parent.ID.String(),
		ParentName:     user.Nickname,
		Avatar:         "",
		AvatarURL:      "",
		Menus:          getMenus,
		LanguageConfig: languageConfig,
		Avatars:        avatars,
		CreatedIndex:   parent.CreatedIndex,
	}, nil
}

func (uc *ParentUseCase) GetParentByUser(ctx context.Context, userID string) (*entity.SParent, error) {
	return uc.ParentRepo.GetByUserID(ctx, userID)
}

func (uc *ParentUseCase) GetAllParents4Search(ctx *gin.Context) ([]entity.SParent, error) {
	userID := ctx.GetString("user_id")
	user, err := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	if user.IsSuperAdmin() {
		parents, err := uc.ParentRepo.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		for i := range parents {
			userParent, _ := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: parents[i].UserID})
			if userParent != nil {
				parents[i].ParentName = userParent.Nickname
			}
		}

		return parents, nil
	}

	orgAdminIds, _ := user.GetManagedOrganizationIDs(DBConn)
	if len(orgAdminIds) == 0 {
		return nil, errors.New("user does not manage any organization")
	}

	// get all student by org
	students, err := uc.StudentRepo.GetByOrganizationID(orgAdminIds[0])
	if err != nil {
		return nil, err
	}

	var parents = make([]entity.SParent, 0)

	// get parent by childid in student
	for _, student := range students {
		parent, err := uc.ParentRepo.GetByUserID(ctx, student.UserID.String())
		if err != nil {
			return nil, err
		}
		if parent != nil {
			userParent, _ := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: parent.UserID})
			if userParent != nil {
				parent.ParentName = userParent.Nickname
			}
			parents = append(parents, *parent)
		}
	}

	uniqueMap := make(map[uuid.UUID]entity.SParent)
	for _, p := range parents {
		uniqueMap[p.ID] = p
	}

	uniqueParents := make([]entity.SParent, 0, len(uniqueMap))
	for _, v := range uniqueMap {
		uniqueParents = append(uniqueParents, v)
	}

	return uniqueParents, nil
}
