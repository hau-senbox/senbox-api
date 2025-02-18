package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
)

type GetFormListUseCase struct {
	*repository.FormRepository
}

func (receiver *GetFormListUseCase) GetFormList(request request.GetFormListRequest) ([]response.GetFormListResponseData, *response.Pagination, error) {
	forms, paging, err := receiver.FormRepository.GetFormList(request)
	if err != nil {
		return nil, nil, err
	}

	var formList []response.GetFormListResponseData
	for _, form := range forms {
		formList = append(formList, response.GetFormListResponseData{
			Id:          form.ID,
			Spreadsheet: form.SpreadsheetUrl,
			Password:    form.Password,
			Note:        form.Note,
			CreatedAt:   form.CreatedAt,
			UpdatedAt:   form.UpdatedAt,
		})
	}

	return formList, paging, nil
}
