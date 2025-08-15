package repository

import (
	"errors"
	"fmt"
	"math"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceRepository struct {
	DBConn                      *gorm.DB
	DefaultRequestPageSize      int
	DefaultOutputSpreadsheetUrl string
}

func (receiver *DeviceRepository) FindDeviceByID(id string) (*entity.SDevice, error) {
	var device entity.SDevice
	err := receiver.DBConn.First(&device, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &device, err
}

func (receiver *DeviceRepository) GetDeviceList(request request.GetListDeviceRequest) ([]entity.SDevice, *response.Pagination, error) {
	var devices []entity.SDevice
	limit := receiver.DefaultRequestPageSize
	if request.Limit != 0 {
		limit = request.Limit
	}
	if request.Page <= 0 {
		request.Page = 1
	}
	var err error
	var count int64
	if request.Keyword != "" {
		err = receiver.DBConn.Raw("SELECT * FROM s_device WHERE device_name LIKE ? OR id LIKE ? AND row_no != ?"+
			"ORDER BY created_at DESC LIMIT ? OFFSET ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%", 0,
			limit, (request.Page-1)*limit).
			Find(&devices).Error
		if err == nil {
			err = receiver.DBConn.Model(&entity.SDevice{}).
				Where("id LIKE ? AND row_no != ?", "%"+request.Keyword+"%", 0).
				Or("id LIKE ?", "%"+request.Keyword+"%").
				Count(&count).Error
		}
	} else {
		err = receiver.DBConn.Where("row_no != ?", 0).Order("created_at desc").Limit(limit).Offset((request.Page - 1) * limit).Find(&devices).Error
		if err == nil {
			err = receiver.DBConn.Model(&entity.SDevice{}).Count(&count).Error
		}
	}

	if int64(request.Page) > count {
		return []entity.SDevice{}, &response.Pagination{
			Page:      request.Page,
			Limit:     limit,
			TotalPage: int(math.Ceil(float64(count) / float64(limit))),
			Total:     count,
		}, errors.New("invalid page number")
	}

	if err != nil {
		return nil, nil, err
	}
	return devices, &response.Pagination{
		Page:      request.Page,
		Limit:     limit,
		TotalPage: int(math.Ceil(float64(count) / float64(limit))),
		Total:     count,
	}, err
}

func (receiver *DeviceRepository) GetDeviceListByUserID(userID string) ([]entity.SDevice, error) {
	var devices []entity.SDevice
	err := receiver.DBConn.Raw(`
    SELECT * FROM s_device 
    WHERE id IN (
        SELECT device_id 
        FROM s_org_devices 
        WHERE user_id = ?
    )`, userID).Scan(&devices).Error
	return devices, err
}

func (receiver *DeviceRepository) GetDeviceListByOrgID(orgID string) ([]entity.SDevice, error) {
	var devices []entity.SDevice
	err := receiver.DBConn.Raw(`
    SELECT * FROM s_device 
    WHERE id IN (
        SELECT device_id 
        FROM s_org_devices 
        WHERE organization_id = ?
    )`, orgID).Scan(&devices).Error
	return devices, err
}

func (receiver *DeviceRepository) DeactivateDevice(id string, deactivateMessage string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("id = ?", id).Updates(map[string]interface{}{"status": value.Inactive, "deactivate_message": deactivateMessage}).Error
}

func (receiver *DeviceRepository) ActivateDevice(id string, deactivateMessage string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("id = ?", id).Updates(map[string]interface{}{"status": value.Active, "deactivate_message": deactivateMessage}).Error
}

func (receiver *DeviceRepository) CreateDevice(req request.RegisterDeviceRequest) (*entity.SDevice, error) {
	input, err := value.GetUserInfoInputTypeFromString(req.InputMode)
	if err != nil {
		return nil, err
	}
	if input == value.UserInfoInputTypeBackOffice {
		return nil, errors.New("invalid input mode for device client")
	}
	device := entity.SDevice{
		ID:         req.DeviceUUID,
		DeviceName: "",
		InputMode:  value.GetInfoInputTypeFromString(req.InputMode),
		Status:     value.DeviceModeT,
		AppVersion: req.AppVersion,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = receiver.DBConn.Create(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, err
}

func (receiver *DeviceRepository) GetDeviceByID(deviceID string) (*entity.SDevice, error) {
	var device entity.SDevice
	err := receiver.DBConn.First(&device, "id = ?", deviceID).Error
	if err != nil {
		return nil, err
	}
	return &device, err
}

func (receiver *DeviceRepository) UpdateDevice(device *entity.SDevice) (*entity.SDevice, error) {
	err := receiver.DBConn.Save(&device).Error
	if err != nil {
		return nil, err
	}
	return device, err
}

func (receiver *DeviceRepository) CopyUserInfoToDevice(device entity.SDevice, req request.RegisterDeviceRequest) error {
	input, err := value.GetUserInfoInputTypeFromString(req.InputMode)
	if err != nil {
		return err
	}
	if input == value.UserInfoInputTypeBackOffice {
		return errors.New("invalid input mode for device client")
	}
	device.InputMode = value.GetInfoInputTypeFromString(req.InputMode)
	device.AppVersion = req.AppVersion

	return receiver.DBConn.Save(&device).Error
}

func (receiver *DeviceRepository) CopyOutputFromDevice(sourceDevice entity.SDevice, targetDevice entity.SDevice, req *request.RegisterDeviceRequest) error {
	targetDevice.AppVersion = req.AppVersion
	return receiver.DBConn.Save(&targetDevice).Error
}

func (receiver *DeviceRepository) SaveDevices(devices []entity.SDevice) error {
	if len(devices) == 0 {
		return nil
	}
	return receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		for _, device := range devices {
			err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"device_name",
					"input_mode",
					"screen_button_type",
					"status",
					"deactivate_message",
					"button_url",
					"note",
					"app_version",
				}),
			}).Save(&device).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (receiver *DeviceRepository) SaveOrUpdateDevices(devices []entity.SDevice) error {
	if len(devices) == 0 {
		return nil
	}
	return receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		for _, device := range devices {
			err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"device_name",
					"input_mode",
					"screen_button_type",
					"status",
					"deactivate_message",
					"button_url",
					"note",
					"app_version",
					"row_no",
				}),
			}).Create(&device).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (receiver *DeviceRepository) UpdateDeviceName(deviceID string, name string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("id = ?", deviceID).Update("device_name", name).Error
}

