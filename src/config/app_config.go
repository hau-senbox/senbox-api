package config

import (
	"sen-global-api/pkg/common"
)

type SenboxFormSubmitBucket struct {
	Domain               string `env-required:"true" yaml:"domain"`
	Region               string `env-required:"true" yaml:"region"`
	BucketName           string `env-required:"true" yaml:"bucket_name"`
	AccessKey            string `env-required:"true" yaml:"access_key"`
	SecretKey            string `env-required:"true" yaml:"secret_key"`
	CloudfrontKeyGroupID string `env-required:"true" yaml:"cloudfront_key_group_id"`
	CloudfrontKeyPath    string `env-required:"true" yaml:"cloudfront_key_path"`
}

type S3 struct {
	SenboxFormSubmitBucket SenboxFormSubmitBucket `env-required:"true" yaml:"senbox-form-submit-bucket"`
}

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
	S3                              S3             `yaml:"s3"`
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
