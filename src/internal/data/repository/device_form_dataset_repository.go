package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sen-global-api/internal/domain/entity"
)

type DeviceFormDatasetRepository struct {
	DBConn *gorm.DB
}

func (r DeviceFormDatasetRepository) Save(datasets []entity.SDeviceFormDataset) {
	r.DBConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"set", "question_date", "question_time", "question_date_time", "question_duration_forward", "question_duration_backward", "question_scale", "question_qr_code", "question_selection", "question_text", "question_count", "question_number", "question_photo", "question_multiple_choice", "question_button_count", "question_single_choice", "question_button_list", "question_message_box", "question_show_pic", "question_button", "question_play_video", "question_qr_code_front", "question_choice_toggle", "question_web"}),
	}).Create(datasets)
}

func (r DeviceFormDatasetRepository) GetDatasetByDeviceId(deviceId string) []entity.SDeviceFormDataset {
	var datasets []entity.SDeviceFormDataset
	r.DBConn.Where("device_id = ?", deviceId).Find(&datasets)
	return datasets
}

func (r DeviceFormDatasetRepository) GetDatasetByDeviceIdAndSet(deviceId string, set string) (entity.SDeviceFormDataset, error) {
	var dataset entity.SDeviceFormDataset
	err := r.DBConn.Where("device_id = ? AND `set` = ?", deviceId, set).First(&dataset).Error
	if err != nil {
		return dataset, err
	}
	return dataset, nil
}
