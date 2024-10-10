package usecase

import (
	"gorm.io/gorm"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
)

type GetCodeCountingsResult struct {
	Codes  []entity.SCodeCounting
	Paging response.Pagination
}

func GetCodeCountings(conn *gorm.DB, rq request.GetCodeCountingsRequest) (GetCodeCountingsResult, error) {
	repo := repository.NewCodeCountingRepository()

	c, p, err := repo.GetCodeCountings(conn, rq)

	return GetCodeCountingsResult{
		Codes:  c,
		Paging: p,
	}, err
}

func UpdateCodeCounting(conn *gorm.DB, rq request.UpdateCodeCountingRequest) error {
	repo := repository.NewCodeCountingRepository()

	return repo.UpdateCodeCounting(conn, rq)
}
