package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
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

func (receiver *UserEntityRepository) GetAllByOrganizationID(organizationID string) ([]entity.SUserEntity, error) {
	var users []entity.SUserEntity
	query := receiver.DBConn.Model(entity.SUserEntity{})
	err := query.
		Preload("Roles").
		Preload("Devices").
		Preload("Organizations").
		Joins("INNER JOIN s_user_organizations ON s_user_entity.id = s_user_organizations.user_id").
		Where("s_user_organizations.organization_id = ?", organizationID).
		Find(&users).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get all users")
	}

	return users, err
}

func (receiver *UserEntityRepository) GetAll() ([]entity.SUserEntity, error) {
	var users []entity.SUserEntity
	query := receiver.DBConn.Model(entity.SUserEntity{})
	err := query.
		Preload("Roles").
		Preload("Devices").
		Preload("Organizations").
		Find(&users).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get all users")
	}

	return users, err
}

func (receiver *UserEntityRepository) GetAllUserAuthorize(userID string) ([]entity.SUserFunctionAuthorize, error) {
	var rights []entity.SUserFunctionAuthorize
	query := receiver.DBConn.Model(entity.SUserFunctionAuthorize{})
	err := query.
		Where("user_id = ?", userID).
		Preload("User").
		Preload("FunctionClaim").
		Preload("FunctionClaimPermission").
		Find(&rights).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAllUserAuthorize: " + err.Error())
		return nil, errors.New("failed to get all rights")
	}

	return rights, err
}

func (receiver *UserEntityRepository) UpdateUserAuthorize(req request.UpdateUserAuthorizeRequest) error {
	// check if user exist
	user, err := receiver.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	// check if function claim exist
	var functionClaim entity.SFunctionClaim
	err = receiver.DBConn.Model(entity.SFunctionClaim{}).Where("id = ?", req.FunctionClaimID).First(&functionClaim).Error
	if err != nil {
		return err
	}

	// check if function claim permission exist
	var functionClaimPermission entity.SFunctionClaimPermission
	err = receiver.DBConn.Model(entity.SFunctionClaimPermission{}).Where(
		"id = ? AND function_claim_id = ?",
		req.FunctionClaimPermissionID,
		req.FunctionClaimID,
	).First(&functionClaimPermission).Error
	if err != nil {
		return err
	}

	err = receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		var userFunctionAuthorize entity.SUserFunctionAuthorize

		// check if user already have function claim
		err = tx.Model(entity.SUserFunctionAuthorize{}).Where(
			"user_id = ? AND function_claim_id = ?",
			req.UserID,
			req.FunctionClaimID,
		).Delete(&userFunctionAuthorize).Error
		if err != nil {
			return errors.New("can't update user authorize")
		}

		// create a new one if user not have function claim
		userFunctionAuthorize = entity.SUserFunctionAuthorize{
			UserID:                    user.ID,
			FunctionClaimID:           functionClaim.ID,
			FunctionClaimPermissionID: functionClaimPermission.ID,
		}
		err = tx.Create(&userFunctionAuthorize).Error
		if err != nil {
			log.Error("UserEntityRepository.UpdateUserAuthorize: " + err.Error())
			return errors.New("failed to create user authorize")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (receiver *UserEntityRepository) DeleteUserAuthorize(req request.DeleteUserAuthorizeRequest) error {
	// check if user exist
	_, err := receiver.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	// check if function claim exist
	var functionClaim entity.SFunctionClaim
	err = receiver.DBConn.Model(entity.SFunctionClaim{}).Where("id = ?", req.FunctionClaimID).First(&functionClaim).Error
	if err != nil {
		return err
	}

	// delete user authorize
	err = receiver.DBConn.Where("user_id = ? AND function_claim_id = ?", req.UserID, req.FunctionClaimID).Delete(&entity.SUserFunctionAuthorize{}).Error
	if err != nil {
		log.Error("UserEntityRepository.DeleteUserAuthorize: " + err.Error())
		return errors.New("failed to delete user authorize")
	}

	return nil
}

func (receiver *UserEntityRepository) BlockUser(userID string) error {
	// check if user exist
	user, err := receiver.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return err
	}

	blocked := !user.IsBlocked
	user.IsBlocked = blocked

	if blocked {
		user.BlockedAt = time.Now()
	} else {
		user.BlockedAt = time.Time{}
	}

	// update user block
	err = receiver.DBConn.Where("id = ?", userID).Save(&user).Error
	if err != nil {
		log.Error("UserEntityRepository.BlockUser: " + err.Error())
		return errors.New("failed to block user")
	}

	return nil
}

func (receiver *UserEntityRepository) GetByID(req request.GetUserEntityByIDRequest) (*entity.SUserEntity, error) {
	var user entity.SUserEntity
	err := receiver.DBConn.
		Preload("Roles").
		Preload("Devices").
		Preload("Organizations").
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
		Preload("Devices").
		Preload("Organizations").
		Where("username = ?", req.Username).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		log.Error("UserEntityRepository.GetByUsername: " + err.Error())
		return nil, errors.New("failed to get user")
	}

	return &user, nil
}

