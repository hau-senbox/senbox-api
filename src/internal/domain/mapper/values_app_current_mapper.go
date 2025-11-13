package mapper

import (
	"sen-global-api/internal/domain/response"
)

func ToGetValuesAppCurrentResponse(value1, value2, value3 string, imageUrl *string) *response.GetValuesAppCurrentResponse {
	return &response.GetValuesAppCurrentResponse{
		Value1:   value1,
		Value2:   value2,
		Value3:   value3,
		ImageUrl: imageUrl,
	}
}
