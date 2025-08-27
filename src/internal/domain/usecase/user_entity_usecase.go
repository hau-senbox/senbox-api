package usecase

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type UserEntityUseCase struct {
	DBConn                  *gorm.DB
	UserBlockSettingUsecase *UserBlockSettingUsecase
	UserImagesUsecase       *UserImagesUsecase
}

func (receiver *UserEntityUseCase) MapUserInfoToResponse(userEntity entity.SUserEntity) (*response.UserEntityResponseV2, error) {
	isDeactive, err := receiver.UserBlockSettingUsecase.GetDeactive4User(userEntity.ID.String())
	if err != nil {
		return nil, err
	}

	avatars, _ := receiver.UserImagesUsecase.GetAvt4Owner(userEntity.ID.String(), value.OwnerRoleUser)

	// Lấy danh sách org mà user thuộc về
	userEntity.Organizations, _ = userEntity.GetOrganizations(receiver.DBConn)
	var orgAdminResp *response.OrganizationAdmin
	if len(userEntity.Organizations) > 0 {
		// Lấy danh sách OrgID mà user là manager
		managedOrgIDs, err := userEntity.GetManagedOrganizationIDs(receiver.DBConn)
		if err != nil {
			orgAdminResp = nil
		}

		// So sánh với các org đã preload, map sang OrganizationAdmin nếu khớp
		if len(managedOrgIDs) > 0 {
			for _, org := range userEntity.Organizations {
				if lo.Contains(managedOrgIDs, org.ID.String()) {
					orgAdminResp = &response.OrganizationAdmin{
						ID:               org.ID.String(),
						OrganizationName: org.OrganizationName,
						Avatar:           org.Avatar,
						AvatarURL:        org.AvatarURL,
						Address:          org.Address,
						Description:      org.Description,
						CreatedAt:        org.CreatedAt,
						UpdatedAt:        org.UpdatedAt,
					}
					break
				}
			}
		}

	}

	userResp := mapper.MapUserEntityToResponseV2(userEntity, isDeactive, avatars, orgAdminResp)
	return &userResp, nil
}
