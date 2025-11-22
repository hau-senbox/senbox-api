package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/consulapi/gateway"
	"sen-global-api/pkg/uploader"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceUsecase struct {
	*repository.DeviceRepository
	*DeviceMenuUseCase
	*repository.ValuesAppCurrentRepository
	*GetImageUseCase
	UserEntityRepository   *repository.UserEntityRepository
	StudentRepo            *repository.StudentApplicationRepository
	ProfileGateway         gateway.ProfileGateway
	OrganizationRepo       *repository.OrganizationRepository
	TeacherRepo            *repository.TeacherApplicationRepository
	StaffRepo              *repository.StaffApplicationRepository
	ValuesAppHistoriesRepo *repository.ValuesAppHistoriesRepository
	ParentRepo             *repository.ParentRepository
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
		cleanDeviceName := strings.ReplaceAll(orgDevice.DeviceName, "NICKNAME", "")

		cleanDeviceName = strings.ReplaceAll(cleanDeviceName, "\t", " ")
		cleanDeviceName = strings.Join(strings.Fields(cleanDeviceName), " ")

		if orgDevice.DeviceNickName != "" {
			cleanDeviceName = fmt.Sprintf("%s %s", cleanDeviceName, orgDevice.DeviceNickName)
		}

		responses = append(responses, response.GetDeviceInfoResponse{
			DeviceName: cleanDeviceName,
		})
	}

	return responses, nil
}

func (receiver *DeviceUsecase) GetOrganizationDeviceInfo4Web(orgID string, deviceID string) (*response.GetDeviceInfoResponse, error) {
	// : Lấy thông tin org device
	orgDevice, err := receiver.GetOrgDeviceByDeviceIdAndOrgID(orgID, deviceID)
	if err != nil {
		return nil, err
	}

	resp := &response.GetDeviceInfoResponse{
		OrganizationID: orgDevice.OrganizationID.String(),
		DeviceName:     orgDevice.DeviceName,
		DeviceNickName: orgDevice.DeviceNickName,
		CreatedIndex:   orgDevice.CreatedIndex,
	}

	// : Lấy menu (không để lỗi menu làm fail hàm)
	if menus, err := receiver.DeviceMenuUseCase.GetByDeviceID(deviceID); err == nil {
		resp.Components = menus.Components
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Nếu là lỗi khác thì log lại để debug, nhưng không return
		log.Printf("GetDeviceMenu error for device %s: %v", deviceID, err)
	}

	currentValuesApp, _ := receiver.getValueAppCurrentByDeviceIDAndOrgID(deviceID, orgID)
	resp.CurrentAppValues = &currentValuesApp

	valueHistories, _ := receiver.getValueHisAppByDeviceIDAndOrgID(deviceID, orgID)
	resp.ValueHistories = valueHistories

	return resp, nil
}

func (receiver *DeviceUsecase) UploadDeviceNickName4Web(orgID string, deviceID string, deviceNickName string) (*entity.SOrgDevices, error) {
	// Validate input
	if orgID == "" || deviceID == "" {
		return nil, errors.New("organization_id and device_id are required")
	}

	// Update device name
	if err := receiver.DeviceRepository.UpdateDeviceNickNameByOrgIDAndDeviceID(orgID, deviceID, deviceNickName); err != nil {
		return nil, err
	}

	// Lấy lại thông tin device sau khi update
	updatedDevice, err := receiver.DeviceRepository.GetOrgDeviceByDeviceIdAndOrgID(orgID, deviceID)
	if err != nil {
		return nil, err
	}

	return updatedDevice, nil
}

func (receiver *DeviceUsecase) DeleteDeviceByOrg(orgID string, deviceID string) error {
	return receiver.DeviceRepository.DeleteDeviceByOrg(orgID, deviceID)
}

func (receiver *DeviceUsecase) GenerateDevicesCode(ctx *gin.Context) {
	devices, _ := receiver.DeviceRepository.GetAllDevices()

	for _, device := range devices {
		_, _ = receiver.ProfileGateway.GenerateDeviceCode(ctx, device.ID, device.CreatedIndex)
	}
}

func (receiver *DeviceUsecase) GetAllPersonalDevices4Web(ctx *gin.Context, deviceCode string) ([]response.GetLoggedDevicesResponse, error) {
	// get all devices
	devices, _ := receiver.DeviceRepository.GetAllDevices()

	logedDevices := make([]response.GetLoggedDevicesResponse, 0, len(devices))
	for _, device := range devices {
		// get all device trong org device
		orgDevices, _ := receiver.DeviceRepository.GetOrgsByDeviceID(device.ID)
		organizationDevices := make([]response.OrganizationDevices, 0)
		for _, orgDevice := range orgDevices {
			code, _ := receiver.ProfileGateway.GetOrganizationCode(ctx, orgDevice.OrganizationID.String())
			organizationDevices = append(organizationDevices, response.OrganizationDevices{
				OrganizationID:         orgDevice.OrganizationID.String(),
				OrganizationName:       orgDevice.Organization.OrganizationNickName,
				OrganizationDeviceCode: code,
			})
		}

		// get device personal code
		deviceCode, _ := receiver.ProfileGateway.GetDeviceCode(ctx, device.ID)

		logedDevices = append(logedDevices, response.GetLoggedDevicesResponse{
			DeviceID:            device.ID,
			DeviceCode:          deviceCode,
			OrganizationDevices: organizationDevices,
		})
	}

	// filter devices by device code
	if deviceCode != "" {
		logedDevices = lo.Filter(logedDevices, func(item response.GetLoggedDevicesResponse, _ int) bool {
			return strings.Contains(strings.ToLower(item.DeviceCode), strings.ToLower(deviceCode))
		})
	}

	return logedDevices, nil
}

