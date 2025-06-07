package repository

import (
	"errors"
	"fmt"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	gofn "github.com/tiendc/gofn"
	"golang.org/x/crypto/bcrypt"
)

type SessionRepository struct {
	*OrganizationRepository
	AuthorizeEncryptKey   string
	TokenExpireTimeInHour time.Duration
}

func (receiver *SessionRepository) VerifyPassword(password string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func (receiver *SessionRepository) VerifyRoleAccesses(user *entity.SUserEntity, roles ...string) bool {
	for _, role := range roles {
		for _, userRole := range user.Roles {
			if strings.ToLower(userRole.Role.String()) == strings.ToLower(role) {
				return true
			}
		}
	}

	return false
}

func (receiver *SessionRepository) GenerateToken(user entity.SUserEntity) (*response.LoginResponseData, error) {
	roles := gofn.MapSliceToMap(user.Roles, func(role entity.SRole) (int64, string) {
		return role.ID, role.Role.String()
	})
	organizations := gofn.MapSliceToMap(user.Organizations, func(organization entity.SOrganization) (int64, string) {
		return organization.ID, organization.OrganizationName
	})

	expirationTime := time.Now().Add(receiver.TokenExpireTimeInHour * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":       user.ID.String(),
		"username":      user.Username,
		"roles":         strings.Join(gofn.MapValues(roles), ", "),
		"organizations": strings.Join(gofn.MapValues(organizations), ", "),
		"exp":           expirationTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(receiver.AuthorizeEncryptKey))
	if err != nil {
		return nil, err
	}

	userOrgs := make([]string, 0)
	for _, organization := range user.Organizations {
		userOrg, err := receiver.GetUserOrgInfo(user.ID.String(), organization.ID)
		if err != nil {
			return nil, err
		}
		if userOrg.IsManager {
			userOrgs = append(userOrgs, organization.OrganizationName)
		}
	}

	isSuperAdmin := gofn.Contain(gofn.MapValues(roles), "SuperAdmin")
	return &response.LoginResponseData{
		UserId:        user.ID.String(),
		Username:      user.Username,
		IsSuperAdmin:  isSuperAdmin,
		Organizations: userOrgs,
		Token:         tokenString,
		Expired:       time.Now().Add(receiver.TokenExpireTimeInHour * time.Hour),
	}, nil
}

func (receiver *SessionRepository) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(receiver.AuthorizeEncryptKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, nil
	}

	return nil, errors.New("invalid token")
}

func (receiver *SessionRepository) ExtractUserIdFromToken(tokenString string) (*string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(receiver.AuthorizeEncryptKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["user_id"].(string); ok {
			return &userId, nil
		}
	}

	return nil, err
}

type TokenData struct {
	UserId        string
	Roles         []string
	Organizations []string
}

func (receiver *SessionRepository) GetDataFromToken(token *jwt.Token) (*TokenData, error) {
	token, err := jwt.Parse(token.Raw, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(receiver.AuthorizeEncryptKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["user_id"].(string)
		roles := claims["roles"].(string)
		organizations := claims["organizations"].(string)

		return &TokenData{
			UserId:        userId,
			Roles:         strings.Split(roles, ", "),
			Organizations: strings.Split(organizations, ", "),
		}, nil
	}

	return nil, err
}

func (receiver *SessionRepository) GenerateTokenByDevice(device entity.SDevice) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)
	tokenClaims["device_uuid"] = device.ID
	tokenClaims["sub"] = 1
	tokenString, err := token.SignedString([]byte(receiver.AuthorizeEncryptKey))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["device_uuid"] = device.ID
	rt, err := refreshToken.SignedString([]byte(receiver.AuthorizeEncryptKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, rt, nil
}

func (receiver *SessionRepository) ExtractDeviceIdFromToken(tokenString string) (*string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(receiver.AuthorizeEncryptKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["device_uuid"].(string); ok {
			return &userId, nil
		}
	}

	return nil, err
}

func (receiver *SessionRepository) GetDeviceIDFromRefreshToken(accessToken string) (string, error) {
	uuid, err := receiver.ExtractDeviceIdFromToken(accessToken)
	if err != nil {
		return "", err
	}

	return *uuid, nil
}

func (receiver *SessionRepository) GeneratePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