func (receiver *DeviceRepository) UpdateDeviceDeactivateMessage(deviceID string, deactivateMessage string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("id = ?", deviceID).Update("deactivate_message", deactivateMessage).Error
}

func (receiver *DeviceRepository) UpdateDeviceNote(deviceID string, note string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("id = ?", deviceID).Update("note", note).Error
}

func (receiver *DeviceRepository) UpdateDeviceMode(deviceID string, deviceMode string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("id = ?", deviceID).Update("status", deviceMode).Error
}

// func (receiver *DeviceRepository) UpdateAppSettingSpreadsheetUrl(deviceID string, spreadsheetUrl string) error {
// 	return receiver.
// 		DBConn.
// 		Model(&entity.SDevice{}).
// 		Where("id = ?", deviceID).
// 		Update("screen_button_value", spreadsheetUrl).
// 		Error
// }

// func (receiver *DeviceRepository) UpdateOutputSpreadsheetUrl(deviceID string, spreadsheetUrl string) error {
// 	return receiver.
// 		DBConn.
// 		Model(&entity.SDevice{}).
// 		Where("id = ?", deviceID).
// 		Update("spreadsheet_id", spreadsheetUrl).
// 		Error
// }

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{
		DBConn:                      db,
		DefaultRequestPageSize:      12,
		DefaultOutputSpreadsheetUrl: "",
	}
}

