## Prerequisites
* [Go](https://go.dev/doc/install) (1.18+)
* [MariaDB](https://mariadb.org/) (10.5+) - [Install](https://computingforgeeks.com/install-mariadb-on-centos-7-centos-8/)
* [Swag](https://github.com/swaggo/swag)

## Create database
- https://mariadb.com/kb/en/create-database/

## Project structure
```bash
src
    ├── config
    │   ├── config.yaml
    ├── credentials
    │   ├── google_service_account.json
    │   ├── uploader_google_service_account.json
```

## Setup Environment variable
```
export SENBOX_API_KEY="978acdfaa332139818f3467664f7b72c"
```


## Configuration
#### config.yaml content
```yaml
config:
  app:
    name: 'Sen-App'
    version: '1.0.0'

  http:
    port: '8003'

  logger:
    log_level: 'debug'
    rollbar_env: 'gopher'


  mysql:
    host: '15.235.184.98'
    port: '3306'
    user: 'sen_master'
    password: '6Yw4rmmaDUkpZtdm'
    database: 'sen_master_db'
    max_conn: 10
    max_idle_conn: 5
    max_lifetime_conn: 1000000

  environment: development

google_config:
  user_credentials_file_path: 'credentials/google_service_account.json'
  uploader_credentials_file_path: 'credentials/uploader_service_account.json'
  scopes:
    - 'https://spreadsheets.google.com/feeds'
    - 'https://www.googleapis.com/auth/spreadsheets'
  spreadsheet_id: '1L0cuLpeOoJlxYCBLcY_DCrDrDUoGSMXRIZvJQqLtg4E'
  first_column: 'K'
  first_row: 12

authorize_encrypt_key: 'senbox-dev-secret-key'
token_expire_duration_in_hour: 1000
default_request_page_size: 12
output_spreadsheet_url: 'https://docs.google.com/spreadsheets/d/1L0cuLpeOoJlxYCBLcY_DCrDrDUoGSMXRIZvJQqLtg4E/edit#gid=753138406'
cron_job_interval: "@every 5m"
default_cron_job_interval_in_minutes: 30
smtp:
  host: 'smtp.gmail.com'
  port: 587
  username: 'senboxapp@gmail.com'
  password: 'rgkwxldypbhtoflb'
messaging:
  service_account: 'credentials/google_service_account.json'
```

## Run localhost
```bash
cd src
go run /cmd/global-api/main.go
```

# Deploy
### Login to server
```
ssh root@103.15.51.114
```

### Change working directory
```
cd go/src/github.com/sen-master-api/src
```

### Run
```bash
./script/run.sh
```