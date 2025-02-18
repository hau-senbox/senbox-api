package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"
)

type SearchFormsUseCase struct {
	FormRepository *repository.FormRepository
}

func (receiver *SearchFormsUseCase) SearchForms(keyword string) ([]response.GetFormListResponseData, error) {
	forms, err := receiver.FormRepository.GetMatchFormByNote(keyword)
	if err != nil {
		return nil, err
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

	return formList, nil
}
