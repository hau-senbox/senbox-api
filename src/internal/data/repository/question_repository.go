package repository

import (
	"encoding/json"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/model"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QuestionRepository struct {
	DBConn *gorm.DB
}

type CreateQuestionParams struct {
	ID               string
	Question         string
	QuestionType     string
	Attributes       string
	Status           string
	Set              string
	EnableOnMobile   value.QuestionForMobile
	QuestionUniqueID *string
	Key              string
	DB               string
}

func (receiver *QuestionRepository) Create(params []CreateQuestionParams) ([]entity.SQuestion, error) {
	if len(params) == 0 {
		return []entity.SQuestion{}, nil
	}
	rawQuestions, err := receiver.unmarshalQuestions(params)
	if err != nil {
		return nil, err
	}

	err = receiver.DBConn.Create(&rawQuestions).Error

	if err != nil {
		return nil, err
	}

	return rawQuestions, nil
}

func (receiver *QuestionRepository) GetMemoryComponentValue(componentName string) (*entity.MemoryComponentValue, error) {
	var componentValue entity.MemoryComponentValue
	err := receiver.DBConn.Where("component_name = ?", componentName).First(&componentValue).Error

	return &componentValue, err
}

func (receiver *QuestionRepository) CreateMemoryComponentValues(values []entity.MemoryComponentValue) error {
	if len(values) == 0 {
		return nil
	}

	err := receiver.DBConn.Model(&entity.MemoryComponentValue{}).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "component_name"}},
			DoNothing: true,
		}).Create(&values).Error
	if err != nil {
		return err
	}

	return nil
}

func (receiver *QuestionRepository) CreateMemoryComponentValuesDuplicate(values []entity.MemoryComponentValue) error {
	if len(values) == 0 {
		return nil
	}

	err := receiver.DBConn.Model(&entity.MemoryComponentValue{}).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "component_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&values).Error
	if err != nil {
		return err
	}

	return nil
}

func (receiver *QuestionRepository) unmarshalQuestions(params []CreateQuestionParams) ([]entity.SQuestion, error) {
	var rawQuestions []entity.SQuestion
	for _, param := range params {
		rawQuestion, err := receiver.unmarshalQuestion(param)
		if err != nil {
			log.Error(err)
			continue
		}
		rawQuestions = append(rawQuestions, *rawQuestion)
	}
	return rawQuestions, nil
}

func (receiver *QuestionRepository) unmarshalQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	qType, err := value.GetQuestionType(param.QuestionType)
	if err != nil {
		return nil, err
	}
	switch qType {
	case value.QuestionDate:
		return receiver.unmarshalDateQuestion(param)
	case value.QuestionTime:
		return receiver.unmarshalTimeQuestion(param)
	case value.QuestionDateTime:
		return receiver.unmarshalDateTimeQuestion(param)
	case value.QuestionDurationForward:
		return receiver.unmarshalDurationForwardQuestion(param)
	case value.QuestionDurationBackward:
		return receiver.unmarshalDurationBackwardQuestion(param)
	case value.QuestionScale:
		return receiver.unmarshalScaleQuestion(param)
	case value.QuestionQRCode:
		return receiver.unmarshalQRCodeQuestion(param)
	case value.QuestionSelection:
		return receiver.unmarshalSelectionQuestion(param)
	case value.QuestionInText:
		return receiver.unmarshalTextQuestion(param)
	case value.QuestionCount:
		return receiver.unmarshalCountQuestion(param)
	case value.QuestionNumber:
		return receiver.unmarshalNumberQuestion(param)
	case value.QuestionPhoto:
		return receiver.unmarshalPhotoQuestion(param)
	case value.QuestionMultipleChoice:
		return receiver.unmarshalMultipleChoiceQuestion(param)
	case value.QuestionSingleChoice:
		return receiver.unmarshalSingleChoiceQuestion(param)
	case value.QuestionButtonCount:
		return receiver.unmarshalButtonCountQuestion(param)
	case value.QuestionButtonList:
		return receiver.unmarshalButtonsQuestion(param)
	case value.QuestionMessageBox:
		return receiver.unmarshalMessageBoxQuestion(param)
	case value.QuestionShowPic:
		return receiver.unmarshalShowPicsQuestion(param)
	case value.QuestionButton:
		return receiver.unmarshalButtonQuestion(param)
	case value.QuestionPlayVideo:
		return receiver.unmarshalPlayVideoQuestion(param)
	case value.QuestionQRCodeFront:
		return receiver.unmarshalQRCodeFrontQuestion(param)
	case value.QuestionChoiceToggle:
		return receiver.unmarshalChoiceToggleQuestion(param)
	case value.QuestionSection:
		return receiver.unmarshalSectionQuestion(param)
	case value.QuestionFormSection:
		return receiver.unmarshalFormSectionQuestion(param)
	case value.QuestionFormSendImmediately:
		return receiver.unmarshalFormSendImmediately(param)
	case value.QuestionSignature:
		return receiver.unmarshalSignatureQuestion(param)
	case value.QuestionWeb:
		return receiver.unmarshalWebQuestion(param)
	case value.QuestionWebUser:
		return receiver.unmarshalWebUserQuestion(param)
	case value.QuestionSignUpPreSetValue1:
		return receiver.unmarshalSignUpPreSetValue1(param)
	case value.QuestionSignUpPreSetValue2:
		return receiver.unmarshalSignUpPreSetValue2(param)
	case value.QuestionSignUpPreSetValue3:
		return receiver.unmarshalSignUpPreSetValue3(param)
	case value.QuestionDraggableList:
		return receiver.unmarshalDraggableListQuestion(param)
	case value.QuestionSendMessage:
		return receiver.unmarshalSendMessageQuestion(param)
	case value.QuestionSendNotification:
		return receiver.unmarshalSendNotification(param)
	case value.QuestionCodeCounting:
		return receiver.unmarshalCodeCountingQuestion(param)
	case value.QuestionRandomizer:
		return receiver.unmarshalRandomizerQuestion(param)
	case value.QuestionDocument:
		return receiver.unmarshalDocumentQuestion(param)
	case value.QuestionQRCodeGenerator:
		return receiver.unmarshalQRCodeGeneratorQuestion(param)
	case value.MemoryText:
		return receiver.unmarshalMemoryTextQuestion(param)
	default:
		return receiver.unmarshalUserQuestion(param)
	}
}

