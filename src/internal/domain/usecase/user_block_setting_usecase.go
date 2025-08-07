package usecase

import (
	"context"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/firebase"
	"time"

	"cloud.google.com/go/firestore"
	"gorm.io/gorm"
)

type UserBlockSettingUsecase struct {
	Repo *repository.UserBlockSettingRepository
}

func NewUserBlockSettingUsecase(repo *repository.UserBlockSettingRepository) *UserBlockSettingUsecase {
	return &UserBlockSettingUsecase{Repo: repo}
}

// GetByUserID returns user block setting by user ID
func (uc *UserBlockSettingUsecase) GetByUserID(userID string) (*response.UserBlockSettingResponse, error) {
	setting, err := uc.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if setting == nil {
		return nil, nil
	}

	return &response.UserBlockSettingResponse{
		ID:              setting.ID,
		UserID:          setting.UserID,
		IsDeactive:      setting.IsDeactive,
		IsViewMessage:   setting.IsViewMessage,
		MessageBox:      setting.MessageBox,
		MessageDeactive: setting.MessageDeactive,
		CreatedAt:       setting.CreatedAt,
		UpdatedAt:       setting.UpdatedAt,
	}, nil
}

// Create a new block setting
func (uc *UserBlockSettingUsecase) Create(req request.UserBlockSettingRequest) error {
	setting := &entity.UserBlockSetting{
		UserID:          req.UserID,
		IsDeactive:      *req.IsDeactive,
		IsViewMessage:   *req.IsViewMessage,
		MessageBox:      req.MessageBox,
		MessageDeactive: req.MessageDeactive,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return uc.Repo.Create(setting)
}

// Update existing block setting
func (uc *UserBlockSettingUsecase) Update(id int, req request.UserBlockSettingRequest) error {
	setting, err := uc.Repo.GetByID(id)
	if err != nil {
		return err
	}

	setting.IsDeactive = *req.IsDeactive
	setting.IsViewMessage = *req.IsViewMessage
	setting.MessageBox = req.MessageBox
	setting.MessageDeactive = req.MessageDeactive
	setting.UpdatedAt = time.Now()

	return uc.Repo.Update(setting)
}

// Upsert (create if not exists, otherwise update)
func (uc *UserBlockSettingUsecase) Upsert(req request.UserBlockSettingRequest) error {
	setting, err := uc.Repo.GetByUserID(req.UserID)

	// Nếu là lỗi "record not found" thì tạo mới
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uc.Create(req)
		}
		// Nếu là lỗi khác → return ra ngoài
		return err
	}

	// Nếu tìm được bản ghi thì update
	setting.IsDeactive = *req.IsDeactive
	setting.IsViewMessage = *req.IsViewMessage
	setting.MessageBox = req.MessageBox
	setting.MessageDeactive = req.MessageDeactive
	setting.UpdatedAt = time.Now()

	if err := uc.Repo.Update(setting); err != nil {
		return err
	}

	// Sau khi cập nhật thành công → ghi Firestore
	return uc.pushToFirestore(req)
}

// Delete block setting by ID
func (uc *UserBlockSettingUsecase) Delete(id int) error {
	return uc.Repo.Delete(id)
}

func (uc *UserBlockSettingUsecase) pushToFirestore(req request.UserBlockSettingRequest) error {
	client := firebase.InitFirestoreClient()
	ctx := context.Background()

	data := map[string]interface{}{
		"user_id":          req.UserID,
		"is_deactive":      req.IsDeactive != nil && *req.IsDeactive,
		"is_view_message":  req.IsViewMessage != nil && *req.IsViewMessage,
		"message_box":      req.MessageBox,
		"message_deactive": req.MessageDeactive,
		"updated_at":       time.Now(),
	}

	// upsert by user id
	_, err := client.Collection("user_block_settings").
		Doc(req.UserID).
		Set(ctx, data, firestore.MergeAll)
	return err
}
