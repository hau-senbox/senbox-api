package usecase

import (
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepartmentMenuOrganizationUseCase struct {
	DepartmentMenuOrganizationRepository *repository.DepartmentMenuOrganizationRepository
	ComponentRepo                        *repository.ComponentRepository
	DeviceRepository                     *repository.DeviceRepository
	OrganizationRepository               *repository.OrganizationRepository
	DepartmentGateway                    gateway.DepartmentGateway
}

func (uc *DepartmentMenuOrganizationUseCase) GetDepartmentMenuOrg4GW(ctx *gin.Context, departmentID, orgID string) ([]response.ComponentResponse, error) {
	// 1. Lấy danh sách menu của giáo viên trong org
	departmentMenusOrg, err := uc.DepartmentMenuOrganizationRepository.GetAllByDepartmentAndOrg(ctx, departmentID, orgID)
	if err != nil {
		return nil, err
	}

	if len(departmentMenusOrg) == 0 {
		return []response.ComponentResponse{}, nil
	}

	// 2. Chuẩn bị danh sách componentID + mapping order
	componentIDs := make([]uuid.UUID, 0, len(departmentMenusOrg))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range departmentMenusOrg {
		compID := cm.ComponentID
		componentIDs = append(componentIDs, compID)
		componentOrderMap[compID] = cm.Order
	}

	// 3. Lấy components theo ID
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
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
			Value: string(comp.Value),
			Order: componentOrderMap[comp.ID],
		}
		menus = append(menus, menu)
	}

	return menus, nil
}

func (uc *DepartmentMenuOrganizationUseCase) GetDepartmentMenuOrg4App(
	ctx *gin.Context,
	req request.GetDepartmentMenuOrganizationRequest,
) (*response.GetDepartmentMenuOrganizationResponse, error) {

	// 1. Kiểm tra device có trong org hay không
	isExist, err := uc.DeviceRepository.CheckDeviceExistInOrganization(req.DeviceID, req.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, nil
	}

	// 2. Lấy danh sách department từ gateway
	departments, err := uc.DepartmentGateway.GetDepartmentsByUser(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Duyệt department
	for _, department := range departments {
		departmentMenusOrg, err := uc.DepartmentMenuOrganizationRepository.
			GetAllByDepartmentAndOrg(ctx, department.ID, req.OrganizationID)
		if err != nil || len(departmentMenusOrg) == 0 {
			continue
		}

		// Chuẩn bị componentIDs + mapping order
		componentIDs := make([]uuid.UUID, 0, len(departmentMenusOrg))
		componentOrderMap := make(map[uuid.UUID]int)
		for _, cm := range departmentMenusOrg {
			componentIDs = append(componentIDs, cm.ComponentID)
			componentOrderMap[cm.ComponentID] = cm.Order
		}

		components, err := uc.ComponentRepo.GetByIDs(componentIDs)
		if err != nil {
			return nil, err
		}

		menus := make([]response.ComponentResponse, 0, len(components))
		for _, comp := range components {
			menus = append(menus, response.ComponentResponse{
				ID:    comp.ID.String(),
				Name:  comp.Name,
				Type:  comp.Type.String(),
				Key:   comp.Key,
				Value: helper.BuildSectionValueMenu(string(comp.Value), comp),
				Order: componentOrderMap[comp.ID],
			})
		}

		orgInfo, _ := uc.OrganizationRepository.GetByID(req.OrganizationID)

		// return ngay department đầu tiên có menu
		return &response.GetDepartmentMenuOrganizationResponse{
			Section:     department.Name + " Menu At " + orgInfo.OrganizationName,
			MenuIconKey: department.Icon,
			Components:  menus,
		}, nil
	}

	return nil, nil // không có department nào có menu
}
