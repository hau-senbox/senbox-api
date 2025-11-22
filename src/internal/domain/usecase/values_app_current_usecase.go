package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/consulapi/gateway"
	"time"

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
	UserRepo       *repository.UserEntityRepository
	HistoriesRepo  *repository.ValuesAppHistoriesRepository
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
		// tao history
		history := &entity.ValuesAppHistories{
			DeviceID:  req.DeviceID,
			Value1:    req.Value1,
			Value2:    req.Value2,
			Value3:    req.Value3,
			ImageKey:  req.ImageKey,
			CreatedAt: time.Now(),
		}
		u.HistoriesRepo.Create(history)
		return u.Repo.Create(newLog)
	} else if err != nil {
		return err
	}

	// đã có -> update
	exist.Value1 = req.Value1
	exist.Value2 = req.Value2
	exist.Value3 = req.Value3
	exist.ImageKey = req.ImageKey

	// tạo history
	history := &entity.ValuesAppHistories{
		DeviceID:  exist.DeviceID,
		Value1:    exist.Value1,
		Value2:    exist.Value2,
		Value3:    exist.Value3,
		ImageKey:  exist.ImageKey,
		CreatedAt: time.Now(),
	}
	u.HistoriesRepo.Create(history)

	return u.Repo.Update(exist)
}

func (u *ValuesAppCurrentUseCase) GetLogedDevices4Student(ctx *gin.Context, studentID string) ([]response.LoggedDevice, error) {
	values, _ := u.Repo.GetAllDevicesByStudentID(studentID)

	logedDevices := make([]response.LoggedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LoggedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}

	return logedDevices, nil
}

func (u *ValuesAppCurrentUseCase) GetLogedDevices4Teacher(ctx *gin.Context, teacherID string) ([]response.LoggedDevice, error) {
	teacher, _ := u.TeacherRepo.GetByID(uuid.MustParse(teacherID))

	if teacher == nil {
		return nil, errors.New("teacher not found")
	}

	values, _ := u.Repo.GetAllDevicesByUserID(teacher.UserID.String())

	logedDevices := make([]response.LoggedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LoggedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}

	return logedDevices, nil
}

func (u *ValuesAppCurrentUseCase) GetLogedDevices4Staff(ctx *gin.Context, staffID string) ([]response.LoggedDevice, error) {
	staff, _ := u.StaffRepo.GetByID(uuid.MustParse(staffID))

	values, _ := u.Repo.GetAllDevicesByUserID(staff.UserID.String())

	logedDevices := make([]response.LoggedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LoggedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}
	return logedDevices, nil
}

func (u *ValuesAppCurrentUseCase) GetLogedDevices4Parent(ctx *gin.Context, parentID string) ([]response.LoggedDevice, error) {
	parent, _ := u.ParentRepo.GetByID(ctx, parentID)

	values, _ := u.Repo.GetAllDevicesByUserID(parent.UserID)

	logedDevices := make([]response.LoggedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LoggedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}
	return logedDevices, nil
}

func (u *ValuesAppCurrentUseCase) GetLogedDevices4User(ctx *gin.Context, userID string) ([]response.LoggedDevice, error) {
	user, _ := u.UserRepo.GetByID(request.GetUserEntityByIDRequest{ID: userID})

	values, _ := u.Repo.GetAllDevicesByUserID(user.ID.String())

	logedDevices := make([]response.LoggedDevice, 0, len(values))
	for _, value := range values {
		// get device personal code
		deviceCode, err := u.ProfileGateway.GetDeviceCode(ctx, value.DeviceID)
		if err != nil {
			return nil, err
		}
		logedDevices = append(logedDevices, response.LoggedDevice{
			DeviceID:   value.DeviceID,
			DeviceCode: deviceCode,
		})
	}
	return logedDevices, nil
}

func (u *ValuesAppCurrentUseCase) GetIsLoggedDevice4User(ctx *gin.Context, userID string) (bool, error) {
	values, _ := u.Repo.GetAllDevicesByUserID(userID)
	if len(values) > 0 {
		return true, nil
	}
	return false, nil
}

func (u *ValuesAppCurrentUseCase) GetIsLoggedDevice4Student(ctx *gin.Context, studentID string) (bool, error) {
	values, _ := u.Repo.GetAllDevicesByStudentID(studentID)
	if len(values) > 0 {
		return true, nil
	}
	return false, nil
}

func (u *ValuesAppCurrentUseCase) GetIsLoggedDevice4Teacher(ctx *gin.Context, teacherID string) (bool, error) {
	values, _ := u.Repo.GetAllDevicesByUserID(teacherID)
	if len(values) > 0 {
		return true, nil
	}
	return false, nil
}

func (u *ValuesAppCurrentUseCase) GetIsLoggedDevice4Staff(ctx *gin.Context, staffID string) (bool, error) {
	values, _ := u.Repo.GetAllDevicesByUserID(staffID)
	if len(values) > 0 {
		return true, nil
	}
	return false, nil
}

func (u *ValuesAppCurrentUseCase) GetIsLoggedDevice4Parent(ctx *gin.Context, parentID string) (bool, error) {
	values, _ := u.Repo.GetAllDevicesByUserID(parentID)
	if len(values) > 0 {
		return true, nil
	}
	return false, nil
}
