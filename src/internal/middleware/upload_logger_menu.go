package middleware

import (
	"bytes"
	"io"
	"sen-global-api/internal/domain/entity"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func MenuUploadLoggerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.FullPath(), "/v1/admin/menu/section") ||
			strings.HasPrefix(c.FullPath(), "/v1/admin/menu/child") ||
			strings.HasPrefix(c.FullPath(), "/v1/admin/menu/student") ||
			strings.HasPrefix(c.FullPath(), "/v1/admin/menu/teacher") ||
			strings.HasPrefix(c.FullPath(), "/v1/admin/menu/staff") ||
			strings.HasPrefix(c.FullPath(), "/v1/admin/menu/parent") {

			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			logEntry := &entity.MenuUploadLog{
				ID:        uuid.New(),
				Endpoint:  c.FullPath(),
				Method:    c.Request.Method,
				Payload:   datatypes.JSON(bodyBytes),
				Status:    "PENDING",
				CreatedAt: time.Now(),
			}

			db.Create(logEntry)
			c.Set("menu_log_id", logEntry.ID)
		}

		c.Next()

		if logID, exists := c.Get("menu_log_id"); exists {
			status := "SUCCESS"
			var errMsg string
			if len(c.Errors) > 0 {
				status = "FAIL"
				errMsg = c.Errors.String()
			}
			db.Model(&entity.MenuUploadLog{}).
				Where("id = ?", logID).
				Updates(map[string]interface{}{
					"status":        status,
					"error_message": errMsg,
				})
		}
	}
}
