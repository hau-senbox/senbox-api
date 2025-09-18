package usecase

import (
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StaffMenuUseCase struct {
	StaffMenuRepo        *repository.StaffMenuRepository
	StaffAppRepo         *repository.StaffApplicationRepository
	ComponentRepo        *repository.ComponentRepository
	UserEntityRepository *repository.UserEntityRepository
	GetUserEntityUseCase *GetUserEntityUseCase
}

func NewStaffMenuUseCase(repo *repository.StaffMenuRepository) *StaffMenuUseCase {
	return &StaffMenuUseCase{
		StaffMenuRepo: repo,
	}
}

func (uc *StaffMenuUseCase) GetByStaffID(ctx *gin.Context, staffID string, isApp bool) (response.GetStaffMenuResponse, error) {
	// B0: Lấy thông tin staff
	staff, err := uc.StaffAppRepo.GetByID(uuid.MustParse(staffID))
	if staff == nil || err != nil {
		return response.GetStaffMenuResponse{}, err
	}

	// B1: Lấy các bản ghi staff_menu
	staffMenus, err := uc.StaffMenuRepo.GetByStaffIDActive(staffID)
	if err != nil {
		return response.GetStaffMenuResponse{}, err
	}

	// B2: Lấy componentID từ staff_menu
	componentIDs := make([]uuid.UUID, 0, len(staffMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, tm := range staffMenus {
		componentIDs = append(componentIDs, tm.ComponentID)
		componentOrderMap[tm.ComponentID] = tm.Order
		componentIsShowMap[tm.ComponentID] = tm.IsShow
	}

	// B3: Lấy danh sách component tương ứng
	appLanguage, _ := ctx.Get("app_language")
	components, err := uc.ComponentRepo.GetByIDsAndLanguage(componentIDs, appLanguage.(uint))
	if err != nil {
		return response.GetStaffMenuResponse{}, err
	}

	// B4: Build danh sách response
	componentResponses := make([]response.ComponentResponse, 0, len(components))
	for _, comp := range components {
		if isApp {
			visible, _ := helper.GetVisibleToValueComponent(comp.Value.String())
			if !visible {
				continue
			}
		}
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

	// get user by user ID from teacher
	user, _ := uc.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{
		ID: staff.UserID.String(),
	})

	return response.GetStaffMenuResponse{
		StaffID:    staffID,
		StaffName:  user.Nickname,
		Components: componentResponses,
	}, nil
}
