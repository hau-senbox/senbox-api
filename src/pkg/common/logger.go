package common

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

func InitLogger(level string) error {

	log.SetFormatter(&log.JSONFormatter{})

	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(log.TraceLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	}

	return nil
}
