package logger

import (
	"io"
	"os"
	"sen-global-api/config"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(config *config.AppConfig) error {
	logLevel := log.InfoLevel
	log.SetLevel(logLevel)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/senbox.log",
		MaxSize:    50, // MB
		MaxBackups: 20,
		MaxAge:     14, // days
		Compress:   true,
	}
	log.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))

	log.SetReportCaller(true)

	log.SetFormatter(&log.TextFormatter{
		PadLevelText:    true,
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return nil
}
