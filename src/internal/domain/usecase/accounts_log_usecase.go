package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
)

type AccountsLogUseCase struct {
	AccountsLogRepository *repository.AccountsLogRepository
}

func (u *AccountsLogUseCase) CreateAccountsLog(req request.CreateAccountsLogRequest) error {
	// check valid type
	if !value.AccountsLogType(req.Type).IsValid() {
		return fmt.Errorf("invalid type: %s", req.Type)
	}

	accountsLog := &entity.AccountsLog{
		Type:           value.AccountsLogType(req.Type),
		UserID:         req.UserID,
		DeviceID:       req.DeviceID,
		OrganizationID: req.OrganizationID,
	}
	return u.AccountsLogRepository.Create(accountsLog)
}
