package usecase

import (
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"

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

func (uc *MessageLanguageUseCase) GetMessageLanguages4GW(typeStr string, typeID string) ([]*response.GetMessageLanguages4GWResponse, error) {
	// Lấy tất cả messages theo type + typeID
	messages, err := uc.messageLanguageRepo.GetByTypeAndTypeID(typeStr, typeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message languages: %w", err)
	}

	// Gom nhóm theo LanguageID
	grouped := make(map[uint]map[string]string)
	for _, msg := range messages {
		if _, ok := grouped[msg.LanguageID]; !ok {
			grouped[msg.LanguageID] = make(map[string]string)
		}
		grouped[msg.LanguageID][msg.Key] = msg.Value
	}

	// Build response
	var responses []*response.GetMessageLanguages4GWResponse
	for langID, msgs := range grouped {
		responses = append(responses, &response.GetMessageLanguages4GWResponse{
			LangID:   langID,
			Contents: msgs,
		})
	}

	return responses, nil
}

func (uc *MessageLanguageUseCase) GetMessageLanguage4GW(typeStr string, typeID string, languageID uint) (*response.GetMessageLanguages4GWResponse, error) {

	// Lấy tất cả messages theo type + typeID + languageID
	messages, err := uc.messageLanguageRepo.GetByTypeAndTypeIDAndLanguage(typeStr, typeID, languageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message languages: %w", err)
	}

	// Build response
	contents := make(map[string]string)
	for _, msg := range messages {
		if msg.TypeID == typeID {
			contents[msg.Key] = msg.Value
		}
	}

	return &response.GetMessageLanguages4GWResponse{
		LangID:   languageID,
		Contents: contents,
	}, nil
}

func (uc *MessageLanguageUseCase) DeleteMessageLanguage4GW(typeStr string, typeID string) error {
	return uc.messageLanguageRepo.DeleteByTypeAndTypeID(typeStr, typeID)
}