func (receiver *UserEntityRepository) GetUserDeviceByID(deviceID string) (*[]entity.SUserDevices, error) {
	var userDevices []entity.SUserDevices
	err := receiver.DBConn.Model(&entity.SUserDevices{}).
		Where("device_id = ?", deviceID).
		Find(&userDevices).Error

	if err != nil {
		log.Error("UserEntityRepository.GetUserDeviceByID: " + err.Error())
		return nil, errors.New("failed to get user")
	}

	return &userDevices, nil
}

func (receiver *UserEntityRepository) CreateUser(req request.CreateUserEntityRequest) error {
	tx := receiver.DBConn.Begin()
	user, err := receiver.GetByUsername(request.GetUserEntityByUsernameRequest{Username: req.Username})

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
	}
	if user != nil {
		return errors.New("user already exist")
	}

	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		log.Error("UserRepository.CreateUser: " + err.Error())
		tx.Rollback()
		return errors.New("failed to create user " + req.Username)
	}

	var setting entity.SSetting
	err = tx.Table("s_setting").Where("type = ?", value.SettingTypeSignUpPresetValue1).First(&setting).Error

	if err != nil {
		log.Error("UserRepository.CreateUser: " + err.Error())
		tx.Rollback()
		return errors.New("failed to create user " + req.Username)
	}

	var signupSetting SignUpFormSetting
	err = json.Unmarshal(setting.Settings, &signupSetting)

	if err != nil {
		log.Error("UserRepository.CreateUser: " + err.Error())
		tx.Rollback()
		return errors.New("failed to create user " + req.Username)
	}

	var userReq = entity.SUserEntity{
		Username: req.Username,
		Nickname: req.Nickname,
		Fullname: signupSetting.SpreadSheetID,
		Birthday: birthday,
		Password: req.Password,
	}
	err = tx.Create(&userReq).Error

	if err != nil {
		log.Error("UserRepository.CreateUser: " + err.Error())
		tx.Rollback()
		return errors.New("failed to create user " + req.Username)
	}

	var organization entity.SOrganization
	err = tx.Model(&entity.SOrganization{}).
		Where("organization_name = 'HOME'").
		Attrs(entity.SOrganization{
			OrganizationName: "HOME",
			Password:         "123",
		}).
		FirstOrCreate(&organization).Error
	if err != nil {
		tx.Rollback()
		log.Error("UserRepository.CreateUser: " + err.Error())
		return fmt.Errorf("failed to link user with default organization")
	}

	err = tx.Table("s_user_organizations").Create(map[string]interface{}{
		"user_id":         userReq.ID.String(),
		"organization_id": organization.ID,
	}).Error
	if err != nil {
		tx.Rollback()
		log.Error("OrganizationRepository.UserJoinOrganization: " + err.Error())
		return errors.New("failed to assign user for organization")
	}

	if req.Role != nil {
		tx.Commit()
		err = receiver.UpdateUserRole(request.UpdateUserRoleRequest{UserID: userReq.ID.String(), Roles: []string{*req.Role}})
		if err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Warnf("Attempted to commit a transaction that is already committed: %v", err)
	}

	return nil
}

