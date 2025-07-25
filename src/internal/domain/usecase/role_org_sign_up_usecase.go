package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
)

type RoleOrgSignUpUseCase struct {
	Repo                 *repository.RoleOrgSignUpRepository
	GetUserEntityUseCase *GetUserEntityUseCase
}

// NewRoleOrgSignUpUseCase: khởi tạo usecase
func NewRoleOrgSignUpUseCase(repo *repository.RoleOrgSignUpRepository) *RoleOrgSignUpUseCase {
	return &RoleOrgSignUpUseCase{
		Repo: repo,
		GetUserEntityUseCase: &GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: repo.DBConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: repo.DBConn},
		},
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

func (receiver *RoleOrgSignUpUseCase) Get4WebAdmin(ctx *gin.Context) ([]entity.SRoleOrgSignUp, error) {
	user, err := receiver.GetUserEntityUseCase.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	// Nếu là SuperAdmin → trả về tất cả roles
	if user.IsSuperAdmin() {
		return receiver.Repo.GetAll()
	}

	// Nếu không phải SuperAdmin → lấy cả Student và Teacher
	roles := make([]entity.SRoleOrgSignUp, 0)

	studentRole, err := receiver.Repo.GetByRoleName(string(value.RoleStudent))
	if err != nil {
		return nil, err
	}
	if studentRole != nil {
		roles = append(roles, *studentRole)
	}

	teacherRole, err := receiver.Repo.GetByRoleName(string(value.RoleTeacher))
	if err != nil {
		return nil, err
	}
	if teacherRole != nil {
		roles = append(roles, *teacherRole)
	}

	return roles, nil
}
