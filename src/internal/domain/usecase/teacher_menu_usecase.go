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
	"gorm.io/gorm"
)

type TeacherMenuUseCase struct {
	TeacherMenuRepo      *repository.TeacherMenuRepository
	TeacherAppRepo       *repository.TeacherApplicationRepository
	ComponentRepo        *repository.ComponentRepository
	UserEntityRepository *repository.UserEntityRepository
	GetUserEntityUseCase *GetUserEntityUseCase
}

func NewTeacherMenuUseCase(repo *repository.TeacherMenuRepository) *TeacherMenuUseCase {
	return &TeacherMenuUseCase{
		TeacherMenuRepo: repo,
	}
}

// Create single teacher menu
func (uc *TeacherMenuUseCase) Create(menu *entity.TeacherMenu) error {
	return uc.TeacherMenuRepo.Create(menu)
}

// Bulk create teacher menus
func (uc *TeacherMenuUseCase) BulkCreate(menus []entity.TeacherMenu) error {
	return uc.TeacherMenuRepo.BulkCreate(menus)
}

func (uc *TeacherMenuUseCase) GetByTeacherID(teacherID string, isApp bool) (response.GetTeacherMenuResponse, error) {
	// B0: Lấy thông tin teacher
	teacher, err := uc.TeacherAppRepo.GetByID(uuid.MustParse(teacherID))
	if teacher == nil || err != nil {
		return response.GetTeacherMenuResponse{}, err
	}

	// B1: Lấy các bản ghi teacher_menu
	teacherMenus, err := uc.TeacherMenuRepo.GetByTeacherIDActive(teacherID)
	if err != nil {
		return response.GetTeacherMenuResponse{}, err
	}

	// B2: Lấy componentID từ teacher_menu
	componentIDs := make([]uuid.UUID, 0, len(teacherMenus))
	componentOrderMap := make(map[uuid.UUID]int)
	componentIsShowMap := make(map[uuid.UUID]bool)

	for _, tm := range teacherMenus {
		componentIDs = append(componentIDs, tm.ComponentID)
		componentOrderMap[tm.ComponentID] = tm.Order
		componentIsShowMap[tm.ComponentID] = tm.IsShow
	}

	// B3: Lấy danh sách component tương ứng
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
	if err != nil {
		return response.GetTeacherMenuResponse{}, err
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
		ID: teacher.UserID.String(),
	})

	return response.GetTeacherMenuResponse{
		TeacherID:   teacherID,
		TeacherName: user.Nickname,
		Components:  componentResponses,
	}, nil
}

// Delete all menus by teacher ID
func (uc *TeacherMenuUseCase) DeleteByTeacherID(teacherID string) error {
	return uc.TeacherMenuRepo.DeleteByTeacherID(teacherID)
}

// Update is_show flag
func (uc *TeacherMenuUseCase) UpdateIsShow(ctx *gin.Context, req request.UpdateTeacherMenuRequest) error {
	user, _ := uc.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)

	// Nếu không phải SuperAdmin → lấy orgIDs mà user đang quản lý
	orgIDs, err := user.GetManagedOrganizationIDs(uc.TeacherAppRepo.GetDB())
	if err != nil {
		return err
	}

	// Kiểm tra teacher có thuộc các tổ chức mà user đang quản lý không
	teacher, _ := uc.TeacherAppRepo.GetByID(uuid.MustParse(req.TeacherID))
	if teacher == nil || !teacher.IsInOrganizations(orgIDs) {
		return errors.New("teacher not found or not in managed organizations")
	}

	return uc.TeacherMenuRepo.UpdateIsShowByTeacherAndComponentID(req.TeacherID, req.ComponentID, *req.IsShow)
}

// Update full record with transaction
func (uc *TeacherMenuUseCase) UpdateMenu(tx *gorm.DB, menu *entity.TeacherMenu) error {
	return uc.TeacherMenuRepo.UpdateWithTx(tx, menu)
}

// Get by teacher + component ID
func (uc *TeacherMenuUseCase) GetByTeacherAndComponent(tx *gorm.DB, teacherID, componentID uuid.UUID) (*entity.TeacherMenu, error) {
	return uc.TeacherMenuRepo.GetByTeacherIDAndComponentID(tx, teacherID, componentID)
}

// Delete all menus globally (dangerous)
func (uc *TeacherMenuUseCase) DeleteAll() error {
	return uc.TeacherMenuRepo.DeleteAll()
}

// Delete by component ID
func (uc *TeacherMenuUseCase) DeleteByComponentID(componentID string) error {
	return uc.TeacherMenuRepo.DeleteByComponentID(componentID)
}