func (receiver *UserEntityRepository) CreateChildForParent(parentID string, req request.CreateChildForParentRequest) error {
	tx := receiver.DBConn.Begin()
	child, err := receiver.GetByUsername(request.GetUserEntityByUsernameRequest{Username: req.Username})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if child == nil {
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			log.Error("UserRepository.CreateChildForParent: " + err.Error())
			tx.Rollback()
			return errors.New("failed to create child " + req.Username)
		}

		var setting entity.SSetting
		err = tx.Table("s_setting").Where("type = ?", value.SettingTypeSignUpPresetValue1).First(&setting).Error

		if err != nil {
			log.Error("UserRepository.CreateChildForParent: " + err.Error())
			tx.Rollback()
			return errors.New("failed to create child " + req.Username)
		}

		var signupSetting SignUpFormSetting
		err = json.Unmarshal(setting.Settings, &signupSetting)

		if err != nil {
			log.Error("UserRepository.CreateChildForParent: " + err.Error())
			tx.Rollback()
			return errors.New("failed to create child " + req.Username)
		}

		if req.Fullname == "" {
			req.Fullname = signupSetting.SpreadSheetID
		}

		child = &entity.SUserEntity{
			Username: req.Username,
			Nickname: req.Nickname,
			Fullname: req.Fullname,
			Birthday: birthday,
			Password: "123",
		}
		err = tx.Create(child).Error

		if err != nil {
			log.Error("UserRepository.CreateChildForParent: " + err.Error())
			tx.Rollback()
			return errors.New("failed to create child " + req.Username)
		}

		var organization entity.SOrganization
		err = tx.Model(&entity.SOrganization{}).
			Where("organization_name = 'SENBOX WAITLIST'").
			Attrs(entity.SOrganization{
				OrganizationName: "SENBOX WAITLIST",
				Password:         "123",
			}).
			FirstOrCreate(&organization).Error
		if err != nil {
			tx.Rollback()
			log.Error("UserRepository.CreateChildForParent: " + err.Error())
			return fmt.Errorf("failed to link child with default organization")
		}

		err = tx.Table("s_user_organizations").Create(map[string]interface{}{
			"user_id":         child.ID.String(),
			"organization_id": organization.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			log.Error("OrganizationRepository.UserJoinOrganization: " + err.Error())
			return errors.New("failed to assign child for organization")
		}

		tx.Commit()
		err = receiver.UpdateUserRole(request.UpdateUserRoleRequest{UserID: child.ID.String(), Roles: []string{"Child"}})
		if err != nil {
			return err
		}

		if err := tx.Commit().Error; err != nil {
			log.Warnf("Attempted to commit a transaction that is already committed: %v", err)
		}
	}

	// check if child is already assign for parent
	var parentChild entity.SUserParentChild
	err = receiver.DBConn.Model(&entity.SUserParentChild{}).Where("parent_id = ? AND child_id = ?", parentID, child.ID).First(&parentChild).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// assign child for parent
	parentChild = entity.SUserParentChild{
		ParentID: uuid.MustParse(parentID),
		ChildID:  child.ID,
	}
	err = receiver.DBConn.Create(&parentChild).Error
	if err != nil {
		log.Error("UserRepository.CreateChildForParent: " + err.Error())
		return errors.New("failed to assign child for parent")
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUser(req request.UpdateUserEntityRequest) error {
	if req.Devices != nil {
		err := receiver.UpdateUserDevice(request.UpdateUserDeviceRequest{UserID: req.ID, Devices: *req.Devices})
		if err != nil {
			return err
		}
	}

	roles := make([]string, 0)
	if req.Roles != nil {
		for _, roleName := range *req.Roles {
			var role entity.SRole
			err := receiver.DBConn.Model(&entity.SRole{}).Where("role = ?", roleName).Find(&role).Error
			if err != nil {
				log.Error("UserRepository.CreateUser: " + err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("role not exist")
				}

				return errors.New("failed to get role")
			}

			roles = append(roles, role.Role.String())
		}
	}

	err := receiver.UpdateUserRole(request.UpdateUserRoleRequest{UserID: req.ID, Roles: roles})
	if err != nil {
		return err
	}

	updatePayload := map[string]interface{}{}
	updatePayload["username"] = req.Username

	if req.Nickname != nil {
		updatePayload["nickname"] = *req.Nickname
	}

	if req.Fullname != nil {
		updatePayload["fullname"] = *req.Fullname
	}

	if req.Phone != nil {
		updatePayload["phone"] = *req.Phone
	}

	if req.Email != nil {
		updatePayload["email"] = *req.Email
	}

	resultUpdate := receiver.DBConn.Model(&entity.SUserEntity{}).Where("id = ?", req.ID).
		Updates(updatePayload)

	if resultUpdate.Error != nil {
		log.Errorf("UserEntityRepository.UpdateUser: %v", resultUpdate.Error)
		return errors.New("failed to update user")
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUserDevice(req request.UpdateUserDeviceRequest) error {
	user, err := receiver.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})

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

	for _, deviceID := range req.Devices {
		// check if device is not exist
		var device entity.SDevice
		err := receiver.DBConn.Table("s_device").Where("id = ?", deviceID).First(&device).Error

		if err != nil {
			log.Error("UserEntityRepository.UpdateUser: " + err.Error())
			return errors.New("device not exist")
		}

		userDeviceResult := receiver.DBConn.Create(&entity.SUserDevices{
			UserID:   user.ID,
			DeviceID: device.ID,
		})
		if userDeviceResult.Error != nil {
			log.Errorf("UserEntityRepository.UpdateUser: %v", userDeviceResult.Error)
			return errors.New("failed to create user device")
		}
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUserRole(req request.UpdateUserRoleRequest) error {
	user, err := receiver.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})

	if err != nil {
		log.Error("UserEntityRepository.UpdateUserRole: " + err.Error())
		return errors.New("failed to get user")
	}

	removeResult := receiver.DBConn.Exec("DELETE FROM s_user_roles WHERE user_id = ?", user.ID)

	if removeResult.Error != nil {
		log.Error("UserEntityRepository.UpdateUserRole: " + removeResult.Error.Error())
		return errors.New("failed to remove user role")
	}

	for _, r := range req.Roles {
		var role entity.SRole
		err = receiver.DBConn.Model(&entity.SRole{}).Where("role = ?", r).First(&role).Error

		if err != nil {
			log.Error("UserEntityRepository.UpdateUserRole: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("role doesn't exist")
			}
			return errors.New("failed to get role")
		}

		// check if role already assigned for user
		var userRole entity.SUserRoles
		err = receiver.DBConn.Model(&entity.SUserRoles{}).Where("user_id = ? AND role_id = ?", user.ID, role.ID).First(&userRole).Error
		if err == nil {
			continue
		}

		result := receiver.DBConn.Create(&entity.SUserRoles{
			UserID: user.ID,
			RoleID: role.ID,
		})

		if result.Error != nil {
			log.Error("UserEntityRepository.UpdateUserRole: " + result.Error.Error())
			return errors.New("failed to assign user role")
		}
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUserRoleClaimPermission(req request.UpdateUserRoleClaimPermissionRequest) error {
	user, err := receiver.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})

	if err != nil {
		log.Error("UserEntityRepository.UpdateUserRoleClaimPermission: " + err.Error())
		return errors.New("failed to get user")
	}

	tx := receiver.DBConn.Begin()

	removeResult := tx.Exec("DELETE FROM s_user_roles WHERE user_id = ? AND role", user.ID)

	if removeResult.Error != nil {
		log.Error("UserEntityRepository.UpdateUserRoleClaimPermission: " + removeResult.Error.Error())
		return errors.New("failed to remove user policy")
	}

	return nil
}

func (receiver *UserEntityRepository) UpdateUserAvatar(userID, key, url string) error {
	err := receiver.DBConn.Model(&entity.SUserEntity{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"avatar":     key,
			"avatar_url": url,
		}).Error
	if err != nil {
		log.Error("UserEntityRepository.UpdateUserAvatar: " + err.Error())
		return errors.New("failed to update user avatar")
	}

	return nil
}

func (receiver *UserEntityRepository) GetAllPreRegisterUser() ([]*entity.SPreRegister, error) {
	var registers []*entity.SPreRegister
	err := receiver.DBConn.Model(&entity.SPreRegister{}).Find(&registers).Error
	if err != nil {
		log.Error("UserEntityRepository.GetAllPreRegisterUser: " + err.Error())
		return nil, errors.New("failed to fetch user pre registers")
	}

	return registers, nil
}

func (receiver *UserEntityRepository) GetPreRegisterUserByEmail(email string) (*entity.SPreRegister, error) {
	var register *entity.SPreRegister
	err := receiver.DBConn.Model(&entity.SPreRegister{}).
		Where("email = ?", email).
		First(&register).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("registered email not found")
		}

		log.Error("UserEntityRepository.GetPreRegisterUserByEmail: " + err.Error())
		return nil, errors.New("failed to fetch user pre register")
	}

	return register, nil
}

