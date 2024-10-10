package repository

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"time"
)

type DeviceRepository struct {
	DBConn                      *gorm.DB
	DefaultRequestPageSize      int
	DefaultOutputSpreadsheetUrl string
}

func (receiver *DeviceRepository) FindDeviceById(id string) (*entity.SDevice, error) {
	var device entity.SDevice
	err := receiver.DBConn.First(&device, "device_id = ?", id).Error
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
		err = receiver.DBConn.Raw("SELECT * FROM s_device WHERE device_name LIKE ? OR device_id LIKE ? "+
			" primary_user_info LIKE ? OR secondary_user_info LIKE ? OR tertiary_user_info LIKE ? AND row_no != ?"+
			"ORDER BY created_at DESC LIMIT ? OFFSET ?", "%"+request.Keyword+"%", "%"+request.Keyword+"%",
			"%"+request.Keyword+"%", "%"+request.Keyword+"%", "%"+request.Keyword+"%", 0,
			limit, (request.Page-1)*limit).
			Find(&devices).Error
		if err == nil {
			err = receiver.DBConn.Model(&entity.SDevice{}).
				Where("device_id LIKE ? AND row_no != ?", "%"+request.Keyword+"%", 0).
				Or("device_id LIKE ?", "%"+request.Keyword+"%").
				Or("attributes LIKE ?", "%"+request.Keyword+"%").
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
		}, errors.New("Invalid page number")
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

func (receiver *DeviceRepository) DeactivateDevice(id string, message string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("device_id = ?", id).Updates(map[string]interface{}{"status": value.Inactive, "message": message}).Error
}

func (receiver *DeviceRepository) ActivateDevice(id string, message string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("device_id = ?", id).Updates(map[string]interface{}{"status": value.Active, "message": message}).Error
}

func (receiver *DeviceRepository) CreateDevice(req request.RegisterDeviceRequest, spreadsheetId string, teacherSpreadsheetId string) (*entity.SDevice, error) {
	input, err := value.GetUserInfoInputTypeFromString(req.InputMode)
	if err != nil {
		return nil, err
	}
	if input == value.UserInfoInputTypeBackOffice {
		return nil, errors.New("invalid input mode for device client")
	}
	device := entity.SDevice{
		DeviceId:             req.DeviceUUID,
		DeviceName:           "",
		PrimaryUserInfo:      req.Primary.Fullname,
		SecondaryUserInfo:    req.Secondary.Fullname,
		TertiaryUserInfo:     req.Tertiary.Fullname,
		InputMode:            value.GetInfoInputTypeFromString(req.InputMode),
		Status:               value.DeviceModeT,
		ProfilePictureUrl:    req.ProfilePictureUrl,
		SpreadsheetId:        spreadsheetId,
		TeacherSpreadsheetId: teacherSpreadsheetId,
		AppVersion:           req.AppVersion,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	err = receiver.DBConn.Create(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, err
}

func (receiver *DeviceRepository) GetDeviceById(deviceId string) (*entity.SDevice, error) {
	var device entity.SDevice
	err := receiver.DBConn.First(&device, "device_id = ?", deviceId).Error
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

func (receiver *DeviceRepository) ReinitDevice(device entity.SDevice, req request.RegisterDeviceRequest) error {
	input, err := value.GetUserInfoInputTypeFromString(req.InputMode)
	if err != nil {
		return err
	}
	if input == value.UserInfoInputTypeBackOffice {
		return errors.New("invalid input mode for device client")
	}
	device.PrimaryUserInfo = req.Primary.Fullname
	device.SecondaryUserInfo = req.Secondary.Fullname
	device.TertiaryUserInfo = req.Tertiary.Fullname
	device.InputMode = value.GetInfoInputTypeFromString(req.InputMode)
	device.ProfilePictureUrl = req.ProfilePictureUrl
	device.AppVersion = req.AppVersion

	return receiver.DBConn.Save(&device).Error
}

func (receiver *DeviceRepository) CopyUserInfoToDevice(device entity.SDevice, req request.RegisterDeviceRequest) error {
	input, err := value.GetUserInfoInputTypeFromString(req.InputMode)
	if err != nil {
		return err
	}
	if input == value.UserInfoInputTypeBackOffice {
		return errors.New("invalid input mode for device client")
	}
	device.PrimaryUserInfo = req.Primary.Fullname
	device.SecondaryUserInfo = req.Secondary.Fullname
	device.TertiaryUserInfo = req.Tertiary.Fullname
	device.InputMode = value.GetInfoInputTypeFromString(req.InputMode)
	device.ProfilePictureUrl = req.ProfilePictureUrl
	device.AppVersion = req.AppVersion

	return receiver.DBConn.Save(&device).Error
}

func (receiver *DeviceRepository) FindByUserInfo(userInfo1 string, userInfo2 string) (*entity.SDevice, error) {
	var device entity.SDevice
	err := receiver.DBConn.First(&device, "primary_user_info LIKE ? AND secondary_user_info LIKE ?", userInfo1, userInfo2).Error
	if err != nil {
		return nil, err
	}
	return &device, err
}

func (receiver *DeviceRepository) CopyOutputFromDevice(sourceDevice entity.SDevice, targetDevice entity.SDevice, req *request.RegisterDeviceRequest) error {
	targetDevice.SpreadsheetId = sourceDevice.SpreadsheetId
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
				Columns:   []clause.Column{{Name: "device_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"device_name", "attributes", "primary_user_info", "secondary_user_info", "screen_button", "status", "profile_picture_url", "spreadsheet_id", "message", "button_url", "note", "app_version", "teacher_spreadsheet_id"}),
			}).Save(&device).Error
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
}

func (receiver *DeviceRepository) UpdateDeviceInfo(device entity.SDevice, version *string, userInfo3 *string) error {
	if version != nil {
		device.AppVersion = *version
	}
	if userInfo3 != nil {
		device.TertiaryUserInfo = *userInfo3
	}

	return receiver.DBConn.Save(&device).Error
}

func (receiver *DeviceRepository) SaveOrUpdateDevices(devices []entity.SDevice) error {
	if len(devices) == 0 {
		return nil
	}
	return receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		for _, device := range devices {
			err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "device_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"device_name", "primary_user_info", "secondary_user_info", "tertiary_user_info", "screen_button_type", "screen_button_value", "status",
					"profile_picture_url", "message", "button_url", "note", "app_version",
					"teacher_spreadsheet_id", "row_no"},
				),
			}).Create(&device).Error
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
}

func (receiver *DeviceRepository) UpdateDeviceName(deviceID string, name string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("device_id = ?", deviceID).Update("device_name", name).Error
}

func (receiver *DeviceRepository) UpdateDeviceMessage(deviceID string, message string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("device_id = ?", deviceID).Update("message", message).Error
}

func (receiver *DeviceRepository) UpdateDeviceNote(deviceID string, note string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("device_id = ?", deviceID).Update("note", note).Error
}

func (receiver *DeviceRepository) UpdateDeviceMode(deviceID string, deviceMode string) error {
	return receiver.DBConn.Model(&entity.SDevice{}).Where("device_id = ?", deviceID).Update("status", deviceMode).Error
}

func (receiver *DeviceRepository) UpdateAppSettingSpreadsheetUrl(deviceID string, spreadsheetUrl string) error {
	return receiver.
		DBConn.
		Model(&entity.SDevice{}).
		Where("device_id = ?", deviceID).
		Update("screen_button_value", spreadsheetUrl).
		Error
}

func (receiver *DeviceRepository) UpdateOutputSpreadsheetUrl(deviceID string, spreadsheetUrl string) error {
	return receiver.
		DBConn.
		Model(&entity.SDevice{}).
		Where("device_id = ?", deviceID).
		Update("spreadsheet_id", spreadsheetUrl).
		Error
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{
		DBConn:                      db,
		DefaultRequestPageSize:      12,
		DefaultOutputSpreadsheetUrl: "",
	}
}

func FindDeviceByDeviceID(deviceID string, conn *gorm.DB) (entity.SDevice, error) {
	var device entity.SDevice
	err := conn.Where("device_id = ?", deviceID).First(&device).Error

	return device, err
}
