package repository

import (
	"gorm.io/gorm"
	"sen-global-api/internal/domain/entity"
)

type OutputRepository struct {
	DBConn *gorm.DB
}

func NewOutputRepository(dbConn *gorm.DB) *OutputRepository {
	return &OutputRepository{DBConn: dbConn}
}

func (r *OutputRepository) GetOutputByValue1AndValue2(value1 string, value2 string) (string, error) {
	type Output struct {
		Value1        string
		Value2        string
		SpreadsheetID string
	}

	var output Output

	_ = r.DBConn.Table("s_output").Where("value1 = ? and value2 = ?", value1, value2).Scan(&output).Error

	return output.SpreadsheetID, nil
}

func (r *OutputRepository) GetTeacherOutputByValue2AndValue3(value2 string, value3 string) (string, error) {
	type Output struct {
		Value1        string
		Value2        string
		SpreadsheetID string
	}

	var output Output

	_ = r.DBConn.Table("s_teacher_output").Where("value2 = ? and value3 = ?", value2, value3).Scan(&output).Error

	return output.SpreadsheetID, nil
}

type CreateOutputParams struct {
	Value1        string `json:"value1" required:"true"`
	Value2        string `json:"value2" required:"true"`
	SpreadsheetID string `json:"spreadsheet_id" required:"true"`
}

func (r *OutputRepository) Create(params CreateOutputParams) (*entity.SOutput, error) {
	var output = &entity.SOutput{
		Value1:        params.Value1,
		Value2:        params.Value2,
		SpreadsheetID: params.SpreadsheetID,
	}

	err := r.DBConn.Create(output).Error

	return output, err
}

func (r *OutputRepository) CreateTeacherOutput(params CreateOutputParams) (*entity.STeacherOutput, error) {
	var output = &entity.STeacherOutput{
		Value2:        params.Value1,
		Value3:        params.Value2,
		SpreadsheetID: params.SpreadsheetID,
	}

	err := r.DBConn.Create(output).Error

	return output, err
}
