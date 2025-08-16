package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceUsecase struct {
	*repository.DeviceRepository
	*DeviceMenuUseCase
}

func NewDeviceUsecase(db *gorm.DB) *GetDeviceByIDUseCase {
	return &GetDeviceByIDUseCase{
		DeviceRepository: &repository.DeviceRepository{DBConn: db},
	}
}

// case device chi active 1 org.
func (receiver *DeviceUsecase) GetDeviceInfoFromOrg(deviceID string) (*response.GetDeviceInfoResponse, error) {
	orgDeviceInfo, err := receiver.GetOrgByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}

	// DeviceName info di theo org ma device dang ky
	res := &response.GetDeviceInfoResponse{
		DeviceName: orgDeviceInfo.DeviceName,
	}

	return res, nil
}

func (receiver *DeviceUsecase) GetDeviceInfoFromOrg4Admin(orgID string, deviceID string) (*response.GetDeviceInfoResponse, error) {
	orgDeviceInfo, err := receiver.GetOrgDeviceByDeviceIdAndOrgID(orgID, deviceID)
	if err != nil {
		return nil, err
	}

	// DeviceName info di theo org ma device dang ky
	res := &response.GetDeviceInfoResponse{
		DeviceName:     orgDeviceInfo.DeviceName,
		CreatedIndex:   orgDeviceInfo.CreatedIndex,
		DeviceNickName: orgDeviceInfo.DeviceNickName,
	}

	return res, nil
}

func (receiver *DeviceUsecase) GetDeviceInfoFromOrg4App(deviceID string) ([]response.GetDeviceInfoResponse, error) {
	orgDevices, err := receiver.GetOrgsByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}

	// Build list response
	responses := make([]response.GetDeviceInfoResponse, 0, len(orgDevices))
	for _, orgDevice := range orgDevices {
		responses = append(responses, response.GetDeviceInfoResponse{
			DeviceName: orgDevice.DeviceName,
			// Nếu có thêm field nào trong GetDeviceInfoResponse, map ở đây
		})
	}

	return responses, nil
}

func (receiver *DeviceUsecase) GetDeviceInfo4Web(orgID string, deviceID string) (*response.GetDeviceInfoResponse, error) {
	// B1: Lấy thông tin org device
	orgDevice, err := receiver.GetOrgDeviceByDeviceIdAndOrgID(orgID, deviceID)
	if err != nil {
		return nil, err
	}

	resp := &response.GetDeviceInfoResponse{
		OrganizationID: orgDevice.OrganizationID.String(),
		DeviceName:     orgDevice.DeviceName,
	}

	// B2: Lấy menu (không để lỗi menu làm fail hàm)
	if menus, err := receiver.DeviceMenuUseCase.GetByDeviceID(deviceID); err == nil {
		resp.Components = menus.Components
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu là lỗi khác thì log lại để debug, nhưng không return
		log.Printf("GetDeviceMenu error for device %s: %v", deviceID, err)
	}

	return resp, nil
}

func (receiver *DeviceUsecase) UploadDeviceName4Web(orgID string, deviceID string, deviceName string) (*entity.SOrgDevices, error) {
	// Validate input
	if orgID == "" || deviceID == "" {
		return nil, errors.New("organization_id and device_id are required")
	}

	// Update device name
	if err := receiver.DeviceRepository.UpdateDeviceNameByOrgIDAndDeviceID(orgID, deviceID, deviceName); err != nil {
		return nil, err
	}

	// Lấy lại thông tin device sau khi update
	updatedDevice, err := receiver.DeviceRepository.GetOrgDeviceByDeviceIdAndOrgID(orgID, deviceID)
	if err != nil {
		return nil, err
	}

	return updatedDevice, nil
}
