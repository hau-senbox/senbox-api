package usecase

import (
	"sen-global-api/internal/data/repository"

	"gorm.io/gorm"
)

type DeviceUsecase struct {
	*repository.DeviceRepository
}

func NewDeviceUsecase(db *gorm.DB) *GetDeviceByIDUseCase {
	return &GetDeviceByIDUseCase{
		DeviceRepository: &repository.DeviceRepository{DBConn: db},
	}
}

// case device chi active 1 org.
// func (receiver *DeviceUsecase) GetOrgActiveByDeviceID(deviceID string) (*entity.SDevice, error) {
// 	orgDevice, err := receiver.GetOrgByDeviceID(deviceID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return receiver.FindDeviceByID(deviceID)
// }
