package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"

	"gorm.io/datatypes"
)

type UploadUserMenuUseCase struct {
	*repository.MenuRepository
	*repository.ComponentRepository
}

func (receiver *UploadUserMenuUseCase) Upload(req request.UploadUserMenuRequest) error {
	tx := receiver.MenuRepository.DBConn.Begin()
	if err := receiver.MenuRepository.DeleteUserMenu(req.UserID, tx); err != nil {
		return err
	}

	if len(req.Components) > 0 {
		err := receiver.ComponentRepository.CreateComponents(&req.Components, tx)
		if err != nil {
			return err
		}

		if err := receiver.MenuRepository.CreateUserMenu(request.CreateUserMenuRequest{
			UserID:     req.UserID,
			Components: req.Components,
		}, tx); err != nil {
			return fmt.Errorf("failed to create user menu: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func componentFromRequest(req request.CreateMenuComponentRequest) (*components.Component, error) {
	var component components.IComponent
	componentType, err := components.GetComponentTypeFromString(req.Type)
	if err != nil {
		return nil, err
	}

	switch componentType {
	case components.ButtonURL:
		component = components.NewButtonURLComponent()
	case components.ButtonForm:
		component = components.NewButtonFormComponent()
	default:
		return nil, errors.New("invalid component type")
	}

	component.SetName(req.Name)
	component.SetKey(req.Key)
	component.SetValue(datatypes.JSON(req.Value))
	component.SetSectionID(req.SectionId)

	if err = component.NormalizeValue(); err != nil {
		return nil, err
	}

	return component.GetComponent(), nil
}

func (receiver *UploadUserMenuUseCase) UploadSectionMenu(req request.UploadSectionMenuRequest) error {

	tx := receiver.MenuRepository.DBConn.Begin()
	for _, item := range req.Components {
		component, err := componentFromRequest(item)
		if err != nil {
			return err
		}

		if err := receiver.MenuRepository.DeleteSectionMenu(component.GetSectionID(), tx); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
