package repository

import (
	"errors"
	"math"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/parameters"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FormRepository struct {
	DBConn                 *gorm.DB
	DefaultRequestPageSize int
}

func (receiver *FormRepository) Create(form *entity.SForm) (*entity.SForm, error) {
	err := receiver.DBConn.Create(form).Error
	if err != nil {
		return nil, err
	}

	return form, nil
}

func (receiver *FormRepository) GetFormById(id uint64) (*entity.SForm, error) {
	var form entity.SForm
	err := receiver.DBConn.Where("id = ?", id).First(&form).Error
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (receiver *FormRepository) SaveForm(request parameters.SaveFormParams) (*entity.SForm, error) {
	form := entity.SForm{
		Note:           request.Note,
		Name:           request.Name,
		SpreadsheetUrl: request.SpreadsheetUrl,
		SpreadsheetId:  request.SpreadsheetId,
		Password:       request.Password,
		Status:         value.Active,
		SheetName:      request.SheetName,
	}
	err := receiver.DBConn.Table("s_form").Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "note"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"spreadsheet_url", "spreadsheet_id", "password", "status",
				"name", "sheet_name",
			}),
		}).Create(&form).Error
	if err != nil {
		return nil, err
	}
	return &form, err
}

func (receiver *FormRepository) DeleteForm(formId uint64) error {
	form := entity.SForm{}
	err := receiver.DBConn.Where("id = ?", formId).First(&form).Error
	if err != nil {
		return err
	}
	err = receiver.DBConn.Delete(&form).Error
	if err != nil {
		return err
	}

	return nil
}

func (receiver *FormRepository) GetFormList(request request.GetFormListRequest) ([]entity.SForm, *response.Pagination, error) {
	var forms []entity.SForm

	limit := receiver.DefaultRequestPageSize
	offset := 0
	if request.Limit != 0 {
		limit = request.Limit
	}
	if request.Page >= 0 {
		offset = request.Page * limit
	} else {
		return []entity.SForm{}, &response.Pagination{
			Page:      request.Page,
			Limit:     limit,
			TotalPage: 0,
			Total:     0,
		}, errors.New("invalid page number")
	}
	var err error
	var count int64
	if request.Keyword != "" {
		err = receiver.DBConn.Where("note LIKE ?", "%"+request.Keyword+"%").Offset(offset).Limit(limit).Find(&forms).Error
		receiver.DBConn.Model(&entity.SForm{}).Where("note LIKE ?", "%"+request.Keyword+"%").Count(&count)
	} else {
		err = receiver.DBConn.Offset(offset).Limit(limit).Find(&forms).Error
		receiver.DBConn.Model(&entity.SForm{}).Count(&count)
	}

	if int64(request.Page) > count {
		return []entity.SForm{}, &response.Pagination{
			Page:      request.Page,
			Limit:     limit,
			TotalPage: int(math.Ceil(float64(count) / float64(limit))),
			Total:     count,
		}, errors.New("invalid page number")
	}

	if err != nil {
		return nil, nil, err
	}

	return forms, &response.Pagination{
		Page:      request.Page,
		Limit:     limit,
		TotalPage: int(math.Ceil(float64(count) / float64(limit))),
		Total:     count,
	}, nil
}

func (receiver *FormRepository) GetFormByQRCode(code string) (*entity.SForm, error) {
	var form entity.SForm
	err := receiver.DBConn.Where("note = ?", code).First(&form).Error
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (receiver *FormRepository) UpdateForm(form *entity.SForm) (*entity.SForm, error) {
	err := receiver.DBConn.Save(form).Error
	if err != nil {
		return nil, err
	}

	return form, nil
}

func (receiver *FormRepository) GetMatchFormByNote(keyword string) ([]entity.SForm, error) {
	var forms []entity.SForm
	err := receiver.DBConn.Where("note LIKE ?", "%"+keyword+"%").Find(&forms).Error
	if err != nil {
		return nil, err
	}
	return forms, err
}

func (receiver *FormRepository) DeleteFormByNote(note string) error {
	err := receiver.DBConn.Where("note = ?", note).Delete(&entity.SForm{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (receiver *FormRepository) FindByCode(code string) (entity.SForm, error) {
	var form entity.SForm
	err := receiver.DBConn.Where("note = ?", code).First(&form).Error
	if err != nil {
		return entity.SForm{}, err
	}
	return form, nil
}
