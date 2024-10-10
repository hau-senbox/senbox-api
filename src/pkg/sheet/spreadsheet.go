package sheet

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"os"
	"sen-global-api/config"
)

type Spreadsheet struct {
	Reader *Reader
	Writer *Writer
}

func NewUserSpreadsheet(config config.AppConfig, contex context.Context) (*Spreadsheet, error) {
	log.Debug(config.Google.UserCredentialsFilePath)

	credentialsInByte, err := os.ReadFile(config.Google.UserCredentialsFilePath)
	if err != nil {
		return nil, err
	}

	jwtConfig, err := google.JWTConfigFromJSON(credentialsInByte, config.Google.Scopes...)
	if err != nil {
		return nil, err
	}
	client := jwtConfig.Client(contex)
	sheetsService, err := sheets.New(client)
	if err != nil {
		return nil, err
	}

	return &Spreadsheet{
		Reader: &Reader{sheetsService: sheetsService},
		Writer: &Writer{sheetsService: sheetsService},
	}, nil
}

func NewUploaderSpreadsheet(config config.AppConfig, contex context.Context) (*Spreadsheet, error) {
	log.Debug(config.Google.UserCredentialsFilePath)

	credentialsInByte, err := os.ReadFile(config.Google.UploaderCredentialsFilePath)
	if err != nil {
		return nil, err
	}

	jwtConfig, err := google.JWTConfigFromJSON(credentialsInByte, config.Google.Scopes...)
	if err != nil {
		return nil, err
	}
	client := jwtConfig.Client(contex)
	sheetsService, err := sheets.New(client)
	if err != nil {
		return nil, err
	}

	return &Spreadsheet{
		Reader: &Reader{sheetsService: sheetsService},
		Writer: &Writer{sheetsService: sheetsService},
	}, nil
}
