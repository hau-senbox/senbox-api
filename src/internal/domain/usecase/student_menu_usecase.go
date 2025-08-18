package usecase

import (
	"errors"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentMenuUseCase struct {
	StudentMenuRepo      *repository.StudentMenuRepository
	StudentAppRepo       *repository.StudentApplicationRepository
	ComponentRepo        *repository.ComponentRepository
	UserEntityRepo       *repository.UserEntityRepository
	GetUserEntityUseCase *GetUserEntityUseCase
}

func NewStudentMenuUseCase(repo *repository.StudentMenuRepository) *StudentMenuUseCase {
	return &StudentMenuUseCase{StudentMenuRepo: repo}
}

func (uc *StudentMenuUseCase) Create(menu *entity.StudentMenu) error {
	return uc.StudentMenuRepo.Create(menu)
}

func (uc *StudentMenuUseCase) BulkCreate(menus []entity.StudentMenu) error {
	return uc.StudentMenuRepo.BulkCreate(menus)
}

func (uc *StudentMenuUseCase) DeleteByStudentID(studentID string) error {
	return uc.StudentMenuRepo.DeleteByStudentID(studentID)
}

func (uc *StudentMenuUseCase) GetByStudentID(studentID string, isApp bool) (response.GetStudentMenuResponse, error) {
	// B0: Lấy thông tin student
	student, err := uc.StudentAppRepo.GetByID(uuid.MustParse(studentID))
	if student == nil || err != nil {
		return response.GetStudentMenuResponse{}, err
	}

	// B1: Lấy các bản ghi student_menu
	studentMenus, err := uc.StudentMenuRepo.GetByStudentIDActive(studentID)
	if err != nil {
		return response.GetStudentMenuResponse{}, err
	}

	// B2: Lấy componentID từ student_menu
	componentIDs := make([]uuid.UUID, 0, len(studentMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, sm := range studentMenus {
		componentIDs = append(componentIDs, sm.ComponentID)
		componentOrderMap[sm.ComponentID] = sm.Order
		componentIsShowMap[sm.ComponentID] = sm.IsShow
	}

	// B3: Lấy danh sách component tương ứng
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
	if err != nil {
		return response.GetStudentMenuResponse{}, err
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

	return response.GetStudentMenuResponse{
		StudentID:   studentID,
		StudentName: student.StudentName,
		Components:  componentResponses,
	}, nil
}

func (uc *StudentMenuUseCase) UpdateIsShowByStudentAndComponentID(ctx *gin.Context, req request.UpdateStudentMenuRequest) error {
	user, _ := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	// Nếu không phải SuperAdmin → lấy orgIDs mà user đang quản lý
	orgIDs, err := user.GetManagedOrganizationIDs(uc.StudentAppRepo.GetDB())
	if err != nil {
		return nil
	}

	// kiem tra student co nam trong orgIDs khong
	student, _ := uc.StudentAppRepo.GetByID(uuid.MustParse(req.StudentID))
	if student == nil || !student.IsInOrganizations(orgIDs) {
		return errors.New("student not found or not in managed organizations")
	}

	return uc.StudentMenuRepo.UpdateIsShowByStudentAndComponentID(req.StudentID, req.ComponentID, *req.IsShow)
}