func (receiver *UserEntityRepository) CreatePreRegisterUser(email string) error {
	// check if email already registered
	var count int64
	err := receiver.DBConn.Model(&entity.SPreRegister{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("UserEntityRepository.CreatePreRegisterUser: " + err.Error())
			return errors.New("failed to fetch user pre register")
		}
	}
	if count > 0 {
		return errors.New("email already registered")
	}

	// create register
	err = receiver.DBConn.Create(&entity.SPreRegister{
		Email: email,
	}).Error
	if err != nil {
		log.Error("UserEntityRepository.CreatePreRegisterUser: " + err.Error())
		return errors.New("failed to create user pre register")
	}

	return err
}

// Teacher

func (receiver *UserEntityRepository) GetAllTeacherFormApplication() ([]*entity.STeacherFormApplication, error) {
	var forms []*entity.STeacherFormApplication
	err := receiver.DBConn.Model(&entity.STeacherFormApplication{}).Find(&forms).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAllTeacherFormApplication: " + err.Error())
		return nil, errors.New("failed to get all application form")
	}

	return forms, nil
}

func (receiver *UserEntityRepository) GetTeacherFormApplicationByID(applicationID int64) (*entity.STeacherFormApplication, error) {
	var form entity.STeacherFormApplication
	err := receiver.DBConn.Model(&entity.STeacherFormApplication{}).
		Where("id = ?", applicationID).
		Preload("User").
		Preload("Organization").
		First(&form).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAllTeacherFormApplication: " + err.Error())
		return nil, errors.New("failed to get application form")
	}

	return &form, nil
}

