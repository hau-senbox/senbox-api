package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
)

type GetRedirectUrlListUseCase struct {
	*repository.RedirectUrlRepository
}

func (receiver *GetRedirectUrlListUseCase) GetList(req request.GetRedirectUrlListRequest) ([]response.GetRedirectUrlListResponseData, *response.Pagination, error) {
	redirectUrls, paging, err := receiver.RedirectUrlRepository.GetList(req)
	if err != nil {
		return nil, nil, err
	}

	var urlListResponseData []response.GetRedirectUrlListResponseData
	for _, url := range redirectUrls {
		urlListResponseData = append(urlListResponseData, response.GetRedirectUrlListResponseData{
			Id:           url.ID,
			QRCode:       url.QRCode,
			TargetUrl:    url.TargetUrl,
			Password:     url.Password,
			Hint:         url.Hint,
			HashPassword: url.HashPassword,
			CreatedAt:    url.CreatedAt,
			UpdatedAt:    url.UpdatedAt,
		})
	}

	return urlListResponseData, paging, nil
}
