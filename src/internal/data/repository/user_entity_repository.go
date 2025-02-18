package repository

import (
	"encoding/json"
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserEntityRepository struct {
	DBConn *gorm.DB
}

func NewUserEntityRepository(dbConn *gorm.DB) *UserEntityRepository {
	return &UserEntityRepository{DBConn: dbConn}
}

func (receiver *UserEntityRepository) GetAll() ([]entity.SUserEntity, error) {
	var users []entity.SUserEntity
	query := receiver.DBConn.Table("s_user_entity")
	err := query.
		Preload("Roles").
		Preload("RolePolicies").
		Preload("Guardians").
		Preload("Devices").
		Preload("Company").
		Preload("UserConfig").
		Find(&users).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get all users")
	}

	return users, err
}

func (receiver *UserEntityRepository) GetByID(req request.GetUserEntityByIdRequest) (*entity.SUserEntity, error) {
	var user entity.SUserEntity
	err := receiver.DBConn.
		Preload("Roles").
		Preload("RolePolicies").
		Preload("Guardians").
		Preload("Devices").
		Preload("Company").
		Preload("UserConfig").
		Where("id = ?", req.ID).
		First(&user).Error
	if err != nil {
		log.Error("UserEntityRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get user")
	}
	return &user, nil
}

func (receiver *UserEntityRepository) GetByUsername(req request.GetUserEntityByUsernameRequest) (*entity.SUserEntity, error) {
	var user entity.SUserEntity
	err := receiver.DBConn.
		Preload("Roles").
		Preload("RolePolicies").
		Preload("Guardians").
		Preload("Devices").
		Preload("Company").
		Preload("UserConfig").
		Where("username = ?", req.Username).
		First(&user).Error
	if err != nil {
		log.Error("UserEntityRepository.GetByUsername: " + err.Error())
		return nil, errors.New("failed to get user")
	}
	return &user, nil
}

func (receiver *UserEntityRepository) GetUserDeviceById(deviceId string) (*[]entity.SUserDevices, error) {
	var userDevices []entity.SUserDevices
	err := receiver.DBConn.Model(&entity.SUserDevices{}).
		Where("device_id = ?", deviceId).
		Find(&userDevices).Error

	if err != nil {
		log.Error("UserEntityRepository.GetUserDeviceById: " + err.Error())
		return nil, errors.New("failed to get user")
	}

	return &userDevices, nil
}

func (receiver *UserEntityRepository) GetChildrenOfGuardian(userId string) (*[]response.UserEntityResponseData, error) {
	var userGuardians []entity.SUserGuardians
	err := receiver.DBConn.Model(&entity.SUserGuardians{}).
		Where("guardian_id = ?", userId).
		Find(&userGuardians).Error

	if err != nil {
		log.Error("UserEntityRepository.GetUserDeviceById: " + err.Error())
		return nil, errors.New("failed to get user")
	}

	var result []response.UserEntityResponseData
	for _, userGuardian := range userGuardians {
		var user entity.SUserEntity
		err := receiver.DBConn.Where("id = ?", userGuardian.UserId).First(&user).Error
		if err != nil {
			log.Error("UserEntityRepository.GetUserDeviceById: " + err.Error())
			return nil, errors.New("failed to get user")
		}

		result = append(result, response.UserEntityResponseData{
			ID:       user.ID.String(),
			Username: user.Username,
		})
	}

	return &result, nil
}

