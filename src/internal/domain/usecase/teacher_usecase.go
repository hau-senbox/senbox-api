package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type TeacherApplicationUseCase struct {
	TeacherRepo repository.TeacherApplicationRepository
}

func NewTeacherApplicationUseCase(repo repository.TeacherApplicationRepository) *TeacherApplicationUseCase {
	return &TeacherApplicationUseCase{TeacherRepo: repo}
}

// Create a new teacher application
func (uc *TeacherApplicationUseCase) Create(teacher *entity.STeacherFormApplication) error {
	teacher.ID = uuid.New()
	teacher.CreatedAt = time.Now()
	return uc.TeacherRepo.Create(teacher)
}

// Get by ID
func (uc *TeacherApplicationUseCase) GetByID(id uuid.UUID) (*entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetByID(id)
}

// Get all applications
func (uc *TeacherApplicationUseCase) GetAll() ([]entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetAll()
}

// Update application
func (uc *TeacherApplicationUseCase) Update(teacher *entity.STeacherFormApplication) error {
	return uc.TeacherRepo.Update(teacher)
}

// Delete by ID
func (uc *TeacherApplicationUseCase) Delete(id uuid.UUID) error {
	return uc.TeacherRepo.Delete(id)
}

// Get by UserID
func (uc *TeacherApplicationUseCase) GetByUserID(userID string) ([]entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetByUserID(userID)
}

// Get by OrganizationID
func (uc *TeacherApplicationUseCase) GetByOrganizationID(orgID string) ([]entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetByOrganizationID(orgID)
}

// Get by list of OrganizationIDs
func (uc *TeacherApplicationUseCase) GetByOrganizationIDs(orgIDs []string) ([]entity.STeacherFormApplication, error) {
	return uc.TeacherRepo.GetByOrganizationIDs(orgIDs)
}

// Check if teacher belongs to one of the given organizations
func (uc *TeacherApplicationUseCase) CheckTeacherBelongsToOrganizations(teacherID uuid.UUID, orgIDs []string) (bool, error) {
	return uc.TeacherRepo.CheckTeacherBelongsToOrganizations(uc.TeacherRepo.GetDB(), teacherID, orgIDs)
}
