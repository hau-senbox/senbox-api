package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"time"

	"github.com/samber/lo"
)

func MapUserEntityToResponseV2(userEntity entity.SUserEntity,
	isDeactive bool,
	avatars []response.Avatar,
	orgAdminResp *response.OrganizationAdmin,
) response.UserEntityResponseV2 {
	// roles
	roleListResponse := make([]response.RoleListResponseData, 0)
	for _, role := range userEntity.Roles {
		roleListResponse = append(roleListResponse, response.RoleListResponseData{
			ID:       role.ID,
			RoleName: role.Role.String(),
		})
	}

	// devices
	deviceListResponse := make([]string, 0)
	for _, device := range userEntity.Devices {
		deviceListResponse = append(deviceListResponse, device.ID)
	}

	// organizations
	organizations := lo.Map(userEntity.Organizations, func(item entity.SOrganization, index int) string {
		return item.ID.String()
	})

	return response.UserEntityResponseV2{
		ID:                userEntity.ID.String(),
		Username:          userEntity.Username,
		Fullname:          userEntity.Fullname,
		Nickname:          userEntity.Nickname,
		Phone:             userEntity.Phone,
		Email:             userEntity.Email,
		Dob:               userEntity.Birthday.Format("2006-01-02"),
		QRLogin:           userEntity.QRLogin,
		Avatar:            userEntity.Avatar,
		AvatarURL:         userEntity.AvatarURL,
		IsBlocked:         userEntity.IsBlocked,
		BlockedAt:         formatDate(userEntity.BlockedAt),
		Organization:      organizations,
		CreatedAt:         userEntity.CreatedAt.Format("2006-01-02"),
		Roles:             &roleListResponse,
		Devices:           &deviceListResponse,
		OrganizationAdmin: orgAdminResp,
		IsDeactive:        isDeactive,
		IsSuperAdmin:      userEntity.IsSuperAdmin(),
		Avatars:           avatars,
	}
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
