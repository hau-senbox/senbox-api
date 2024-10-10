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
