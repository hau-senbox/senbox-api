package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeacherMenuUseCase struct {
	TeacherMenuRepo *repository.TeacherMenuRepository
	DB              *gorm.DB
}

func NewTeacherMenuUseCase(repo *repository.TeacherMenuRepository, db *gorm.DB) *TeacherMenuUseCase {
	return &TeacherMenuUseCase{
		TeacherMenuRepo: repo,
		DB:              db,
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

// Get all menus by teacher ID
func (uc *TeacherMenuUseCase) GetByTeacherID(teacherID string) ([]entity.TeacherMenu, error) {
	return uc.TeacherMenuRepo.GetByTeacherID(teacherID)
}

// Delete all menus by teacher ID
func (uc *TeacherMenuUseCase) DeleteByTeacherID(teacherID string) error {
	return uc.TeacherMenuRepo.DeleteByTeacherID(teacherID)
}

// Update is_show flag
func (uc *TeacherMenuUseCase) UpdateIsShow(teacherID, componentID string, isShow bool) error {
	return uc.TeacherMenuRepo.UpdateIsShowByTeacherAndComponentID(teacherID, componentID, isShow)
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

// Bulk replace menus by teacher ID (transactional)
func (uc *TeacherMenuUseCase) ReplaceMenusForTeacher(teacherID uuid.UUID, newMenus []entity.TeacherMenu) error {
	return uc.DB.Transaction(func(tx *gorm.DB) error {
		// Xoá toàn bộ menu cũ
		if err := uc.TeacherMenuRepo.DeleteByTeacherID(teacherID.String()); err != nil {
			return err
		}
		// Tạo mới
		for _, menu := range newMenus {
			menu.TeacherID = teacherID
			if err := uc.TeacherMenuRepo.CreateWithTx(tx, &menu); err != nil {
				return err
			}
		}
		return nil
	})
}
