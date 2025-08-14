package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"

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
		DeviceName: orgDeviceInfo.DeviceName,
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
