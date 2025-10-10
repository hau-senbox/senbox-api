package usecase

import (
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ParentMenuUseCase struct {
	Repo          *repository.ParentMenuRepository
	ComponentRepo *repository.ComponentRepository
	UserRepo      *repository.UserEntityRepository
}

func NewParentMenuUseCase(repo *repository.ParentMenuRepository) *ParentMenuUseCase {
	return &ParentMenuUseCase{Repo: repo}
}

func (uc *ParentMenuUseCase) Create(menu *entity.ParentMenu) error {
	return uc.Repo.Create(menu)
}

func (uc *ParentMenuUseCase) BulkCreate(menus []entity.ParentMenu) error {
	return uc.Repo.BulkCreate(menus)
}

func (uc *ParentMenuUseCase) DeleteByParentID(parentID string) error {
	return uc.Repo.DeleteByParentID(parentID)
}

func (uc *ParentMenuUseCase) GetByParentID(ctx *gin.Context, parentID string, userID string) (response.GetParentMenuResponse, error) {

	parentMenus, err := uc.Repo.GetByParentIDActive(parentID)
	if err != nil {
		return response.GetParentMenuResponse{}, err
	}

	// B1: Lấy tất cả ComponentID từ parentMenus
	componentIDs := make([]uuid.UUID, 0, len(parentMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, pm := range parentMenus {
		componentIDs = append(componentIDs, pm.ComponentID)
		componentOrderMap[pm.ComponentID] = pm.Order
		componentIsShowMap[pm.ComponentID] = pm.IsShow
	}

	// B2: Lấy danh sách Component theo IDs
	appLanguage, _ := ctx.Get("app_language")
	components, err := uc.ComponentRepo.GetByIDsAndLanguage(componentIDs, appLanguage.(uint))
	if err != nil {
		return response.GetParentMenuResponse{}, err
	}

	// B3: Map sang ComponentParentResponse
	componentResponses := make([]response.ComponentResponse, 0, len(components))
	for _, comp := range components {
		componentResponses = append(componentResponses, response.ComponentResponse{
			ID:     comp.ID.String(),
			Name:   comp.Name,
			Type:   comp.Type.String(),
			Key:    comp.Key,
			Value:  helper.BuildSectionValueMenu(string(comp.Value), comp),
			Order:  componentOrderMap[comp.ID],
			IsShow: componentIsShowMap[comp.ID],
		})
	}

	user, _ := uc.UserRepo.GetByID(request.GetUserEntityByIDRequest{
		ID: userID,
	})

	return response.GetParentMenuResponse{
		ParentID:   parentID,
		ParentName: user.Nickname,
		Components: componentResponses,
	}, nil
}
