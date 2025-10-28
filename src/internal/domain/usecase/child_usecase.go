package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChildUseCase struct {
	dbConn                   *gorm.DB
	childRepo                *repository.ChildRepository
	userRepo                 *repository.UserEntityRepository
	componentRepo            *repository.ComponentRepository
	childMenuRepo            *repository.ChildMenuRepository
	roleOrgRepo              *repository.RoleOrgSignUpRepository
	getUserEntityUseCase     *GetUserEntityUseCase
	userImagesUsecase        *UserImagesUsecase
	languagesConfigUsecase   *LanguagesConfigUsecase
	languageSettingRepo      *repository.LanguageSettingRepository
	parentRepo               *repository.ParentRepository
	parentChildsRepo         *repository.ParentChildsRepository
	profileGateway           gateway.ProfileGateway
	generateOwnerCodeUseCase GenerateOwnerCodeUseCase
}

func NewChildUseCase(
	dbConn *gorm.DB,
	childRepo *repository.ChildRepository,
	userRepo *repository.UserEntityRepository,
	componentRepo *repository.ComponentRepository,
	childMenuRepo *repository.ChildMenuRepository,
	roleOrgRepo *repository.RoleOrgSignUpRepository,
	getUserEntityUseCase *GetUserEntityUseCase,
	languagesConfigUsecase *LanguagesConfigUsecase,
	userImagesUsecase *UserImagesUsecase,
	languageSettingRepo *repository.LanguageSettingRepository,
	parentRepo *repository.ParentRepository,
	parentChildsRepo *repository.ParentChildsRepository,
	profileGateway gateway.ProfileGateway,
	gengenerateOwnerCodeUseCase GenerateOwnerCodeUseCase,
) *ChildUseCase {
	return &ChildUseCase{
		dbConn:                   dbConn,
		childRepo:                childRepo,
		userRepo:                 userRepo,
		componentRepo:            componentRepo,
		childMenuRepo:            childMenuRepo,
		roleOrgRepo:              roleOrgRepo,
		getUserEntityUseCase:     getUserEntityUseCase,
		userImagesUsecase:        userImagesUsecase,
		languagesConfigUsecase:   languagesConfigUsecase,
		languageSettingRepo:      languageSettingRepo,
		parentRepo:               parentRepo,
		parentChildsRepo:         parentChildsRepo,
		profileGateway:           profileGateway,
		generateOwnerCodeUseCase: gengenerateOwnerCodeUseCase,
	}
}

func (uc *ChildUseCase) CreateChild(req request.CreateChildRequest, ctx *gin.Context) error {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		return errors.New("unauthorized: user_id not found in context")
	}

	var userID uuid.UUID
	switch v := userIDRaw.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return errors.New("invalid user_id format")
		}
		userID = parsed
	default:
		return errors.New("invalid user_id type in context")
	}

	childID := uuid.New()
	child := &entity.SChild{
		ID:        childID,
		ChildName: req.ChildName,
		Age:       req.Age,
		ParentID:  userID,
	}

	// Bắt đầu Transaction
	tx := uc.dbConn.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// ---- Create Child ----
	if err := uc.childRepo.WithTx(tx).Create(child); err != nil {
		tx.Rollback()
		return fmt.Errorf("create child failed: %w", err)
	}

	// ---- Tạo menu cho child ----
	childRoleOrg, _ := uc.roleOrgRepo.WithTx(tx).GetByRoleName(string(value.RoleChild))
	if childRoleOrg != nil {
		comps, _ := uc.componentRepo.WithTx(tx).GetBySectionID(childRoleOrg.ID.String())

		for index, component := range comps {
			visible, _ := helper.GetVisibleToValueComponent(string(component.Value))

			newComponent := &components.Component{
				ID:        uuid.New(),
				Name:      component.Name,
				Type:      component.Type,
				Key:       component.Key,
				SectionID: component.SectionID,
				Value:     component.Value,
			}

			if err := uc.componentRepo.WithTx(tx).Create(newComponent); err != nil {
				tx.Rollback()
				return fmt.Errorf("create component failed: %w", err)
			}

			childMenu := &entity.ChildMenu{
				ID:          uuid.New(),
				ChildID:     childID,
				ComponentID: newComponent.ID,
				Order:       index,
				IsShow:      true,
				Visible:     visible,
			}

			if err := uc.childMenuRepo.WithTx(tx).Create(childMenu); err != nil {
				tx.Rollback()
				return fmt.Errorf("create child menu failed: %w", err)
			}
		}
	}

	// ---- Tạo Parent - Child mapping ----
	if err := uc.parentRepo.WithTx(tx).Create(ctx, &entity.SParent{ID: uuid.New(), UserID: userID.String()}); err != nil {
		tx.Rollback()
		return fmt.Errorf("create parent-child failed: %w", err)
	}

	if err := uc.parentChildsRepo.WithTx(tx).Create(ctx, &entity.SParentChilds{ParentID: userID.String(), ChildID: childID.String()}); err != nil {
		tx.Rollback()
		return fmt.Errorf("create parent-child failed: %w", err)
	}

	// generate child code
	uc.generateOwnerCodeUseCase.GenerateChildCode(ctx, childID.String())

	// Commit nếu tất cả OK
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (uc *ChildUseCase) UpdateChild(req request.UpdateChildRequest, ctx *gin.Context) error {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		return errors.New("unauthorized: user_id not found in context")
	}

	var userID uuid.UUID
	switch v := userIDRaw.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return errors.New("invalid user_id format")
		}
		userID = parsed
	default:
		return errors.New("invalid user_id type in context")
	}

	childID, err := uuid.Parse(req.ID)
	if err != nil {
		return errors.New("invalid child_id format")
	}

	child := &entity.SChild{
		ID:        childID,
		ChildName: req.ChildName,
		Age:       req.Age,
		ParentID:  userID,
	}

	return uc.childRepo.Update(child)
}

