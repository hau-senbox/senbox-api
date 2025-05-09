package repository

import (
	"encoding/json"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DeviceComponentValuesRepository struct {
	DBConn *gorm.DB
}

func NewDeviceComponentValuesRepository(dbConn *gorm.DB) *DeviceComponentValuesRepository {
	return &DeviceComponentValuesRepository{DBConn: dbConn}
}

func (receiver *DeviceComponentValuesRepository) GetByOrganization(req request.GetDeviceComponentValuesByOrganizationRequest) (*entity.SDeviceComponentValues, error) {
	var deviceComponentValues entity.SDeviceComponentValues
	err := receiver.DBConn.Where("organization_id = ?", req.ID).First(&deviceComponentValues).Error
	if err != nil {
		return nil, err
	}
	return &deviceComponentValues, nil
}

func (receiver *DeviceComponentValuesRepository) GetByDevice(req request.GetDeviceComponentValuesByDeviceRequest) (*entity.SDeviceComponentValues, error) {
	var deviceComponentValues entity.SDeviceComponentValues
	err := receiver.DBConn.Where("id = ? AND organization_id = ?", req.ID, nil).First(&deviceComponentValues).Error
	if err != nil {
		return nil, err
	}
	return &deviceComponentValues, nil
}

func (receiver *DeviceComponentValuesRepository) SaveByOrganization(req request.SaveDeviceComponentValuesByOrganizationRequest) error {
	setting, _ := receiver.GetByOrganization(request.GetDeviceComponentValuesByOrganizationRequest{ID: req.Organization})

	jsonSetting, err := json.Marshal(req.Settings)
	if err != nil {
		return err
	}

	if setting == nil {
		organizationID := int64(req.Organization)
		result := receiver.DBConn.Create(&entity.SDeviceComponentValues{
			Setting:   datatypes.JSON(string(jsonSetting)),
			OrganizationId: &organizationID,
		})

		return result.Error
	} else {
		setting.Setting = datatypes.JSON(string(jsonSetting))

		return receiver.DBConn.Save(&setting).Error
	}
}

func (receiver *DeviceComponentValuesRepository) SaveByDevice(req request.SaveDeviceComponentValuesByDeviceRequest) error {
	setting, err := receiver.GetByDevice(request.GetDeviceComponentValuesByDeviceRequest{ID: *req.ID})
	if err != nil {
		return err
	}

	jsonSetting, err := json.Marshal(req.Settings)
	if err != nil {
		return err
	}

	if setting == nil {
		result := receiver.DBConn.Create(&entity.SDeviceComponentValues{
			Setting: datatypes.JSON(string(jsonSetting)),
		})

		return result.Error
	} else {
		setting.Setting = datatypes.JSON(string(jsonSetting))

		return receiver.DBConn.Save(&setting).Error
	}
}
