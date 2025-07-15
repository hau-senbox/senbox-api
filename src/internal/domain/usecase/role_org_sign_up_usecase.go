package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type RoleOrgSignUpUseCase struct {
	Repo *repository.RoleOrgSignUpRepository
}

// NewRoleOrgSignUpUseCase: khởi tạo usecase
func NewRoleOrgSignUpUseCase(repo *repository.RoleOrgSignUpRepository) *RoleOrgSignUpUseCase {
	return &RoleOrgSignUpUseCase{
		Repo: repo,
	}
}

// Execute: Gọi UpdateOrCreate từ repository
func (uc *RoleOrgSignUpUseCase) UpdateOrCreateExecute(role *entity.SRoleOrgSignUp) error {
	return uc.Repo.UpdateOrCreate(role)
}

// Execute: Gọi UpdateOrCreate từ repository
func (uc *RoleOrgSignUpUseCase) GetAll() ([]entity.SRoleOrgSignUp, error) {
	return uc.Repo.GetAll()
}

// GetByRoleName: Gọi GetByRoleName từ repository
func (uc *RoleOrgSignUpUseCase) GetByRoleName(roleName string) (*entity.SRoleOrgSignUp, error) {
	return uc.Repo.GetByRoleName(roleName)
}
