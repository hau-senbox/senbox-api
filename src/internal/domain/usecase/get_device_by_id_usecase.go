package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type GetDeviceByIDUseCase struct {
	*repository.DeviceRepository
}

func NewGetDeviceByIDUseCase(db *gorm.DB) *GetDeviceByIDUseCase {
	return &GetDeviceByIDUseCase{
		DeviceRepository: &repository.DeviceRepository{DBConn: db},
	}
}

func (receiver *GetDeviceByIDUseCase) Get(deviceID string) (*entity.SDevice, error) {
	return receiver.FindDeviceByID(deviceID)
}
