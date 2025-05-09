package repository

import (
	"errors"
	"math"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceRepository struct {
	DBConn                      *gorm.DB
	DefaultRequestPageSize      int
	DefaultOutputSpreadsheetUrl string
}

func (receiver *DeviceRepository) FindDeviceById(id string) (*entity.SDevice, error) {
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

func (receiver *DeviceRepository) GetDeviceById(deviceId string) (*entity.SDevice, error) {
	var device entity.SDevice
	err := receiver.DBConn.First(&device, "id = ?", deviceId).Error
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

func (receiver *DeviceRepository) GetDevicesByUserId(userId string) (*[]entity.SDevice, error) {
	var userDevices []entity.SUserDevices
	err := receiver.DBConn.Table("s_user_devices").Where("user_id = ?", userId).Find(&userDevices).Error

	if err != nil {
		log.Error("DeviceRepository.GetDevicesByUserId: " + err.Error())
		return nil, err
	}

	var devices []entity.SDevice
	for _, userDevice := range userDevices {
		device, err := receiver.FindDeviceById(userDevice.DeviceId)

		if err != nil {
			log.Error("DeviceRepository.GetDevicesByUserId: " + err.Error())
			return nil, err
		}

		devices = append(devices, *device)
	}

	return &devices, nil
}

func (receiver *DeviceRepository) CheckUserDeviceExist(req request.RegisteringDeviceForUser) error {
	queryCheck := receiver.DBConn.Table("s_user_devices").Where("user_id = ? AND device_id = ?", req.UserId, req.DeviceId)

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

func (receiver *DeviceRepository) CheckDeviceLimitation(userId string) error {
	queryCheck := receiver.DBConn.Table("s_user_devices").Where("user_id = ?", userId)

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
	var deviceId *string
	queryCheck := receiver.DBConn.Table("s_user_devices").Where("user_id = ? AND device_id = ?", user.ID, req.DeviceUUID)

	var userDevice *entity.SUserDevices
	err := queryCheck.First(&userDevice).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	err = receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		if err := receiver.CheckDeviceLimitation(user.ID.String()); err != nil {
			return err
		}

		if userDevice.DeviceId != "" {
			// add new user_device
			userDeviceResult := receiver.DBConn.Create(&entity.SUserDevices{
				UserId:   user.ID,
				DeviceId: userDevice.DeviceId,
			})

			if userDeviceResult.Error != nil {
				log.Error("UserEntityRepository.UpdateUser: " + userDeviceResult.Error.Error())
				return errors.New("failed to register device for user")
			}

			deviceId = &userDevice.DeviceId
			return nil
		}

		var device *entity.SDevice
		// check if device already exist
		device, _ = receiver.GetDeviceById(req.DeviceUUID)
		if device == nil {
			device, err = receiver.CreateDevice(req)
			if err != nil {
				log.Error("UserEntityRepository.UpdateUser: " + err.Error())
				return errors.New("failed to create new device")
			}
		}

		// add new user_device
		userDeviceResult := receiver.DBConn.Create(&entity.SUserDevices{
			UserId:   user.ID,
			DeviceId: device.ID,
		})

		if userDeviceResult.Error != nil {
			log.Error("UserEntityRepository.UpdateUser: " + userDeviceResult.Error.Error())
			return errors.New("failed to register device for user")
		}

		deviceId = &device.ID
		return nil
	})

	if deviceId != nil {
		return deviceId, nil
	}

	if err != nil {
		return nil, err
	}

	return &userDevice.DeviceId, nil
}
