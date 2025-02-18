package usecase

import (
	"sen-global-api/pkg/job"
)

type SynchronizeUseCase struct {
}

func (receiver SynchronizeUseCase) StartSync(timeMachine job.TimeMachine, appSetting *AppSettings) {
	var formInterval uint64 = 0
	var formInterval2 uint64 = 0
	var formInterval3 uint64 = 0
	var formInterval4 uint64 = 0
	var redirectInterval uint64 = 0
	var devicesInterval uint64 = 0
	var toDosInterval uint64 = 0
	if appSetting != nil {
		if appSetting.Form != nil {
			if appSetting.Form.AutoImport && appSetting.Form.Interval > 0 {
				formInterval = appSetting.Form.Interval
			}
		}
		if appSetting.Form2 != nil {
			if appSetting.Form2.AutoImport && appSetting.Form2.Interval > 0 {
				formInterval2 = appSetting.Form2.Interval
			}
		}
		if appSetting.Form3 != nil {
			if appSetting.Form3.AutoImport && appSetting.Form3.Interval > 0 {
				formInterval3 = appSetting.Form3.Interval
			}
		}
		if appSetting.Form4 != nil {
			if appSetting.Form4.AutoImport && appSetting.Form4.Interval > 0 {
				formInterval4 = appSetting.Form4.Interval
			}
		}
		if appSetting.Url != nil {
			if appSetting.Url.AutoImport && appSetting.Url.Interval > 0 {
				redirectInterval = appSetting.Url.Interval
			}
		}
		if appSetting.SyncDevices != nil {
			if appSetting.SyncDevices.AutoImport && appSetting.SyncDevices.Interval > 0 {
				devicesInterval = appSetting.SyncDevices.Interval
			}
		}
		if appSetting.SyncToDos != nil {
			if appSetting.SyncToDos.AutoImport && appSetting.SyncToDos.Interval > 0 {
				toDosInterval = appSetting.SyncToDos.Interval
			}
		}
	}

	timeMachine.Start(formInterval, redirectInterval, devicesInterval, toDosInterval, formInterval2, formInterval3, formInterval4, 0)
}

func (receiver SynchronizeUseCase) StopSync(timeMachine job.TimeMachine) {
	timeMachine.Stop()
}
