package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type UserImageUsecase struct {
	Repo *repository.UserImageRepository
}

// Tạo ảnh người dùng
func (u *UserImageUsecase) Create(userImage *entity.SUserImage) error {
	return u.Repo.Create(userImage)
}

// Lấy ảnh theo ID
func (u *UserImageUsecase) GetByID(id int) (*entity.SUserImage, error) {
	return u.Repo.FindByID(id)
}

// Lấy tất cả ảnh
func (u *UserImageUsecase) GetAll() ([]entity.SUserImage, error) {
	return u.Repo.FindAll()
}

// Cập nhật ảnh
func (u *UserImageUsecase) Update(userImage *entity.SUserImage) error {
	return u.Repo.Update(userImage)
}

// Xóa ảnh
func (u *UserImageUsecase) Delete(id int) error {
	return u.Repo.Delete(id)
}
