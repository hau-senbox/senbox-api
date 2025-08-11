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

type StudentBlockSettingUsecase struct {
	Repo *repository.StudentBlockSettingRepository
}

func NewStudentBlockSettingUsecase(repo *repository.StudentBlockSettingRepository) *StudentBlockSettingUsecase {
	return &StudentBlockSettingUsecase{Repo: repo}
}

// GetByStudentID returns student block setting by student ID
func (uc *StudentBlockSettingUsecase) GetByStudentID(studentID string) (*response.StudentBlockSettingResponse, error) {
	setting, err := uc.Repo.GetByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	if setting == nil {
		return nil, nil
	}

	return &response.StudentBlockSettingResponse{
		ID:              setting.ID,
		StudentID:       setting.StudentID,
		IsDeactive:      setting.IsDeactive,
		IsViewMessage:   setting.IsViewMessage,
		MessageBox:      setting.MessageBox,
		MessageDeactive: setting.MessageDeactive,
		CreatedAt:       setting.CreatedAt,
		UpdatedAt:       setting.UpdatedAt,
	}, nil
}

// Create a new block setting
func (uc *StudentBlockSettingUsecase) Create(req request.StudentBlockSettingRequest) error {
	setting := &entity.StudentBlockSetting{
		StudentID:       req.StudentID,
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
func (uc *StudentBlockSettingUsecase) Update(id int, req request.StudentBlockSettingRequest) error {
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
func (uc *StudentBlockSettingUsecase) Upsert(req request.StudentBlockSettingRequest) error {
	setting, err := uc.Repo.GetByStudentID(req.StudentID)

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
func (uc *StudentBlockSettingUsecase) Delete(id int) error {
	return uc.Repo.Delete(id)
}

func (uc *StudentBlockSettingUsecase) pushToFirestore(req request.StudentBlockSettingRequest) error {
	client := firebase.InitFirestoreClient()
	ctx := context.Background()

	data := map[string]interface{}{
		"student_id":       req.StudentID,
		"is_deactive":      req.IsDeactive != nil && *req.IsDeactive,
		"is_view_message":  req.IsViewMessage != nil && *req.IsViewMessage,
		"message_box":      req.MessageBox,
		"message_deactive": req.MessageDeactive,
		"updated_at":       time.Now(),
	}

	// upsert by student id
	_, err := client.Collection("student_block_settings").
		Doc(req.StudentID).
		Set(ctx, data, firestore.MergeAll)
	return err
}

func (uc *StudentBlockSettingUsecase) GetDeactive4Student(studentID string) (bool, error) {
	return uc.Repo.GetIsDeactiveByStudentID(studentID)
}
