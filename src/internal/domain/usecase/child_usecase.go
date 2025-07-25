package usecase

import (
	"errors"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChildUseCase struct {
	childRepo            *repository.ChildRepository
	userRepo             *repository.UserEntityRepository
	componentRepo        *repository.ComponentRepository
	childMenuRepo        *repository.ChildMenuRepository
	roleOrgRepo          *repository.RoleOrgSignUpRepository
	getUserEntityUseCase *GetUserEntityUseCase
}

func NewChildUseCase(
	childRepo *repository.ChildRepository,
	userRepo *repository.UserEntityRepository,
	componentRepo *repository.ComponentRepository,
	childMenuRepo *repository.ChildMenuRepository,
	roleOrgRepo *repository.RoleOrgSignUpRepository,
	getUserEntityUseCase *GetUserEntityUseCase,
) *ChildUseCase {
	return &ChildUseCase{
		childRepo:            childRepo,
		userRepo:             userRepo,
		componentRepo:        componentRepo,
		childMenuRepo:        childMenuRepo,
		roleOrgRepo:          roleOrgRepo,
		getUserEntityUseCase: getUserEntityUseCase,
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

	err := uc.childRepo.Create(child)

	if err == nil {
		//tao child menu
		childRoleOrg, _ := uc.roleOrgRepo.GetByRoleName(string(value.RoleChild))
		if childRoleOrg != nil {
			components, _ := uc.componentRepo.GetBySectionID(childRoleOrg.ID.String())

			for index, component := range components {
				visible, _ := helper.GetVisibleToValueComponent(string(component.Value))
				err := uc.childMenuRepo.Create(&entity.ChildMenu{
					ID:          uuid.New(),
					ChildID:     childID,
					ComponentID: component.ID,
					Order:       index,
					IsShow:      true,
					Visible:     visible,
				})
				if err != nil {
					continue
				}
			}
		}
	}

	return err
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

	// Lấy thông tin parent
	// parent, err := uc.userRepo.GetByID(request.GetUserEntityByIDRequest{
	// 	ID: child.ParentID.String(),
	// })
	// if err != nil {
	// 	return nil, err
	// }

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

	// Build danh sách ComponentChildResponse
	menus := make([]response.ComponentResponse, 0)
	for _, comp := range components {
		menu := response.ComponentResponse{
			ID:     comp.ID.String(),
			Name:   comp.Name,
			Type:   comp.Type.String(),
			Key:    comp.Key,
			Value:  string(comp.Value),
			Order:  componentOrderMap[comp.ID],
			IsShow: componentIsShowMap[comp.ID],
		}
		menus = append(menus, menu)
	}

	// lay qr profile form
	childRoleOrg, err := uc.roleOrgRepo.GetByRoleName(string(value.RoleChild))
	if err != nil {
		return nil, err
	}
	formProfile := childRoleOrg.OrgProfile + ":" + child.ID.String()
	// Trả về kết quả
	return &response.ChildResponse{
		ChildID:       child.ID.String(),
		ChildName:     child.ChildName,
		Avatar:        "", // Nếu bạn có trường Avatar trong DB thì lấy thêm ở đây
		AvatarURL:     "", // Có thể generate từ link
		QrFormProfile: formProfile,
		// Parent:    *parent,
		Menus: menus,
	}, nil
}

func (receiver *ChildUseCase) GetAll4Search(ctx *gin.Context) ([]entity.SChild, error) {
	user, err := receiver.getUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	// Nếu là SuperAdmin → trả về tất cả child
	if user.IsSuperAdmin() {
		return receiver.childRepo.GetAll()
	}

	// Nếu không phải SuperAdmin → return nil
	return nil, nil
}
