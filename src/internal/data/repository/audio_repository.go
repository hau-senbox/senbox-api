package repository

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/entity"
)

type AudioRepository struct {
	DBConn *gorm.DB
}

func NewAudioRepository(dbConn *gorm.DB) *AudioRepository {
	return &AudioRepository{DBConn: dbConn}
}

func (receiver *AudioRepository) GetAllByIDs(ids []int) ([]entity.SAudio, error) {
	var audios []entity.SAudio
	err := receiver.DBConn.Model(entity.SAudio{}).Find(&audios).Where("id IN (?)", ids).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.SAudio{}, nil
		}
		log.Error("AudioRepository.GetAllByIDs: " + err.Error())
		return nil, errors.New("failed to get audios")
	}

	return audios, err
}

func (receiver *AudioRepository) GetAllByName(audioName string) ([]entity.SAudio, error) {
	var audios []entity.SAudio
	err := receiver.DBConn.Model(entity.SAudio{}).Find(&audios).Where("audio_name = ", audioName).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.SAudio{}, nil
		}
		log.Error("AudioRepository.GetAllByName: " + err.Error())
		return nil, errors.New("failed to get audios")
	}

	return audios, err
}

func (receiver *AudioRepository) GetByID(id uint64) (*entity.SAudio, error) {
	var audios entity.SAudio
	err := receiver.DBConn.Model(&entity.SAudio{}).Where("id = ?", id).First(&audios).Error
	if err != nil {
		log.Error("AudioRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get audio")
	}

	return &audios, nil
}

func (receiver *AudioRepository) GetByKey(key string) (*entity.SAudio, error) {
	var audios entity.SAudio
	err := receiver.DBConn.Model(&entity.SAudio{}).Where("`key` = ?", key).First(&audios).Error
	if err != nil {
		log.Error("AudioRepository.GetByKey: " + err.Error())
		return nil, errors.New("failed to get audio")
	}

	return &audios, nil
}

func (receiver *AudioRepository) CreateAudios(audios []entity.SAudio) error {
	if err := receiver.DBConn.Model(entity.SAudio{}).Create(&audios).Error; err != nil {
		log.Error("AudioRepository.CreateAudios: " + err.Error())
		return errors.New("failed to create audios")
	}

	return nil
}

func (receiver *AudioRepository) CreateAudio(audio entity.SAudio) error {
	if err := receiver.DBConn.Model(entity.SAudio{}).Create(&audio).Error; err != nil {
		log.Error("AudioRepository.CreateAudios: " + err.Error())
		return errors.New("failed to create audios")
	}

	return nil
}

func (receiver *AudioRepository) DeleteAudio(id uint64) error {
	if err := receiver.DBConn.Model(entity.SAudio{}).Where("id = ?", id).Delete(&entity.SAudio{}).Error; err != nil {
		log.Error("AudioRepository.DeleteAudio: " + err.Error())
		return errors.New("failed to delete audio")
	}

	return nil
}