func (receiver *QuestionRepository) unmarshalDateQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Attr struct {
		Value string `json:"value"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate date format of attr.Value
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalTimeQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Attr struct {
		Value string `json:"value"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	//TODO: validate time format of attr.Value
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalDateTimeQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Attr struct {
		Value string `json:"value"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalDurationForwardQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalDurationBackwardQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Attr struct {
		Value string `json:"value"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalScaleQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Attr struct {
		Number int `json:"number"`
		Steps  int `json:"steps"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalQRCodeQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalSelectionQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Option struct {
		Name string `json:"name"`
	}
	type Attr struct {
		Options []Option `json:"options"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) SaveQuestions(questions []entity.SQuestion) ([]entity.SQuestion, error) {
	err := receiver.DBConn.Table("s_question").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "question_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"question", "question_type", "attributes", "status", "updated_at"}),
		}).Create(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, err
}

func (receiver *QuestionRepository) unmarshalTextQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
		Key:              param.Key,
		DB:               param.DB,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalCountQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
		Key:              param.Key,
		DB:               param.DB,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalNumberQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalPhotoQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) FindByID(id string) (*entity.SQuestion, error) {
	var question entity.SQuestion
	err := receiver.DBConn.Table("s_question").Where("question_id = ?", id).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, err
}

func (receiver *QuestionRepository) unmarshalMultipleChoiceQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Option struct {
		Name string `json:"name"`
	}
	type Attr struct {
		Options []Option `json:"options"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalButtonsQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalButtonQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalMessageBoxQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalShowPicsQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalButtonCountQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) GetQuestionsByFormID(id uint64) ([]model.FormQuestionItem, error) {
	var questions []model.FormQuestionItem
	//goland:noinspection ALL
	rows, err := receiver.DBConn.Raw("SELECT s_question.id as id,"+
		"s_question.question_type as question_type, s_question.attributes as attributes, s_question.status as status, "+
		"s_question.created_at as created_at, s_question.updated_at as updated_at, s_form_question.order as `order`, "+
		"s_form_question.answer_required as answer_required, s_form_question.answer_remember as answer_remember, s_question.question as question,"+
		"s_question.enable_on_mobile as enable_on_mobile, s_question.question_unique_id as question_unique_id, "+
		"`s_question`.`key` as `key`, s_question.db as db "+
		"FROM s_question RIGHT JOIN s_form_question ON s_form_question.question_id = s_question.id WHERE s_form_question.form_id = ? AND s_question.status = ? ORDER BY `order` ASC", id, value.Active).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var question model.FormQuestionItem
		err := receiver.DBConn.ScanRows(rows, &question)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, err
}

func (receiver *QuestionRepository) GetAllQuestions() ([]entity.SQuestion, error) {
	var questions []entity.SQuestion
	err := receiver.DBConn.Table("s_question").Find(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, err
}

func (receiver *QuestionRepository) unmarshalSingleChoiceQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Option struct {
		Name string `json:"name"`
	}
	type Attr struct {
		Options []Option `json:"options"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) DeleteQuestionsFormNote(note string) error {
	return receiver.DBConn.Exec("DELETE s FROM s_question s INNER JOIN s_form_question fq ON fq.question_id = s.id INNER JOIN s_form f ON f.id = fq.form_id WHERE f.note = ?", note).Error
}