func (receiver *DeviceUsecase) GetPersonalDeviceInfo4Web(ctx *gin.Context, deviceID string) (*response.GetPersonalDeviceInfoResponse, error) {
	// get device info
	device, _ := receiver.DeviceRepository.GetDeviceByID(deviceID)
	if device == nil {
		return nil, errors.New("device not found")
	}
	// get device personal code
	deviceCode, _ := receiver.ProfileGateway.GetDeviceCode(ctx, deviceID)

	values, err := receiver.ValuesAppCurrentRepository.GetAllByDeviceID(deviceID)
	if err != nil {
		return &response.GetPersonalDeviceInfoResponse{
			DeviceID:       deviceID,
			DeviceCode:     deviceCode,
			Students:       make([]response.StudentResponse, 0),
			Teachers:       make([]response.TeacherResponse, 0),
			Parents:        make([]response.ParentResponse, 0),
			Staffs:         make([]response.StaffResponse, 0),
			ValueHistories: make([]*response.GetValuesAppResponse, 0),
		}, nil
	}

	// get students by device id
	students := make([]response.StudentResponse, 0)
	for _, value := range values {
		if studentID, err := uuid.Parse(value.Value1); err == nil && studentID != uuid.Nil {
			student, _ := receiver.StudentRepo.GetByID(studentID)
			code, _ := receiver.ProfileGateway.GetStudentCode(ctx, studentID.String())
			if student != nil {
				students = append(students, response.StudentResponse{
					StudentID:   studentID.String(),
					StudentName: student.StudentName,
					Code:        code,
				})
			}
		}
	}
	// get organizationDevices
	organizationDevices := make([]response.OrganizationDevices, 0)
	for _, value := range values {
		if orgID, err := uuid.Parse(value.Value2); err == nil && orgID != uuid.Nil {
			org, _ := receiver.OrganizationRepo.GetByID(orgID.String())
			if org != nil {
				code, _ := receiver.ProfileGateway.GetOrganizationCode(ctx, orgID.String())
				organizationDevices = append(organizationDevices, response.OrganizationDevices{
					OrganizationID:         orgID.String(),
					OrganizationName:       org.OrganizationName,
					OrganizationDeviceCode: code,
				})
			}
		}
	}

	// get teachers
	teachers := make([]response.TeacherResponse, 0)
	for _, value := range values {
		// get user tu value 3
		if userID, err := uuid.Parse(value.Value3); err == nil && userID != uuid.Nil {
			user, _ := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID.String()})
			teacher, _ := receiver.TeacherRepo.GetByUserID(userID.String())
			code, _ := receiver.ProfileGateway.GetTeacherCode(ctx, teacher.ID.String())
			if teacher.ID != uuid.Nil && user != nil {
				teachers = append(teachers, response.TeacherResponse{
					TeacherID:   teacher.ID.String(),
					TeacherName: user.Nickname,
					Code:        code,
				})
			}
		}
	}

	// get staffs
	staffs := make([]response.StaffResponse, 0)
	for _, value := range values {
		if userID, err := uuid.Parse(value.Value3); err == nil && userID != uuid.Nil {
			user, _ := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID.String()})
			staff, _ := receiver.StaffRepo.GetByUserID(userID.String())
			code, _ := receiver.ProfileGateway.GetStaffCode(ctx, staff.ID.String())
			if staff.ID != uuid.Nil && user != nil {
				staffs = append(staffs, response.StaffResponse{
					StaffID:   staff.ID.String(),
					StaffName: user.Nickname,
					Code:      code,
				})
			}
		}
	}

	// get parents tu students
	parents := make([]response.ParentResponse, 0)
	if len(students) > 0 {
		for _, student := range students {
			std, _ := receiver.StudentRepo.GetByID(uuid.MustParse(student.StudentID))
			user, _ := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: std.UserID.String()})
			if user != nil {
				parent, _ := receiver.ParentRepo.GetByUserID(ctx, std.UserID.String())
				code, _ := receiver.ProfileGateway.GetParentCode(ctx, parent.ID.String())
				parents = append(parents, response.ParentResponse{
					ParentID:   parent.ID.String(),
					ParentName: parent.ParentName,
					Code:       code,
				})
			}
		}
	}
	// get value histories
	valueHistories, err := receiver.getValueHisAppByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}
	return mapper.ToGetPersonalDeviceInfoResponse(device, deviceCode, organizationDevices, teachers, students, parents, staffs, valueHistories), nil
}

