package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	"github.com/gin-gonic/gin"
)

type GetUserEntityUseCase struct {
	*repository.UserEntityRepository
	*repository.OrganizationRepository
	*repository.ChildRepository
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

func (receiver *GetUserEntityUseCase) GetAllUsers4Search(ctx *gin.Context) ([]entity.SUserEntity, error) {
	user, err := receiver.GetCurrentUserWithOrganizations(ctx)
	if err != nil {
		return nil, err
	}

	// Nếu là SuperAdmin → trả về tất cả user
	if user.IsSuperAdmin() {
		allUsers, err := receiver.UserEntityRepository.GetAll()
		if err != nil {
			return nil, err
		}

		result := make([]entity.SUserEntity, 0, len(allUsers))
		for _, u := range allUsers {
			if !u.IsSuperAdmin() {
				result = append(result, u)
			}
		}

		return result, nil
	}

	// Nếu không phải SuperAdmin → lấy danh sách org mà user quản lý
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

	// Lấy user thuộc các tổ chức mà user này quản lý
	users, err := receiver.UserEntityRepository.GetUsersByOrganizationIDs(orgIDsManaged)
	if err != nil {
		return nil, err
	}

	// Lọc ra chính user đang truy cập khỏi kết quả
	result := make([]entity.SUserEntity, 0, len(users))
	for _, u := range users {
		if u.ID != user.ID && !u.IsSuperAdmin() {
			result = append(result, u)
		}
	}

	return result, nil
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
