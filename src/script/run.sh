#!/usr/bin/env bash
echo "Start backup database sen_master_db..."
mysqldump -u sen_master sen_master_db > sen_master_db.sql
tar -zcvf "sen_master_db.sql.tar.gz" sen_master_db.sql
echo "Finished backup database sen_master_db..."

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Configure GO env
    set GOOS=linux
    kill -9 $(lsof -t -i:8003)
elif [[ "$OSTYPE" == "darwin"* ]]; then
  # Configure GO env
  set GOOS=darwin
  lsof -t -i tcp:8003 | xargs kill
else
    echo "Unknown OS"
    exit 1
fi

go get

swag init -g cmd/global-api/main.go

# Compile the internal
go build -o sen-global-api ./cmd/global-api/main.go

# Run the internal in the background
rm -rf logs
nohup ./sen-global-api config/config.yaml > /dev/null 2>&1&