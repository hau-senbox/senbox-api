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
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var componentID string
	isNewComp := false

	// Xử lý component nếu có
	if req.Component.Name != "" || (req.Component.ID != nil && *req.Component.ID != uuid.Nil) {
		valueJSON, err := json.Marshal(req.Component.Value)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("marshal component value fail: %w", err)
		}

		component := &components.Component{
			Name:  req.Component.Name,
			Type:  components.ComponentType(req.Component.Type),
			Key:   req.Component.Key,
			Value: datatypes.JSON(valueJSON),
		}

		if req.Component.ID != nil && *req.Component.ID != uuid.Nil {
			// Update component
			component.ID = *req.Component.ID
			existingComp, err := u.ComponentRepo.GetByID(component.ID.String())
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return fmt.Errorf("query component fail: %w", err)
			}
			if existingComp != nil {
				if err := u.ComponentRepo.UpdateWithTx(tx, component); err != nil {
					tx.Rollback()
					return fmt.Errorf("update component fail: %w", err)
				}
				isNewComp = false
			} else {
				tx.Rollback()
				return fmt.Errorf("component ID not found, cannot update")
			}
		} else {
			// Create component mới
			component.ID = uuid.New()
			if err := u.ComponentRepo.CreateWithTx(tx, component); err != nil {
				tx.Rollback()
				return fmt.Errorf("create component fail: %w", err)
			}
			isNewComp = true
		}

		componentID = component.ID.String()
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

	if componentID != "" {
		setting.ComponentID = componentID
	}

	existingSetting, err := u.Repo.GetByDeviceID(req.DeviceID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("query organization setting fail: %w", err)
	}

	if existingSetting != nil {
		setting.ID = existingSetting.ID
		// Merge dữ liệu
		if req.IsViewMessageBox == false && req.IsShowMessage == false && req.MessageBox == "" {
			setting.IsViewMessageBox = existingSetting.IsViewMessageBox
			setting.IsShowMessage = existingSetting.IsShowMessage
			setting.MessageBox = existingSetting.MessageBox
		}

		if req.IsShowSpecialBtn == false {
			setting.IsShowSpecialBtn = existingSetting.IsShowSpecialBtn
		}

		if req.IsDeactiveApp == false && req.MessageDeactiveApp == "" {
			setting.IsDeactiveApp = existingSetting.IsDeactiveApp
			setting.MessageDeactiveApp = existingSetting.MessageDeactiveApp
		}

		if req.IsDeactiveTopMenu == false && req.MessageTopMenu == "" && req.TopMenuPassword == "" {
			setting.IsDeactiveTopMenu = existingSetting.IsDeactiveTopMenu
			setting.MessageTopMenu = existingSetting.MessageTopMenu
			setting.TopMenuPasswod = existingSetting.TopMenuPasswod
		}

		// Merge component
		if !isNewComp && componentID == "" {
			setting.ComponentID = existingSetting.ComponentID
		}
		if err := u.Repo.UpdateWithTx(tx, setting); err != nil {
			tx.Rollback()
			return fmt.Errorf("update organization setting fail: %w", err)
		}
	} else {
		if err := u.Repo.CreateWithTx(tx, setting); err != nil {
			tx.Rollback()
			return fmt.Errorf("create organization setting fail: %w", err)
		}
	}

	return tx.Commit().Error
}

func (u *OrganizationSettingUsecase) GetOrgSetting(deviceID string) (response.OrgSettingResponse, error) {
	// Lấy thông tin OrgSetting
	orgSetting, err := u.Repo.GetByDeviceID(deviceID)
	if err != nil {
		return response.OrgSettingResponse{}, err
	}

	// Lấy danh sách components
	component, _ := u.ComponentRepo.GetByID(orgSetting.ComponentID)

	resp := mapper.MapOrgSettingToResponse(orgSetting, component)

	return resp, nil
}
