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
	"strconv"
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
	UserEntityRepository *repository.UserEntityRepository
	StudentRepo          *repository.StudentApplicationRepository
	ProfileGateway       gateway.ProfileGateway
	OrganizationRepo     *repository.OrganizationRepository
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

func (receiver *DeviceUsecase) GetDeviceInfo4Web(orgID string, deviceID string) (*response.GetDeviceInfoResponse, error) {
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

	// get values app current
	if values, err := receiver.ValuesAppCurrentRepository.FindByDeviceID(deviceID); err == nil {
		//get image url
		url, _ := receiver.GetImageUseCase.GetUrlByKey(values.ImageKey, uploader.UploadPrivate)
		userNickName := ""
		studentName := ""
		organizationName := ""

		// get user nick name tu value 1
		if userID, err := uuid.Parse(values.Value1); err == nil && userID != uuid.Nil {
			user, _ := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID.String()})
			if user != nil {
				userNickName = user.Nickname
			}
		}

		// get orginzation name tu value 2
		if orgID, err := uuid.Parse(values.Value2); err == nil && orgID != uuid.Nil {
			org, _ := receiver.OrganizationRepo.GetByID(orgID.String())
			if org != nil {
				organizationName = org.OrganizationNickName
			}
		}

		// get student tu value 3
		if studentID, err := uuid.Parse(values.Value3); err == nil && studentID != uuid.Nil {
			student, _ := receiver.StudentRepo.GetByID(studentID)
			if student != nil {
				studentName = student.StudentName
			}
		}

		resp.CurrentAppValues = mapper.ToGetValuesAppCurrentResponse(userNickName, organizationName, studentName, url)
	}

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
			organizationDevices = append(organizationDevices, response.OrganizationDevices{
				OrganizationID:         orgDevice.OrganizationID.String(),
				OrganizationName:       orgDevice.Organization.OrganizationNickName,
				OrganizationDeviceCode: "O.D." + strconv.Itoa(orgDevice.CreatedIndex),
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
