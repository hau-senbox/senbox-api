package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/pkg/monitor"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type UpdateOutputTemplateSettingForTeacherUseCase struct {
	*repository.SettingRepository
	config.AppConfig
}

func (receiver *UpdateOutputTemplateSettingForTeacherUseCase) Execute(req request.UpdateOutputTemplateRequest) error {
	//Download spreadsheet file
	pwd, err := os.Getwd()
	if err != nil {
		monitor.SendMessageViaTelegram(fmt.Sprintf("Error getting current directory: %s", err))
		return err
	}
	srv, err := drive.NewService(context.Background(), option.WithCredentialsFile(pwd+"/credentials/google_service_account.json"))
	if err != nil {
		log.Debug("Unable to access Drive API:", err)
	}

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		log.Error("failed to parse spreadsheet id from sync devices sheet")
		return errors.New("invalid spreadsheet url")
	}

	spreadsheetID := match[1]

	// Export the Google Spreadsheet as a CSV file
	resp, err := srv.Files.Export(spreadsheetID, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet").Download()
	if err != nil {
		monitor.SendMessageViaTelegram(fmt.Sprintf("Error downloading spreadsheet: %s", err))
		return err
	}
	defer resp.Body.Close()

	// Create the output file
	file, err := os.Create(pwd + "/config/output_template_teacher.xlsx")
	if err != nil {
		monitor.SendMessageViaTelegram(fmt.Sprintf("Error creating output file: %s", err))
		return err
	}
	defer file.Close()

	// Copy the response body to the output file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		monitor.SendMessageViaTelegram(fmt.Sprintf("Error copying spreadsheet to output file: %s", err))
		return err
	}

	//Update setting
	err = receiver.UpdateOutputTemplateSettingForTeacher(spreadsheetID)
	if err != nil {
		return err
	}

	return nil
}
