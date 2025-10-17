package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"time"
)

type PreRegisterUseCase struct {
	*repository.UserEntityRepository
	*repository.FormRepository
}

func (receiver *PreRegisterUseCase) CreatePreRegister(req request.CreatePreRegisterRequest) error {
	now := time.Now()

	entity := &entity.SPreRegister{
		Email:      req.Email,
		DeviceID:   req.DeviceID,
		DeviceName: req.DeviceName,
		FormQR:     req.FormQr,
		CreatedAt:  &now,
	}

	return receiver.UserEntityRepository.CreatePreRegisterUser(entity)
}

func (receiver *PreRegisterUseCase) GetAllPreRegister4Web() ([]*response.GetAllPreRegister4Web, error) {
	preRegisters, err := receiver.UserEntityRepository.GetAllPreRegister()
	if err != nil {
		return nil, err
	}

	var result = make([]*response.GetAllPreRegister4Web, 0)

	// tim form theo form id
	for _, preRegister := range preRegisters {
		form, _ := receiver.FormRepository.GetFormByQRCode(preRegister.FormQR)

		formName := ""
		formQr := ""
		formSheet := ""
		createdAt := ""

		if form != nil {
			formName = form.Name
			formQr = form.Note
			formSheet = form.SpreadsheetUrl
		}

		if preRegister.CreatedAt != nil {
			createdAt = preRegister.CreatedAt.Format("2006-01-02 15:04:05")
		}
		result = append(result, &response.GetAllPreRegister4Web{
			Email:      preRegister.Email,
			DeviceID:   preRegister.DeviceID,
			DeviceName: preRegister.DeviceName,
			FormName:   formName,
			FormQr:     formQr,
			FormSheet:  formSheet,
			CreatedAt:  createdAt,
		})
	}

	return result, nil
}

func (receiver *PreRegisterUseCase) GetAllPreRegister4App() ([]*response.GetAllPreRegister4App, error) {
	preRegisters, err := receiver.UserEntityRepository.GetAllPreRegister()
	if err != nil {
		return nil, err
	}

	var result = make([]*response.GetAllPreRegister4App, 0)

	// tim form theo form id
	for _, preRegister := range preRegisters {
		form, _ := receiver.FormRepository.GetFormByQRCode(preRegister.FormQR)

		formName := ""
		formQr := ""
		formSheet := ""
		createdAt := ""

		if form != nil {
			formName = form.Name
			formQr = form.Note
			formSheet = form.SpreadsheetUrl
		}

		if preRegister.CreatedAt != nil {
			createdAt = preRegister.CreatedAt.Format("2006-01-02 15:04:05")
		}

		result = append(result, &response.GetAllPreRegister4App{
			Email:      preRegister.Email,
			DeviceID:   preRegister.DeviceID,
			DeviceName: preRegister.DeviceName,
			FormName:   formName,
			FormQr:     formQr,
			FormSheet:  formSheet,
			CreatedAt:  createdAt,
		})
	}

	return result, nil
}
