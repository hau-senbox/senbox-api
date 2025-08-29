package usecase

import (
	"context"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
)

type TeacherMenuOrganizationUseCase struct {
	TeacherMenuOrganizationRepository *repository.TeacherMenuOrganizationRepository
	ComponentRepo                     *repository.ComponentRepository
	TeacherRepo                       *repository.TeacherApplicationRepository
	DeviceRepository                  *repository.DeviceRepository
}

// Lấy danh sách menu component của giáo viên theo organization
func (uc *TeacherMenuOrganizationUseCase) GetTeacherMenuOrg4Admin(ctx context.Context, teacherID, orgID string) ([]response.ComponentResponse, error) {
	// 1. Lấy danh sách menu của giáo viên trong org
	teacherMenus, err := uc.TeacherMenuOrganizationRepository.GetAllByTeacherAndOrg(ctx, teacherID, orgID)
	if err != nil {
		return nil, err
	}

	if len(teacherMenus) == 0 {
		return []response.ComponentResponse{}, nil
	}

	// 2. Chuẩn bị danh sách componentID + mapping order
	componentIDs := make([]uuid.UUID, 0, len(teacherMenus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range teacherMenus {
		compID, err := uuid.Parse(cm.ComponentID)
		if err != nil {
			continue // skip nếu ComponentID không hợp lệ
		}
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

func (uc *TeacherMenuOrganizationUseCase) GetTeacherMenuOrg4App(ctx context.Context, req request.GetTeacherOrganizationMenuRequest) ([]response.ComponentResponse, error) {
	// kiem tra device dang co nam trong org hay khong
	isExist, err := uc.DeviceRepository.CheckDeviceExistInOrganization(req.DeviceID, req.OrganizationID)
	if err != nil {
		return []response.ComponentResponse{}, err
	}
	if !isExist {
		return []response.ComponentResponse{}, nil
	}
	// get teacher by user ID
	teacher, err := uc.TeacherRepo.GetByUserID(req.UserID)
	if err != nil {
		return nil, err
	}

	// Lấy danh sách menu của giáo viên trong org
	teacherMenus, err := uc.TeacherMenuOrganizationRepository.GetAllByTeacherAndOrg(ctx, teacher.ID.String(), req.OrganizationID)
	if err != nil {
		return nil, err
	}

	if len(teacherMenus) == 0 {
		return []response.ComponentResponse{}, nil
	}

	// 2. Chuẩn bị danh sách componentID + mapping order
	componentIDs := make([]uuid.UUID, 0, len(teacherMenus))
	componentOrderMap := make(map[uuid.UUID]int)

	for _, cm := range teacherMenus {
		compID, err := uuid.Parse(cm.ComponentID)
		if err != nil {
			continue // skip nếu ComponentID không hợp lệ
		}
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
			Value: helper.BuildSectionValueMenu(string(comp.Value), comp),
			Order: componentOrderMap[comp.ID],
		}
		menus = append(menus, menu)
	}

	return menus, nil
}
