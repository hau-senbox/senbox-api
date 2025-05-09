package middleware

import (
	"net/http"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/monitor"

	"github.com/go-errors/errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func RecoveryHandler(c *gin.Context, err any) {
	goErr := errors.Wrap(err, 2).ErrorStack()
	monitor.SendMessageViaTelegram("URGENT: SENBOX is panicking: ", goErr)
	log.Error(goErr)
	c.AbortWithStatusJSON(500, response.FailedResponse{
		Code:  http.StatusInternalServerError,
		Error: "Internal server error",
	})
}
