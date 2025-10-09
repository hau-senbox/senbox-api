package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ParentUseCase struct {
	UserRepo               *repository.UserEntityRepository
	ParentMenuRepo         *repository.ParentMenuRepository
	ComponentRepo          *repository.ComponentRepository
	LanguagesConfigUsecase *LanguagesConfigUsecase
	UserImagesUsecase      *UserImagesUsecase
	LanguageSettingRepo    *repository.LanguageSettingRepository
	ParentRepo             *repository.ParentRepository
	ParentChildsRepo       *repository.ParentChildsRepository
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

func (uc *ParentUseCase) CreateParent(ctx *gin.Context, childID string, userID string) error {

	parent := &entity.SParent{
		ID:     uuid.New(),
		UserID: userID,
	}

	// tao parent
	err := uc.ParentRepo.Create(parent)
	if err != nil {
		return err
	}

	// tao parent childs
	parentChilds := &entity.SParentChilds{
		ParentID: parent.ID.String(),
		ChildID:  childID,
	}

	err = uc.ParentChildsRepo.Create(ctx, parentChilds)
	if err != nil {
		return err
	}
	return nil
}
