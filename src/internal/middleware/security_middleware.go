package middleware

import (
	"errors"
	"net/http"
	"sen-global-api/internal/data/repository"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type SecuredMiddleware struct {
	SessionRepository repository.SessionRepository
}

func (receiver SecuredMiddleware) Secured() gin.HandlerFunc {
	return func(context *gin.Context) {
		authorizationHeader := context.GetHeader("Authorization")

		// get app language
		appLanguage := uint(1) // default
		if header := context.GetHeader("X-App-Language"); header != "" {
			log.Info("App language from header:", header)
			if val, err := strconv.Atoi(header); err == nil {
				appLanguage = uint(val)
			}
		}
		context.Writer.Header().Set("X-App-Language", strconv.Itoa(int(appLanguage)))
		context.Set("app_language", appLanguage)

		// check header
		if len(authorizationHeader) == 0 {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// parse token
		tokenString := strings.Split(authorizationHeader, " ")[1]
		token, err := receiver.SessionRepository.ValidateToken(tokenString)

		if err != nil {
			// phân biệt lỗi hết hạn (expired) với lỗi khác
			if errors.Is(err, jwt.ErrTokenExpired) {
				context.AbortWithStatus(http.StatusUnauthorized) // 401
			} else {
				context.AbortWithStatus(http.StatusForbidden) // 403
			}
			return
		}

		if token.Valid {
			userID, err := receiver.SessionRepository.ExtractUserIDFromToken(tokenString)
			if err != nil {
				context.AbortWithStatus(http.StatusForbidden)
				return
			}
			context.Set("user_id", *userID)
			context.Set("token", tokenString)
			context.Next()
		} else {
			context.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func (receiver SecuredMiddleware) ValidateSuperAdminRole() gin.HandlerFunc {
	return func(context *gin.Context) {
		authorizationHeader := context.GetHeader("Authorization")
		if len(authorizationHeader) == 0 {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]
		token, err := receiver.SessionRepository.ValidateToken(tokenString)
		if err != nil {
			context.AbortWithStatus(http.StatusForbidden)
		} else if token.Valid {
			tokenData, err := receiver.SessionRepository.GetDataFromToken(token)
			if err != nil {
				context.AbortWithStatus(http.StatusForbidden)
			}
			if lo.Contains(tokenData.Roles, "SuperAdmin") {
				context.Set("user_id", tokenData.UserID)
				context.Next()
			} else {
				context.AbortWithStatus(http.StatusForbidden)
			}
		} else {
			log.Info("Token is not valid")
			context.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
