package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
)

type ParentUseCase struct {
	UserRepo       *repository.UserEntityRepository
	ParentMenuRepo *repository.ParentMenuRepository
	ComponentRepo  *repository.ComponentRepository
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

	return &response.ParentResponseBase{
		ParentID:   parentID,
		ParentName: parent.Nickname,
		Avatar:     "",
		AvatarURL:  "",
		Menus:      menus,
		CustomID:   parent.CustomID,
	}, nil
}
