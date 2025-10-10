package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type RoleOrgSignUpRepository struct {
	DBConn *gorm.DB
}

func (r *RoleOrgSignUpRepository) Create(roleOrg *entity.SRoleOrgSignUp) error {
	return r.DBConn.Create(roleOrg).Error
}

// UpdateOrCreate: cập nhật nếu tồn tại theo role_name, nếu không thì tạo mới
func (r *RoleOrgSignUpRepository) UpdateOrCreate(role *entity.SRoleOrgSignUp) error {
	var existing entity.SRoleOrgSignUp

	// Tìm theo role_name
	err := r.DBConn.Where("role_name = ?", role.RoleName).First(&existing).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Không tìm thấy => tạo mới
			return r.DBConn.Create(role).Error
		}
		// Lỗi khác
		return err
	}

	// Cập nhật nếu đã tồn tại
	role.ID = existing.ID // giữ nguyên ID
	return r.DBConn.Model(&existing).Updates(role).Error
}

func (r *RoleOrgSignUpRepository) GetAll() ([]entity.SRoleOrgSignUp, error) {
	var roles []entity.SRoleOrgSignUp
	err := r.DBConn.Find(&roles).Error
	return roles, err
}

func (r *RoleOrgSignUpRepository) GetByRoleName(roleName string) (*entity.SRoleOrgSignUp, error) {
	var role entity.SRoleOrgSignUp
	err := r.DBConn.Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleOrgSignUpRepository) GetByID(id string) (*entity.SRoleOrgSignUp, error) {
	var role entity.SRoleOrgSignUp

	err := r.DBConn.
		Where("id = ?", id).
		First(&role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &role, nil
}

func (r *RoleOrgSignUpRepository) WithTx(tx *gorm.DB) *RoleOrgSignUpRepository {
	return &RoleOrgSignUpRepository{DBConn: tx}
}
