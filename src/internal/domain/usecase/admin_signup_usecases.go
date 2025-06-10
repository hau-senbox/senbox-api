package usecase

import (
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/parameters"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AdminSignUpUseCases struct {
	SettingRepository *repository.SettingRepository
	FormRepository    *repository.FormRepository
	SpreadsheetReader *sheet.Reader
	config.AppConfig
	ImportFormsUseCase *ImportFormsUseCase
}

func (c *AdminSignUpUseCases) UpdateSignUpButton1(name string, value string) error {
	_, err := c.SettingRepository.UpdateSignUpButton1(name, value)

	return err
}

func (c *AdminSignUpUseCases) UpdateSignUpButton2(name string, value string) error {
	_, err := c.SettingRepository.UpdateSignUpButton2(name, value)
	return err
}

func (c *AdminSignUpUseCases) UpdateSignUpButton3(name string, value string) error {
	_, err := c.SettingRepository.UpdateSignUpButton3(name, value)
	return err
}

func (c *AdminSignUpUseCases) UpdateSignUpButton4(name string, value string) error {
	_, err := c.SettingRepository.UpdateSignUpButton4(name, value)
	return err
}

func (c *AdminSignUpUseCases) UpdateSignUpButton5(name string, value string) error {
	_, err := c.SettingRepository.UpdateSignUpButton5(name, value)
	return err
}

func (c *AdminSignUpUseCases) UpdateRegistrationForm(url string) error {
	f, err := c.importForm(url, "SENBOX.ORG/[SIGNUP]:default", "Sign Up 0")
	if err != nil {
		return err
	}
	_, err = c.SettingRepository.UpdateRegistrationForm(f.ID, url)
	return err
}

func (c *AdminSignUpUseCases) UpdateSignUpButtonConfiguration(url string) error {
	f, err := c.importForm(url, "SENBOX.ORG/[SIGNUP-BUTTON]:Configuration", "Configuration")
	if err != nil {
		return err
	}
	_, err = c.SettingRepository.UpdateSignUpButtonConfiguration(f.ID, url)
	return err
}

func (c *AdminSignUpUseCases) UpdateRegistrationSubmission(url string) error {
	_, err := c.SettingRepository.UpdateRegistrationSubmission(url)
	return err
}

func (c *AdminSignUpUseCases) UpdateRegistrationPreset2(url string) error {
	_, err := c.SettingRepository.UpdateRegistrationPreset2(url)
	return err
}

func (c *AdminSignUpUseCases) UpdateRegistrationPreset1(url string) error {
	_, err := c.SettingRepository.UpdateRegistrationPreset1(url)
	return err
}

func (c *AdminSignUpUseCases) importForm(spreadsheetUrl, note, sheetNameToRead string) (entity.SForm, error) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(spreadsheetUrl)

	if len(match) < 2 {
		log.Error("Import Sign Up Form Invalid spreadsheet url: ", spreadsheetUrl)
		return entity.SForm{}, fmt.Errorf("import sign up form invalid spreadsheet url: %s", spreadsheetUrl)
	}

	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestImportForm()
	values, err := c.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     sheetNameToRead + `!` + c.Google.FirstColumn + strconv.Itoa(c.Google.FirstRow-1) + `:Q`,
	})
	if err != nil || values == nil {
		log.Error(fmt.Sprintf("Error reading spreadsheet: %s - note : %s", err.Error(), note))
		return entity.SForm{}, err
	}

	var rawQuestions = make([]parameters.RawQuestion, 0)
	var formName = ""
	for index, row := range values {
		if index == 0 && len(row) >= 1 && cap(row) >= 1 {
			formName = row[0].(string)
			continue
		} else if len(row) >= 4 && cap(row) >= 4 && index > 1 && row[1].(string) != "" {
			if row[1].(string) == "" {
				continue
			}

			additionalInfo := ""
			if len(row) >= 6 {
				additionalInfo = row[5].(string)
			}
			required := "false"
			if len(row) >= 5 {
				required = row[4].(string)
			}
			item := parameters.RawQuestion{
				// ID:        strings.ToUpper(note) + "_" + spreadsheetId + "_" + row[0].(string),
				QuestionId:        uuid.NewString(),
				Question:          row[2].(string),
				Type:              row[1].(string),
				Attributes:        strings.ReplaceAll(row[3].(string), "\n", ""),
				AnswerRequired:    required,
				AdditionalOptions: additionalInfo,
				Status:            "1",
				RowNumber:         index + 1,
			}
			rawQuestions = append(rawQuestions, item)
		}
	}

	f, msg, err := c.ImportFormsUseCase.CreateSignUpForm(parameters.SaveFormParams{
		Note:           note,
		Name:           formName,
		SpreadsheetUrl: spreadsheetUrl,
		SpreadsheetId:  spreadsheetId,
		Password:       "",
		RawQuestions:   rawQuestions,
		SheetName:      sheetNameToRead,
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error creating form: %s - note : %s", err.Error(), note))
		return entity.SForm{}, err
	}

	log.Warning(msg)

	return *f, nil
}
