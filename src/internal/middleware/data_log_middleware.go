package middleware

import (
	"bytes"
	"io"
	"sen-global-api/internal/domain/entity"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func DataLogMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// đọc body request (payload)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			// reset body để handler khác vẫn đọc được
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// capture response bằng custom writer
		respBody := &bytes.Buffer{}
		writer := &bodyLogWriter{body: respBody, ResponseWriter: c.Writer}
		c.Writer = writer

		// chạy tiếp handler
		c.Next()

		// lấy status và error nếu có
		status := "success"
		var errMsg *string
		if len(c.Errors) > 0 {
			status = "error"
			msg := c.Errors.String()
			errMsg = &msg
		}

		payload := bodyBytes
		if len(payload) == 0 {
			payload = []byte("{}")
		}

		// tạo log
		log := entity.DataLog{
			ID:           uuid.New(),
			Endpoint:     c.FullPath(),
			Method:       c.Request.Method,
			Payload:      datatypes.JSON(payload),
			Response:     datatypes.JSON(respBody.Bytes()),
			Status:       status,
			ErrorMessage: errMsg,
			CreatedAt:    start,
		}

		// save DB
		db.Create(&log)
	}
}

// custom writer để capture response
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // copy response vào buffer
	return w.ResponseWriter.Write(b) // ghi response ra client
}
