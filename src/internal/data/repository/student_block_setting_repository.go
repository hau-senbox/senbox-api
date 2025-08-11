package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type StudentBlockSettingRepository struct {
	DBConn *gorm.DB
}

func (r *StudentBlockSettingRepository) Create(setting *entity.StudentBlockSetting) error {
	return r.DBConn.Create(setting).Error
}

func (r *StudentBlockSettingRepository) Update(setting *entity.StudentBlockSetting) error {
	return r.DBConn.Save(setting).Error
}

func (r *StudentBlockSettingRepository) GetByStudentID(studentID string) (*entity.StudentBlockSetting, error) {
	var setting entity.StudentBlockSetting
	err := r.DBConn.Where("student_id = ?", studentID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *StudentBlockSettingRepository) GetByID(id int) (*entity.StudentBlockSetting, error) {
	var setting entity.StudentBlockSetting
	err := r.DBConn.First(&setting, id).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *StudentBlockSettingRepository) Delete(id int) error {
	return r.DBConn.Delete(&entity.StudentBlockSetting{}, id).Error
}

func (r *StudentBlockSettingRepository) GetIsDeactiveByStudentID(studentID string) (bool, error) {
	var isDeactive bool
	err := r.DBConn.Model(&entity.StudentBlockSetting{}).
		Select("is_deactive").
		Where("student_id = ?", studentID).
		Scan(&isDeactive).Error
	if err != nil {
		return false, err
	}
	return isDeactive, nil
}