func (receiver *UserEntityRepository) CreateUser(req request.CreateUserEntityRequest) error {
	receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		user, _ := receiver.GetByUsername(request.GetUserEntityByUsernameRequest{Username: req.Username})

		if user != nil {
			return errors.New("user already existed")
		}

		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			log.Error("UserRepository.CreateUser: " + err.Error())
			return errors.New("failed to create user " + req.Username)
		}

		var setting entity.SSetting
		err = receiver.DBConn.Table("s_setting").Where("type = ?", value.SettingTypeSignUpPresetValue1).First(&setting).Error

		if err != nil {
			log.Error("UserRepository.CreateUser: " + err.Error())
			return errors.New("failed to create user " + req.Username)
		}

		var signupSetting SignUpFormSetting
		err = json.Unmarshal([]byte(setting.Settings), &signupSetting)

		if err != nil {
			log.Error("UserRepository.CreateUser: " + err.Error())
			return errors.New("failed to create user " + req.Username)
		}

		var userReq = entity.SUserEntity{
			Username:  req.Username,
			Fullname:  signupSetting.SpreadSheetId,
			Birthday:  birthday,
			CompanyId: 1,
			Password:  req.Password,
		}
		userResult := receiver.DBConn.Create(&userReq)

		if userResult.Error != nil {
			log.Error("UserRepository.CreateUser: " + userResult.Error.Error())
			return errors.New("failed to create user " + req.Username)
		}

		if req.Guardians != nil {
			for _, guardian := range *req.Guardians {
				userGuardianResult := receiver.DBConn.Create(&entity.SUserGuardians{
					UserId:     userReq.ID,
					GuardianId: uuid.MustParse(guardian),
				})
				if userGuardianResult.Error != nil {
					log.Error("UserRepository.CreateUser: " + userGuardianResult.Error.Error())
					return errors.New("failed to create user guardian")
				}
			}
		}

		if req.Roles != nil {
			roles := make([]uint, 0)
			for _, roleName := range *req.Roles {
				var role entity.SRole
				err := receiver.DBConn.Model(&entity.SRole{}).Where("role_name = ?", roleName).Find(&role).Error
				if err != nil {
					log.Error("UserRepository.CreateUser: " + err.Error())
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("role does not exist")
					}

					return errors.New("failed to get role")
				}

				if role.ID == 0 {
					return errors.New("role does not exist")
				}

				roles = append(roles, uint(role.ID))
			}
			err := receiver.UpdateUserRole(request.UpdateUserRoleRequest{UserId: userReq.ID.String(), Roles: roles})
			if err != nil {
				return err
			}
		}

		if req.Policies != nil {
			err := receiver.UpdateUserRolePolicy(request.UpdateUserRolePolicyRequest{UserId: userReq.ID.String(), Policies: *req.Policies})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func (receiver *UserEntityRepository) UpdateUser(req request.UpdateUserEntityRequest) error {
	if req.Devices != nil {
		err := receiver.UpdateUserDevice(request.UpdateUserDeviceRequest{UserId: req.ID, Devices: *req.Devices})
		if err != nil {
			return err
		}
	}

	if req.Guardians != nil {
		err := receiver.UpdateUserGuardian(request.UpdateUserGuardianRequest{UserId: req.ID, Guardians: *req.Guardians})
		if err != nil {
			return err
		}
	}

	if req.Roles != nil {
		roles := make([]uint, 0)
		for _, roleName := range *req.Roles {
			var role entity.SRole
			err := receiver.DBConn.Model(&entity.SRole{}).Where("role_name = ?", roleName).Find(&role).Error
			if err != nil {
				log.Error("UserRepository.CreateUser: " + err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("role not exist")
				}

				return errors.New("failed to get role")
			}

			roles = append(roles, uint(role.ID))
		}
		err := receiver.UpdateUserRole(request.UpdateUserRoleRequest{UserId: req.ID, Roles: roles})
		if err != nil {
			return err
		}
	}

	if req.Policies != nil {
		err := receiver.UpdateUserRolePolicy(request.UpdateUserRolePolicyRequest{UserId: req.ID, Policies: *req.Policies})
		if err != nil {
			return err
		}
	}

	updatePayload := map[string]interface{}{}
	updatePayload["username"] = req.Username

	if req.Fullname != nil {
		updatePayload["fullname"] = *req.Fullname
	}

	if req.Phone != nil {
		updatePayload["phone"] = *req.Phone
	}

	if req.Email != nil {
		updatePayload["email"] = *req.Email
	}

	if req.UserConfig != nil {
		configId := *req.UserConfig
		updatePayload["user_config_id"] = configId
	}

	log.Info("PAYLOAD: ", updatePayload)

	resultUpdate := receiver.DBConn.Model(&entity.SUserEntity{}).Where("id = ?", req.ID).
		Updates(updatePayload)

	if resultUpdate.Error != nil {
		log.Errorf("UserEntityRepository.UpdateUser: %v", resultUpdate.Error)
		return errors.New("failed to update user")
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUserDevice(req request.UpdateUserDeviceRequest) error {
	user, err := receiver.GetByID(request.GetUserEntityByIdRequest{ID: req.UserId})

	if err != nil {
		log.Error("UserEntityRepository.UpdateUserRole: " + err.Error())
		return errors.New("failed to get user")
	}

	// delete if exist
	deleteDevices := receiver.DBConn.Where("user_id = ?", user.ID).Delete(&entity.SUserDevices{})

	if deleteDevices.Error != nil {
		log.Error("UserEntityRepository.UpdateUser: " + deleteDevices.Error.Error())
		return errors.New("failed to delete user device")
	}

	for _, deviceId := range req.Devices {
		// check if device is not exist
		var device entity.SDevice
		err := receiver.DBConn.Table("s_device").Where("id = ?", deviceId).First(&device).Error

		if err != nil {
			log.Error("UserEntityRepository.UpdateUser: " + err.Error())
			return errors.New("device not exist")
		}

		userDeviceResult := receiver.DBConn.Create(&entity.SUserDevices{
			UserId:   user.ID,
			DeviceId: device.ID,
		})
		if userDeviceResult.Error != nil {
			log.Errorf("UserEntityRepository.UpdateUser: %v", userDeviceResult.Error)
			return errors.New("failed to create user device")
		}
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUserGuardian(req request.UpdateUserGuardianRequest) error {
	user, err := receiver.GetByID(request.GetUserEntityByIdRequest{ID: req.UserId})

	if err != nil {
		log.Error("UserEntityRepository.UpdateUserRole: " + err.Error())
		return errors.New("failed to get user")
	}

	// delete if exist
	deleteGuardians := receiver.DBConn.Where("user_id = ?", req.UserId).Delete(&entity.SUserGuardians{})

	if deleteGuardians.Error != nil {
		log.Errorf("UserEntityRepository.UpdateUser: %v", deleteGuardians.Error)
		return errors.New("failed to delete user guardian")
	}

	for _, guardianId := range req.Guardians {
		// check guardian user is not exist
		guardian, err := receiver.GetByID(request.GetUserEntityByIdRequest{ID: guardianId})

		if err != nil {
			log.Error("UserEntityRepository.UpdateUser: " + err.Error())
			return errors.New("guardian user not exist")
		}

		userGuardianResult := receiver.DBConn.Create(&entity.SUserGuardians{
			UserId:     user.ID,
			GuardianId: guardian.ID,
		})

		if userGuardianResult.Error != nil {
			log.Errorf("UserEntityRepository.UpdateUser: %v", userGuardianResult.Error)
			return errors.New("failed to create user guardian")
		}
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUserRole(req request.UpdateUserRoleRequest) error {
	user, err := receiver.GetByID(request.GetUserEntityByIdRequest{ID: req.UserId})

	if err != nil {
		log.Error("UserEntityRepository.UpdateUserRole: " + err.Error())
		return errors.New("failed to get user")
	}

	removeResult := receiver.DBConn.Exec("DELETE FROM s_user_roles WHERE user_id = ?", user.ID)

	if removeResult.Error != nil {
		log.Error("UserEntityRepository.UpdateUserRole: " + removeResult.Error.Error())
		return errors.New("failed to remove user role")
	}

	for _, roleId := range req.Roles {
		var role entity.SRole
		err = receiver.DBConn.Model(&entity.SRole{}).Where("id = ?", roleId).First(&role).Error

		if err != nil {
			log.Error("UserEntityRepository.UpdateUserRole: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("role doesn't exist")
			}
			return errors.New("failed to get role")
		}

		result := receiver.DBConn.Create(&entity.SUserRoles{
			UserId: user.ID,
			RoleId: role.ID,
		})

		if result.Error != nil {
			log.Error("UserEntityRepository.UpdateUserRole: " + result.Error.Error())
			return errors.New("failed to assign user role")
		}
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUserRolePolicy(req request.UpdateUserRolePolicyRequest) error {
	user, err := receiver.GetByID(request.GetUserEntityByIdRequest{ID: req.UserId})

	if err != nil {
		log.Error("UserEntityRepository.UpdateUserRolePolicy: " + err.Error())
		return errors.New("failed to get user")
	}

	removeResult := receiver.DBConn.Exec("DELETE FROM s_user_policies WHERE user_id = ?", user.ID)

	if removeResult.Error != nil {
		log.Error("UserEntityRepository.UpdateUserRolePolicy: " + removeResult.Error.Error())
		return errors.New("failed to remove user policy")
	}

	for _, policyId := range req.Policies {
		var policy entity.SRolePolicy
		err = receiver.DBConn.Where("id = ?", policyId).First(&policy).Error

		if err != nil {
			log.Error("UserEntityRepository.UpdateUserRolePolicy: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("policy doesn't exist")
			}
			return errors.New("failed to get policy")
		}

		result := receiver.DBConn.Create(&entity.SUserPolicies{
			UserId:   user.ID,
			PolicyId: policy.ID,
		})

		if result.Error != nil {
			log.Error("UserEntityRepository.UpdateUserRolePolicy: " + result.Error.Error())
			return errors.New("failed to assign user policy")
		}
	}

	return nil
}