func (receiver *UserEntityRepository) ApproveTeacherFormApplication(applicationID int64) error {
	form, err := receiver.GetTeacherFormApplicationByID(applicationID)
	if err != nil {
		log.Error("UserEntityRepository.ApproveTeacherFormApplication: " + err.Error())
		return err
	}

	if !form.ApprovedAt.IsZero() {
		return errors.New("teacher has already been approved")
	}

	tx := receiver.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update application status
	err = tx.Model(&entity.STeacherFormApplication{}).
		Where("id = ?", applicationID).
		Updates(map[string]interface{}{
			"status":      value.Approved,
			"approved_at": time.Now(),
		}).Error
	if err != nil {
		tx.Rollback()
		log.Error("Update application status failed: " + err.Error())
		return errors.New("failed to approve teacher")
	}

	tx.Rollback()
	return errors.New("this function is not implemented yet")

	//// Create organization
	//organization := entity.STeacherFormApplication{
	//	TeacheranizationName: form.TeacheranizationName,
	//}
	//err = tx.Create(&organization).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Create organization failed: " + err.Error())
	//	return errors.New("failed to create organization")
	//}
	//
	//// Create user organization mapping (Manager)
	//userTeacher := entity.SUserTeacher{
	//	UserID:             form.UserID,
	//	TeacheranizationID: organization.ID,
	//	UserNickName:       "Manager",
	//	IsManager:          true,
	//}
	//err = tx.Create(&userTeacher).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Create user organization mapping failed: " + err.Error())
	//	return errors.New("failed to create user-organization relationship")
	//}
	//
	//// Add user to organization (if needed)
	//err = tx.Table("s_user_organizations").Create(map[string]interface{}{
	//	"user_id":         form.UserID.String(),
	//	"organization_id": organization.ID,
	//}).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("insert into s_user_organizations failed: " + err.Error())
	//	return errors.New("failed to join organization")
	//}
	//
	//// Assign admin role to user
	//userRole := entity.SUserRoles{
	//	UserID: form.UserID,
	//	RoleID: 5, // admin
	//}
	//err = tx.Create(&userRole).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Assign admin role failed: " + err.Error())
	//	return errors.New("failed to assign admin role to user")
	//}
	//
	//// Commit the transaction
	//if err := tx.Commit().Error; err != nil {
	//	log.Error("Transaction commit failed: " + err.Error())
	//	return errors.New("failed to commit transaction")
	//}
	//
	//return nil
}

