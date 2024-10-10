package repository

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"time"
)

type SessionRepository struct {
	AuthorizeEncryptKey   string
	TokenExpireTimeInHour time.Duration
}

func (receiver *SessionRepository) VerifyPassword(password string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func (receiver *SessionRepository) GenerateToken(user entity.SUser) (*response.LoginResponseData, error) {
	roles := user.Role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.UserId,
		"username": user.Username,
		"roles":    roles,
	})
	tokenString, err := token.SignedString([]byte(receiver.AuthorizeEncryptKey))

	if err != nil {
		return nil, err
	}
	return &response.LoginResponseData{
		UserId:   user.UserId,
		UserName: user.Username,
		Token:    tokenString,
		Expired:  time.Now().Add(receiver.TokenExpireTimeInHour * time.Hour),
	}, nil
}

func (receiver *SessionRepository) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(receiver.AuthorizeEncryptKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userId, ok := claims["userId"].(string); ok {
			return &userId, nil
		}
	}

	return nil, err
}

func (receiver *SessionRepository) GetRoleFromToken(token *jwt.Token) (uint, string, error) {
	token, err := jwt.Parse(token.Raw, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(receiver.AuthorizeEncryptKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"].(string)
		roles := claims["roles"].(float64)

		return uint(roles), userId, nil
	}

	return 0, "", err
}

func (receiver *SessionRepository) GenerateTokenByDevice(device entity.SDevice) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)
	tokenClaims["device_uuid"] = device.DeviceId
	tokenClaims["sub"] = 1
	tokenString, err := token.SignedString([]byte(receiver.AuthorizeEncryptKey))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["device_uuid"] = device.DeviceId
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
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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