func (receiver *QuestionRepository) unmarshalQRCodeFrontQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalPlayVideoQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalChoiceToggleQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Option struct {
		Name string `json:"name"`
	}
	type Attr struct {
		Options []Option `json:"options"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalSectionQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalUserQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		Set:              param.Set,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalFormSectionQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil

}

func (receiver *QuestionRepository) unmarshalFormSendImmediately(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) GetQuestionsByIDs(IDs []string) ([]entity.SQuestion, error) {
	var questions []entity.SQuestion
	err := receiver.DBConn.Where("id IN (?)", IDs).Find(&questions).Error
	return questions, err
}

func (receiver *QuestionRepository) unmarshalSignatureQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:             uuid.MustParse(param.ID),
		Question:       param.Question,
		QuestionType:   param.QuestionType,
		Attributes:     datatypes.JSON(param.Attributes),
		Status:         status,
		EnableOnMobile: param.EnableOnMobile,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalWebUserQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalWebQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalSignUpPreSetValue1(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalSignUpPreSetValue2(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalSignUpPreSetValue3(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalDraggableListQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Option struct {
		Name string `json:"name"`
	}
	type Attr struct {
		Options []Option `json:"options"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:             uuid.MustParse(param.ID),
		Question:       param.Question,
		QuestionType:   param.QuestionType,
		Attributes:     datatypes.JSON(param.Attributes),
		Status:         status,
		EnableOnMobile: param.EnableOnMobile,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalSendNotification(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Option struct {
		Name string `json:"name"`
	}
	type Attr struct {
		Options []Option `json:"options"`
	}
	var attr = Attr{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalSendMessageQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	type Msg struct {
		Email          []string `json:"email"`
		Value3         []string `json:"value3"`
		ShowMessageBox bool     `json:"showMessageBox"`
	}
	type Messaging struct {
		Data Msg `json:"messaging"`
	}
	var attr = Messaging{}
	err := json.Unmarshal([]byte(param.Attributes), &attr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(attr)
	//TODO: validate time format of attr.Value
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
	}

	return &question, err
}

func (receiver *QuestionRepository) unmarshalCodeCountingQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:             uuid.MustParse(param.ID),
		Question:       param.Question,
		QuestionType:   param.QuestionType,
		Attributes:     datatypes.JSON(param.Attributes),
		Status:         status,
		EnableOnMobile: param.EnableOnMobile,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalRandomizerQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:             uuid.MustParse(param.ID),
		Question:       param.Question,
		QuestionType:   param.QuestionType,
		Attributes:     datatypes.JSON(param.Attributes),
		Status:         status,
		EnableOnMobile: param.EnableOnMobile,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalDocumentQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:             uuid.MustParse(param.ID),
		Question:       param.Question,
		QuestionType:   param.QuestionType,
		Attributes:     datatypes.JSON(param.Attributes),
		Status:         status,
		EnableOnMobile: param.EnableOnMobile,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalQRCodeGeneratorQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:             uuid.MustParse(param.ID),
		Question:       param.Question,
		QuestionType:   param.QuestionType,
		Attributes:     datatypes.JSON(param.Attributes),
		Status:         status,
		EnableOnMobile: param.EnableOnMobile,
	}

	return &question, nil
}

func (receiver *QuestionRepository) unmarshalMemoryTextQuestion(param CreateQuestionParams) (*entity.SQuestion, error) {
	status, err := value.GetStatusFromString(param.Status)
	if err != nil {
		return nil, err
	}
	var question = entity.SQuestion{
		ID:               uuid.MustParse(param.ID),
		Question:         param.Question,
		QuestionType:     param.QuestionType,
		Attributes:       datatypes.JSON(param.Attributes),
		Status:           status,
		EnableOnMobile:   param.EnableOnMobile,
		QuestionUniqueID: param.QuestionUniqueID,
		Key:              param.Key,
		DB:               param.DB,
	}

	return &question, nil
}