func (receiver *UserEntityRepository) BlockTeacherFormApplication(applicationID int64) error {
	form, err := receiver.GetTeacherFormApplicationByID(applicationID)

	if err != nil {
		log.Error("UserEntityRepository.ApproveTeacherFromApplication: " + err.Error())
		return errors.New("failed to get application form")
	}

	if !form.ApprovedAt.IsZero() {
		return errors.New("teacher has not been approved")
	}

	err = receiver.DBConn.Model(&entity.STeacherFormApplication{}).Where("id = ?", applicationID).
		Updates(map[string]interface{}{"status": value.Blocked}).Error

	if err != nil {
		log.Error("UserEntityRepository.BlockTeacherFromApplication: " + err.Error())
		return errors.New("failed to block teacher")
	}

	return nil
}

func (receiver *UserEntityRepository) CreateTeacherFormApplication(req request.CreateTeacherFormApplicationRequest) error {
	result := receiver.DBConn.Create(&entity.STeacherFormApplication{
		UserID:         uuid.MustParse(req.UserID),
		OrganizationID: uuid.MustParse(req.OrganizationID),
	})

	if result.Error != nil {
		log.Error("UserEntityRepository.CreateTeacherFormApplication: " + result.Error.Error())
		return errors.New("failed to create teacher form application")
	}

	return nil
}

// Staff

func (receiver *UserEntityRepository) GetAllStaffFormApplication() ([]*entity.SStaffFormApplication, error) {
	var forms []*entity.SStaffFormApplication
	err := receiver.DBConn.Model(&entity.SStaffFormApplication{}).Find(&forms).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAllStaffFormApplication: " + err.Error())
		return nil, errors.New("failed to get all application form")
	}

	return forms, nil
}

func (receiver *UserEntityRepository) GetStaffFormApplicationByID(applicationID int64) (*entity.SStaffFormApplication, error) {
	var form entity.SStaffFormApplication
	err := receiver.DBConn.Model(&entity.SStaffFormApplication{}).
		Where("id = ?", applicationID).
		Preload("User").
		Preload("Organization").
		First(&form).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAllStaffFormApplication: " + err.Error())
		return nil, errors.New("failed to get application form")
	}

	return &form, nil
}

