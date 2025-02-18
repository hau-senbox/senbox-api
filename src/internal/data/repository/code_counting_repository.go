package repository

import (
	"encoding/json"
	"errors"
	"math"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CodeCountingRepository struct {
}

func NewCodeCountingRepository() *CodeCountingRepository {
	return &CodeCountingRepository{}
}

func (receiver *CodeCountingRepository) Create(codeCounting *entity.SCodeCounting, db *gorm.DB) error {
	return db.Create(codeCounting).Error
}

func (receiver *CodeCountingRepository) Update(codeCounting *entity.SCodeCounting, db *gorm.DB) error {
	return db.Save(codeCounting).Error
}

func (receiver *CodeCountingRepository) FindByToken(token string, db *gorm.DB) (*entity.SCodeCounting, error) {
	var codeCounting entity.SCodeCounting
	err := db.Where("token = ?", token).First(&codeCounting).Error
	return &codeCounting, err
}

func (receiver *CodeCountingRepository) CreateForQuestion(question entity.SQuestion, db *gorm.DB) (string, error) {
	var att response.QuestionAttributes
	err := json.Unmarshal(question.Attributes, &att)
	if err != nil {
		log.Error(err)
		return "", err
	}

	var existing entity.SCodeCounting
	result := db.Where("token = ? AND deleted_at IS NULL", att.Value).
		Order("current_value desc").
		First(&existing)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// insert new
		codeCounting := entity.SCodeCounting{
			Token:        att.Value,
			CurrentValue: 1,
		}

		err = db.Create(&codeCounting).Error
		if err != nil {
			log.Error(err)
			return "", err
		}

		return codeCounting.Token + strconv.Itoa(codeCounting.CurrentValue), nil
	}

	codeCounting := entity.SCodeCounting{
		Token:        att.Value,
		CurrentValue: existing.CurrentValue + 1,
	}

	err = db.Create(&codeCounting).Error
	if err != nil {
		log.Error(err)
		return "", err
	}

	//Remove other records
	err = db.Where("token = ? AND id != ?", att.Value, codeCounting.ID).Delete(&entity.SCodeCounting{}).Error
	if err != nil {
		log.Error(err)
		return "", err
	}

	return att.Value + strconv.Itoa(codeCounting.CurrentValue), nil
}

func (receiver *CodeCountingRepository) CreateForQuestionWithID(questionId string, db *gorm.DB) (string, error) {
	var q entity.SQuestion

	err := db.Where("question_id = ?", questionId).First(&q).Error
	if err != nil {
		log.Error(err)
		return "", err
	}

	return receiver.CreateForQuestion(q, db)
}

func (receiver *CodeCountingRepository) ResetCodeCounting(req request.ResetCodeCountingRequest, db *gorm.DB) error {
	//Delete all
	err := db.Where("token = ?", req.Prefix).Delete(&entity.SCodeCounting{}).Error
	if err != nil {
		log.Error(err)
		return err
	}

	//Insert new
	codeCounting := entity.SCodeCounting{
		Token:        req.Prefix,
		CurrentValue: req.ResetTo,
	}

	return db.Create(&codeCounting).Error
}

func (receiver *CodeCountingRepository) GetCodeCountings(conn *gorm.DB, rq request.GetCodeCountingsRequest) ([]entity.SCodeCounting, response.Pagination, error) {
	//Select records in s_code_counting table where token starts with rq.Prefix and the value is maximum group by token
	var result []entity.SCodeCounting
	var paging response.Pagination
	if rq.PerPage == 0 {
		rq.PerPage = 12
	}
	if rq.PageNo == 0 {
		rq.PageNo = 1
	}

	// Case empty prefix searching keyword
	if rq.Keyword == "" {
		err := conn.
			Where("deleted_at IS NULL").
			Limit(rq.PerPage).
			Offset(rq.PerPage * (rq.PageNo - 1)).
			Order("id desc").
			Group("token").
			Find(&result).
			Error
		if err != nil {
			log.Error(err)
			return []entity.SCodeCounting{}, paging, err
		}
		// Count total records
		err = conn.Model(&entity.SCodeCounting{}).Where("deleted_at IS NULL").Group("token").Count(&paging.Total).Error
		if err != nil {
			log.Error(err)
			return []entity.SCodeCounting{}, paging, err
		}
		// Count total pages
		paging.TotalPage = int(math.Ceil(float64(paging.Total) / float64(rq.PerPage)))
	} else {
		err := conn.
			Limit(rq.PerPage).Offset(rq.PerPage*(rq.PageNo-1)).Where("token like ? AND deleted_at IS NULL", "%"+rq.Keyword+"%").Order("id desc").Group("token").Find(&result).Error
		if err != nil {
			log.Error(err)
			return []entity.SCodeCounting{}, paging, err
		}
		// Count total records
		err = conn.Model(&entity.SCodeCounting{}).Where("token like ? and deleted_at IS NULL", "%"+rq.Keyword+"%").Group("token").Count(&paging.Total).Error
		if err != nil {
			log.Error(err)
			return []entity.SCodeCounting{}, paging, err
		}
		// Count total pages
		paging.TotalPage = int(math.Ceil(float64(paging.Total) / float64(rq.PerPage)))
	}

	paging.Limit = rq.PerPage
	paging.Page = rq.PageNo

	return result, paging, nil
}

func (receiver *CodeCountingRepository) UpdateCodeCounting(conn *gorm.DB, rq request.UpdateCodeCountingRequest) error {
	//Update records in s_code_counting table where token = rq.Prefix
	var codeCounting entity.SCodeCounting
	codeCounting.ID = rq.ID
	codeCounting.CurrentValue = rq.ResetTo

	err := conn.Exec("UPDATE s_code_counting SET current_value = ? WHERE id = ?", codeCounting.CurrentValue, codeCounting.ID).Error
	if err != nil {
		log.Error(err)
		return err
	}

	//Delete records in s_code_counting table where token = rq.Prefix except rq.ID
	err = conn.Where("token = ?", codeCounting.Token).Not("id = ?", rq.ID).Delete(&entity.SCodeCounting{}).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
