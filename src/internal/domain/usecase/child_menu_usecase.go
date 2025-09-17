package usecase

import (
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
)

type ChildMenuUseCase struct {
	Repo          *repository.ChildMenuRepository
	ComponentRepo *repository.ComponentRepository
	ChildRepo     *repository.ChildRepository
}

func NewChildMenuUseCase(repo *repository.ChildMenuRepository) *ChildMenuUseCase {
	return &ChildMenuUseCase{Repo: repo}
}

func (uc *ChildMenuUseCase) Create(menu *entity.ChildMenu) error {
	return uc.Repo.Create(menu)
}

func (uc *ChildMenuUseCase) BulkCreate(menus []entity.ChildMenu) error {
	return uc.Repo.BulkCreate(menus)
}

func (uc *ChildMenuUseCase) DeleteByChildID(childID string) error {
	return uc.Repo.DeleteByChildID(childID)
}

func (uc *ChildMenuUseCase) GetByChildID(childID string, isApp bool) (response.GetChildMenuResponse, error) {
	child, err := uc.ChildRepo.GetByID(childID)
	if child == nil || err != nil {
		return response.GetChildMenuResponse{}, err
	}
	childMenus, err := uc.Repo.GetByChildIDActive(childID)
	if err != nil {
		return response.GetChildMenuResponse{}, err
	}

	// B1: Lấy tất cả ComponentID từ childMenus
	componentIDs := make([]uuid.UUID, 0, len(childMenus))
	componentOrderMap := make(map[uuid.UUID]int)   // lưu order theo ComponentID
	componentIsShowMap := make(map[uuid.UUID]bool) // lưu is_show theo ComponentID

	for _, cm := range childMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
		componentIsShowMap[cm.ComponentID] = cm.IsShow
	}

	// B2: Lấy danh sách Component theo IDs
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
	if err != nil {
		return response.GetChildMenuResponse{}, err
	}

	// B3: Map sang ComponentChildResponse
	componentResponses := make([]response.ComponentResponse, 0, len(components))
	for _, comp := range components {
		if isApp {
			visible, _ := helper.GetVisibleToValueComponent(comp.Value.String())
			if !visible {
				continue
			}
		}
		componentResponses = append(componentResponses, response.ComponentResponse{
			ID:       comp.ID.String(),
			Name:     comp.Name,
			Type:     comp.Type.String(),
			Key:      comp.Key,
			Value:    helper.BuildSectionValueMenu(string(comp.Value), comp),
			Order:    componentOrderMap[comp.ID],
			IsShow:   componentIsShowMap[comp.ID],
			Language: comp.Language,
		})
	}

	return response.GetChildMenuResponse{
		ChildID:    childID,
		ChildName:  child.ChildName,
		Components: componentResponses,
	}, nil
}

func (uc *ChildMenuUseCase) UpdateIsShowByChildAndComponentID(req request.UpdateChildMenuRequest) error {
	return uc.Repo.UpdateIsShowByChildAndComponentID(req.ChildID, req.ComponentID, *req.IsShow)
}
