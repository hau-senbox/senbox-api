package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/firebase"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrganizationSettingUsecase struct {
	Repo             *repository.OrganizationSettingRepository
	ComponentRepo    *repository.ComponentRepository
	OrganizationRepo *repository.OrganizationRepository
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
	if req.Component.Name != "" || req.Component.ID != "" {
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

		if req.Component.ID != "" {
			// Update component
			component.ID = uuid.MustParse(req.Component.ID)
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
			componentID = uuid.New().String()
			component.ID = uuid.MustParse(componentID)
			if err := u.ComponentRepo.CreateWithTx(tx, component); err != nil {
				tx.Rollback()
				return fmt.Errorf("create component fail: %w", err)
			}
			isNewComp = true
		}
	}

	// Upsert organization setting
	setting := &entity.OrganizationSetting{
		OrganizationID: req.OrganizationID,
		DeviceID:       req.DeviceID,
	}

	if req.IsViewMessageBox != nil {
		setting.IsViewMessageBox = *req.IsViewMessageBox
	}
	if req.IsShowMessage != nil {
		setting.IsShowMessage = *req.IsShowMessage
	}
	if req.IsDeactiveApp != nil {
		setting.IsDeactiveApp = *req.IsDeactiveApp
	}
	if req.IsDeactiveTopMenu != nil {
		setting.IsDeactiveTopMenu = *req.IsDeactiveTopMenu
	}
	if req.IsShowSpecialBtn != nil {
		setting.IsShowSpecialBtn = *req.IsShowSpecialBtn
	}

	if componentID != "" {
		setting.ComponentID = componentID
	}

	if req.IsShowSpecialBtn != nil {
		setting.IsShowSpecialBtn = *req.IsShowSpecialBtn
	}

	// text
	if req.MessageBox != nil {
		setting.MessageBox = *req.MessageBox
	}
	if req.MessageDeactiveApp != nil {
		setting.MessageDeactiveApp = *req.MessageDeactiveApp
	}
	if req.MessageTopMenu != nil {
		setting.MessageTopMenu = *req.MessageTopMenu
	}
	if req.TopMenuPassword != nil {
		setting.TopMenuPassword = *req.TopMenuPassword
	}

	existingSetting, err := u.Repo.GetByDeviceID(req.DeviceID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("query organization setting fail: %w", err)
	}

	if existingSetting != nil {
		setting.ID = existingSetting.ID
		if req.IsViewMessageBox != nil {
			existingSetting.IsViewMessageBox = *req.IsViewMessageBox
		}
		if req.IsShowMessage != nil {
			existingSetting.IsShowMessage = *req.IsShowMessage
		}

		if req.MessageBox != nil {
			existingSetting.MessageBox = *req.MessageBox
		}

		if req.IsShowSpecialBtn != nil {
			existingSetting.IsShowSpecialBtn = *req.IsShowSpecialBtn
		}

		if req.IsDeactiveApp != nil {
			existingSetting.IsDeactiveApp = *req.IsDeactiveApp
		}

		if req.MessageDeactiveApp != nil {
			existingSetting.MessageDeactiveApp = *req.MessageDeactiveApp
		}

		if req.IsDeactiveTopMenu != nil {
			existingSetting.IsDeactiveTopMenu = *req.IsDeactiveTopMenu
		}

		if req.MessageTopMenu != nil {
			existingSetting.MessageTopMenu = *req.MessageTopMenu
		}

		if req.TopMenuPassword != nil {
			existingSetting.TopMenuPassword = *req.TopMenuPassword
		}

		// Merge component
		if isNewComp && componentID != "" {
			existingSetting.ComponentID = componentID
		}

		if err := u.Repo.UpdateWithTx(tx, existingSetting); err != nil {
			tx.Rollback()
			return fmt.Errorf("update organization setting fail: %w", err)
		}
	} else {
		setting.CreatedAt = time.Now()
		if err := u.Repo.CreateWithTx(tx, setting); err != nil {
			tx.Rollback()
			return fmt.Errorf("create organization setting fail: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Push lên Firestore sau khi commit thành công
	// get setting by device id
	settingFirestore, _ := u.GetOrgSetting(req.DeviceID)
	if err := u.pushToFirestore(settingFirestore); err != nil {
		// Nếu push Firestore fail thì log lại nhưng không rollback DB
		log.Printf("pushToFirestore error: %v", err)
		return err
	}

	return nil
}

func (u *OrganizationSettingUsecase) GetOrgSetting(deviceID string) (response.OrgSettingResponse, error) {
	// Lấy thông tin OrgSetting
	orgSetting, err := u.Repo.GetByDeviceID(deviceID)
	if err != nil {
		return response.OrgSettingResponse{}, err
	}

	// Lấy danh sách components
	component, _ := u.ComponentRepo.GetByID(orgSetting.ComponentID)

	// ger org info
	orgInfo, _ := u.OrganizationRepo.GetByID(orgSetting.OrganizationID)
	resp := mapper.MapOrgSettingToResponse(orgSetting, component, orgInfo.OrganizationName)

	return resp, nil
}

func (uc *OrganizationSettingUsecase) pushToFirestore(setting response.OrgSettingResponse) error {
	client := firebase.InitFirestoreClient()
	ctx := context.Background()

	// Convert Component.Value (struct) -> map[string]interface{}
	var valueMap map[string]interface{}
	if b, err := json.Marshal(setting.Component.Value); err == nil {
		_ = json.Unmarshal(b, &valueMap)
	} else {
		return fmt.Errorf("failed to marshal component value: %w", err)
	}

	// Build data map
	data := map[string]interface{}{
		"device_id":            setting.DeviceID,
		"is_view_message_box":  setting.IsViewMessageBox,
		"is_show_message":      setting.IsShowMessage,
		"message_box":          setting.MessageBox,
		"is_show_special_btn":  setting.IsShowSpecialBtn,
		"is_deactive_app":      setting.IsDeactiveApp,
		"message_deactive_app": setting.MessageDeactiveApp,
		"is_deactive_top_menu": setting.IsDeactiveTopMenu,
		"message_top_menu":     setting.MessageTopMenu,
		"top_menu_password":    setting.TopMenuPassword,
		"component": map[string]interface{}{
			"name":  setting.Component.Name,
			"type":  setting.Component.Type,
			"key":   setting.Component.Key,
			"value": valueMap, // <-- giữ nguyên tất cả field trong struct
		},
		"updated_at": time.Now(),
	}

	// Upsert theo device_id
	docID := setting.DeviceID
	_, err := client.Collection("device_settings").
		Doc(docID).
		Set(ctx, data, firestore.MergeAll)

	return err
}

// UploadOrgSettingNewsDevice uploads organization setting news for device & portal
func (u *OrganizationSettingUsecase) UploadOrgSettingNewsDevice(req request.UploadOrgSettingDeviceNewsRequest) error {
	// check exist by org id
	exist, _ := u.Repo.GetSettingNewsByOrganizationID(req.OrganizationID)

	if exist != nil {
		// Update
		exist.IsPublishedDevice = req.IsPublishedDevice
		exist.MessageDeviceNews = req.MessageDeviceNews

		return u.Repo.UpdateSettingNews(exist)
	}

	// Create
	newSetting := &entity.OrganizationNewsSetting{
		OrganizationID:    req.OrganizationID,
		IsPublishedDevice: req.IsPublishedDevice,
		MessageDeviceNews: req.MessageDeviceNews,
	}

	return u.Repo.CreateSettingNews(newSetting)
}

func (u *OrganizationSettingUsecase) UploadOrgSettingNewsPortal(req request.UploadOrgSettingPortalNewsRequest) error {
	// check exist by org id
	exist, _ := u.Repo.GetSettingNewsByOrganizationID(req.OrganizationID)

	if exist != nil {
		// Update
		exist.IsPublishedPortal = req.IsPublishedPortal
		exist.MessagePortalNews = req.MessagePortalNews

		return u.Repo.UpdateSettingNews(exist)
	}

	// Create
	newSetting := &entity.OrganizationNewsSetting{
		OrganizationID:    req.OrganizationID,
		IsPublishedPortal: req.IsPublishedPortal,
		MessagePortalNews: req.MessagePortalNews,
	}

	return u.Repo.CreateSettingNews(newSetting)
}

func (u *OrganizationSettingUsecase) GetOrgSettingNews(orgID string) (*response.OrgSettingNewsResponse, error) {
	setting, err := u.Repo.GetSettingNewsByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, nil
	}

	return &response.OrgSettingNewsResponse{
		OrganizationID:    setting.OrganizationID,
		IsPublishedDevice: setting.IsPublishedDevice,
		MessageDeviceNews: setting.MessageDeviceNews,
		IsPublishedPortal: setting.IsPublishedPortal,
		MessagePortalNews: setting.MessagePortalNews,
	}, nil
}
