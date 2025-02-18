package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetCompanyUseCase struct {
	*repository.CompanyRepository
}

func (receiver *GetCompanyUseCase) GetCompanyById(id uint) (*entity.SCompany, error) {
	return receiver.CompanyRepository.GetByID(id)
}
