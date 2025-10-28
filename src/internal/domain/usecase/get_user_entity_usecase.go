package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/consulapi/gateway"

	"github.com/gin-gonic/gin"
)

type GetUserEntityUseCase struct {
	*repository.UserEntityRepository
	*repository.OrganizationRepository
	*repository.ChildRepository
	*UserBlockSettingUsecase
	*UserImagesUsecase
	gateway.ProfileGateway
}

func (receiver *GetUserEntityUseCase) GetUserByID(req request.GetUserEntityByIDRequest) (*entity.SUserEntity, error) {
	return receiver.UserEntityRepository.GetByID(req)
}

func (receiver *GetUserEntityUseCase) GetUserByUsername(req request.GetUserEntityByUsernameRequest) (*entity.SUserEntity, error) {
	return receiver.UserEntityRepository.GetByUsername(req)
}

func (receiver *GetUserEntityUseCase) GetAllUsers() ([]entity.SUserEntity, error) {
	return receiver.UserEntityRepository.GetAll()
}

func (receiver *GetUserEntityUseCase) GetAllByOrganization(organizationID string) ([]entity.SUserEntity, error) {
	return receiver.UserEntityRepository.GetAllByOrganizationID(organizationID)
}

func (receiver *GetUserEntityUseCase) GetUserOrgInfo(userID, organization string) (*entity.SUserOrg, error) {
	return receiver.OrganizationRepository.GetUserOrgInfo(userID, organization)
}

func (receiver *GetUserEntityUseCase) GetAllOrgManagerInfo(organization string) (*[]entity.SUserOrg, error) {
	return receiver.OrganizationRepository.GetAllOrgManagerInfo(organization)
}

func (receiver *GetUserEntityUseCase) GetAllUserAuthorize(userID string) ([]entity.SUserFunctionAuthorize, error) {
	return receiver.UserEntityRepository.GetAllUserAuthorize(userID)
}

func (receiver *GetUserEntityUseCase) GetAllUsers4Search(ctx *gin.Context) ([]response.UserResponse, error) {
	user, err := receiver.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	var users []entity.SUserEntity

	// Nếu là SuperAdmin → lấy tất cả user (trừ SuperAdmin khác)
	if user.IsSuperAdmin() {
		allUsers, err := receiver.UserEntityRepository.GetAll()
		if err != nil {
			return nil, err
		}
		for _, u := range allUsers {
			if !u.IsSuperAdmin() {
				users = append(users, u)
			}
		}
	} else {
		// Nếu không phải SuperAdmin → lấy user theo org mà user quản lý
		if len(user.Organizations) == 0 {
			return nil, errors.New("user does not belong to any organization")
		}

		orgIDsManaged, err := user.GetManagedOrganizationIDs(receiver.UserEntityRepository.GetDB())
		if err != nil {
			return nil, err
		}
		if len(orgIDsManaged) == 0 {
			return nil, errors.New("user does not manage any organization")
		}

		users, err = receiver.UserEntityRepository.GetUsersByOrganizationIDs(orgIDsManaged)
		if err != nil {
			return nil, err
		}

		// Loại bỏ chính user đang login & SuperAdmin
		filtered := make([]entity.SUserEntity, 0, len(users))
		for _, u := range users {
			if u.ID != user.ID && !u.IsSuperAdmin() {
				filtered = append(filtered, u)
			}
		}
		users = filtered
	}

	// Map sang response
	responses := make([]response.UserResponse, 0, len(users))
	for _, u := range users {
		isDeactive, _ := receiver.UserBlockSettingUsecase.GetDeactive4User(u.ID.String())
		avatar, _ := receiver.UserImagesUsecase.GetAvtIsMain4Owner(u.ID.String(), value.OwnerRoleUser)
		code, _ := receiver.ProfileGateway.GetUserCode(ctx, u.ID.String())
		responses = append(responses, response.UserResponse{
			ID:           u.ID.String(),
			Username:     u.Username,
			Nickname:     u.Nickname,
			IsDeactive:   isDeactive,
			Avatar:       avatar,
			CreatedIndex: u.CreatedIndex,
			Code:         code,
			LanguageKeys: []string{"vietnamese", "english"},
		})
	}

	return responses, nil
}

func (receiver *GetUserEntityUseCase) GetAllParents4Search(ctx *gin.Context) ([]entity.SUserEntity, error) {
	user, err := receiver.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	// Nếu ko là SuperAdmin -> return nil
	if !user.IsSuperAdmin() {
		return nil, nil
	}

	return receiver.ChildRepository.GetAllParents()
}

func (receiver *GetUserEntityUseCase) GetCurrentUserWithOrganizations(ctx *gin.Context) (*entity.SUserEntity, error) {
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		return nil, errors.New("user ID not found in context")
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return nil, errors.New("invalid user ID format")
	}

	user, err := receiver.UserEntityRepository.GetByIDWithOrganizations(request.GetUserEntityByIDRequest{ID: userIDStr})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
