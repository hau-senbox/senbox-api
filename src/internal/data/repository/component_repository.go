package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ComponentRepository struct {
	DBConn *gorm.DB
}

func NewComponentRepository(dbConn *gorm.DB) *ComponentRepository {
	return &ComponentRepository{DBConn: dbConn}
}

func (receiver *ComponentRepository) GetByID(componentID string) (*components.Component, error) {
	var component components.Component
	err := receiver.DBConn.Where("id = ?", componentID).First(&component).Error
	if err != nil {
		log.Error("ComponentRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get role")
	}

	return &component, nil
}

func (receiver *ComponentRepository) GetAllComponentKey() (*[]string, error) {
	var keys []string
	err := receiver.DBConn.Model(&components.Component{}).
		Where("`key` != ?", "").
		Distinct("key").
		Pluck("key", &keys).Error
	if err != nil {
		log.Error("ComponentRepository.GetAllComponentKey: " + err.Error())
		return nil, errors.New("failed to get component keys")
	}

	return &keys, nil
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

func (receiver *ComponentRepository) CreateComponent(req *request.CreateMenuComponentRequest, tx *gorm.DB) error {
	component, err := componentFromRequest(*req)
	if err != nil {
		return err
	}

	if tx == nil {
		err = receiver.DBConn.Create(&components.Component{
			Name:  component.GetName(),
			Type:  component.GetType(),
			Key:   component.GetKey(),
			Value: component.GetValue(),
		}).Error

		if err != nil {
			log.Error("ComponentRepository.CreateComponent: " + err.Error())
			return errors.New("failed to create component")
		}

		req.ID = component.GetID()
		return nil
	}

	err = tx.Create(&components.Component{
		Name:  component.GetName(),
		Type:  component.GetType(),
		Key:   component.GetKey(),
		Value: component.GetValue(),
	}).Error

	if err != nil {
		tx.Rollback()
		log.Error("ComponentRepository.CreateComponent: " + err.Error())
		return errors.New("failed to create component")
	}

	req.ID = component.GetID()
	return nil
}

func (receiver *ComponentRepository) CreateComponents(request *[]request.CreateMenuComponentRequest, tx *gorm.DB) error {
	var componentList []*components.Component
	for _, req := range *request {
		component, err := componentFromRequest(req)
		if err != nil {
			return err
		}

		componentList = append(componentList, component)
	}

	if tx == nil {
		err := receiver.DBConn.Create(&componentList).Error

		if err != nil {
			log.Error("ComponentRepository.CreateComponent: " + err.Error())
			return errors.New("failed to create components")
		}

		for i, component := range componentList {
			(*request)[i].ID = component.GetID()
		}

		return nil
	}

	err := tx.Create(&componentList).Error
	if err != nil {
		tx.Rollback()
		log.Error("ComponentRepository.CreateComponent: " + err.Error())
		return errors.New("failed to create components")
	}

	for i, component := range componentList {
		(*request)[i].ID = component.GetID()
	}

	return nil
}

func (receiver *ComponentRepository) UpdateComponent(req request.UpdateComponentRequest, tx *gorm.DB) error {
	if tx != nil {
		err := tx.Delete(&components.Component{}, req.ID).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = receiver.CreateComponent(&request.CreateMenuComponentRequest{
			Name:  req.Name,
			Type:  req.Type,
			Key:   req.Key,
			Value: req.Value,
		}, tx)
		if err != nil {
			log.Error("ComponentRepository.UpdateComponent: " + err.Error())
			return errors.New("failed to update component")
		}

		return nil
	}

	err := receiver.DBConn.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&components.Component{}, req.ID).Error
		if err != nil {
			return err
		}

		err = receiver.CreateComponent(&request.CreateMenuComponentRequest{
			Name:  req.Name,
			Type:  req.Type,
			Key:   req.Key,
			Value: req.Value,
		}, tx)

		return err
	})

	if err != nil {
		log.Error("ComponentRepository.UpdateComponent: " + err.Error())
		return errors.New("failed to update component")
	}

	return nil
}

func (receiver *ComponentRepository) DeleteComponent(componentID string, tx *gorm.DB) error {
	if tx == nil {
		err := receiver.DBConn.Delete(&components.Component{}, componentID).Error

		if err != nil {
			log.Error("ComponentRepository.DeleteComponent: " + err.Error())
			return errors.New("failed to delete component")
		}

		return nil
	}

	err := tx.Delete(&components.Component{}, componentID).Error

	if err != nil {
		tx.Rollback()
		log.Error("ComponentRepository.DeleteComponent: " + err.Error())
		return errors.New("failed to delete component")
	}

	return nil
}

func (receiver *ComponentRepository) DeleteBySectionID(sectionID string, tx *gorm.DB) error {
	if tx == nil {
		err := receiver.DBConn.
			Where("section_id = ?", sectionID).
			Delete(&components.Component{}).Error

		if err != nil {
			log.Error("ComponentRepository.DeleteBySectionID: " + err.Error())
			return errors.New("failed to delete components by section_id")
		}
		return nil
	}

	err := tx.
		Where("section_id = ?", sectionID).
		Delete(&components.Component{}).Error

	if err != nil {
		_ = tx.Rollback()
		log.Error("ComponentRepository.DeleteBySectionID: " + err.Error())
		return errors.New("failed to delete components by section_id")
	}
	return nil
}

func (receiver *ComponentRepository) GetAllByKey(key string) ([]components.Component, error) {
	var componentsList []components.Component

	err := receiver.DBConn.
		Where("`key` = ?", key).
		Find(&componentsList).Error

	if err != nil {
		log.Error("ComponentRepository.GetAllByKey: " + err.Error())
		return nil, errors.New("failed to get components by key")
	}

	return componentsList, nil
}

func (receiver *ComponentRepository) GetBySectionID(sectionID string) (*components.Component, error) {
	var component components.Component

	err := receiver.DBConn.
		Where("`section_id` = ?", sectionID).
		First(&component).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error("ComponentRepository.GetByKey: " + err.Error())
		return nil, errors.New("failed to get component by sectionID")
	}

	return &component, nil
}
