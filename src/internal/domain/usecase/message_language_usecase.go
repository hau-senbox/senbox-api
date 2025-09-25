package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	"gorm.io/gorm"
)

type MessageLanguageUseCase struct {
	messageLanguageRepo *repository.MessageLanguageRepository
}

func NewMessageLanguageUseCase(messageLanguageRepo *repository.MessageLanguageRepository) *MessageLanguageUseCase {
	return &MessageLanguageUseCase{
		messageLanguageRepo: messageLanguageRepo,
	}
}

func (uc *MessageLanguageUseCase) UploadMessageLanguage(req request.UploadMessageLanguageRequest) error {
	// check exist
	existMessage, err := uc.messageLanguageRepo.GetByTypeAndKeyAndLanguageAndTypeID(req.Type, req.Key, req.LanguageID, req.TypeID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existMessage != nil {
		// update
		existMessage.Value = req.Value
		return uc.messageLanguageRepo.Update(existMessage)
	} else {
		// create new
		newMessage := &entity.MessageLanguage{
			Type:       req.Type,
			Key:        req.Key,
			Value:      req.Value,
			LanguageID: req.LanguageID,
		}
		return uc.messageLanguageRepo.Create(newMessage)
	}
}

func (uc *MessageLanguageUseCase) UploadMessageLanguages(req request.UploadMessageLanguagesRequest) error {
	for _, msgReq := range req.MessageLanguages {
		// check exist
		existMessage, err := uc.messageLanguageRepo.GetByTypeAndKeyAndLanguageAndTypeID(
			msgReq.Type, msgReq.Key, msgReq.LanguageID, msgReq.TypeID,
		)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if existMessage != nil {
			// update
			existMessage.Value = msgReq.Value
			if err := uc.messageLanguageRepo.Update(existMessage); err != nil {
				return err
			}
		} else {
			// create new
			newMessage := &entity.MessageLanguage{
				TypeID:     msgReq.TypeID,
				Type:       msgReq.Type,
				Key:        msgReq.Key,
				Value:      msgReq.Value,
				LanguageID: msgReq.LanguageID,
			}
			if err := uc.messageLanguageRepo.Create(newMessage); err != nil {
				return err
			}
		}
	}

	return nil
}