func (receiver *UserEntityRepository) ApproveStaffFormApplication(applicationID int64) error {
	form, err := receiver.GetStaffFormApplicationByID(applicationID)
	if err != nil {
		log.Error("UserEntityRepository.ApproveStaffFormApplication: " + err.Error())
		return err
	}

	if !form.ApprovedAt.IsZero() {
		return errors.New("staff has already been approved")
	}

	tx := receiver.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update application status
	err = tx.Model(&entity.SStaffFormApplication{}).
		Where("id = ?", applicationID).
		Updates(map[string]interface{}{
			"status":      value.Approved,
			"approved_at": time.Now(),
		}).Error
	if err != nil {
		tx.Rollback()
		log.Error("Update application status failed: " + err.Error())
		return errors.New("failed to approve staff")
	}

	tx.Rollback()
	return errors.New("this function is not implemented yet")

	//// Create organization
	//organization := entity.SStaffFormApplication{
	//	StaffanizationName: form.StaffanizationName,
	//}
	//err = tx.Create(&organization).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Create organization failed: " + err.Error())
	//	return errors.New("failed to create organization")
	//}
	//
	//// Create user organization mapping (Manager)
	//userStaff := entity.SUserStaff{
	//	UserID:             form.UserID,
	//	StaffanizationID: organization.ID,
	//	UserNickName:       "Manager",
	//	IsManager:          true,
	//}
	//err = tx.Create(&userStaff).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Create user organization mapping failed: " + err.Error())
	//	return errors.New("failed to create user-organization relationship")
	//}
	//
	//// Add user to organization (if needed)
	//err = tx.Table("s_user_organizations").Create(map[string]interface{}{
	//	"user_id":         form.UserID.String(),
	//	"organization_id": organization.ID,
	//}).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("insert into s_user_organizations failed: " + err.Error())
	//	return errors.New("failed to join organization")
	//}
	//
	//// Assign admin role to user
	//userRole := entity.SUserRoles{
	//	UserID: form.UserID,
	//	RoleID: 5, // admin
	//}
	//err = tx.Create(&userRole).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Assign admin role failed: " + err.Error())
	//	return errors.New("failed to assign admin role to user")
	//}
	//
	//// Commit the transaction
	//if err := tx.Commit().Error; err != nil {
	//	log.Error("Transaction commit failed: " + err.Error())
	//	return errors.New("failed to commit transaction")
	//}
	//
	//return nil
}

func (receiver *UserEntityRepository) BlockStaffFormApplication(applicationID int64) error {
	form, err := receiver.GetStaffFormApplicationByID(applicationID)

	if err != nil {
		log.Error("UserEntityRepository.ApproveStaffFromApplication: " + err.Error())
		return errors.New("failed to get application form")
	}

	if !form.ApprovedAt.IsZero() {
		return errors.New("staff has not been approved")
	}

	err = receiver.DBConn.Model(&entity.SStaffFormApplication{}).Where("id = ?", applicationID).
		Updates(map[string]interface{}{"status": value.Blocked}).Error

	if err != nil {
		log.Error("UserEntityRepository.BlockStaffFromApplication: " + err.Error())
		return errors.New("failed to block staff")
	}

	return nil
}

func (receiver *UserEntityRepository) CreateStaffFormApplication(req request.CreateStaffFormApplicationRequest) error {
	result := receiver.DBConn.Create(&entity.SStaffFormApplication{
		UserID:         uuid.MustParse(req.UserID),
		OrganizationID: uuid.MustParse(req.OrganizationID),
	})

	if result.Error != nil {
		log.Error("UserEntityRepository.CreateStaffFormApplication: " + result.Error.Error())
		return errors.New("failed to create staff form application")
	}

	return nil
}

// Student

func (receiver *UserEntityRepository) GetAllStudentFormApplication() ([]*entity.SStudentFormApplication, error) {
	var forms []*entity.SStudentFormApplication
	err := receiver.DBConn.Model(&entity.SStudentFormApplication{}).Find(&forms).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAllStudentFormApplication: " + err.Error())
		return nil, errors.New("failed to get all application form")
	}

	return forms, nil
}

