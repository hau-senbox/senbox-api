package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrganizationSettingUsecase struct {
	Repo          *repository.OrganizationSettingRepository
	ComponentRepo *repository.ComponentRepository
}

func NewOrganizationSettingUsecase(repo *repository.OrganizationSettingRepository) *OrganizationSettingUsecase {
	return &OrganizationSettingUsecase{Repo: repo}
}

// CreateOrganizationSetting tạo mới
func (u *OrganizationSettingUsecase) CreateOrganizationSetting(setting *entity.OrganizationSetting) error {
	return u.Repo.Create(setting)
}

// GetOrganizationSettingByID lấy theo ID
func (u *OrganizationSettingUsecase) GetOrganizationSettingByID(id uint) (*entity.OrganizationSetting, error) {
	return u.Repo.GetByID(id)
}

// UpdateOrganizationSetting cập nhật
func (u *OrganizationSettingUsecase) UpdateOrganizationSetting(setting *entity.OrganizationSetting) error {
	return u.Repo.Update(setting)
}

// DeleteOrganizationSetting xóa theo ID
func (u *OrganizationSettingUsecase) DeleteOrganizationSetting(id uint) error {
	return u.Repo.Delete(id)
}

// ListOrganizationSettings lấy tất cả
func (u *OrganizationSettingUsecase) ListOrganizationSettings() ([]entity.OrganizationSetting, error) {
	return u.Repo.List()
}

// upload org setting
func (u *OrganizationSettingUsecase) UploadOrgSetting(req request.UploadOrgSettingRequest) error {
	tx := u.Repo.DBConn.Begin()
	rolledBack := false

	var componentID string

	// Nếu request có component thì xử lý upsert component
	if req.Component.Name != "" || req.Component.ID != nil {
		var cid uuid.UUID
		valueJSON, err := json.Marshal(req.Component.Value)

		if err != nil {
			return err
		}

		component := &components.Component{
			Name:  req.Component.Name,
			Type:  components.ComponentType(req.Component.Type),
			Key:   req.Component.Key,
			Value: datatypes.JSON([]byte(valueJSON)),
		}

		if req.Component.ID != nil && *req.Component.ID != uuid.Nil {
			// Update component
			cid = *req.Component.ID
			component.ID = cid

			existingComponent, err := u.ComponentRepo.GetByID(cid.String())
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.Error("rollback by error query component:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("query component fail: %w", err)
			}

			if existingComponent != nil {
				if err := u.ComponentRepo.UpdateWithTx(tx, component); err != nil {
					logrus.Error("rollback by error update component:", err)
					tx.Rollback()
					rolledBack = true
					return fmt.Errorf("update component fail: %w", err)
				}
			} else {
				logrus.Error("rollback by error create component (from non-existent id)")
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("component ID not found, cannot update")
			}
		} else {
			// Create component
			cid = uuid.New()
			component.ID = cid
			if err := u.ComponentRepo.CreateWithTx(tx, component); err != nil {
				logrus.Error("rollback by error create component:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("create component fail: %w", err)
			}
		}

		componentID = cid.String()
	}

	// Upsert organization setting
	setting := &entity.OrganizationSetting{
		OrganizationID:     req.OrganizationID,
		DeviceID:           req.DeviceID,
		IsViewMessageBox:   req.IsViewMessageBox,
		IsShowMessage:      req.IsShowMessage,
		MessageBox:         req.MessageBox,
		IsShowSpecialBtn:   req.IsShowSpecialBtn,
		IsDeactiveApp:      req.IsDeactiveApp,
		MessageDeactiveApp: req.MessageDeactiveApp,
		IsDeactiveTopMenu:  req.IsDeactiveTopMenu,
		MessageTopMenu:     req.MessageTopMenu,
		TopMenuPasswod:     req.TopMenuPassword,
	}

	// Nếu có component thì set ComponentID
	if componentID != "" {
		setting.ComponentID = componentID
	}

	existingSetting, err := u.Repo.GetByDeviceIdAndOrgId(req.DeviceID, req.OrganizationID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error("rollback by error query org setting:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("query organization setting fail: %w", err)
	}

	if existingSetting != nil {
		setting.ID = existingSetting.ID
		setting.ComponentID = existingSetting.ComponentID
		if err := u.Repo.UpdateWithTx(tx, setting); err != nil {
			logrus.Error("rollback by error update org setting:", err)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("update organization setting fail: %w", err)
		}
	} else {
		if err := u.Repo.CreateWithTx(tx, setting); err != nil {
			logrus.Error("rollback by error create org setting:", err)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("create organization setting fail: %w", err)
		}
	}

	if !rolledBack {
		return tx.Commit().Error
	}
	return nil
}

func (u *OrganizationSettingUsecase) GetOrgSetting(deviceID string, orgID string) (response.OrgSettingResponse, error) {
	// Lấy thông tin OrgSetting
	orgSetting, err := u.Repo.GetByDeviceIdAndOrgId(deviceID, orgID)
	if err != nil {
		return response.OrgSettingResponse{}, err
	}

	// Lấy danh sách components
	component, _ := u.ComponentRepo.GetByID(orgSetting.ComponentID)

	resp := mapper.MapOrgSettingToResponse(orgSetting, component)

	return resp, nil
}
