package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"sen-global-api/internal/domain/entity"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func GeneralLoggerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// đọc body request
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// lấy tất cả params từ gin.Context
		params := make(map[string]string)
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}

		// convert sang JSON
		paramsJSON, _ := json.Marshal(params)

		// lấy tất cả headers
		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			// header có thể nhiều value, gộp lại bằng ","
			headers[k] = strings.Join(v, ",")
		}
		headersJSON, _ := json.Marshal(headers)

		// tạo log entry ban đầu
		logEntry := &entity.DataLog{
			ID:        uuid.New(),
			Endpoint:  c.FullPath(),
			Method:    c.Request.Method,
			Payload:   datatypes.JSON(bodyBytes),
			Status:    "PENDING",
			CreatedAt: time.Now(),
			Params:    datatypes.JSON(paramsJSON),
			Headers:   datatypes.JSON(headersJSON),
		}
		db.Create(logEntry)
		c.Set("general_log_id", logEntry.ID)

		// wrap writer để capture response
		respBody := &bytes.Buffer{}
		writer := &bodyLogWriter{body: respBody, ResponseWriter: c.Writer}
		c.Writer = writer

		// tiếp tục request
		c.Next()

		// update log sau khi handler xong
		if logID, exists := c.Get("general_log_id"); exists {
			status := "SUCCESS"
			var errMsg string
			if len(c.Errors) > 0 {
				status = "FAIL"
				errMsg = c.Errors.String()
			}

			db.Model(&entity.DataLog{}).
				Where("id = ?", logID).
				Updates(map[string]interface{}{
					"status":        status,
					"error_message": errMsg,
					"response":      datatypes.JSON(respBody.Bytes()),
				})
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // copy response vào buffer
	return w.ResponseWriter.Write(b) // ghi response ra client
}
