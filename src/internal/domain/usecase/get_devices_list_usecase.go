package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
)

type GetDeviceListUseCase struct {
	*repository.DeviceRepository
}

func (receiver *GetDeviceListUseCase) GetDeviceList(request request.GetListDeviceRequest) ([]entity.SDevice, *response.Pagination, error) {
	return receiver.DeviceRepository.GetDeviceList(request)
}

func (receiver *GetDeviceListUseCase) GetDeviceListByUserID(userID string) ([]entity.SDevice, error) {
	return receiver.DeviceRepository.GetDeviceListByUserID(userID)
}

func (receiver *GetDeviceListUseCase) GetDeviceListByOrgID(orgID string) ([]entity.SDevice, error) {
	return receiver.DeviceRepository.GetDeviceListByOrgID(orgID)
}
