package usecase

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"net/smtp"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"
)

type SendEmailUseCase struct {
	config.SMTPConfig
	*repository.SettingRepository
	*sheet.Writer
}

func (receiver *SendEmailUseCase) SendEmail(target string, subject string, content string, device entity.SDevice) error {
	_, err := mail.ParseAddress(target)
	if err != nil {
		return err
	}

	smtpServer := receiver.Host
	auth := smtp.PlainAuth(
		"",
		receiver.Username,
		receiver.Password,
		smtpServer,
	)

	from := mail.Address{Name: "SENBOX", Address: receiver.Username}
	to := mail.Address{Name: target, Address: target}
	title := subject

	body := content

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = title

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err = smtp.SendMail(
		smtpServer+":"+strconv.Itoa(receiver.Port),
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)

	receiver.logHistory(target, subject, device)

	return err
}

func (receiver *SendEmailUseCase) logHistory(target string, subject string, device entity.SDevice) {
	setting, err := receiver.GetEmailSettings()
	if err != nil {
		return
	}

	if setting == nil {
		return
	}

	type EmailSetting struct {
		SpreadsheetID string `json:"spreadsheet_id"`
	}

	var emailSettings *EmailSetting = nil
	err = json.Unmarshal([]byte(setting.Settings), &emailSettings)

	if err != nil {
		return
	}

	log := make([][]interface{}, 0)
	log = append(log, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
	log = append(log, []interface{}{device.ID})
	log = append(log, []interface{}{device.DeviceName})
	log = append(log, []interface{}{target})
	log = append(log, []interface{}{nil})
	log = append(log, []interface{}{nil})
	log = append(log, []interface{}{nil})
	log = append(log, []interface{}{subject})

	_, err = receiver.WriteRanges(sheet.WriteRangeParams{
		Range:     "History!K11",
		Dimension: "COLUMNS",
		Rows:      log,
	}, emailSettings.SpreadsheetID)

	if err != nil {
		return
	}
}

func (receiver *SendEmailUseCase) SendMessage(subject string, bccList []string, body string) error {
	smtpServer := receiver.Host
	auth := smtp.PlainAuth(
		"",
		receiver.Username,
		receiver.Password,
		smtpServer,
	)

	from := mail.Address{Name: "SENBOX", Address: receiver.Username}

	header := make(map[string]string)
	header["From"] = from.String()
	header["Bcc"] = strings.Join(bccList, ",")
	header["Subject"] = subject

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpServer+":"+strconv.Itoa(receiver.Port),
		auth,
		from.Address,
		bccList,
		[]byte(message),
	)

	return err
}
