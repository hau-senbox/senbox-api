package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
)

func ToGetValuesAppCurrentResponse(values *entity.ValuesAppCurrent) *response.GetValuesAppCurrentResponse {
	return &response.GetValuesAppCurrentResponse{
		Value1:   values.Value1,
		Value2:   values.Value2,
		Value3:   values.Value3,
		ImageKey: values.ImageKey,
	}
}
