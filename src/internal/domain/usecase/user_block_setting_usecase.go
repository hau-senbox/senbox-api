package usecase

import (
	"context"
	"errors"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/firebase"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBlockSettingUsecase struct {
	Repo        *repository.UserBlockSettingRepository
	TeacherRepo *repository.TeacherApplicationRepository
	StaffRepo   *repository.StaffApplicationRepository
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

func (uc *UserBlockSettingUsecase) GetByUserID4App(userID string) (*response.UserBlockSettingResponse4App, error) {
	setting, err := uc.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if setting == nil {
		return nil, nil
	}

	return &response.UserBlockSettingResponse4App{
		ID:              setting.ID,
		IsNeedToUpdate:  setting.IsNeedToUpdate,
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
	return uc.pushToFirestore(setting)
}

// Delete block setting by ID
func (uc *UserBlockSettingUsecase) Delete(id int) error {
	return uc.Repo.Delete(id)
}

func (uc *UserBlockSettingUsecase) pushToFirestore(setting *entity.UserBlockSetting) error {
	client := firebase.InitFirestoreClient()
	ctx := context.Background()

	data := map[string]interface{}{
		"user_id":           setting.UserID,
		"is_deactive":       setting.IsDeactive,
		"is_view_message":   setting.IsViewMessage,
		"message_box":       setting.MessageBox,
		"message_deactive":  setting.MessageDeactive,
		"is_need_to_update": setting.IsNeedToUpdate,
		"updated_at":        time.Now(),
	}

	// upsert by user id
	collection := "user_block_settings"
	// get app config
	if config.IsDevMode() {
		collection = "dev_user_block_settings"
	}
	_, err := client.Collection(collection).
		Doc(setting.UserID).
		Set(ctx, data, firestore.MergeAll)
	return err
}

func (uc *UserBlockSettingUsecase) GetDeactive4User(userID string) (bool, error) {
	return uc.Repo.GetIsDeactiveByUserID(userID)
}

func (uc *UserBlockSettingUsecase) GetDeactive4Teacher(teacherID string) (bool, error) {
	// Parse teacherID sang UUID
	tid, err := uuid.Parse(teacherID)
	if err != nil {
		return false, err
	}

	// Lấy thông tin teacher
	teacher, err := uc.TeacherRepo.GetByID(tid)
	if err != nil {
		return false, err
	}
	if teacher == nil {
		return false, nil
	}

	// Lấy is_deactive từ user_id
	isDeactive, err := uc.Repo.GetIsDeactiveByUserID(teacher.UserID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return isDeactive, nil

}

func (uc *UserBlockSettingUsecase) GetDeactive4Staff(staffID string) (bool, error) {
	// Parse teacherID sang UUID
	stid, err := uuid.Parse(staffID)
	if err != nil {
		return false, err
	}

	// Lấy thông tin teacher
	staff, err := uc.StaffRepo.GetByID(stid)
	if err != nil {
		return false, err
	}
	if staff == nil {
		return false, nil
	}

	// Lấy is_deactive từ user_id
	isDeactive, err := uc.Repo.GetIsDeactiveByUserID(staff.UserID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return isDeactive, nil
}

func (uc *UserBlockSettingUsecase) OnIsNeedToUpdate() error {

	err := uc.Repo.OnIsNeedToUpdate()
	if err != nil {
		return err
	}

	//

	// update firestore
	// err = uc.pushToFirestore(request.UserBlockSettingRequest{
	// 	UserID:         "all",
	// 	IsNeedToUpdate: true,
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (uc *UserBlockSettingUsecase) OffIsNeedToUpdateByUser(userID string) error {
	return uc.Repo.OffIsNeedToUpdate(userID)
}

func (uc *UserBlockSettingUsecase) MigrateFirestore() error {
	// get all user block settings
	userBlockSettings, err := uc.Repo.GetAll()
	if err != nil {
		return err
	}

	for _, userBlockSetting := range userBlockSettings {
		uc.pushToFirestore(&userBlockSetting)
	}

	return nil
}
