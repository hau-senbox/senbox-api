#!/bin/bash
echo "Backup database...in 60 seconds"
curl -X GET http://localhost:8003/infra/backup > ~/go/src/sen-master-api/src/logs/backup.log 2>&1
