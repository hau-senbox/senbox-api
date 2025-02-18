package usecase

import (
	"sen-global-api/pkg/job"
	"sen-global-api/pkg/sheet"

	"github.com/hashicorp/consul/api"
)

var AdminSpreadsheetClient *sheet.Spreadsheet
var TheTimeMachine *job.TimeMachine = nil
var ConsulClient *api.Client = nil
