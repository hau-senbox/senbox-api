package usecase

import (
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
)

type DeviceMenuUseCase struct {
	Repo          *repository.DeviceMenuRepository
	ComponentRepo *repository.ComponentRepository
	DeviceRepo    *repository.DeviceRepository
}

func NewDeviceMenuUseCase(repo *repository.DeviceMenuRepository) *DeviceMenuUseCase {
	return &DeviceMenuUseCase{Repo: repo}
}

func (uc *DeviceMenuUseCase) Create(menu *entity.SDeviceMenuV2) error {
	return uc.Repo.Create(menu)
}

func (uc *DeviceMenuUseCase) BulkCreate(menus []entity.SDeviceMenuV2) error {
	return uc.Repo.BulkCreate(menus)
}

func (uc *DeviceMenuUseCase) DeleteByDeviceID(deviceID string) error {
	return uc.Repo.DeleteByDeviceID(deviceID)
}

func (uc *DeviceMenuUseCase) GetByDeviceID(deviceID string) (response.GetDeviceMenuResponse, error) {
	// B1: Lấy thông tin device
	device, err := uc.DeviceRepo.FindDeviceByID(deviceID)
	if device == nil || err != nil {
		return response.GetDeviceMenuResponse{}, err
	}

	// B2: Lấy tất cả menu đang active
	deviceMenus, err := uc.Repo.GetByDeviceIDActive(deviceID)
	if err != nil {
		return response.GetDeviceMenuResponse{}, err
	}

	// B3: Lấy tất cả ComponentID từ deviceMenus
	componentIDs := make([]uuid.UUID, 0, len(deviceMenus))
	componentOrderMap := make(map[uuid.UUID]int)   // lưu order theo ComponentID
	componentIsShowMap := make(map[uuid.UUID]bool) // lưu is_show theo ComponentID

	for _, dm := range deviceMenus {
		componentIDs = append(componentIDs, dm.ComponentID)
		componentOrderMap[dm.ComponentID] = dm.Order
		componentIsShowMap[dm.ComponentID] = dm.IsShow
	}

	// B4: Lấy danh sách Component theo IDs
	components, err := uc.ComponentRepo.GetByIDs(componentIDs)
	if err != nil {
		return response.GetDeviceMenuResponse{}, err
	}

	// B5: Map sang ComponentResponse
	componentResponses := make([]response.ComponentResponse, 0, len(components))
	for _, comp := range components {
		componentResponses = append(componentResponses, response.ComponentResponse{
			ID:     comp.ID.String(),
			Name:   comp.Name,
			Type:   comp.Type.String(),
			Key:    comp.Key,
			Value:  helper.BuildSectionValueMenu(string(comp.Value), comp),
			Order:  componentOrderMap[comp.ID],
			IsShow: componentIsShowMap[comp.ID],
		})
	}

	return response.GetDeviceMenuResponse{
		DeviceID:   deviceID,
		DeviceName: device.DeviceName, // giả sử entity.Device có DeviceName
		Components: componentResponses,
	}, nil
}
