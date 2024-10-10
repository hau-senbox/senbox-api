package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sen-global-api/internal/domain/entity"
)

type SecureAppMiddleware struct {
	DB *gorm.DB
}

func NewSecureAppMiddleware(db *gorm.DB) SecureAppMiddleware {
	return SecureAppMiddleware{
		DB: db,
	}
}

func (receiver SecureAppMiddleware) Secure() gin.HandlerFunc {
	return func(context *gin.Context) {
		appKeyArgs := context.GetHeader("App-Key")
		if len(appKeyArgs) == 0 {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}

		var appKey entity.SAppKey
		err := receiver.DB.Where("app_key = ?", appKeyArgs).First(&appKey).Error
		if err != nil {
			context.AbortWithStatus(http.StatusForbidden)
			return
		}
		if appKey.ID == 0 {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		context.Next()
	}
}
