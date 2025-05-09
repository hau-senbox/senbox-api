package middleware

import (
	"net/http"
	"sen-global-api/internal/data/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type SecuredMiddleware struct {
	SessionRepository repository.SessionRepository
}

func (receiver SecuredMiddleware) Secured() gin.HandlerFunc {
	return func(context *gin.Context) {
		authorizationHeader := context.GetHeader("Authorization")
		if len(authorizationHeader) == 0 {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authorizationHeader, " ")[1]
		token, err := receiver.SessionRepository.ValidateToken(tokenString)
		if err != nil {
			context.AbortWithStatus(http.StatusForbidden)
		} else if token.Valid {
			userId, err := receiver.SessionRepository.ExtractUserIdFromToken(tokenString)
			if err != nil {
				context.AbortWithStatus(http.StatusForbidden)
			}
			context.Set("user_id", *userId)
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
			roles, userId, err := receiver.SessionRepository.GetRoleFromToken(token)
			if err != nil {
				context.AbortWithStatus(http.StatusForbidden)
			}
			if lo.Contains(roles, "SuperAdmin") {
				context.Set("user_id", userId)
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
