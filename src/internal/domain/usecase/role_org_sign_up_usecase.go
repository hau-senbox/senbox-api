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

// Execute: Gọi UpdateOrCreate từ repository
func (uc *RoleOrgSignUpUseCase) Get4App() ([]entity.SRoleOrgSignUp, error) {
	roles, err := uc.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	// app không lay Staff
	var filtered []entity.SRoleOrgSignUp
	for _, role := range roles {
		if role.RoleName != string(value.RoleStaff) {
			filtered = append(filtered, role)
		}
	}

	return filtered, nil
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

	// Nếu là SuperAdmin → chỉ lấy role Child
	if user.IsSuperAdmin() {
		childRole, err := receiver.Repo.GetByRoleName(string(value.RoleChild))
		if err != nil {
			return nil, err
		}
		if childRole != nil {
			return []entity.SRoleOrgSignUp{*childRole}, nil
		}
		return []entity.SRoleOrgSignUp{}, nil
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

	staffRole, err := receiver.Repo.GetByRoleName(string(value.RoleStaff))
	if err != nil {
		return nil, err
	}
	if staffRole != nil {
		roles = append(roles, *staffRole)
	}

	return roles, nil
}
