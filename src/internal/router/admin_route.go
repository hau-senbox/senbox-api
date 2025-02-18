package router

import (
	"encoding/json"
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/usecase/infrastructure"
	"sen-global-api/internal/middleware"
	"sen-global-api/pkg/job"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"time"

	firebase "firebase.google.com/go/v4"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func setupAdminRoutes(engine *gin.Engine, dbConn *gorm.DB, config config.AppConfig, userSpreadsheet *sheet.Spreadsheet, uploaderSpreadsheet *sheet.Spreadsheet, fcm *firebase.App) {
	usecase.AdminSpreadsheetClient = userSpreadsheet
	usecase.TheTimeMachine = job.New()
	sessionRepository := repository.SessionRepository{
		AuthorizeEncryptKey: config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}
	formRepo := &repository.FormRepository{DBConn: dbConn, DefaultRequestPageSize: config.DefaultRequestPageSize}

	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}
	settingRepository := &repository.SettingRepository{DBConn: dbConn}

	importFormsUseCase := &usecase.ImportFormsUseCase{
		FormRepository:                  formRepo,
		QuestionRepository:              &repository.QuestionRepository{DBConn: dbConn},
		FormQuestionRepository:          &repository.FormQuestionRepository{DBConn: dbConn},
		SpreadsheetReader:               uploaderSpreadsheet.Reader,
		SpreadsheetWriter:               uploaderSpreadsheet.Writer,
		SettingRepository:               settingRepository,
		DefaultCronJobIntervalInMinutes: config.DefaultCronJobIntervalInMinutes,
		TimeMachine:                     usecase.TheTimeMachine,
		AppConfig:                       config,
	}
	importUrlsUseCase := &usecase.ImportRedirectUrlsUseCase{
		RedirectUrlRepository: &repository.RedirectUrlRepository{
			DBConn:                 dbConn,
			DefaultRequestPageSize: config.DefaultRequestPageSize,
		},
		SpreadsheetReader: uploaderSpreadsheet.Reader,
		SpreadsheetWriter: uploaderSpreadsheet.Writer,
		SettingRepository: settingRepository,
		TimeMachine:       usecase.TheTimeMachine,
	}

	deviceRepository := &repository.DeviceRepository{DBConn: dbConn, DefaultRequestPageSize: config.DefaultRequestPageSize, DefaultOutputSpreadsheetUrl: config.OutputSpreadsheetUrl}

	syncDevicesUseCase := &usecase.SyncDevicesUseCase{
		SettingRepository: settingRepository,
		DeviceRepository:  deviceRepository,
		DeviceFormDatasetRepository: &repository.DeviceFormDatasetRepository{
			DBConn: dbConn,
		},
		UserEntityRepository:  &repository.UserEntityRepository{DBConn: dbConn},
		UserConfigRepository:  &repository.UserConfigRepository{DBConn: dbConn},
		Reader:                uploaderSpreadsheet.Reader,
		Writer:                uploaderSpreadsheet.Writer,
		TimeMachine:           usecase.TheTimeMachine,
		UserSpreadsheetReader: userSpreadsheet.Reader,
		UserSpreadsheetWriter: userSpreadsheet.Writer,
	}

	v1 := engine.Group("/v1/admin")
	{
		userRepository := &repository.UserRepository{DBConn: dbConn}
		userController := &controller.UserController{
			GetUsersUseCase: usecase.GetUsersUseCase{
				UserRepository: userRepository,
			},
			CreateUserUseCase: usecase.CreateUserUseCase{
				UserRepository: userRepository,
			},
			ChangePasswordUseCase: usecase.ChangePasswordUseCase{
				SessionRepository: &sessionRepository,
				UserRepository:    userRepository,
			},
		}
		loginController := &controller.LoginController{DBConn: dbConn,
			AuthorizeUseCase: usecase.AuthorizeUseCase{
				UserRepository:    &repository.UserRepository{DBConn: dbConn},
				SessionRepository: sessionRepository,
			},
		}
		v1.POST("/login", loginController.Login)

		v1.GET("/users", secureMiddleware.ValidateAdminRole(), userController.GetUserList)
		v1.PUT("/me/new-password", secureMiddleware.ValidateAdminRole(), userController.ChangePassword)

		v1.POST("/user/create", secureMiddleware.ValidateAdminRole(), userController.CreateUser)
		form := controller.FormController{
			SaveFormUseCase: usecase.SaveFormUseCase{
				FormRepository:         formRepo,
				QuestionRepository:     &repository.QuestionRepository{DBConn: dbConn},
				FormQuestionRepository: &repository.FormQuestionRepository{DBConn: dbConn},
				SpreadsheetReader:      userSpreadsheet.Reader,
			},
			DeleteFormUseCase: usecase.DeleteFormUseCase{
				FormRepository: formRepo,
			},
			GetFormListUseCase: usecase.GetFormListUseCase{
				FormRepository: formRepo,
			},
			UpdateFormUseCase: usecase.UpdateFormUseCase{
				FormRepository:         formRepo,
				QuestionRepository:     &repository.QuestionRepository{DBConn: dbConn},
				FormQuestionRepository: &repository.FormQuestionRepository{DBConn: dbConn},
				SpreadsheetReader:      userSpreadsheet.Reader,
			},
			SearchFormsUseCase: usecase.SearchFormsUseCase{
				FormRepository: formRepo,
			},
			ImportFormsUseCase: importFormsUseCase,
		}
		v1.POST("/form/create", secureMiddleware.ValidateAdminRole(), form.CreateForm)

		v1.GET("/form/list", secureMiddleware.ValidateAdminRole(), form.GetFormList)

		v1.DELETE("/form/delete/:id", secureMiddleware.ValidateAdminRole(), form.DeleteForm)

		v1.GET("/forms/search", secureMiddleware.ValidateAdminRole(), form.SearchForms)

		v1.PUT("/form/:id", secureMiddleware.ValidateAdminRole(), form.UpdateForm)

		v1.POST("/forms/import", secureMiddleware.ValidateAdminRole(), form.ImportForms)

		v1.POST("/forms2/import", secureMiddleware.ValidateAdminRole(), form.ImportForms2)

		v1.POST("/forms3/import", secureMiddleware.ValidateAdminRole(), form.ImportForms3)

		v1.POST("/forms4/import", secureMiddleware.ValidateAdminRole(), form.ImportForms4)

		v1.POST("/forms/partially/import", middleware.NewSecureAppMiddleware(dbConn).Secure(), form.ImportFormsPartially)

		v1.POST("/forms/signup", form.ImportSignUpForms)

		deviceController := &controller.DeviceController{
			DBConn: dbConn,
			UpdateDeviceSheetUseCase: &usecase.UpdateDeviceSheetUseCase{
				DeviceRepository: deviceRepository,
			},
			RegisterDeviceUseCase: &usecase.RegisterDeviceUseCase{
				UserRepository:    &repository.UserRepository{DBConn: dbConn},
				DeviceRepository:  deviceRepository,
				SessionRepository: &sessionRepository,
				Writer:            userSpreadsheet.Writer,
				Reader:            userSpreadsheet.Reader,
			},
			GetDeviceByIdUseCase: &usecase.GetDeviceByIdUseCase{
				DeviceRepository: deviceRepository,
			},
			GetDeviceListUseCase: &usecase.GetDeviceListUseCase{
				DeviceRepository: deviceRepository,
			},
			UpdateDeviceUseCase: &usecase.UpdateDeviceUseCase{
				DeviceRepository:  deviceRepository,
				SettingRepository: settingRepository,
				Writer:            userSpreadsheet.Writer,
			},
			GetSettingMessageUseCase: &usecase.GetSettingMessageUseCase{DeviceRepository: deviceRepository, Reader: userSpreadsheet.Reader},
		}

		v1.GET("/device/message/:device_id", secureMiddleware.ValidateAdminRole(), deviceController.GetSettingMessageV2)

		v1.PUT("/device/deactivate/:device_id", secureMiddleware.ValidateAdminRole(), deviceController.DeactivateDevice)

		v1.PUT("/device/activate/:device_id", secureMiddleware.ValidateAdminRole(), deviceController.ActivateDevice)

		v1.PUT("/device/:device_id/update", secureMiddleware.ValidateAdminRole(), deviceController.UpdateDevice)

		v1.PUT("/device/:device_id/updatev2", secureMiddleware.ValidateAdminRole(), deviceController.UpdateDeviceV2)
	}
	redirectUrl := engine.Group("/v1/admin/redirect-url")
	{
		redirectUrlRepository := &repository.RedirectUrlRepository{DBConn: dbConn, DefaultRequestPageSize: config.DefaultRequestPageSize}
		redirectController := &controller.RedirectUrlController{
			SaveRedirectUrlUseCase: &usecase.SaveRedirectUrlUseCase{
				RedirectUrlRepository: redirectUrlRepository,
			},
			GetRedirectUrlListUseCase: &usecase.GetRedirectUrlListUseCase{
				RedirectUrlRepository: redirectUrlRepository,
			},
			DeleteRedirectUrlUseCase: &usecase.DeleteRedirectUrlUseCase{
				RedirectUrlRepository: redirectUrlRepository,
			},
			UpdateRedirectUrlUseCase: &usecase.UpdateRedirectUrlUseCase{
				RedirectUrlRepository: redirectUrlRepository,
			},
			GetRedirectUrlByQRCodeUseCase: nil,
			ImportRedirectUrlsUseCase:     importUrlsUseCase,
		}
		redirectUrl.POST("/create", secureMiddleware.ValidateAdminRole(), redirectController.CreateRedirectUrl)

		redirectUrl.GET("/list", secureMiddleware.ValidateAdminRole(), redirectController.GetRedirectUrlList)

		redirectUrl.DELETE("/:id", secureMiddleware.ValidateAdminRole(), redirectController.DeleteRedirectUrl)

		redirectUrl.PUT("/:id", secureMiddleware.ValidateAdminRole(), redirectController.UpdateRedirectUrl)

		redirectUrl.POST("/import", secureMiddleware.ValidateAdminRole(), redirectController.ImportRedirectUrls)
		//Partially import
		redirectUrl.POST("/import/partially", middleware.NewSecureAppMiddleware(dbConn).Secure(), redirectController.ImportPartiallyRedirectUrls)
	}

	todo := engine.Group("/v1/admin/todo")
	{
		controller := controller.NewImportToDoListController(config, dbConn, uploaderSpreadsheet.Reader, uploaderSpreadsheet.Writer, usecase.TheTimeMachine)
		todo.POST("/import", secureMiddleware.ValidateAdminRole(), controller.ImportTodos)
		todo.POST("/import/partially", middleware.NewSecureAppMiddleware(dbConn).Secure(), controller.ImportPartiallyTodos)
	}

	system := engine.Group("/v1/admin/settings", secureMiddleware.ValidateAdminRole())
	{
		systemController := &controller.SettingController{
			GetSettingsUseCase: &usecase.GetSettingsUseCase{
				SettingRepository: settingRepository,
			},
			UpdateOutputSubmissionSettingUseCase: &usecase.UpdateOutputSubmissionSettingUseCase{
				SettingRepository: settingRepository,
			},
			UpdateOutputSummarySettingUseCase: &usecase.UpdateOutputSummarySettingUseCase{
				SettingRepository: settingRepository,
			},
			SyncDevicesUseCase: syncDevicesUseCase,
			UpdateEmailHistorySettingUseCase: &usecase.UpdateEmailHistorySettingUseCase{
				SettingRepository: settingRepository,
			},
			UpdateOutputTemplateSettingUseCase: &usecase.UpdateOutputTemplateSettingUseCase{
				SettingRepository: settingRepository,
				AppConfig:         config,
			},
			UpdateOutputTemplateSettingForTeacherUseCase: &usecase.UpdateOutputTemplateSettingForTeacherUseCase{
				SettingRepository: settingRepository,
				AppConfig:         config,
			},
			AdminSignUpUseCases: &usecase.AdminSignUpUseCases{
				SettingRepository: settingRepository,
				FormRepository:    formRepo,
				SpreadsheetReader: uploaderSpreadsheet.Reader,
				AppConfig:         config,
				ImportFormsUseCase: &usecase.ImportFormsUseCase{
					FormRepository:                  formRepo,
					QuestionRepository:              &repository.QuestionRepository{DBConn: dbConn},
					FormQuestionRepository:          &repository.FormQuestionRepository{DBConn: dbConn},
					SpreadsheetReader:               uploaderSpreadsheet.Reader,
					SpreadsheetWriter:               uploaderSpreadsheet.Writer,
					SettingRepository:               settingRepository,
					DefaultCronJobIntervalInMinutes: 0,
					TimeMachine:                     nil,
					AppConfig:                       config,
				},
			},
			UpdateSettingNameUseCase:    usecase.NewUpdateSettingNameUseCase(dbConn),
			UpdateApiDistributorUseCase: usecase.NewUpdateApiDistributorUseCase(dbConn, userSpreadsheet.Reader, userSpreadsheet.Writer),
		}
		system.GET("/", systemController.GetSettings)

		system.POST("/output-sheet", systemController.UpdateOutputSubmissionSettings)

		system.POST("/output-summary", systemController.UpdateOutputSummarySettings)

		system.POST("/sync-devices", systemController.SyncDevices)

		system.POST("/email-history", systemController.UpdateEmailHistorySettings)

		system.POST("/output-template", systemController.UpdateOutputTemplateSettings)

		system.POST("/output-template-teacher", systemController.UpdateOutputTemplateSettingsForTeacher)

		system.POST("/sign-up-button-1", systemController.UpdateSignUpButton1)

		system.POST("/sign-up-button-2", systemController.UpdateSignUpButton2)

		system.POST("/sign-up-button-3", systemController.UpdateSignUpButton3)

		system.POST("/sign-up-button-4", systemController.UpdateSignUpButton4)

		system.POST("/sign-up-button-5", systemController.UpdateSignUpButton5)

		system.POST("/sign-up-button-configuration", systemController.UpdateSignUpButtonConfiguration)

		system.POST("/registration-form", systemController.UpdateRegistrationForm)

		system.POST("/registration-submission", systemController.UpdateRegistrationSubmission)

		system.POST("/registration-preset-2", systemController.UpdateRegistrationPreset2)

		system.POST("/registration-preset-1", systemController.UpdateRegistrationPreset1)

		system.POST("/api-distributer", systemController.UpdateAPIDistributor)

		usecase.DBConn = dbConn
		usecase.FirebaseApp = fcm
		system.POST("/code-counting-data", systemController.UpdateCodeCountingData)

		system.POST("/label/name", systemController.SetSettingNames)

		system.POST("/logo-refresh-interval", systemController.SetupLogoRefreshInterval)
		system.GET("/logo-refresh-interval", systemController.GetLogoRefreshInterval)
	}

	monitoring := engine.Group("/v1/admin/monitor")
	{
		monitoringController := &controller.MonitoringController{}

		monitoring.GET("/google-api", monitoringController.GetGoogleAPIMonitoring)
	}

	controller.DBConn = dbConn
	codeCounter := engine.Group("/v1/admin/code-counting", secureMiddleware.ValidateAdminRole())
	{
		codeCounter.GET("/list", controller.GetCodeCounterList)
		codeCounter.PUT("/update", controller.UpdateCodeCounter)
	}

	executor := &TimeMachineSubscriber{
		ImportFormsUseCase:        importFormsUseCase,
		ImportRedirectUrlsUseCase: importUrlsUseCase,
		SyncDevicesUseCase:        syncDevicesUseCase,
		SettingRepository:         settingRepository,
		ImportToDoListUseCase:     usecase.NewImportToDoListUseCase(config, dbConn, uploaderSpreadsheet.Reader, uploaderSpreadsheet.Writer, usecase.TheTimeMachine),
		SyncSubmissionUseCase: &usecase.SyncSubmissionUseCase{
			SubmissionRepository: &repository.SubmissionRepository{
				DBConn: dbConn,
			},
			// DeviceRepository:      deviceRepository,
			FormRepository:        &repository.FormRepository{DBConn: dbConn},
			QuestionRepository:    &repository.QuestionRepository{DBConn: dbConn},
			SettingRepository:     settingRepository,
			UserSpreadsheetReader: userSpreadsheet.Reader,
			UserSpreadsheetWriter: userSpreadsheet.Writer,
			SendEmailUseCase: &usecase.SendEmailUseCase{
				SMTPConfig:        config.SMTP,
				SettingRepository: settingRepository,
				Writer:            userSpreadsheet.Writer,
			},
			GetSettingMessageUseCase: &usecase.GetSettingMessageUseCase{
				DeviceRepository: deviceRepository,
				Reader:           userSpreadsheet.Reader,
			},
			UserEntityRepository: &repository.UserEntityRepository{
				DBConn: dbConn,
			},
		},
	}
	//usecase.TheTimeMachine.Start(config.CronJobInterval)
	getSettingsUseCase := &usecase.GetSettingsUseCase{
		SettingRepository: settingRepository,
	}
	// appSetting, _ := getSettingsUseCase.GetSettings()
	var formInterval uint64 = 0
	var formInterval2 uint64 = 0
	var formInterval3 uint64 = 0
	var formInterval4 uint64 = 0
	var redirectInterval uint64 = 0
	var devicesInterval uint64 = 0
	var toDosInterval uint64 = 0
	// if appSetting != nil {
	// 	if appSetting.Form != nil {
	// 		if appSetting.Form.AutoImport == true && appSetting.Form.Interval > 0 {
	// 			formInterval = appSetting.Form.Interval
	// 		}
	// 	}
	// 	if appSetting.Form2 != nil {
	// 		if appSetting.Form2.AutoImport == true && appSetting.Form2.Interval > 0 {
	// 			formInterval2 = appSetting.Form2.Interval
	// 		}
	// 	}
	// 	if appSetting.Form3 != nil {
	// 		if appSetting.Form3.AutoImport == true && appSetting.Form3.Interval > 0 {
	// 			formInterval3 = appSetting.Form3.Interval
	// 		}
	// 	}
	// 	if appSetting.Form4 != nil {
	// 		if appSetting.Form4.AutoImport == true && appSetting.Form4.Interval > 0 {
	// 			formInterval4 = appSetting.Form4.Interval
	// 		}
	// 	}
	// 	if appSetting.Url != nil {
	// 		if appSetting.Url.AutoImport == true && appSetting.Url.Interval > 0 {
	// 			redirectInterval = appSetting.Url.Interval
	// 		}
	// 	}
	// 	if appSetting.SyncDevices != nil {
	// 		if appSetting.SyncDevices.AutoImport == true && appSetting.SyncDevices.Interval > 0 {
	// 			devicesInterval = appSetting.SyncDevices.Interval
	// 		}
	// 	}
	// 	if appSetting.SyncToDos != nil {
	// 		if appSetting.SyncToDos.AutoImport == true && appSetting.SyncToDos.Interval > 0 {
	// 			toDosInterval = appSetting.SyncToDos.Interval
	// 		}
	// 	}
	// }

	infra := engine.Group("/infra")
	{
		infra.GET("/backup", infrastructure.BackupDatabase())
	}

	sync := engine.Group("v1/admin/sync")
	{
		synchronizeController := &controller.SyncController{
			SynchronizeUseCase: &usecase.SynchronizeUseCase{},
			GetSettingsUseCase: getSettingsUseCase,
		}

		sync.POST("/start", secureMiddleware.ValidateAdminRole(), synchronizeController.StartSync)
		sync.POST("/stop", secureMiddleware.ValidateAdminRole(), synchronizeController.StopSync)
	}

	deviceComponentValues := engine.Group("/v1/admin/device-component-values")
	{
		deviceComponentValuesRepository := &repository.DeviceComponentValuesRepository{DBConn: dbConn}
		deviceComponentValuesController := &controller.DeviceComponentValuesController{
			GetDeviceComponentValuesUseCase: &usecase.GetDeviceComponentValuesUseCase{
				DeviceComponentValuesRepository: deviceComponentValuesRepository,
			},
			SaveDeviceComponentValuesUseCase: &usecase.SaveDeviceComponentValuesUseCase{
				DeviceComponentValuesRepository: deviceComponentValuesRepository,
			},
		}

		deviceComponentValues.GET("/company/:company_id", deviceComponentValuesController.GetDeviceComponentValuesByCompany)

		deviceComponentValues.GET("/device/:company_id", deviceComponentValuesController.GetDeviceComponentValuesByDevice)

		deviceComponentValues.POST("/company", secureMiddleware.ValidateAdminRole(), deviceComponentValuesController.SaveDeviceComponentValuesByCompany)

		deviceComponentValues.POST("/device", secureMiddleware.ValidateAdminRole(), deviceComponentValuesController.SaveDeviceComponentValuesByCompany)
	}

	usecase.TheTimeMachine.Start(formInterval, redirectInterval, devicesInterval, toDosInterval, formInterval2, formInterval3, formInterval4, 0)

	usecase.TheTimeMachine.SubscribeFormsExec(executor)
	usecase.TheTimeMachine.SubscribeForms2Exec(executor)
	usecase.TheTimeMachine.SubscribeForms3Exec(executor)
	usecase.TheTimeMachine.SubscribeForms4Exec(executor)
	usecase.TheTimeMachine.SubscribeUrlsExec(executor)
	usecase.TheTimeMachine.SubscribeSyncDevicesExec(executor)
	usecase.TheTimeMachine.SubscribeSyncToDosExec(executor)
	usecase.TheTimeMachine.SubscribeGoogleAPIRequestMonitorExec(executor)
	usecase.TheTimeMachine.SubscribeSyncSubmissionsExec(executor)
}

