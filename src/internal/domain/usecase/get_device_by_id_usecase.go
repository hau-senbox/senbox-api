package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type GetDeviceByIdUseCase struct {
	*repository.DeviceRepository
}

func NewGetDeviceByIdUseCase(db *gorm.DB) *GetDeviceByIdUseCase {
	return &GetDeviceByIdUseCase{
		DeviceRepository: &repository.DeviceRepository{DBConn: db},
	}
}

func (receiver *GetDeviceByIdUseCase) Get(deviceId string) (*entity.SDevice, error) {
	return receiver.FindDeviceById(deviceId)
}
