package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type ChangePasswordUseCase struct {
	*repository.SessionRepository
	*repository.UserRepository
}

func (c *ChangePasswordUseCase) ValidateCurrentPassword(userId string, oldPassword string) (entity.SUser, error) {
	user, err := c.UserRepository.FindUserById(&userId)
	if err != nil {
		return entity.SUser{}, err
	}

	err = c.SessionRepository.VerifyPassword(oldPassword, user.Password)
	if err != nil {
		return entity.SUser{}, err
	}

	return *user, nil
}

func (c *ChangePasswordUseCase) VerifyNewPassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	return nil
}

func (c *ChangePasswordUseCase) ChangePassword(user entity.SUser, newPassword string) error {
	hashedPassword, err := c.SessionRepository.GeneratePassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = c.UserRepository.UpdateUser(&user)
	if err != nil {
		return err
	}

	return nil
}