func (receiver *DeviceUsecase) getValueHisAppByDeviceID(deviceID string) ([]*response.GetValuesAppResponse, error) {
	valueHistories := make([]*response.GetValuesAppResponse, 0)
	values, err := receiver.ValuesAppHistoriesRepo.GetByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}
	for _, value := range values {
		//get image url
		url, _ := receiver.GetImageUseCase.GetUrlByKey(value.ImageKey, uploader.UploadPrivate)
		userNickName := ""
		studentName := ""
		organizationName := ""

		// get student tu value 1
		if studentID, err := uuid.Parse(value.Value1); err == nil && studentID != uuid.Nil {
			student, _ := receiver.StudentRepo.GetByID(studentID)
			if student != nil {
				studentName = student.StudentName
			}
		}

		// get orginzation name tu value 2
		if orgID, err := uuid.Parse(value.Value2); err == nil && orgID != uuid.Nil {
			org, _ := receiver.OrganizationRepo.GetByID(orgID.String())
			if org != nil {
				organizationName = org.OrganizationName
			}
		}

		// get user tu value 3
		if userID, err := uuid.Parse(value.Value3); err == nil && userID != uuid.Nil {
			user, _ := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID.String()})
			if user != nil {
				userNickName = user.Nickname
			}
		}
		valueHistories = append(valueHistories, &response.GetValuesAppResponse{
			Value1:   studentName,
			Value2:   organizationName,
			Value3:   userNickName,
			ImageKey: value.ImageKey,
			ImageUrl: url,
		})
	}
	return valueHistories, nil
}

func (receiver *DeviceUsecase) getValueHisAppByDeviceIDAndOrgID(deviceID string, orgID string) ([]*response.GetValuesAppResponse, error) {
	valueHistories := make([]*response.GetValuesAppResponse, 0)
	values, err := receiver.ValuesAppHistoriesRepo.GetByDeviceIDAndOrgID(deviceID, orgID)
	if err != nil {
		return valueHistories, err
	}
	for _, value := range values {
		//get image url
		url, _ := receiver.GetImageUseCase.GetUrlByKey(value.ImageKey, uploader.UploadPrivate)
		userNickName := ""
		studentName := ""
		organizationName := ""

		// get student tu value 1
		if studentID, err := uuid.Parse(value.Value1); err == nil && studentID != uuid.Nil {
			student, _ := receiver.StudentRepo.GetByID(studentID)
			if student != nil {
				studentName = student.StudentName
			}
		}

		// get orginzation name tu value 2
		if orgID, err := uuid.Parse(value.Value2); err == nil && orgID != uuid.Nil {
			org, _ := receiver.OrganizationRepo.GetByID(orgID.String())
			if org != nil {
				organizationName = org.OrganizationName
			}
		}

		// get user tu value 3
		if userID, err := uuid.Parse(value.Value3); err == nil && userID != uuid.Nil {
			user, _ := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID.String()})
			if user != nil {
				userNickName = user.Nickname
			}
		}
		valueHistories = append(valueHistories, &response.GetValuesAppResponse{
			Value1:   studentName,
			Value2:   organizationName,
			Value3:   userNickName,
			ImageKey: value.ImageKey,
			ImageUrl: url,
		})
	}
	return valueHistories, nil
}

func (receiver *DeviceUsecase) getValueAppCurrentByDeviceIDAndOrgID(deviceID string, orgID string) (response.GetValuesAppResponse, error) {

	valueHistory, err := receiver.ValuesAppCurrentRepository.FindByDeviceIDAndOrgID(deviceID, orgID)
	if err != nil {
		return response.GetValuesAppResponse{}, err
	}
	//get image url
	url, _ := receiver.GetImageUseCase.GetUrlByKey(valueHistory.ImageKey, uploader.UploadPrivate)
	userNickName := ""
	studentName := ""
	organizationName := ""

	// get student tu value 1
	if studentID, err := uuid.Parse(valueHistory.Value1); err == nil && studentID != uuid.Nil {
		student, _ := receiver.StudentRepo.GetByID(studentID)
		if student != nil {
			studentName = student.StudentName
		}
	}

	// get orginzation name tu value 2
	if orgID, err := uuid.Parse(valueHistory.Value2); err == nil && orgID != uuid.Nil {
		org, _ := receiver.OrganizationRepo.GetByID(orgID.String())
		if org != nil {
			organizationName = org.OrganizationName
		}
	}

	// get user tu value 3
	if userID, err := uuid.Parse(valueHistory.Value3); err == nil && userID != uuid.Nil {
		user, _ := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID.String()})
		if user != nil {
			userNickName = user.Nickname
		}
	}
	return response.GetValuesAppResponse{
		Value1:   studentName,
		Value2:   organizationName,
		Value3:   userNickName,
		ImageKey: valueHistory.ImageKey,
		ImageUrl: url,
	}, nil
}
