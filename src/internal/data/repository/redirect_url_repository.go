package repository

import (
	"errors"
	"math"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RedirectUrlRepository struct {
	DBConn                 *gorm.DB
	DefaultRequestPageSize int
}

func (receiver *RedirectUrlRepository) Save(redirectUrl entity.SRedirectUrl) (*entity.SRedirectUrl, error) {
	err := receiver.DBConn.Create(&redirectUrl).Error
	if err != nil {
		return nil, err
	}
	return &redirectUrl, nil
}

func (receiver *RedirectUrlRepository) GetList(req request.GetRedirectUrlListRequest) ([]entity.SRedirectUrl, *response.Pagination, error) {
	var urls []entity.SRedirectUrl

	limit := receiver.DefaultRequestPageSize
	offset := 0
	if req.Limit != 0 {
		limit = req.Limit
	}
	if req.Page >= 0 {
		offset = req.Page * limit
	} else {
		return []entity.SRedirectUrl{}, &response.Pagination{
			Page:      req.Page,
			Limit:     limit,
			TotalPage: 0,
			Total:     0,
		}, errors.New("invalid page number")
	}
	var err error
	var count int64
	if req.Keyword != "" {
		err = receiver.DBConn.Where("qr_code LIKE ?", "%"+req.Keyword+"%").Offset(offset).Limit(limit).Find(&urls).Error
		receiver.DBConn.Model(&entity.SRedirectUrl{}).Where("qr_code LIKE ?", "%"+req.Keyword+"%").Count(&count)
	} else {
		err = receiver.DBConn.Offset(offset).Limit(limit).Find(&urls).Error
		receiver.DBConn.Model(&entity.SRedirectUrl{}).Count(&count)
	}

	if int64(req.Page) > count {
		return []entity.SRedirectUrl{}, &response.Pagination{
			Page:      req.Page,
			Limit:     limit,
			TotalPage: int(math.Ceil(float64(count) / float64(limit))),
			Total:     count,
		}, errors.New("invalid page number")
	}

	if err != nil {
		return nil, nil, err
	}

	return urls, &response.Pagination{
		Page:      req.Page,
		Limit:     limit,
		TotalPage: int(math.Ceil(float64(count) / float64(limit))),
		Total:     count,
	}, nil
}

func (receiver *RedirectUrlRepository) Delete(id uint64) error {
	return receiver.DBConn.Delete(&entity.SRedirectUrl{}, id).Error
}

func (receiver *RedirectUrlRepository) GetById(id uint64) (*entity.SRedirectUrl, error) {
	var url entity.SRedirectUrl
	err := receiver.DBConn.First(&url, id).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (receiver *RedirectUrlRepository) Update(form *entity.SRedirectUrl) error {
	return receiver.DBConn.Save(form).Error
}

func (receiver *RedirectUrlRepository) GetByQRCode(qrCode string) (*entity.SRedirectUrl, error) {
	var url entity.SRedirectUrl
	err := receiver.DBConn.Where("qr_code = ?", qrCode).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (receiver *RedirectUrlRepository) SaveRedirectUrl(qrCode string, targetUrl string, password string, status string, hint string, hashPwd *string) error {
	spreadSheetStatus, err := value.GetImportSpreadsheetStatusFromString(status)
	if err != nil {
		return err
	}
	switch spreadSheetStatus {
	case value.ImportSpreadsheetStatusNew:
		return receiver.saveNewRedirectUrl(qrCode, targetUrl, password, hint, hashPwd)
	case value.ImportSpreadsheetStatusDeleted:
		return receiver.DBConn.Where("qr_code = ?", qrCode).Delete(&entity.SRedirectUrl{}).Error
	case value.ImportSpreadsheetStatusSkip:
		return nil
	default:
		return errors.New("invalid status")
	}
}

func (receiver *RedirectUrlRepository) saveNewRedirectUrl(code string, url string, password string, hint string, hashPwd *string) error {
	var pwd *string = nil
	if password != "" {
		pwd = &password
	}
	redirectUrl := entity.SRedirectUrl{
		QRCode:       code,
		TargetUrl:    url,
		Password:     pwd,
		Hint:         hint,
		HashPassword: hashPwd,
	}
	return receiver.DBConn.Table("s_redirect_url").Clauses(
		clause.OnConflict{Columns: []clause.Column{{Name: "qr_code"}},
			DoUpdates: clause.AssignmentColumns([]string{"target_url", "password", "hint", "hash_password"}),
		}).Create(&redirectUrl).Error
}
