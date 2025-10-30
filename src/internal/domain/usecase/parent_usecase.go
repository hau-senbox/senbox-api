package usecase

import (
	"context"
	"errors"
	"sen-global-api/internal/cache/caching"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParentUseCase struct {
	DBConn                   *gorm.DB
	UserRepo                 *repository.UserEntityRepository
	ParentMenuRepo           *repository.ParentMenuRepository
	ComponentRepo            *repository.ComponentRepository
	LanguagesConfigUsecase   *LanguagesConfigUsecase
	UserImagesUsecase        *UserImagesUsecase
	LanguageSettingRepo      *repository.LanguageSettingRepository
	ParentRepo               *repository.ParentRepository
	ParentChildsRepo         *repository.ParentChildsRepository
	StudentRepo              *repository.StudentApplicationRepository
	ProfileGateway           gateway.ProfileGateway
	UserBlockSettingUsecase  *UserBlockSettingUsecase
	GenerateOwnerCodeUseCase GenerateOwnerCodeUseCase
	CachingService           caching.CachingService
}

func NewParentUseCase(
	userRepo *repository.UserEntityRepository,
	parentMenuRepo *repository.ParentMenuRepository,
	componentRepo *repository.ComponentRepository,
	languagesConfigUsecase *LanguagesConfigUsecase,
	userImagesUsecase *UserImagesUsecase,
	languageSettingRepo *repository.LanguageSettingRepository,
	parentRepo *repository.ParentRepository,
	parentChildsRepo *repository.ParentChildsRepository,
	studentRepo *repository.StudentApplicationRepository,
	profileGateway gateway.ProfileGateway,
	userBlockSettingUsecase *UserBlockSettingUsecase,
	generateOwnerCodeUseCase GenerateOwnerCodeUseCase,
	cachingService caching.CachingService,
) *ParentUseCase {
	return &ParentUseCase{
		UserRepo:                 userRepo,
		ParentMenuRepo:           parentMenuRepo,
		ComponentRepo:            componentRepo,
		LanguagesConfigUsecase:   languagesConfigUsecase,
		UserImagesUsecase:        userImagesUsecase,
		LanguageSettingRepo:      languageSettingRepo,
		ParentRepo:               parentRepo,
		ParentChildsRepo:         parentChildsRepo,
		StudentRepo:              studentRepo,
		ProfileGateway:           profileGateway,
		UserBlockSettingUsecase:  userBlockSettingUsecase,
		GenerateOwnerCodeUseCase: generateOwnerCodeUseCase,
		CachingService:           cachingService,
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

func (uc *ParentUseCase) GetAllParents4Search(ctx *gin.Context) ([]response.ParentResponse, error) {
	userID := ctx.GetString("user_id")
	user, err := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	var parents []entity.SParent

	if user.IsSuperAdmin() {
		// SuperAdmin → lấy tất cả
		parents, err = uc.ParentRepo.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		for i := range parents {
			userParent, _ := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: parents[i].UserID})
			if userParent != nil {
				parents[i].ParentName = userParent.Nickname
			}
		}
	} else {
		// Non-SuperAdmin → lấy theo org
		orgAdminIds, _ := user.GetManagedOrganizationIDs(DBConn)
		if len(orgAdminIds) == 0 {
			return nil, errors.New("user does not manage any organization")
		}

		// Lấy học sinh theo org
		students, err := uc.StudentRepo.GetByOrganizationID(orgAdminIds[0])
		if err != nil {
			return nil, err
		}

		var parentList []entity.SParent
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
				parentList = append(parentList, *parent)
			}
		}

		// unique
		uniqueMap := make(map[uuid.UUID]entity.SParent)
		for _, p := range parentList {
			uniqueMap[p.ID] = p
		}
		parents = make([]entity.SParent, 0, len(uniqueMap))
		for _, v := range uniqueMap {
			parents = append(parents, v)
		}
	}

	return uc.mapParentEntitiesToResponse(ctx, parents), nil
}

func (uc *ParentUseCase) mapParentEntitiesToResponse(ctx *gin.Context, parents []entity.SParent) []response.ParentResponse {
	res := make([]response.ParentResponse, 0, len(parents))
	for _, p := range parents {
		userEntity, _ := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{
			ID: p.UserID,
		})
		isDeactive, _ := uc.UserBlockSettingUsecase.GetDeactive4User(p.ID.String())
		avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(p.ID.String(), value.OwnerRoleParent)
		code, _ := uc.ProfileGateway.GetParentCode(ctx, p.ID.String())
		res = append(res, response.ParentResponse{
			ParentID:         p.ID.String(),
			ParentName:       userEntity.Nickname,
			IsDeactive:       isDeactive,
			CreatedIndex:     p.CreatedIndex,
			UserCreatedIndex: userEntity.CreatedIndex,
			Avatar:           avatar,
			Code:             code,
			LanguageKeys:     []string{"vietnamese", "english"},
		})
	}

	return res
}

func (uc *ParentUseCase) GenerateParentCode(ctx *gin.Context) {
	// get all parents
	parents, err := uc.ParentRepo.GetAll(ctx)
	if err != nil {
		return
	}

	for _, pr := range parents {
		// call profile gateway to generate parent code
		_, _ = uc.ProfileGateway.GenerateParentCode(ctx, pr.ID.String(), pr.CreatedIndex)
	}
}

func (uc *ParentUseCase) GetParentByUser4Gw(ctx *gin.Context, userID string) (*response.GetParent4Gateway, error) {
	parent, err := uc.ParentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent not found")
	}

	// get avts
	avatar, _ := uc.UserImagesUsecase.GetAvtIsMain4Owner(parent.ID.String(), value.OwnerRoleParent)

	// get user
	user, _ := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: parent.UserID})
	code, _ := uc.ProfileGateway.GetParentCode(ctx, parent.ID.String())

	res := &response.GetParent4Gateway{
		ParentID:   parent.ID.String(),
		ParentName: user.Nickname,
		Avatar:     avatar,
		Code:       code,
	}

	_ = uc.CachingService.SetParentByUserCacheKey(ctx, userID, res)
	return res, nil
}
