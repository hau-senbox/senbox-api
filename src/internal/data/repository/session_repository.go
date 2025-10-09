package repository

import (
	"errors"
	"fmt"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"strings"
	"time"

	"github.com/tiendc/gofn"

	"github.com/golang-jwt/jwt/v4"
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

func (receiver *SessionRepository) VerifyPassword4LoginQr(password string, hashedPassword string, loginType value.LoginType) error {
	if loginType == value.ForScan {
		if password != hashedPassword {
			return errors.New("invalid QR password")
		}
	}
	if loginType == value.ForRegister {
		return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	}
	return nil
}

func (receiver *SessionRepository) VerifyRoleAccesses(user *entity.SUserEntity, roles ...string) bool {
	for _, role := range roles {
		for _, userRole := range user.Roles {
			if strings.EqualFold(userRole.Role.String(), role) {
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
	organizations := gofn.MapSliceToMap(user.Organizations, func(organization entity.SOrganization) (string, string) {
		return organization.ID.String(), organization.OrganizationName
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
		userOrg, err := receiver.GetUserOrgInfo(user.ID.String(), organization.ID.String())
		if err != nil {
			return nil, err
		}
		if userOrg.IsManager {
			userOrgs = append(userOrgs, organization.OrganizationName)
		}
	}

	isSuperAdmin := gofn.Contain(gofn.MapValues(roles), "SuperAdmin")
	return &response.LoginResponseData{
		UserID:        user.ID.String(),
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

	// if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	// 	return token, nil
	// }

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := claims.Valid(); err != nil {
			return nil, fmt.Errorf("token is expired or not valid: %w", err)
		}
		return token, nil
	}

	return nil, errors.New("invalid token")
}

func (receiver *SessionRepository) ExtractUserIDFromToken(tokenString string) (*string, error) {
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
		if userID, ok := claims["user_id"].(string); ok {
			return &userID, nil
		}
	}

	return nil, err
}

type TokenData struct {
	UserID        string
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
		userID := claims["user_id"].(string)
		roles := claims["roles"].(string)
		organizations := claims["organizations"].(string)

		return &TokenData{
			UserID:        userID,
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

func (receiver *SessionRepository) ExtractDeviceIDFromToken(tokenString string) (*string, error) {
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
		if userID, ok := claims["device_uuid"].(string); ok {
			return &userID, nil
		}
	}

	return nil, err
}

func (receiver *SessionRepository) GetDeviceIDFromRefreshToken(accessToken string) (string, error) {
	uuid, err := receiver.ExtractDeviceIDFromToken(accessToken)
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

func (receiver *SessionRepository) ExtractUserIDIgnoreExp(tokenString string) (*string, error) {
	// parse token nhưng không validate exp
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("cannot parse token: %w", err)
	}

	// cast claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(string); ok {
			return &userID, nil
		}
		return nil, errors.New("user_id not found in token claims")
	}

	return nil, errors.New("invalid token claims")
}
