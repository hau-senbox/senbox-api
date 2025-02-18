package repository

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	DBConn *gorm.DB
}

func (receiver *UserRepository) FindUserByUsername(username string) *entity.SUser {
	var result entity.SUser
	if err := receiver.DBConn.Where("username = ?", username).First(&result).Error; err != nil {
		log.Debug(err)
		return nil
	}

	return &result
}

func (receiver *UserRepository) FindUserById(id *string) (*entity.SUser, error) {
	var result entity.SUser
	if err := receiver.DBConn.Where("user_id = ?", id).First(&result).Error; err != nil {
		log.Debug(err)
		return nil, err
	}

	return &result, nil
}

func (receiver *UserRepository) GetUsers(role value.Role) ([]*entity.SUser, error) {
	var result []*entity.SUser
	err := receiver.DBConn.Table("s_user").Where("role & ? = ?", role, role).Find(&result).Order("updated_at DESC").Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (receiver *UserRepository) SaveUser(user *entity.SUser) (*entity.SUser, error) {
	err := receiver.DBConn.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (receiver *UserRepository) UpdateUser(e *entity.SUser) error {
	log.Debug("UPDATE s_user SET password = ''", e.Password, "' WHERE user_id = ''", e.UserId, "''")
	err := receiver.DBConn.Exec("UPDATE s_user SET password = ? WHERE user_id = ?", e.Password, e.UserId).Error
	log.Error(err)

	return err
}