func FindDeviceByDeviceID(deviceID string, conn *gorm.DB) (entity.SDevice, error) {
	var device entity.SDevice
	err := conn.Where("id = ?", deviceID).First(&device).Error

	return device, err
}

func (receiver *DeviceRepository) CheckUserDeviceExist(req request.RegisteringDeviceForUser) error {
	queryCheck := receiver.DBConn.Table("s_user_devices").Where("user_id = ? AND device_id = ?", req.UserID, req.DeviceID)

	var userDevice *entity.SUserDevices
	err := queryCheck.First(&userDevice).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		log.Error("UserEntityRepository.CheckUserDeviceExist: " + err.Error())
		return errors.New("failed to get user device")
	}

	if userDevice != nil {
		return errors.New("device already exist and assigned to user")
	}

	return nil
}

func (receiver *DeviceRepository) CheckOrgDeviceExist(req request.RegisteringDeviceForOrg) error {
	queryCheck := receiver.DBConn.Table("s_org_devices").Where("organization_id = ? AND device_id = ?", req.OrgID, req.DeviceID)

	var userDevice entity.SOrgDevices
	err := queryCheck.First(&userDevice).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("UserEntityRepository.CheckOrgDeviceExist: " + err.Error())
		return errors.New("failed to get org device")
	}

	if userDevice.DeviceID != "" {
		return errors.New("device already exist and assigned to organization")
	}

	return nil
}

func (receiver *DeviceRepository) CheckDeviceLimitation(userID string) error {
	queryCheck := receiver.DBConn.Table("s_user_devices").Where("user_id = ?", userID)

	var deviceCount int64
	err := queryCheck.Count(&deviceCount).Error

	if err != nil {
		log.Error("UserEntityRepository.CheckDeviceLimitation: " + err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return errors.New("failed to get device count")
	}

	if deviceCount > 1 {
		return errors.New("device limitation reached")
	}

	return nil
}

func (receiver *DeviceRepository) RegisteringDeviceForUser(user *entity.SUserEntity, req request.RegisterDeviceRequest) (*string, error) {
	var deviceID *string
	device, err := receiver.GetDeviceByID(req.DeviceUUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	err = receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		if err := receiver.CheckDeviceLimitation(user.ID.String()); err != nil {
			return err
		}

		if device != nil {
			// add new user_device
			userDeviceResult := receiver.DBConn.Create(&entity.SUserDevices{
				UserID:   user.ID,
				DeviceID: device.ID,
			})

			if userDeviceResult.Error != nil {
				log.Error("UserEntityRepository.RegisteringDeviceForUser: " + userDeviceResult.Error.Error())
				return errors.New("failed to register device for user")
			}

			deviceID = &device.ID
			return nil
		}

		var device *entity.SDevice
		// check if device already exist
		device, _ = receiver.GetDeviceByID(req.DeviceUUID)
		if device == nil {
			device, err = receiver.CreateDevice(req)
			if err != nil {
				log.Error("UserEntityRepository.RegisteringDeviceForUser: " + err.Error())
				return errors.New("failed to create new device")
			}
		}

		// add new user_device
		userDeviceResult := receiver.DBConn.Create(&entity.SUserDevices{
			UserID:   user.ID,
			DeviceID: device.ID,
		})

		if userDeviceResult.Error != nil {
			log.Error("UserEntityRepository.RegisteringDeviceForUser: " + userDeviceResult.Error.Error())
			return errors.New("failed to register device for user")
		}

		deviceID = &device.ID
		return nil
	})

	if err != nil {
		return nil, err
	}

	return deviceID, nil
}