func (receiver *UserEntityRepository) GetStudentFormApplicationByID(applicationID int64) (*entity.SStudentFormApplication, error) {
	var form entity.SStudentFormApplication
	err := receiver.DBConn.Model(&entity.SStudentFormApplication{}).
		Where("id = ?", applicationID).
		Preload("User").
		Preload("Organization").
		First(&form).Error

	if err != nil {
		log.Error("UserEntityRepository.GetAllStudentFormApplication: " + err.Error())
		return nil, errors.New("failed to get application form")
	}

	return &form, nil
}

func (receiver *UserEntityRepository) ApproveStudentFormApplication(applicationID int64) error {
	form, err := receiver.GetStudentFormApplicationByID(applicationID)
	if err != nil {
		log.Error("UserEntityRepository.ApproveStudentFormApplication: " + err.Error())
		return err
	}

	if !form.ApprovedAt.IsZero() {
		return errors.New("student has already been approved")
	}

	tx := receiver.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update application status
	err = tx.Model(&entity.SStudentFormApplication{}).
		Where("id = ?", applicationID).
		Updates(map[string]interface{}{
			"status":      value.Approved,
			"approved_at": time.Now(),
		}).Error
	if err != nil {
		tx.Rollback()
		log.Error("Update application status failed: " + err.Error())
		return errors.New("failed to approve student")
	}

	tx.Rollback()
	return errors.New("this function is not implemented yet")

	//// Create organization
	//organization := entity.SStudentFormApplication{
	//	StudentanizationName: form.StudentanizationName,
	//}
	//err = tx.Create(&organization).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Create organization failed: " + err.Error())
	//	return errors.New("failed to create organization")
	//}
	//
	//// Create user organization mapping (Manager)
	//userStudent := entity.SUserStudent{
	//	UserID:             form.UserID,
	//	StudentanizationID: organization.ID,
	//	UserNickName:       "Manager",
	//	IsManager:          true,
	//}
	//err = tx.Create(&userStudent).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Create user organization mapping failed: " + err.Error())
	//	return errors.New("failed to create user-organization relationship")
	//}
	//
	//// Add user to organization (if needed)
	//err = tx.Table("s_user_organizations").Create(map[string]interface{}{
	//	"user_id":         form.UserID.String(),
	//	"organization_id": organization.ID,
	//}).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("insert into s_user_organizations failed: " + err.Error())
	//	return errors.New("failed to join organization")
	//}
	//
	//// Assign admin role to user
	//userRole := entity.SUserRoles{
	//	UserID: form.UserID,
	//	RoleID: 5, // admin
	//}
	//err = tx.Create(&userRole).Error
	//if err != nil {
	//	tx.Rollback()
	//	log.Error("Assign admin role failed: " + err.Error())
	//	return errors.New("failed to assign admin role to user")
	//}
	//
	//// Commit the transaction
	//if err := tx.Commit().Error; err != nil {
	//	log.Error("Transaction commit failed: " + err.Error())
	//	return errors.New("failed to commit transaction")
	//}
	//
	//return nil
}

func (receiver *UserEntityRepository) BlockStudentFormApplication(applicationID int64) error {
	form, err := receiver.GetStudentFormApplicationByID(applicationID)

	if err != nil {
		log.Error("UserEntityRepository.ApproveStudentFromApplication: " + err.Error())
		return errors.New("failed to get application form")
	}

	if !form.ApprovedAt.IsZero() {
		return errors.New("student has not been approved")
	}

	err = receiver.DBConn.Model(&entity.SStudentFormApplication{}).Where("id = ?", applicationID).
		Updates(map[string]interface{}{"status": value.Blocked}).Error

	if err != nil {
		log.Error("UserEntityRepository.BlockStudentFromApplication: " + err.Error())
		return errors.New("failed to block student")
	}

	return nil
}

func (receiver *UserEntityRepository) CreateStudentFormApplication(entity *entity.SStudentFormApplication) error {
	return receiver.DBConn.Create(entity).Error
}
