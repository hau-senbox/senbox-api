package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ValuesAppCurrentUseCase struct {
	Repo           *repository.ValuesAppCurrentRepository
	DeviceRepo     *repository.DeviceRepository
	TeacherRepo    *repository.TeacherApplicationRepository
	StaffRepo      *repository.StaffApplicationRepository
	ParentRepo     *repository.ParentRepository
	ProfileGateway gateway.ProfileGateway
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

func (u *ValuesAppCurrentUseCase) GetLogedDevices4Student(ctx *gin.Context, studentID string) ([]response.LogedDevice, error) {
	values, _ := u.Repo.GetAllDevicesByStudentID(studentID)

	logedDevices := make([]response.LogedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LogedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}

	return logedDevices, nil
}

func (u *ValuesAppCurrentUseCase) GetLogedDevices4Teacher(ctx *gin.Context, teacherID string) ([]response.LogedDevice, error) {
	teacher, _ := u.TeacherRepo.GetByID(uuid.MustParse(teacherID))

	values, _ := u.Repo.GetAllDevicesByUserID(teacher.UserID.String())

	logedDevices := make([]response.LogedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LogedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}

	return logedDevices, nil
}

func (u *ValuesAppCurrentUseCase) GetLogedDevices4Staff(ctx *gin.Context, staffID string) ([]response.LogedDevice, error) {
	staff, _ := u.StaffRepo.GetByID(uuid.MustParse(staffID))

	values, _ := u.Repo.GetAllDevicesByUserID(staff.UserID.String())

	logedDevices := make([]response.LogedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LogedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}
	return logedDevices, nil
}

func (u *ValuesAppCurrentUseCase) GetLogedDevices4Parent(ctx *gin.Context, parentID string) ([]response.LogedDevice, error) {
	parent, _ := u.ParentRepo.GetByID(ctx, parentID)

	values, _ := u.Repo.GetAllDevicesByUserID(parent.UserID)

	logedDevices := make([]response.LogedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LogedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}
	return logedDevices, nil
}