type TimeMachineSubscriber struct {
	*usecase.ImportFormsUseCase
	*usecase.ImportRedirectUrlsUseCase
	*usecase.SyncDevicesUseCase
	*repository.SettingRepository
	*usecase.ImportToDoListUseCase
	*usecase.SyncSubmissionUseCase
}

func (t *TimeMachineSubscriber) ExecuteSyncDevices() {
	log.Debug("Start sync devices")
	monitor.SendMessageViaTelegram("Start sync devices by cron job")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	deviceSetting, err := t.SettingRepository.GetSyncDevicesSettings()

	if err != nil {
		log.Error(err)
	} else {
		log.Debug(deviceSetting)
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(deviceSetting.Settings), &importSetting)
		if err != nil {
			log.Error(err)
		} else {
			var importSetting ImportSetting
			err = json.Unmarshal([]byte(deviceSetting.Settings), &importSetting)
			if err != nil {
				log.Error(err)
			} else {
				err = t.SyncDevicesUseCase.SyncDevices(request.SyncDevicesRequest{
					SpreadsheetUrl: importSetting.SpreadSheetUrl,
					AutoImport:     importSetting.AutoImport,
					Interval:       importSetting.Interval,
				})
				if err != nil {
					log.Error(err)
					monitor.SendMessageViaTelegram("Error sync devices by cron job Err:", err.Error())
				} else {
					monitor.SendMessageViaTelegram("Finish sync devices by cron job")
				}
			}
		}
	}
}