func (receiver *DeviceRepository) RegisteringDeviceForOrg(org *entity.SOrganization, req request.RegisterDeviceRequest) (*string, error) {
	var deviceID *string
	device, err := receiver.GetDeviceByID(req.DeviceUUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Get device name
	listDeviceByOrg, _ := receiver.GetDeviceListByOrgID(org.ID.String())
	count := len(listDeviceByOrg) // số thiết bị hiện có
	deviceName := fmt.Sprintf("%s - DEVICE -[%d]", org.OrganizationName, count+1)

	err = receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		if device != nil {
			// add new user_device
			userDeviceResult := receiver.DBConn.Create(&entity.SOrgDevices{
				OrganizationID: org.ID,
				DeviceID:       device.ID,
				DeviceName:     deviceName,
			})

			if userDeviceResult.Error != nil {
				log.Error("UserEntityRepository.RegisteringDeviceForOrg: " + userDeviceResult.Error.Error())
				return errors.New("failed to register device for user")
			}

			deviceID = &device.ID
			return nil
		}

		var device *entity.SDevice
		// check if device already exist
		device, _ = receiver.GetDeviceByID(req.DeviceUUID)
		if device == nil {
			device, err = receiver.CreateDevice(req)
			if err != nil {
				log.Error("UserEntityRepository.RegisteringDeviceForOrg: " + err.Error())
				return errors.New("failed to create new device")
			}
		}

		// add new user_device
		userDeviceResult := receiver.DBConn.Create(&entity.SOrgDevices{
			OrganizationID: org.ID,
			DeviceID:       device.ID,
			DeviceName:     deviceName,
		})

		if userDeviceResult.Error != nil {
			log.Error("UserEntityRepository.RegisteringDeviceForOrg: " + userDeviceResult.Error.Error())
			return errors.New("failed to register device for user")
		}

		deviceID = &device.ID
		return nil
	})

	if err != nil {
		return nil, err
	}

	return deviceID, nil
}

func (r *DeviceRepository) GetOrgIDsByDeviceID(deviceID string) ([]uuid.UUID, error) {
	var orgDevices []entity.SOrgDevices
	if err := r.DBConn.
		Select("organization_id").
		Where("device_id = ?", deviceID).
		Find(&orgDevices).Error; err != nil {
		return nil, err
	}

	orgIDs := make([]uuid.UUID, 0, len(orgDevices))
	for _, od := range orgDevices {
		orgIDs = append(orgIDs, od.OrganizationID)
	}

	return orgIDs, nil
}

func (r *DeviceRepository) GetOrgByDeviceID(deviceID string) (*entity.SOrgDevices, error) {
	// Gia su chi lay 1 org theo device id (case 1 device chi active 1 org)
	var orgDevices *entity.SOrgDevices
	if err := r.DBConn.
		Where("device_id = ?", deviceID).
		First(&orgDevices).Error; err != nil {
		return nil, err
	}

	return orgDevices, nil
}

func (r *DeviceRepository) GetOrgsByDeviceID(deviceID string) ([]entity.SOrgDevices, error) {
	var orgDevices []entity.SOrgDevices
	if err := r.DBConn.
		Where("device_id = ?", deviceID).
		Find(&orgDevices).Error; err != nil {
		return nil, err
	}

	return orgDevices, nil
}

func (r *DeviceRepository) GetOrgDeviceByDeviceIdAndOrgID(orgID string, deviceID string) (*entity.SOrgDevices, error) {
	var orgDevices *entity.SOrgDevices
	if err := r.DBConn.
		Where("organization_id = ? AND device_id = ?", orgID, deviceID).
		First(&orgDevices).Error; err != nil {
		return nil, err
	}

	return orgDevices, nil
}

func (r *DeviceRepository) UpdateDeviceNameByOrgIDAndDeviceID(orgID string, deviceID string, deviceName string) error {
	var orgDevice entity.SOrgDevices
	if err := r.DBConn.
		Where("organization_id = ? AND device_id = ?", orgID, deviceID).
		First(&orgDevice).Error; err != nil {
		return err
	}

	orgDevice.DeviceName = deviceName
	return r.DBConn.Save(&orgDevice).Error
}
