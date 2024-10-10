package config

import (
	"sen-global-api/pkg/common"
)

type GoogleConfig struct {
	UserCredentialsFilePath     string   `env-required:"true" yaml:"user_credentials_file_path" env:"GOOGLE_CREDENTIALS_USER_FILE_PATH"`
	UploaderCredentialsFilePath string   `env-required:"true" yaml:"uploader_credentials_file_path" env:"GOOGLE_CREDENTIALS_UPLOADER_FILE_PATH"`
	Scopes                      []string `env-required:"true" yaml:"scopes" env:"GOOGLE_SCOPES"`
	SpreadsheetId               string   `env-required:"true" yaml:"spreadsheet_id" env:"GOOGLE_SPREADSHEET_ID"`
	FirstColumn                 string   `env-required:"true" yaml:"first_column" env:"GOOGLE_FIRST_COLUMN"`
	FirstRow                    int      `env-required:"true" yaml:"first_row" env:"GOOGLE_FIRST_ROW"`
}

type SMTPConfig struct {
	Host     string `env-required:"true" yaml:"host" env:"SMTP_HOST"`
	Port     int    `env-required:"true" yaml:"port" env:"SMTP_PORT"`
	Username string `env-required:"true" yaml:"username" env:"SMTP_USERNAME"`
	Password string `env-required:"true" yaml:"password" env:"SMTP_PASSWORD"`
}

type Messaging struct {
	ServiceAccount string `env-required:"true" yaml:"service_account"`
}

type AppConfig struct {
	Config                          *common.Config `yaml:"config"`
	Google                          *GoogleConfig  `yaml:"google_config"`
	AuthorizeEncryptKey             string         `env-required:"true" yaml:"authorize_encrypt_key" env:"AUTHORIZE_ENCRYPT_KEY"`
	TokenExpireDurationInHour       int            `env-required:"true" yaml:"token_expire_duration_in_hour" env:"TOKEN_EXPIRE_DURATION_IN_HOUR"`
	DefaultRequestPageSize          int            `env-required:"true" yaml:"default_request_page_size" env:"DEFAULT_REQUEST_PAGE_SIZE"`
	OutputSpreadsheetUrl            string         `env-required:"true" yaml:"output_spreadsheet_url" env:"OUTPUT_SPREADSHEET_URL"`
	CronJobInterval                 string         `env-required:"true" yaml:"cron_job_interval" env:"CRON_JOB_INTERVAL"`
	DefaultCronJobIntervalInMinutes uint8          `env-required:"true" yaml:"default_cron_job_interval_in_minutes" env:"DEFAULT_CRON_JOB_INTERVAL"`
	SMTP                            SMTPConfig     `yaml:"smtp"`
	Messaging                       Messaging      `yaml:"messaging"`
}
