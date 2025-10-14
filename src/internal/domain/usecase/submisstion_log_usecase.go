package usecase

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm" //
)

type SubmissionLogUseCase struct {
	SubmissionLogRepository *repository.SubmissionLogsRepository
	UserRepository          *repository.UserEntityRepository
	AuthorizeEncryptKey     string
}

func NewSubmissionLogUseCase(cfg config.AppConfig, db *gorm.DB) *SubmissionLogUseCase {
	return &SubmissionLogUseCase{
		SubmissionLogRepository: &repository.SubmissionLogsRepository{
			DBConn: db,
		},
		UserRepository: &repository.UserEntityRepository{
			DBConn: db,
		},
		AuthorizeEncryptKey: cfg.AuthorizeEncryptKey,
	}
}

func (u *SubmissionLogUseCase) GetSubmissionsFormLogs(
	ctx *gin.Context,
	startStr, endStr, qrCode, customID string,
	page, limit int,
) (map[string]interface{}, error) {

	var startTime, endTime *time.Time
	now := time.Now()

	if startStr != "" {
		st, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format")
		}
		startTime = &st
	}

	if endStr != "" {
		et, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format")
		}
		endTime = &et
	}

	if startTime != nil && endTime == nil {
		endTime = &now
	}

	data, total, err := u.SubmissionLogRepository.GetSubmissionsFormLogs(ctx, startTime, endTime, qrCode, page, limit)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for _, result := range data {
		var payload map[string]interface{}
		if err := json.Unmarshal(result.Payload, &payload); err != nil {
			continue
		}

		var headers map[string]string
		if err := json.Unmarshal(result.Headers, &headers); err != nil {
			continue
		}

		qr := ""
		deviceID := ""
		if val, ok := payload["qr_code"].(string); ok {
			qr = val
		}
		if val, ok := payload["device_id"].(string); ok {
			deviceID = val
		}

		var userID string
		if tokenStr, ok := headers["Authorization"]; ok && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			claims := jwt.MapClaims{}
			token, _ := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(u.AuthorizeEncryptKey), nil
			})
			if token != nil && token.Valid {
				if id, ok := claims["user_id"].(string); ok {
					userID = id
				}
			}
		}

		var custom string
		var nickname string
		if userID != "" {
			userReq := request.GetUserEntityByIDRequest{ID: userID}
			user, err := u.UserRepository.GetByID(userReq)
			if err == nil && user != nil {
				custom = user.CustomID
				nickname = user.Nickname
			}

			if customID != "" && user.CustomID != customID {
				continue
			}
		}

		results = append(results, map[string]interface{}{
			"nickname":   nickname,
			"qr_code":    qr,
			"device_id":  deviceID,
			"user_id":    userID,
			"custom_id":  custom,
			"created_at": result.CreatedAt,
		})
	}

	totalPages := (total + limit - 1) / limit
	pagination := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"hasNext":    page < totalPages,
		"hasPrev":    page > 1,
	}

	return map[string]interface{}{
		"data":       results,
		"pagination": pagination,
	}, nil

}

func (u *SubmissionLogUseCase) GetSubmissionsFormLogsBySubmit(
	ctx *gin.Context,
	startStr, endStr, qrCode, customID string,
	page, limit int,
) (map[string]interface{}, error) {

	var startTime, endTime *time.Time
	now := time.Now()

	if startStr != "" {
		st, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format")
		}
		startTime = &st
	}

	if endStr != "" {
		et, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format")
		}
		endTime = &et
	}

	if startTime != nil && endTime == nil {
		endTime = &now
	}

	data, total, err := u.SubmissionLogRepository.GetSubmissionsFormLogsBySubmit(ctx, startTime, endTime, qrCode, customID, page, limit)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for _, result := range data {
		var payload map[string]interface{}
		if err := json.Unmarshal(result.Payload, &payload); err != nil {
			continue
		}

		var headers map[string]string
		if err := json.Unmarshal(result.Headers, &headers); err != nil {
			continue
		}

		qr := ""
		if val, ok := payload["qr_code"].(string); ok {
			qr = val
		}

		customID := ""
		if val, ok := payload["user_custom_id"].(string); ok {
			customID = val
		}

		var userID string
		if tokenStr, ok := headers["Authorization"]; ok && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			claims := jwt.MapClaims{}
			token, _ := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(u.AuthorizeEncryptKey), nil
			})
			if token != nil && token.Valid {
				if id, ok := claims["user_id"].(string); ok {
					userID = id
				}
			}
		}

		var nickname string
		if userID != "" {
			userReq := request.GetUserEntityByIDRequest{ID: userID}
			user, err := u.UserRepository.GetByID(userReq)
			if err == nil && user != nil {
				nickname = user.Nickname
			}

			if customID != "" && user.CustomID != customID {
				continue
			}
		}

		var matchedAnswers []string

		if ansArr, ok := payload["answers"].([]interface{}); ok {
			for _, a := range ansArr {
				if ansMap, ok := a.(map[string]interface{}); ok {
					if answer, ok := ansMap["answer"].(string); ok {
						matched, _ := regexp.MatchString(`^[0-9]{4}-[0-9]{2}-[0-9]{2}( [0-9]{2}:[0-9]{2}:[0-9]{2})?$`, answer)
						if matched {
							matchedAnswers = append(matchedAnswers, answer)
						}
					}
				}
			}
		}

		results = append(results, map[string]interface{}{
			"nickname":        nickname,
			"qr_code":         qr,
			"user_id":         userID,
			"custom_id":       customID,
			"created_at":      result.CreatedAt,
			"matched_answers": matchedAnswers,
		})
	}

	totalPages := (total + limit - 1) / limit
	pagination := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"hasNext":    page < totalPages,
		"hasPrev":    page > 1,
	}

	return map[string]interface{}{
		"data":       results,
		"pagination": pagination,
	}, nil

}
