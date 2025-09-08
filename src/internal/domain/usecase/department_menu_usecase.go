package usecase

import (
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
)

type DepartmentMenuUseCase struct {
	DepartmentMenuRepository *repository.DepartmentMenuRepository
	ComponentRepository      *repository.ComponentRepository
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

	// B4: Build danh sách response
	componentResponses := make([]response.ComponentResponse, 0, len(components))
	for _, comp := range components {
		componentResponses = append(componentResponses, response.ComponentResponse{
			ID:    comp.ID.String(),
			Name:  comp.Name,
			Type:  comp.Type.String(),
			Key:   comp.Key,
			Value: string(comp.Value),
		})
	}

	return response.GetDepartmentMenuResponse{
		Components: componentResponses,
	}, nil
}

func (uc *DepartmentMenuUseCase) GetDepartmentMenu4App(departmentID string) (response.GetDepartmentMenuResponse, error) {

	departmentMenus, err := uc.DepartmentMenuRepository.GetByDepartmentID(departmentID)
	if err != nil {
		return response.GetDepartmentMenuResponse{}, err
	}

	// B1: Lấy tất cả ComponentID từ departmentMenus
	componentIDs := make([]uuid.UUID, 0, len(departmentMenus))
	componentOrderMap := make(map[uuid.UUID]int) // lưu order theo ComponentID

	for _, cm := range departmentMenus {
		componentIDs = append(componentIDs, cm.ComponentID)
		componentOrderMap[cm.ComponentID] = cm.Order
	}

	// B2: Lấy danh sách Component theo IDs
	components, err := uc.ComponentRepository.GetByIDs(componentIDs)
	if err != nil {
		return response.GetDepartmentMenuResponse{}, err
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

	return response.GetDepartmentMenuResponse{
		Components: componentResponses,
	}, nil
}