func (t *TimeMachineSubscriber) ExecuteSyncUrls() {
	log.Debug("Start sync urls")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	urlSetting, err := t.SettingRepository.GetUrlSettings()
	if err != nil {
		log.Error(err)
	} else {
		log.Debug(urlSetting)
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(urlSetting.Settings), &importSetting)
		if err != nil {
			log.Error(err)
		} else {
			err = t.ImportRedirectUrlsUseCase.SyncUrls(request.ImportRedirectUrlsRequest{
				SpreadsheetUrl: importSetting.SpreadSheetUrl,
				AutoImport:     importSetting.AutoImport,
				Interval:       importSetting.Interval,
			})
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (t *TimeMachineSubscriber) ExecuteSyncForms() {
	log.Debug("Start sync forms")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	log.Debug("TimeMachineSubscriber: ExecuteSyncForms")
	formSettings, err := t.SettingRepository.GetFormSettings()
	if err != nil {
		log.Error(err)
		return
	} else {
		log.Debug(formSettings)
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(formSettings.Settings), &importSetting)
		if err != nil {
			log.Error(err)
		} else {
			err = t.ImportFormsUseCase.SyncForms(request.ImportFormRequest{
				SpreadsheetUrl: importSetting.SpreadSheetUrl,
				AutoImport:     importSetting.AutoImport,
				Interval:       importSetting.Interval,
			})
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (t *TimeMachineSubscriber) ExecuteSyncForms2() {
	log.Debug("Start sync forms")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	log.Debug("TimeMachineSubscriber: ExecuteSyncForms")
	formSettings, err := t.SettingRepository.GetFormSettings2()
	if err != nil {
		log.Error(err)
		return
	} else {
		log.Debug(formSettings)
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(formSettings.Settings), &importSetting)
		if err != nil {
			log.Error(err)
		} else {
			err = t.ImportFormsUseCase.SyncForms(request.ImportFormRequest{
				SpreadsheetUrl: importSetting.SpreadSheetUrl,
				AutoImport:     importSetting.AutoImport,
				Interval:       importSetting.Interval,
			})
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (t *TimeMachineSubscriber) ExecuteSyncForms3() {
	log.Debug("Start sync forms")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	log.Debug("TimeMachineSubscriber: ExecuteSyncForms")
	formSettings, err := t.SettingRepository.GetFormSettings3()
	if err != nil {
		log.Error(err)
		return
	} else {
		log.Debug(formSettings)
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(formSettings.Settings), &importSetting)
		if err != nil {
			log.Error(err)
		} else {
			err = t.ImportFormsUseCase.SyncForms(request.ImportFormRequest{
				SpreadsheetUrl: importSetting.SpreadSheetUrl,
				AutoImport:     importSetting.AutoImport,
				Interval:       importSetting.Interval,
			})
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (t *TimeMachineSubscriber) ExecuteSyncForms4() {
	log.Debug("Start sync forms")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	log.Debug("TimeMachineSubscriber: ExecuteSyncForms")
	formSettings, err := t.SettingRepository.GetFormSettings4()
	if err != nil {
		log.Error(err)
		return
	} else {
		log.Debug(formSettings)
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(formSettings.Settings), &importSetting)
		if err != nil {
			log.Error(err)
		} else {
			err = t.ImportFormsUseCase.SyncForms(request.ImportFormRequest{
				SpreadsheetUrl: importSetting.SpreadSheetUrl,
				AutoImport:     importSetting.AutoImport,
				Interval:       importSetting.Interval,
			})
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (t *TimeMachineSubscriber) ExecuteSyncTodos() {
	log.Debug("Start sync devices")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	deviceSetting, err := t.SettingRepository.GetSyncToDosSettings()

	if err != nil {
		log.Error(err)
	} else {
		log.Debug(deviceSetting)
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(deviceSetting.Settings), &importSetting)
		if err != nil {
			log.Error(err)
		} else {
			var importSetting ImportSetting
			err = json.Unmarshal([]byte(deviceSetting.Settings), &importSetting)
			if err != nil {
				log.Error(err)
			} else {
				err = t.ImportToDoListUseCase.ImportToDoList(request.ImportFormRequest{
					SpreadsheetUrl: importSetting.SpreadSheetUrl,
					AutoImport:     importSetting.AutoImport,
					Interval:       importSetting.Interval,
				})
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}

// register, import 1 todo, import 1 form, screen button, top button
func (t *TimeMachineSubscriber) ExecuteGoogleAPIRequestMonitor() {
	if (monitor.TotalRequestInitDevice +
		monitor.TotalRequestImportToDo +
		monitor.TotalRequestImportForm +
		monitor.TotalRequestGETScreenButton +
		monitor.TotalRequestGETTopButton) > 300 {
		monitor.SendMessageViaTelegram(
			"Number of request for register: "+strconv.Itoa(monitor.TotalRequestInitDevice),
			"Number of request for import 1 todo: "+strconv.Itoa(monitor.TotalRequestImportToDo),
			"Number of request for import 1 form: "+strconv.Itoa(monitor.TotalRequestImportForm),
			"Number of request for screen button: "+strconv.Itoa(monitor.TotalRequestGETScreenButton),
			"Number of request for top button: "+strconv.Itoa(monitor.TotalRequestGETTopButton),
		)
	}

	monitor.ResetGoogleAPIRequestMonitor()
}

func (t *TimeMachineSubscriber) ExecuteSubmissionSync() {
	log.Info("Start interval sync submission")
	t.SyncSubmissionUseCase.Execute()
}
