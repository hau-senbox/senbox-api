package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
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

	// upsert compoent
	var componentID uuid.UUID

	component := &components.Component{
		Name:  req.Component.Name,
		Type:  components.ComponentType(req.Component.Type),
		Key:   req.Component.Key,
		Value: datatypes.JSON([]byte(req.Component.Value)),
	}

	if req.Component.ID != nil && *req.Component.ID != uuid.Nil {
		// Nếu có ID truyền lên
		componentID = *req.Component.ID
		component.ID = componentID

		existingComponent, err := u.ComponentRepo.GetByID(componentID.String())
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error("rollback by error query component:", err)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("query component fail: %w", err)
		}

		if existingComponent != nil {
			// Update
			if err := u.ComponentRepo.UpdateWithTx(tx, component); err != nil {
				logrus.Error("rollback by error update component:", err)
				tx.Rollback()
				rolledBack = true
				return fmt.Errorf("update component fail: %w", err)
			}
		} else {
			logrus.Error("rollback by error create component (from non-existent id):", err)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("create component fail (Component ID wrong): %w", err)
		}
	} else {
		// Tạo mới
		component.ID = uuid.New()
		componentID = component.ID
		if err := u.ComponentRepo.CreateWithTx(tx, component); err != nil {
			logrus.Error("rollback by error create component:", err)
			tx.Rollback()
			rolledBack = true
			return fmt.Errorf("create component fail: %w", err)
		}
	}

	// Upsert organization setting
	setting := &entity.OrganizationSetting{
		OrganizationID:    req.OrganizationID,
		ComponentID:       componentID.String(),
		IsViewMessage:     req.IsViewMessage,
		IsShowOrgNews:     req.IsShowOrgNews,
		IsDeactiveTopMenu: req.IsDeactiveTopMenu,
		IsShowSpecialBtn:  req.IsShowSpecialBtn,
		MessageBox:        req.MessageBox,
		MessageTopMenu:    req.MessageTopMenu,
	}

	existingSetting, err := u.Repo.GetByOrgID(req.OrganizationID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error("rollback by error query org setting:", err)
		tx.Rollback()
		rolledBack = true
		return fmt.Errorf("query organization setting fail: %w", err)
	}

	if existingSetting != nil {
		setting.ID = existingSetting.ID // Giữ nguyên ID cũ khi update
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

func (u *OrganizationSettingUsecase) GetOrgSetting(orgID string) (*response.OrgSettingResponse, error) {
	// Lấy thông tin OrgSetting
	orgSetting, err := u.Repo.GetByOrgID(orgID)
	if err != nil {
		return nil, err
	}

	// Lấy danh sách components
	components, err := u.ComponentRepo.GetByID(orgSetting.ComponentID)
	if err != nil {
		return nil, err
	}

	// Mapping sang response
	resp := &response.OrgSettingResponse{
		ID:                orgSetting.ID.String(),
		OrganizationID:    orgSetting.OrganizationID,
		IsViewMessage:     orgSetting.IsViewMessage,
		IsShowOrgNews:     orgSetting.IsShowOrgNews,
		IsDeactiveTopMenu: orgSetting.IsDeactiveTopMenu,
		IsShowSpecialBtn:  orgSetting.IsShowSpecialBtn,
		MessageBox:        orgSetting.MessageBox,
		MessageTopMenu:    orgSetting.MessageTopMenu,
		Component:         components,
	}

	return resp, nil
}
