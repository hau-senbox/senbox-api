package usecase

import (
	"errors"
	"github.com/samber/lo"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
)

type GetMenuUseCase struct {
	MenuRepository         *repository.MenuRepository
	UserEntityRepository   *repository.UserEntityRepository
	OrganizationRepository *repository.OrganizationRepository
}

func (receiver *GetMenuUseCase) GetSuperAdminMenu() ([]menu.SuperAdminMenu, error) {
	return receiver.MenuRepository.GetSuperAdminMenu()
}

func (receiver *GetMenuUseCase) GetOrgMenu(orgID string) ([]menu.OrgMenu, error) {
	org, err := receiver.OrganizationRepository.GetByID(orgID)
	if err != nil {
		return nil, err
	}

	return receiver.MenuRepository.GetOrgMenu(org.ID.String())
}

func (receiver *GetMenuUseCase) GetStudentMenu(userID string) ([]menu.UserMenu, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	present := lo.ContainsBy(user.Roles, func(role entity.SRole) bool {
		return role.Role == entity.Student
	})
	if !present {
		return nil, errors.New("failed to get student menu")
	}

	return receiver.MenuRepository.GetUserMenu(user.ID.String())
}

func (receiver *GetMenuUseCase) GetTeacherMenu(userID string) ([]menu.UserMenu, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	present := lo.ContainsBy(user.Roles, func(role entity.SRole) bool {
		return role.Role == entity.Teacher
	})
	if !present {
		return nil, errors.New("failed to get teacher menu")
	}

	return receiver.MenuRepository.GetUserMenu(user.ID.String())
}

func (receiver *GetMenuUseCase) GetUserMenu(userID string) ([]menu.UserMenu, error) {
	user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: userID})
	if err != nil {
		return nil, err
	}

	return receiver.MenuRepository.GetUserMenu(user.ID.String())
}
