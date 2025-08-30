package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthorizeUseCase struct {
	*repository.UserEntityRepository
	*repository.DeviceRepository
	repository.SessionRepository
	*repository.OrganizationRepository
	UserEntityUseCase      *UserEntityUseCase
	DBConn                 *gorm.DB
	ManageUserLoginUseCase *ManageUserLoginUseCase
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

	// get organization_admin
	var orgAdminResp *response.OrganizationAdmin
	if len(user.Organizations) > 0 {
		// Lấy danh sách OrgID mà user là manager
		managedOrgIDs, err := user.GetManagedOrganizationIDs(receiver.DBConn)
		if err != nil {
			orgAdminResp = nil
		}

		// So sánh với các org đã preload, map sang OrganizationAdmin nếu khớp
		if len(managedOrgIDs) > 0 {
			for _, org := range user.Organizations {
				if lo.Contains(managedOrgIDs, org.ID.String()) {
					orgAdminResp = &response.OrganizationAdmin{
						ID:               org.ID.String(),
						OrganizationName: org.OrganizationName,
						Avatar:           org.Avatar,
						AvatarURL:        org.AvatarURL,
						Address:          org.Address,
						Description:      org.Description,
						CreatedAt:        org.CreatedAt,
						UpdatedAt:        org.UpdatedAt,
					}
					break
				}
			}
		}

	}
	token.OrganizationAdmin = orgAdminResp
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

	// reqRegisterDevice := request.RegisterDeviceRequest{
	// 	DeviceUUID: req.DeviceUUID,
	// 	InputMode:  string(value.InfoInputTypeBarcode),
	// }

	// if err := receiver.CheckUserDeviceExist(request.RegisteringDeviceForUser{
	// 	UserID:   user.ID.String(),
	// 	DeviceID: req.DeviceUUID,
	// }); err == nil {
	// 	_, err = receiver.RegisteringDeviceForUser(user, reqRegisterDevice)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// check user device login
	if req.DeviceUUID != "" {
		err := receiver.ManageUserLoginUseCase.ManageUserDeviceLogin(user.ID.String(), req.DeviceUUID)
		if err != nil {
			log.Error("AuthorizeUseCase.UserLoginUsecase.HandleUserDeviceLogin: " + err.Error())
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

func (receiver AuthorizeUseCase) SwitchToOrganizationAdmin(organizationID string) (*response.SwitchToOrganizationResponse, error) {
	// lấy thông tin manager theo organizationID
	manager, err := receiver.OrganizationRepository.GetManagerByOrganizationID(organizationID)
	if err != nil {
		log.Error("AuthorizeUseCase.SwitchToOrganizationAdmin: " + err.Error())
		return nil, errors.New("failed to get organization manager")
	}

	// sinh token cho manager
	tokenData, err := receiver.GenerateToken(manager.User) // manager.User là SUserEntity
	if err != nil {
		log.Error("AuthorizeUseCase.SwitchToOrganizationAdmin.GenerateToken: " + err.Error())
		return nil, errors.New("cannot generate token for manager")
	}

	// map sang response
	user, _ := receiver.UserEntityUseCase.MapUserInfoToResponse(manager.User)

	return &response.SwitchToOrganizationResponse{
		Token:   tokenData.Token,
		Expired: tokenData.Expired,
		User:    *user,
	}, nil
}

func (receiver AuthorizeUseCase) UserLogoutUsecase(c *gin.Context, req request.UserLogoutReqeust) error {

	// get user id by context
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		return errors.New("user_id not found")
	}
	// xoa user device login
	err := receiver.ManageUserLoginUseCase.ManageUserDeviceLogout(userIDRaw.(string), req.DeviceID)
	if err != nil {
		log.Error("AuthorizeUseCase.UserLogoutUsecase.ManageUserDeviceLogout: " + err.Error())
		return err
	}
	return nil
}
