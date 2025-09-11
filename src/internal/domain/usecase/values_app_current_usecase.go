package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	"gorm.io/gorm"
)

type ValuesAppCurrentUseCase struct {
	Repo       *repository.ValuesAppCurrentRepository
	DeviceRepo *repository.DeviceRepository
}

func (u *ValuesAppCurrentUseCase) Upload(req request.UploadValuesAppCurrentRequest) error {
	// tìm record theo DeviceID
	exist, err := u.Repo.FindByDeviceID(req.DeviceID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// chưa có -> tạo mới
		newLog := &entity.ValuesAppCurrent{
			DeviceID: req.DeviceID,
			Value1:   req.Value1,
			Value2:   req.Value2,
			Value3:   req.Value3,
			ImageKey: req.ImageKey,
		}
		return u.Repo.Create(newLog)
	} else if err != nil {
		return err
	}

	// đã có -> update
	exist.Value1 = req.Value1
	exist.Value2 = req.Value2
	exist.Value3 = req.Value3
	exist.ImageKey = req.ImageKey

	return u.Repo.Update(exist)
}
