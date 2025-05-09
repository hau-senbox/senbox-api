package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/sheet"
	"strings"
)

type UpdateDeviceUseCase struct {
	*repository.DeviceRepository
	*repository.SettingRepository
	*sheet.Writer
}

func (receiver *UpdateDeviceUseCase) UpdateDevice(deviceId string, req request.UpdateDeviceRequest) (*entity.SDevice, error) {
	device, err := receiver.GetDeviceById(deviceId)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		device.DeviceName = *req.Name
	}

	if req.DeactivateMessage != nil {
		device.DeactivateMessage = *req.DeactivateMessage
	}

	if req.ButtonType != nil {
		bttType := value.ScreenButtonType_Scan
		if strings.ToLower(*req.ButtonType) == "list" {
			bttType = value.ScreenButtonType_List
		}
		device.ScreenButtonType = bttType
	}

	dv, err := receiver.DeviceRepository.UpdateDevice(device)
	if err != nil {
		return nil, err
	}

	return dv, nil
}

func (receiver *UpdateDeviceUseCase) UpdateDeviceV2(deviceId string, req request.UpdateDeviceRequestV2) (*entity.SDevice, error) {
	device, err := receiver.GetDeviceById(deviceId)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		device.DeviceName = *req.Name
	}
	if req.Note != nil {
		device.Note = *req.Note
	}
	if req.Status != nil {
		_, err := value.GetDeviceModeFromString(*req.Status)
		if err != nil {
			return nil, err
		}
		device.Status = value.DeviceMode(*req.Status)
	}
	if req.OutputSpreadsheetUrl != nil {
		// device.SpreadsheetId = *req.OutputSpreadsheetUrl
	}
	if req.ButtonUrl != nil {
		device.ButtonUrl = *req.ButtonUrl
	}
	if req.Message != nil {
		device.DeactivateMessage = *req.Message
	}
	if req.UserInfo != nil {
		// device.PrimaryUserInfo = reverseNormalize(req.UserInfo.UserInfo1Prefix, req.UserInfo.UserInfo1, req.UserInfo.UserInfo1ID)
		// device.SecondaryUserInfo = reverseNormalize(req.UserInfo.UserInfo2Prefix, req.UserInfo.UserInfo2, req.UserInfo.UserInfo2ID)
		// device.TertiaryUserInfo = reverseNormalize(req.UserInfo.UserInfo3Prefix, req.UserInfo.UserInfo3, req.UserInfo.UserInfo3ID)
		device.InputMode = value.InfoInputTypeBackOffice
	}

	if req.ScreenButton != nil {
		bttType := value.ScreenButtonType_Scan
		if strings.ToLower(req.ScreenButton.ButtonType) == "list" {
			bttType = value.ScreenButtonType_List
		}
		device.ScreenButtonType = bttType
	}

	dv, err := receiver.DeviceRepository.UpdateDevice(device)
	if err != nil {
		return nil, err
	}
	return dv, nil
}
