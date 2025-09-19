package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
)

type DeviceMenuUseCase struct {
	Repo                *repository.DeviceMenuRepository
	ComponentRepo       *repository.ComponentRepository
	DeviceRepo          *repository.DeviceRepository
	LanguageSettingRepo *repository.LanguageSettingRepository
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

	// Gom components theo language_id
	menusByLang := make(map[uint][]response.ComponentResponse)
	langMap := make(map[uint]entity.LanguageSetting)

	for _, comp := range components {
		menu := response.ComponentResponse{
			ID:         comp.ID.String(),
			Name:       comp.Name,
			Type:       comp.Type.String(),
			Key:        comp.Key,
			Value:      string(comp.Value),
			Order:      componentOrderMap[comp.ID],
			IsShow:     componentIsShowMap[comp.ID],
			LanguageID: comp.LanguageID,
		}

		menusByLang[comp.LanguageID] = append(menusByLang[comp.LanguageID], menu)

		// nếu chưa có languageID trong cache -> query DB
		if _, ok := langMap[comp.LanguageID]; !ok {
			langSetting, err := uc.LanguageSettingRepo.GetByID(comp.LanguageID)
			if err != nil {
				return response.GetDeviceMenuResponse{}, err
			}
			if langSetting != nil {
				langMap[comp.LanguageID] = *langSetting
			}
		}
	}

	// Build []GetMenus4Web
	getMenus := make([]response.GetMenus4Web, 0, len(menusByLang))
	for langID, comps := range menusByLang {
		getMenus = append(getMenus, response.GetMenus4Web{
			Language: langMap[langID],
			Menus:    comps,
		})
	}

	return response.GetDeviceMenuResponse{
		DeviceID:   deviceID,
		DeviceName: device.DeviceName, // giả sử entity.Device có DeviceName
		Components: getMenus,
	}, nil
}