func (uc *ChildUseCase) GetByID(childID string) (*entity.SChild, error) {
	return uc.childRepo.GetByID(childID)
}

func (uc *ChildUseCase) GetAll() ([]entity.SChild, error) {
	return uc.childRepo.GetAll()
}

func (uc *ChildUseCase) GetByID4WebAdmin(childID string) (*response.ChildResponse, error) {
	// Lấy thông tin child
	child, err := uc.childRepo.GetByID(childID)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return nil, errors.New("child not found")
	}

	// Lấy danh sách ChildMenu
	childMenus, err := uc.childMenuRepo.GetByChildID(childID)
	if err != nil {
		return nil, err
	}

	// Tạo danh sách componentID để lấy Component
	componentIDs := make([]uuid.UUID, 0, len(childMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, cm := range childMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
		componentIsShowMap[cm.ComponentID] = cm.IsShow
	}

	// Lấy tất cả components theo danh sách ID
	components, err := uc.componentRepo.GetByIDs(componentIDs)
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

		// nếu chưa có languageID trong cache -> query DB
		if _, ok := langMap[comp.LanguageID]; !ok {
			langSetting, err := uc.languageSettingRepo.GetByID(comp.LanguageID)
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

	// lay qr profile form
	childRoleOrg, err := uc.roleOrgRepo.GetByRoleName(string(value.RoleChild))
	if err != nil {
		return nil, err
	}
	formProfile := childRoleOrg.OrgProfile + ":" + child.ID.String()

	// get languages config
	languageConfig, _ := uc.languagesConfigUsecase.GetLanguagesConfigByOwnerNoCtx(childID, value.OwnerRoleLangChild)

	// get avts
	avatars, _ := uc.userImagesUsecase.GetAvt4Owner(childID, value.OwnerRoleChild)

	// Trả về kết quả
	return &response.ChildResponse{
		ChildID:       child.ID.String(),
		ChildName:     child.ChildName,
		Avatar:        "", // Nếu bạn có trường Avatar trong DB thì lấy thêm ở đây
		AvatarURL:     "", // Có thể generate từ link
		QrFormProfile: formProfile,
		// Parent:    *parent,
		Menus:          getMenus,
		LanguageConfig: languageConfig,
		Avatars:        avatars,
		CreatedIndex:   child.CreatedIndex,
	}, nil
}

func (receiver *ChildUseCase) GetAll4Search(ctx *gin.Context) ([]response.ChildrenResponse, error) {
	user, err := receiver.getUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	var children []entity.SChild

	// Nếu là SuperAdmin → lấy tất cả
	if user.IsSuperAdmin() {
		children, err = receiver.childRepo.GetAll()
		if err != nil {
			return nil, err
		}
	} else {
		// Không phải SuperAdmin → return rỗng
		return []response.ChildrenResponse{}, nil
	}

	// Map sang response
	var res []response.ChildrenResponse
	for _, c := range children {
		avatar, _ := receiver.userImagesUsecase.GetAvtIsMain4Owner(c.ID.String(), value.OwnerRoleChild)
		code, _ := receiver.profileGateway.GetChildCode(ctx, c.ID.String())
		res = append(res, response.ChildrenResponse{
			ChildID:      c.ID.String(),
			ChildName:    c.ChildName,
			CreatedIndex: c.CreatedIndex,
			Avatar:       avatar,
			Code:         code,
			LanguageKeys: []string{"vietnamese", "english"},
		})
	}

	return res, nil
}

func (uc *ChildUseCase) GetParentIDByChildID(childID string) (string, error) {
	return uc.childRepo.GetParentIDByChildID(childID)
}

func (uc *ChildUseCase) IsParentOfChild(userID string, childID string) (bool, error) {
	parentID, err := uc.childRepo.GetParentIDByChildID(childID)
	if err != nil {
		return false, err
	}

	if parentID == "" {
		return false, nil
	}

	return parentID == userID, nil
}

func (uc *ChildUseCase) GenerateChildCode(ctx *gin.Context) {
	// get all children
	children, err := uc.childRepo.GetAll()
	if err != nil {
		return
	}

	for _, child := range children {
		// call profile gateway to generate children code
		_, _ = uc.profileGateway.GenerateChildCode(ctx, child.ID.String(), child.CreatedIndex)
	}
}
