package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthorizeUseCase struct {
	*repository.UserEntityRepository
	*repository.DeviceRepository
	repository.SessionRepository
}

func (receiver AuthorizeUseCase) LoginInputDao(req request.UserLoginRequest) (*response.LoginResponseData, error) {
	user, _ := receiver.GetByUsername(request.GetUserEntityByUsernameRequest{Username: req.Username})
	if user == nil {
		log.Info("No user has username matches", req.Username)
		return nil, errors.New("user not found")
	}

	err := receiver.VerifyPassword(req.Password, user.Password)

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !receiver.VerifyRoleAccesses(user, "SuperAdmin", "Admin", "Staff") {
		return nil, errors.New("you don't have access to login")
	}

	token, err := receiver.GenerateToken(*user)
	if err != nil {
		return nil, errors.New("cannot generate token")
	}

	//authMiddleware := jwtauth.JwtMiddleware()
	//token := authMiddleware.TokenGen(user.UserID)
	return token, nil
}

func (receiver AuthorizeUseCase) UserLoginUsecase(req request.UserLoginFromDeviceReqest, loginType value.LoginType) (*response.LoginResponseData, error) {
	user, err := receiver.GetByUsername(request.GetUserEntityByUsernameRequest{Username: req.Username})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	reqRegisterDevice := request.RegisterDeviceRequest{
		DeviceUUID: req.DeviceUUID,
		InputMode:  string(value.InfoInputTypeBarcode),
	}

	if err := receiver.CheckUserDeviceExist(request.RegisteringDeviceForUser{
		UserID:   user.ID.String(),
		DeviceID: req.DeviceUUID,
	}); err == nil {
		_, err = receiver.RegisteringDeviceForUser(user, reqRegisterDevice)
		if err != nil {
			return nil, err
		}
	}

	err = receiver.VerifyPassword4LoginQr(req.Password, user.Password, loginType)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, err := receiver.GenerateToken(*user)
	if err != nil {
		return nil, errors.New("cannot generate token")
	}

	//authMiddleware := jwtauth.JwtMiddleware()
	//token := authMiddleware.TokenGen(user.UserID)
	return token, nil
}
