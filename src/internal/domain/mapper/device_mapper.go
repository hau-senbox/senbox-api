package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
)

func ToGetPersonalDeviceInfoResponse(
	device *entity.SDevice,
	deviceCode string,
	organizations []response.OrganizationDevices,
	teachers []response.TeacherResponse,
	students []response.StudentResponse,
	parents []response.ParentResponse,
	staffs []response.StaffResponse,
	valueHistories []*response.GetValuesAppCurrentResponse,
) *response.GetPersonalDeviceInfoResponse {
	return &response.GetPersonalDeviceInfoResponse{
		DeviceID:       device.ID,
		DeviceCode:     deviceCode,
		Organizations:  organizations,
		Teachers:       teachers,
		Students:       students,
		Parents:        parents,
		Staffs:         staffs,
		ValueHistories: valueHistories,
	}
}
