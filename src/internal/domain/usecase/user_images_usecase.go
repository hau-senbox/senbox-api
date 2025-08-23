package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/uploader"
	"sort"
	"time"

	"github.com/google/uuid"
)

type UserImagesUsecase struct {
	Repo               *repository.UserImagesRepository
	ImageRepo          *repository.ImageRepository
	DeleteImageUsecase *DeleteImageUseCase
	GetImageUseCase    *GetImageUseCase
}

// UploadAvt tạo avatar mới cho Owner
func (uc *UserImagesUsecase) UploadAvt(req request.UploadAvatarRequest) error {
	// validate UUID
	ownerID, err := uuid.Parse(req.OwnerID)
	if err != nil {
		return fmt.Errorf("invalid owner_id: %w", err)
	}

	// validate OwnerRole
	switch req.OwnerRole {
	case value.OwnerRoleUser,
		value.OwnerRoleTeacher,
		value.OwnerRoleStaff,
		value.OwnerRoleStudent,
		value.OwnerRoleChild:
	default:
		return fmt.Errorf("invalid owner_role: %s", req.OwnerRole)
	}

	// map sang entity
	userImage := entity.UserImages{
		OwnerID:   ownerID,
		OwnerRole: req.OwnerRole,
		ImageID:   req.ImageID,
		Index:     req.Index,
		Feature:   string(value.ImageFeatureAvatar),
	}

	// kiểm tra xem đã có avatar ở index này chưa
	userImageExist, err := uc.Repo.GetByOwnerAndIndex(req.OwnerID, string(req.OwnerRole), req.Index)
	if err == nil && userImageExist != nil {
		// tìm image metadata cũ
		img, err := uc.ImageRepo.GetByID(uint64(userImageExist.ImageID))
		if err == nil && img != nil {
			if delErr := uc.DeleteImageUsecase.DeleteImage(img.Key); delErr != nil {
				return fmt.Errorf("failed to delete old image: %w", delErr)
			}
		}

		// update record cũ với ImageID mới
		userImageExist.ImageID = req.ImageID
		userImageExist.UpdatedAt = time.Now()

		if err := uc.Repo.Update(userImageExist); err != nil {
			return fmt.Errorf("failed to update user_image record: %w", err)
		}

		return nil
	}

	// tạo record mới
	return uc.Repo.Create(&userImage)
}

// Create
func (uc *UserImagesUsecase) Create(userImage *entity.UserImages) error {
	return uc.Repo.Create(userImage)
}

// Get by ID
func (uc *UserImagesUsecase) GetByID(id string) (*entity.UserImages, error) {
	return uc.Repo.GetByID(id)
}

// Get all
func (uc *UserImagesUsecase) GetAll() ([]entity.UserImages, error) {
	return uc.Repo.GetAll()
}

// Update
func (uc *UserImagesUsecase) Update(userImage *entity.UserImages) error {
	return uc.Repo.Update(userImage)
}

// Delete
func (uc *UserImagesUsecase) Delete(id string) error {
	return uc.Repo.Delete(id)
}

func (uc *UserImagesUsecase) GetAvt4Owner(ownerID string, ownerRole value.OwnerRole) ([]response.Avatar, error) {
	// Lấy danh sách ảnh từ DB
	userImages, err := uc.Repo.GetAvtByOwnerRole(ownerID, string(ownerRole), value.ImageFeatureAvatar)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for owner %s: %w", ownerID, err)
	}

	// Map sang response
	avatars := make([]response.Avatar, 0, len(userImages))
	for _, img := range userImages {
		// get img key by img id
		imageEntity, _ := uc.ImageRepo.GetByID(img.ImageID)
		// get img url
		url, _ := uc.GetImageUseCase.GetUrlByKey(imageEntity.Key, uploader.UploadPrivate)
		avatars = append(avatars, response.Avatar{
			ImageID:  img.ImageID,
			ImageKey: imageEntity.Key,
			Index:    img.Index,
			IsMain:   img.IsMain,
			ImageUrl: *url,
		})
	}

	// Sort by Index (tăng dần)
	sort.Slice(avatars, func(i, j int) bool {
		return avatars[i].Index < avatars[j].Index
	})

	return avatars, nil
}

func (uc *UserImagesUsecase) UpdateIsMain(request request.UpdateIsMainAvatar) error {
	// Tìm ảnh theo ownerID, ownerRole và index
	userImage, err := uc.Repo.GetByOwnerRoleAndIndex(request.OwnerID, string(request.OwnerRole), request.Index)
	if err != nil {
		return fmt.Errorf("failed to find user image: %w", err)
	}
	if userImage == nil {
		return fmt.Errorf("user image not found for owner=%s role=%s index=%d", request.OwnerID, request.OwnerRole, request.Index)
	}

	if userImage.IsMain {
		// Nếu đang true thì chuyển về false
		userImage.IsMain = false
		if err := uc.Repo.Update(userImage); err != nil {
			return fmt.Errorf("failed to update user image: %w", err)
		}
	} else {
		// Nếu đang false thì reset tất cả về false trước
		if err := uc.Repo.ResetIsMain(request.OwnerID, string(request.OwnerRole)); err != nil {
			return fmt.Errorf("failed to reset is_main: %w", err)
		}
		// Rồi set ảnh này thành true
		userImage.IsMain = true
		if err := uc.Repo.Update(userImage); err != nil {
			return fmt.Errorf("failed to update user image: %w", err)
		}
	}

	if err := uc.Repo.Update(userImage); err != nil {
		return fmt.Errorf("failed to update user image: %w", err)
	}

	return nil
}

func (uc *UserImagesUsecase) GetImg4Ownewr(ownerID string, ownerRole value.OwnerRole) (*entity.SImage, error) {
	userImage, err := uc.Repo.GetByOwnerAndRoleIsMain(ownerID, string(ownerRole))
	if err != nil {
		return nil, err
	}
	// get img
	return uc.ImageRepo.GetByID(userImage.ImageID)
}

func (uc *UserImagesUsecase) DeleteUserAvatar(request request.DeleteUserAvatarRequest) error {
	// B1: Tìm ảnh theo ownerID, ownerRole và index
	userImage, err := uc.Repo.GetByOwnerRoleAndIndex(request.OwnerID, string(request.OwnerRole), request.Index)
	if err != nil {
		return fmt.Errorf("failed to find user image: %w", err)
	}
	if userImage == nil {
		return fmt.Errorf("user image not found for owner=%s role=%s index=%d", request.OwnerID, request.OwnerRole, request.Index)
	}

	// tìm image metadata cũ
	img, err := uc.ImageRepo.GetByID(uint64(userImage.ImageID))
	if err == nil && img != nil {
		if delErr := uc.DeleteImageUsecase.DeleteImage(img.Key); delErr != nil {
			return fmt.Errorf("failed to delete old image: %w", delErr)
		}
	}

	// xoa user image
	if err := uc.Repo.Delete(userImage.ID.String()); err != nil {
		return fmt.Errorf("failed to reset is_main: %w", err)
	}

	return nil
}
