package repository

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"strings"
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

func (receiver *UserRepository) RegisterDeviceWithUsers(req request.RegisterDeviceRequest) (*entity.SUser, *entity.SUser, *entity.SUser, error) {
	var result []entity.SUser
	var primaryUser = entity.SUser{}
	err := receiver.DBConn.Where("username = ?", req.Primary.Email).First(&primaryUser).Error
	if err != nil {
		primaryUser = entity.SUser{
			UserId:   uuid.New().String(),
			Username: strings.ToLower(req.Primary.Email),
			Fullname: req.Primary.Fullname,
			Phone:    req.Primary.Phone,
			Email:    req.Primary.Email,
		}
		err = receiver.DBConn.Create(&primaryUser).Error
		if err != nil {
			return nil, nil, nil, err
		}
		result = append(result, primaryUser)
	}

	var secondary = entity.SUser{}
	err = receiver.DBConn.Where("username = ?", req.Secondary.Email).First(&secondary).Error
	if err != nil {
		secondary = entity.SUser{
			UserId:   uuid.New().String(),
			Username: strings.ToLower(req.Secondary.Email),
			Fullname: req.Secondary.Fullname,
			Phone:    req.Secondary.Phone,
			Email:    req.Secondary.Email,
		}
		err = receiver.DBConn.Create(&secondary).Error
		if err != nil {
			return nil, nil, nil, err
		}
		result = append(result, secondary)
	}

	var tertiary = entity.SUser{}
	err = receiver.DBConn.Where("username = ?", req.Tertiary.Email).First(&tertiary).Error
	if err != nil {
		tertiary = entity.SUser{
			UserId:   uuid.New().String(),
			Username: strings.ToLower(req.Tertiary.Email),
			Fullname: req.Tertiary.Fullname,
			Phone:    req.Tertiary.Phone,
			Email:    req.Tertiary.Email,
		}
		err = receiver.DBConn.Create(&tertiary).Error
		if err != nil {
			return nil, nil, nil, err
		}
		result = append(result, tertiary)
	}
	return &primaryUser, &secondary, &tertiary, nil
}

func (receiver *UserRepository) UpdateUser(e *entity.SUser) error {
	log.Debug("UPDATE s_user SET password = ''", e.Password, "' WHERE user_id = ''", e.UserId, "''")
	err := receiver.DBConn.Exec("UPDATE s_user SET password = ? WHERE user_id = ?", e.Password, e.UserId).Error
	log.Error(err)

	return err
}
