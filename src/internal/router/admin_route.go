package router

import (
	"encoding/json"
	"sen-global-api/config"
	"sen-global-api/helper"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/usecase/infrastructure"
	"sen-global-api/internal/middleware"
	"sen-global-api/pkg/job"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"sen-global-api/pkg/uploader"
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
		OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		AuthorizeEncryptKey:    config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}
	formRepo := &repository.FormRepository{DBConn: dbConn, DefaultRequestPageSize: config.DefaultRequestPageSize}

	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}
	settingRepository := &repository.SettingRepository{DBConn: dbConn}

	s3Provider := uploader.NewS3Provider(
		config.S3.SenboxFormSubmitBucket.AccessKey,
		config.S3.SenboxFormSubmitBucket.SecretKey,
		config.S3.SenboxFormSubmitBucket.BucketName,
		config.S3.SenboxFormSubmitBucket.Region,
		config.S3.SenboxFormSubmitBucket.Domain,
		config.S3.SenboxFormSubmitBucket.CloudfrontKeyGroupID,
		config.S3.SenboxFormSubmitBucket.CloudfrontKeyPath,
	)

	importFormsUseCase := &usecase.ImportFormsUseCase{
		FormRepository:                  formRepo,
		QuestionRepository:              &repository.QuestionRepository{DBConn: dbConn},
		FormQuestionRepository:          &repository.FormQuestionRepository{DBConn: dbConn},
		SpreadsheetReader:               uploaderSpreadsheet.Reader,
		SpreadsheetWriter:               uploaderSpreadsheet.Writer,
		SettingRepository:               settingRepository,
		RoleOrgSignUpRepo:               &repository.RoleOrgSignUpRepository{DBConn: dbConn},
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

	v1 := engine.Group("/v1/admin")
	{
		loginController := &controller.LoginController{DBConn: dbConn,
			AuthorizeUseCase: usecase.AuthorizeUseCase{
				UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
				DeviceRepository:     &repository.DeviceRepository{DBConn: dbConn},
				SessionRepository:    sessionRepository,
				DBConn:               dbConn,
				UpdateUserEntityUseCase: &usecase.UpdateUserEntityUseCase{
					UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
				},
			},
		}
		v1.POST("/login", loginController.Login)

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
		v1.POST("/form/create", secureMiddleware.ValidateSuperAdminRole(), form.CreateForm)

		v1.GET("/form/list", secureMiddleware.ValidateSuperAdminRole(), form.GetFormList)

		v1.DELETE("/form/delete/:id", secureMiddleware.ValidateSuperAdminRole(), form.DeleteForm)

		v1.GET("/forms/search", secureMiddleware.ValidateSuperAdminRole(), form.SearchForms)

		v1.PUT("/form/:id", secureMiddleware.ValidateSuperAdminRole(), form.UpdateForm)

		v1.POST("/forms/import", secureMiddleware.ValidateSuperAdminRole(), form.ImportForms)

		v1.POST("/forms2/import", secureMiddleware.ValidateSuperAdminRole(), form.ImportForms2)

		v1.POST("/forms3/import", secureMiddleware.ValidateSuperAdminRole(), form.ImportForms3)

		v1.POST("/forms4/import", secureMiddleware.ValidateSuperAdminRole(), form.ImportForms4)

		v1.POST("/forms/partially/import", middleware.NewSecureAppMiddleware(dbConn).Secure(), form.ImportFormsPartially)

		v1.POST("/forms/signup", form.ImportSignUpForms)

		deviceController := &controller.DeviceController{
			DBConn: dbConn,
			UpdateDeviceSheetUseCase: &usecase.UpdateDeviceSheetUseCase{
				DeviceRepository: deviceRepository,
			},
			RegisterDeviceUseCase: &usecase.RegisterDeviceUseCase{
				DeviceRepository:  deviceRepository,
				SessionRepository: &sessionRepository,
				Writer:            userSpreadsheet.Writer,
				Reader:            userSpreadsheet.Reader,
			},
			GetDeviceByIDUseCase: &usecase.GetDeviceByIDUseCase{
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
		}

		v1.PUT("/device/deactivate/:device_id", secureMiddleware.ValidateSuperAdminRole(), deviceController.DeactivateDevice)

		v1.PUT("/device/activate/:device_id", secureMiddleware.ValidateSuperAdminRole(), deviceController.ActivateDevice)

		v1.PUT("/device/:device_id/update", secureMiddleware.ValidateSuperAdminRole(), deviceController.UpdateDevice)

		v1.PUT("/device/:device_id/updatev2", secureMiddleware.ValidateSuperAdminRole(), deviceController.UpdateDeviceV2)
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
		redirectUrl.POST("/create", secureMiddleware.ValidateSuperAdminRole(), redirectController.CreateRedirectUrl)

		redirectUrl.GET("/list", secureMiddleware.ValidateSuperAdminRole(), redirectController.GetRedirectUrlList)

		redirectUrl.DELETE("/:id", secureMiddleware.ValidateSuperAdminRole(), redirectController.DeleteRedirectUrl)

		redirectUrl.PUT("/:id", secureMiddleware.ValidateSuperAdminRole(), redirectController.UpdateRedirectUrl)

		redirectUrl.POST("/import", secureMiddleware.ValidateSuperAdminRole(), redirectController.ImportRedirectUrls)
		//Partially import
		redirectUrl.POST("/import/partially", middleware.NewSecureAppMiddleware(dbConn).Secure(), redirectController.ImportPartiallyRedirectUrls)
	}

	todo := engine.Group("/v1/admin/todo")
	{
		todoController := controller.NewImportToDoListController(config, dbConn, uploaderSpreadsheet.Reader, uploaderSpreadsheet.Writer, usecase.TheTimeMachine)
		todo.POST("/import", secureMiddleware.ValidateSuperAdminRole(), todoController.ImportTodos)
		todo.POST("/import/partially", middleware.NewSecureAppMiddleware(dbConn).Secure(), todoController.ImportPartiallyTodos)
	}

	system := engine.Group("/v1/admin/settings", secureMiddleware.ValidateSuperAdminRole())
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
			LanguageSettingUseCase: &usecase.LanguageSettingUseCase{
				LanguageSettingRepository: &repository.LanguageSettingRepository{DBConn: dbConn},
			},
		}
		system.GET("/", systemController.GetSettings)

		system.POST("/output-sheet", systemController.UpdateOutputSubmissionSettings)

		system.POST("/output-summary", systemController.UpdateOutputSummarySettings)

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

		// language setting
		system.POST("/language", systemController.UploadLanguageSetting)
		system.GET("/language", systemController.GetLanguageSettings)
	}

	monitoring := engine.Group("/v1/admin/monitor")
	{
		monitoringController := &controller.MonitoringController{}

		monitoring.GET("/google-api", monitoringController.GetGoogleAPIMonitoring)
	}

	controller.DBConn = dbConn
	codeCounter := engine.Group("/v1/admin/code-counting", secureMiddleware.ValidateSuperAdminRole())
	{
		codeCounter.GET("/list", controller.GetCodeCounterList)
		codeCounter.PUT("/update", controller.UpdateCodeCounter)
	}

	roleUseCase := usecase.NewRoleOrgSignUpUseCase(&repository.RoleOrgSignUpRepository{
		DBConn: dbConn,
	})

	roleOrgSignUpController := controller.NewRoleOrgSignUpController(roleUseCase)
	roleSignUp := engine.Group("/v1/admin/role-sign-up", secureMiddleware.Secured())
	{
		roleSignUp.GET("", roleOrgSignUpController.Get4AdminWeb)
		roleSignUp.POST("", secureMiddleware.ValidateSuperAdminRole(), roleOrgSignUpController.CreateOrUpdate)
	}

	//menu
	menuController := &controller.MenuController{
		GetMenuUseCase: &usecase.GetMenuUseCase{
			MenuRepository:          &repository.MenuRepository{DBConn: dbConn},
			UserEntityRepository:    &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository:  &repository.OrganizationRepository{DBConn: dbConn},
			DeviceRepository:        &repository.DeviceRepository{DBConn: dbConn},
			RoleOrgSignUpRepository: &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			FormRepository:          &repository.FormRepository{DBConn: dbConn},
			SubmissionRepository:    &repository.SubmissionRepository{DBConn: dbConn},
			ComponentRepository:     &repository.ComponentRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
			OrganizationMenuTemplateRepository: &repository.OrganizationMenuTemplateRepository{DBConn: dbConn},
			StudentAppRepo:                     &repository.StudentApplicationRepository{DB: dbConn},
			SuperAdminEmergencyMenuRepo:        &repository.SuperAdminEmergencyMenuRepository{DBConn: dbConn},
			OrganizationEmergencyMenuRepo:      &repository.OrganizationEmergencyMenuRepository{DBConn: dbConn},
		},
		UploadSectionMenuUseCase: &usecase.UploadSectionMenuUseCase{
			MenuRepository:                     &repository.MenuRepository{DBConn: dbConn},
			ComponentRepository:                &repository.ComponentRepository{DBConn: dbConn},
			ChildMenuRepository:                &repository.ChildMenuRepository{DBConn: dbConn},
			ChildRepository:                    &repository.ChildRepository{DB: dbConn},
			RoleOrgSignUpRepository:            &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			StudentMenuRepository:              &repository.StudentMenuRepository{DBConn: dbConn},
			StudentApplicationRepository:       &repository.StudentApplicationRepository{DB: dbConn},
			OrganizationMenuTemplateRepository: &repository.OrganizationMenuTemplateRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
			TeacherApplicationRepository:         &repository.TeacherApplicationRepository{DBConn: dbConn},
			TeacherMenuRepository:                &repository.TeacherMenuRepository{DBConn: dbConn},
			StaffMenuRepository:                  &repository.StaffMenuRepository{DBConn: dbConn},
			StaffApplicationRepository:           &repository.StaffApplicationRepository{DBConn: dbConn},
			DeviceMenuRepository:                 &repository.DeviceMenuRepository{DBConn: dbConn},
			ParentMenuRepository:                 &repository.ParentMenuRepository{DBConn: dbConn},
			TeacherMenuOrganizationRepository:    &repository.TeacherMenuOrganizationRepository{DBConn: dbConn},
			DepartmentMenuRepository:             &repository.DepartmentMenuRepository{DBConn: dbConn},
			DepartmentMenuOrganizationRepository: &repository.DepartmentMenuOrganizationRepository{DBConn: dbConn},
			SuperAdminEmergencyMenuRepository:    &repository.SuperAdminEmergencyMenuRepository{DBConn: dbConn},
			OrganizationEmergencyMenuRepository:  &repository.OrganizationEmergencyMenuRepository{DBConn: dbConn},
		},
		ChildMenuUseCase: &usecase.ChildMenuUseCase{
			Repo:          &repository.ChildMenuRepository{DBConn: dbConn},
			ComponentRepo: &repository.ComponentRepository{DBConn: dbConn},
			ChildRepo:     &repository.ChildRepository{DB: dbConn},
		},
		StudentMenuUseCase: &usecase.StudentMenuUseCase{
			StudentMenuRepo: &repository.StudentMenuRepository{DBConn: dbConn},
			StudentAppRepo:  &repository.StudentApplicationRepository{DB: dbConn},
			ComponentRepo:   &repository.ComponentRepository{DBConn: dbConn},
			UserEntityRepo:  &repository.UserEntityRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
		},
		StudentApplicationUseCase: &usecase.StudentApplicationUseCase{
			StudentAppRepo:  &repository.StudentApplicationRepository{DB: dbConn},
			StudentMenuRepo: &repository.StudentMenuRepository{DBConn: dbConn},
			ComponentRepo:   &repository.ComponentRepository{DBConn: dbConn},
			RoleOrgRepo:     &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
		},
		TeacherMenuUseCase: &usecase.TeacherMenuUseCase{
			TeacherMenuRepo:      &repository.TeacherMenuRepository{DBConn: dbConn},
			TeacherAppRepo:       &repository.TeacherApplicationRepository{DBConn: dbConn},
			ComponentRepo:        &repository.ComponentRepository{DBConn: dbConn},
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
		},
		DeviceMenuUseCase: &usecase.DeviceMenuUseCase{
			Repo:          &repository.DeviceMenuRepository{DBConn: dbConn},
			ComponentRepo: &repository.ComponentRepository{DBConn: dbConn},
			DeviceRepo:    &repository.DeviceRepository{DBConn: dbConn},
		},
		TeacherMenuOrganizationUseCase: &usecase.TeacherMenuOrganizationUseCase{
			TeacherMenuOrganizationRepository: &repository.TeacherMenuOrganizationRepository{DBConn: dbConn},
			ComponentRepo:                     &repository.ComponentRepository{DBConn: dbConn},
		},
	}

	menu := engine.Group("/v1/admin/menu", secureMiddleware.Secured(), middleware.MenuUploadLoggerMiddleware(dbConn))
	{
		menu.GET("/section", menuController.GetSectionMenu4WebAdmin)
		menu.POST("/section", menuController.UploadSectionMenu)
		menu.POST("/section/student", menuController.UploadStudentMenu)
		menu.POST("/section/teacher", menuController.UploadTeacherMenu)
		menu.POST("/section/staff", menuController.UploadStaffMenu)
		menu.POST("/section/child", secureMiddleware.ValidateSuperAdminRole(), menuController.UploadChildMenu)
		menu.POST("/section/parent", secureMiddleware.ValidateSuperAdminRole(), menuController.UploadParentMenu)
		menu.GET("/section/device/:id", menuController.GetDeviceMenu4Admin)
		menu.POST("/section/device", menuController.UploadDeviceSectionMenu)
		menu.DELETE("/section/:id", menuController.DeleteSectionMenu)
		menu.GET("/child/:id", menuController.GetChildMenuByChildID)
		menu.PUT("/child", menuController.UpdateIsShowChildMenu)
		menu.PUT("/student", menuController.UpdateIsShowStudentMenu)
		menu.PUT("/teacher", menuController.UpdateIsShowTeacherMenu)
		menu.POST("/section/teacher/organization", menuController.UploadTeacherMenuOrganization)
		menu.GET("/section/teacher/:teacher_id/organization/:organization_id", menuController.GetTeacherMenuOrganization4Admin)

		// get user menu
		menu.GET("/user/:id", menuController.GetUserMenu)
		// super admin menu
		menu.GET("", secureMiddleware.ValidateSuperAdminRole(), menuController.GetSuperAdminMenu)
		menu.POST("/top", secureMiddleware.ValidateSuperAdminRole(), menuController.UploadSuperAdminMenuTop)
		menu.POST("/bottom", secureMiddleware.ValidateSuperAdminRole(), menuController.UploadSuperAdminMenuBottom)
		// organization admin menu
		menu.GET("/organization/:id", menuController.GetOrgMenu)
		menu.POST("/organization/top", menuController.UploadOrganizationAdminMenuTop)
		menu.POST("/organization/bottom", menuController.UploadOrganizationAdminMenuBottom)

		// emergency menu
		menu.POST("/emergency", menuController.UploadEmergencyMenu)
		menu.GET("/emergency", menuController.GetEmergencyMenu4WebAdmin)
	}

	// user
	studentUseCase := usecase.NewStudentApplicationUseCase(
		&repository.StudentApplicationRepository{DB: dbConn},
		&repository.StudentMenuRepository{DBConn: dbConn},
		&repository.ComponentRepository{DBConn: dbConn},
		&repository.RoleOrgSignUpRepository{DBConn: dbConn},
		&usecase.GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		&repository.OrganizationRepository{DBConn: dbConn},
		&usecase.StudentBlockSettingUsecase{
			Repo: &repository.StudentBlockSettingRepository{DBConn: dbConn},
		},
		&usecase.LanguagesConfigUsecase{
			Repo: &repository.LanguagesConfigRepository{DBConn: dbConn},
		},
		&usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	)

	childUseCase := usecase.NewChildUseCase(
		&repository.ChildRepository{DB: dbConn},
		&repository.UserEntityRepository{DBConn: dbConn},
		&repository.ComponentRepository{DBConn: dbConn},
		&repository.ChildMenuRepository{DBConn: dbConn},
		&repository.RoleOrgSignUpRepository{DBConn: dbConn},
		&usecase.GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		&usecase.LanguagesConfigUsecase{
			Repo: &repository.LanguagesConfigRepository{DBConn: dbConn},
		},
		&usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	)

	staffUsecase := usecase.NewStaffApplicationUseCase(
		&repository.StaffApplicationRepository{DBConn: dbConn},
		&repository.StaffMenuRepository{DBConn: dbConn},
		&repository.ComponentRepository{DBConn: dbConn},
		&repository.RoleOrgSignUpRepository{DBConn: dbConn},
		&repository.OrganizationRepository{DBConn: dbConn},
		&usecase.GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		&repository.UserEntityRepository{DBConn: dbConn},
		&usecase.LanguagesConfigUsecase{
			Repo: &repository.LanguagesConfigRepository{DBConn: dbConn},
		},
		&usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	)

	userEntityController := &controller.UserEntityController{
		ChildUseCase:              childUseCase,
		StudentApplicationUseCase: studentUseCase,
		StaffApplicationUseCase:   staffUsecase,
		GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			ChildRepository:        &repository.ChildRepository{DB: dbConn},
		},
		TeacherApplicationUseCase: &usecase.TeacherApplicationUseCase{
			TeacherRepo: &repository.TeacherApplicationRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
			TeacherMenuRepo:      &repository.TeacherMenuRepository{DBConn: dbConn},
			ComponentRepo:        &repository.ComponentRepository{DBConn: dbConn},
			RoleOrgRepo:          &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			LanguagesConfigUsecase: &usecase.LanguagesConfigUsecase{
				Repo: &repository.LanguagesConfigRepository{DBConn: dbConn},
			},
			UserImagesUsecase: &usecase.UserImagesUsecase{
				Repo:      &repository.UserImagesRepository{DBConn: dbConn},
				ImageRepo: &repository.ImageRepository{DBConn: dbConn},
				GetImageUseCase: &usecase.GetImageUseCase{
					UploadProvider:  s3Provider,
					ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				},
			},
		},
		UpdateUserEntityUseCase: &usecase.UpdateUserEntityUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		UserBlockSettingUsecase: &usecase.UserBlockSettingUsecase{
			Repo:        &repository.UserBlockSettingRepository{DBConn: dbConn},
			TeacherRepo: &repository.TeacherApplicationRepository{DBConn: dbConn},
			StaffRepo:   &repository.StaffApplicationRepository{DBConn: dbConn},
		},
		ParentUseCase: &usecase.ParentUseCase{
			UserRepo:       &repository.UserEntityRepository{DBConn: dbConn},
			ParentMenuRepo: &repository.ParentMenuRepository{DBConn: dbConn},
			ComponentRepo:  &repository.ComponentRepository{DBConn: dbConn},
			LanguagesConfigUsecase: &usecase.LanguagesConfigUsecase{
				Repo: &repository.LanguagesConfigRepository{DBConn: dbConn},
			},
			UserImagesUsecase: &usecase.UserImagesUsecase{
				Repo:      &repository.UserImagesRepository{DBConn: dbConn},
				ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			},
		},
		StudentBlockSettingUsecase: &usecase.StudentBlockSettingUsecase{
			Repo: &repository.StudentBlockSettingRepository{DBConn: dbConn},
		},
		UploadImageUseCase: &usecase.UploadImageUseCase{
			UploadProvider:  s3Provider,
			ImageRepository: &repository.ImageRepository{DBConn: dbConn},
		},
		UserImagesUsecase: &usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			DeleteImageUsecase: &usecase.DeleteImageUseCase{
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				UploadProvider:  s3Provider,
			},
			GetImageUseCase: &usecase.GetImageUseCase{
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				UploadProvider:  s3Provider,
			},
		},
		OwnerAssignUseCase: &usecase.OwnerAssignUseCase{
			TeacherRepo:    &repository.TeacherApplicationRepository{DBConn: dbConn},
			StaffRepo:      &repository.StaffApplicationRepository{DBConn: dbConn},
			StudentRepo:    &repository.StudentApplicationRepository{DB: dbConn},
			UserEntityRepo: &repository.UserEntityRepository{DBConn: dbConn},
			UserImagesUsecase: &usecase.UserImagesUsecase{
				Repo:      &repository.UserImagesRepository{DBConn: dbConn},
				ImageRepo: &repository.ImageRepository{DBConn: dbConn},
				GetImageUseCase: &usecase.GetImageUseCase{
					UploadProvider:  s3Provider,
					ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				},
			},
		},
	}

	userBlockSettingController := &controller.BlockSettingController{
		UserBlockUsecase: &usecase.UserBlockSettingUsecase{
			Repo:        &repository.UserBlockSettingRepository{DBConn: dbConn},
			TeacherRepo: &repository.TeacherApplicationRepository{DBConn: dbConn},
			StaffRepo:   &repository.StaffApplicationRepository{DBConn: dbConn},
		},
		StudentBlockUsecase: &usecase.StudentBlockSettingUsecase{
			Repo: &repository.StudentBlockSettingRepository{DBConn: dbConn},
		},
	}

	userSettingController := &controller.UserSettingController{
		UserSettingUsecase: &usecase.UserSettingUseCase{
			Repo: &repository.UserSettingRepository{DBConn: dbConn},
		},
	}

	user := engine.Group("/v1/admin/user", secureMiddleware.Secured())
	{
		user.GET("/search", userEntityController.SearchUser4WebAdmin)
		user.GET("/child/:id", userEntityController.GetChild4WebAdmin)
		user.GET("/student/:id", userEntityController.GetStudent4WebAdmin)
		user.GET("/teacher/:id", userEntityController.GetTeacher4WebAdmin)
		user.GET("/staff/:id", userEntityController.GetStaff4WebAdmin)
		user.GET("/parent/:id", secureMiddleware.ValidateSuperAdminRole(), userEntityController.GetParent4WebAdmin)
		user.POST("/student/add-custom-id", userEntityController.AddCustomID2Student)
		user.POST("/add-custom-id", userEntityController.AddCustomID2User)
		block := user.Group("/block")
		{
			block.GET("/:user_id", userBlockSettingController.GetByUserID)
			block.POST("", userBlockSettingController.UpsertUserBlockSetting)
			block.POST("/student", userBlockSettingController.UpsertStudentBlockSetting)
			block.GET("/student/:student_id", userBlockSettingController.GetByStudentID)
		}
		// avatar
		user.POST("/avatar", userEntityController.UploadAvatarV2)
		user.PUT("/avatar/is-main", userEntityController.UpdateIsMain)
		user.DELETE("/avatar", userEntityController.DeleteUserAvatar)

		// setting
		user.POST("/setting", secureMiddleware.ValidateSuperAdminRole(), userSettingController.UploadUserSetting)

		// assign department list
		user.GET("/organization/:organization_id/assigned-owners", userEntityController.GetListOwner2Assign)

		// re logimn
		user.PUT("/set-relogin", secureMiddleware.ValidateSuperAdminRole(), userEntityController.SetReLogin)
	}

	// application
	sheetsService, _ := helper.GetSheetsService("credentials/uploader_service_account.json")

	syncDataUsecase := &usecase.SyncDataUsecase{
		SheetService:       sheetsService,
		SubmissionRepo:     &repository.SubmissionRepository{DBConn: dbConn},
		SyncQueueRepo:      &repository.SyncQueueRepository{DBConn: dbConn},
		SettingRepository:  &repository.SettingRepository{DBConn: dbConn},
		ImportFormsUseCase: importFormsUseCase,
	}

	// call schedule sync cron job
	syncDataUsecase.StartAutoSyncScheduler()
	syncDataUsecase.StartAutoSyncForm2Scheduler()

	applicationController := &controller.ApplicationController{
		StaffAppUsecase: &usecase.StaffApplicationUseCase{
			StaffAppRepo:  &repository.StaffApplicationRepository{DBConn: dbConn},
			StaffMenuRepo: &repository.StaffMenuRepository{DBConn: dbConn},
			ComponentRepo: &repository.ComponentRepository{DBConn: dbConn},
			RoleOrgRepo:   &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
			OrganizationRepo:     &repository.OrganizationRepository{DBConn: dbConn},
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		StudentAppUsecase: &usecase.StudentApplicationUseCase{
			StudentAppRepo:  &repository.StudentApplicationRepository{DB: dbConn},
			StudentMenuRepo: &repository.StudentMenuRepository{DBConn: dbConn},
			ComponentRepo:   &repository.ComponentRepository{DBConn: dbConn},
			RoleOrgRepo:     &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
			OrganizationRepo: &repository.OrganizationRepository{DBConn: dbConn},
		},
		TeacherAppUsecase: &usecase.TeacherApplicationUseCase{
			TeacherRepo: &repository.TeacherApplicationRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
			TeacherMenuRepo:      &repository.TeacherMenuRepository{DBConn: dbConn},
			ComponentRepo:        &repository.ComponentRepository{DBConn: dbConn},
			RoleOrgRepo:          &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			OrganizationRepo:     &repository.OrganizationRepository{DBConn: dbConn},
		},
		SyncDataUsecase: syncDataUsecase,
	}

	application := engine.Group("/v1/admin/application", secureMiddleware.Secured())
	{
		// student application
		application.GET("/student", applicationController.GetAllStudentApplications)
		application.GET("/student/:id", applicationController.GetDetailStudentApplication)
		application.PUT("/student/approve/:id", applicationController.ApproveStudentApplication)
		application.PUT("/student/block/:id", applicationController.BlockStudentApplication)

		// teacher application
		application.GET("/teacher", applicationController.GetAllTeacherApplications)
		application.GET("/teacher/:id", applicationController.GetDetailTeacherApplication)
		application.PUT("/teacher/approve/:id", applicationController.ApproveTeacherApplication)
		application.PUT("/teacher/block/:id", applicationController.BlockTeacherApplication)

		// staff application
		application.GET("/staff", applicationController.GetAllStaffApplications)
		application.GET("/staff/:id", applicationController.GetDetailStaffApplication)
		application.PUT("/staff/approve/:id", applicationController.ApproveStaffApplication)
		application.PUT("/staff/block/:id", applicationController.BlockStaffApplication)
	}

	// organization

	orgController := &controller.OrganizationController{
		OrganizationSettingUsecase: &usecase.OrganizationSettingUsecase{
			Repo:          &repository.OrganizationSettingRepository{DBConn: dbConn},
			ComponentRepo: &repository.ComponentRepository{DBConn: dbConn},
		},
		AuthorizeUseCase: &usecase.AuthorizeUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			SessionRepository:      sessionRepository,
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			UserEntityUseCase: &usecase.UserEntityUseCase{
				DBConn: dbConn,
				UserBlockSettingUsecase: &usecase.UserBlockSettingUsecase{
					Repo: &repository.UserBlockSettingRepository{DBConn: dbConn},
				},
				UserImagesUsecase: &usecase.UserImagesUsecase{
					Repo:      &repository.UserImagesRepository{DBConn: dbConn},
					ImageRepo: &repository.ImageRepository{DBConn: dbConn},
					GetImageUseCase: &usecase.GetImageUseCase{
						UploadProvider:  s3Provider,
						ImageRepository: &repository.ImageRepository{DBConn: dbConn},
					},
				},
			},
		},
	}

	deviceController := &controller.DeviceController{
		DBConn: dbConn,
		UpdateDeviceSheetUseCase: &usecase.UpdateDeviceSheetUseCase{
			DeviceRepository: deviceRepository,
		},
		RegisterDeviceUseCase: &usecase.RegisterDeviceUseCase{
			DeviceRepository:  deviceRepository,
			SessionRepository: &sessionRepository,
			SettingRepository: &repository.SettingRepository{DBConn: dbConn},
			Writer:            userSpreadsheet.Writer,
			Reader:            userSpreadsheet.Reader,
		},
		GetDeviceByIDUseCase: &usecase.GetDeviceByIDUseCase{
			DeviceRepository: deviceRepository,
		},
		GetDeviceListUseCase: &usecase.GetDeviceListUseCase{
			DeviceRepository: deviceRepository,
		},
		UpdateDeviceUseCase: &usecase.UpdateDeviceUseCase{
			DeviceRepository: deviceRepository,
		},
		FindDeviceFromRequestCase: &usecase.FindDeviceFromRequestCase{
			DeviceRepository:  deviceRepository,
			SessionRepository: &sessionRepository,
		},
		GetFormByIDUseCase: &usecase.GetFormByIDUseCase{
			FormRepository: formRepo,
		},
		TakeNoteUseCase: &usecase.TakeNoteUseCase{
			DeviceRepository: deviceRepository,
		},
		RefreshAccessTokenUseCase: &usecase.RefreshAccessTokenUseCase{
			SessionRepository: &sessionRepository,
			DeviceRepository:  deviceRepository,
		},
		DiscoverUseCase: &usecase.DiscoverUseCase{
			DeviceRepository: deviceRepository,
		},
		DeviceSignUpUseCases: &usecase.DeviceSignUpUseCases{
			SettingRepository: &repository.SettingRepository{DBConn: dbConn},
			FormRepository:    formRepo,
			GetQuestionsByFormUseCase: &usecase.GetQuestionsByFormUseCase{
				QuestionRepository:     &repository.QuestionRepository{DBConn: dbConn},
				CodeCountingRepository: repository.NewCodeCountingRepository(),
				DB:                     dbConn,
			},
		},
		GetRecentSubmissionByFormIDUseCase:     usecase.NewGetRecentSubmissionByFormIDUseCase(dbConn),
		GetSubmissionByConditionUseCase:        usecase.NewGetSubmissionByConditionUseCase(dbConn),
		GetTotalNrSubmissionByConditionUseCase: usecase.NewGetTotalNrSubmissionByConditionUseCase(dbConn),
		GetSubmission4MemoriesFormUseCase:      usecase.NewGetSubmission4MemoriesFormUseCase(dbConn),
		RegisterFcmDeviceUseCase:               usecase.NewRegisterFcmDeviceUseCase(dbConn, fcm),
		SendNotificationUseCase:                usecase.NewSendNotificationUseCase(dbConn, fcm),
		ResetCodeCountingUseCase:               usecase.NewResetCodeCountingUseCase(dbConn),
		GetUserFromTokenUseCase: &usecase.GetUserFromTokenUseCase{
			UserEntityRepository: repository.UserEntityRepository{DBConn: dbConn},
			SessionRepository:    sessionRepository,
		},
		GetUserDeviceUseCase: &usecase.GetUserDeviceUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		GetDeviceStatusUseCase: &usecase.GetDeviceStatusUseCase{
			DeviceRepository: deviceRepository,
		},
		OrgDeviceRegistrationUseCase: &usecase.OrgDeviceRegistrationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			DeviceRepository:       deviceRepository,
		},
		GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		ChildUseCase: childUseCase,
		DeviceUsecase: &usecase.DeviceUsecase{
			DeviceRepository: &repository.DeviceRepository{DBConn: dbConn},
			DeviceMenuUseCase: &usecase.DeviceMenuUseCase{
				Repo:          &repository.DeviceMenuRepository{DBConn: dbConn},
				ComponentRepo: &repository.ComponentRepository{DBConn: dbConn},
				DeviceRepo:    &repository.DeviceRepository{DBConn: dbConn},
			},
			ValuesAppCurrentRepository: &repository.ValuesAppCurrentRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	}

	org := engine.Group("/v1/admin/organization", secureMiddleware.Secured())
	{
		org.GET("/setting/device/:device_id", orgController.GetOrgSetting)
		org.POST("/setting/device", orgController.UploadOrgSetting)
		org.GET("/:organization_id/device", deviceController.GetAllDeviceByOrgID)
		org.GET("/:organization_id/device/:device_id", deviceController.GetDevice4Web)
		org.POST("/:organization_id/setting/device/news", orgController.UploadOrgSettingNewsDevice)
		org.POST("/:organization_id/setting/portal/news", orgController.UploadOrgSettingNewsPortal)
		org.GET("/:organization_id/setting/news", orgController.GetOrgSettingNews)
		org.PUT("/:organization_id/device/:device_id", deviceController.UploadDeviceByOrg4Web)
		// switch to organization admin
		org.GET("/switch/:organization_id", secureMiddleware.ValidateSuperAdminRole(), orgController.SwitchToOrganizationAdmin)
		// delete device by org
		org.DELETE("/:organization_id/device/:device_id", deviceController.DeleteDeviceByOrgID)
	}

	sync := engine.Group("/v1/admin/sync", secureMiddleware.ValidateSuperAdminRole())
	{
		sync.POST("/form", applicationController.SyncDataDemoV3)
		sync.GET("/form/check-status", applicationController.CheckStatusSyncQueue)
		sync.GET("/form/sync-queues", applicationController.GetAllSycnQueue)
	}

	// languages config

	languagesConfigController := &controller.LanguagesConfigController{
		LanguagesConfigUsecase: &usecase.LanguagesConfigUsecase{
			Repo: &repository.LanguagesConfigRepository{DBConn: dbConn},
		},
	}

	languagesConfig := engine.Group("/v1/admin/languages-config", secureMiddleware.Secured())
	{
		languagesConfig.GET("", languagesConfigController.GetByOwner)
		languagesConfig.POST("", languagesConfigController.UploadLanguagesConfig)
	}

	executor := &TimeMachineSubscriber{
		ImportFormsUseCase:        importFormsUseCase,
		ImportRedirectUrlsUseCase: importUrlsUseCase,
		SettingRepository:         settingRepository,
		ImportToDoListUseCase:     usecase.NewImportToDoListUseCase(config, dbConn, uploaderSpreadsheet.Reader, uploaderSpreadsheet.Writer, usecase.TheTimeMachine),
	}
	//usecase.TheTimeMachine.Start(config.CronJobInterval)
	// appSetting, _ := getSettingsUseCase.GetSettings()
	var formInterval uint64 = 0
	var formInterval2 uint64 = 0
	var formInterval3 uint64 = 0
	var formInterval4 uint64 = 0
	var redirectInterval uint64 = 0
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

		deviceComponentValues.GET("/organization/:organization_id", deviceComponentValuesController.GetDeviceComponentValuesByOrganization)

		deviceComponentValues.GET("/device/:organization_id", deviceComponentValuesController.GetDeviceComponentValuesByDevice)

		deviceComponentValues.POST("/organization", secureMiddleware.ValidateSuperAdminRole(), deviceComponentValuesController.SaveDeviceComponentValuesByOrganization)

		deviceComponentValues.POST("/device", secureMiddleware.ValidateSuperAdminRole(), deviceComponentValuesController.SaveDeviceComponentValuesByOrganization)
	}

	usecase.TheTimeMachine.Start(formInterval, redirectInterval, toDosInterval, formInterval2, formInterval3, formInterval4)

	usecase.TheTimeMachine.SubscribeFormsExec(executor)
	// usecase.TheTimeMachine.SubscribeForms2Exec(executor)
	usecase.TheTimeMachine.SubscribeForms3Exec(executor)
	usecase.TheTimeMachine.SubscribeForms4Exec(executor)
	usecase.TheTimeMachine.SubscribeUrlsExec(executor)
	usecase.TheTimeMachine.SubscribeSyncDevicesExec(executor)
	usecase.TheTimeMachine.SubscribeSyncToDosExec(executor)
	usecase.TheTimeMachine.SubscribeGoogleAPIRequestMonitorExec(executor)
}

type TimeMachineSubscriber struct {
	*usecase.ImportFormsUseCase
	*usecase.ImportRedirectUrlsUseCase
	*repository.SettingRepository
	*usecase.ImportToDoListUseCase
}

func (t *TimeMachineSubscriber) ExecuteSyncUrls() {
	log.Debug("Start sync urls")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	urlSetting, err := t.GetUrlSettings()
	if err != nil {
		log.Error(err)
	} else {
		log.Debug(urlSetting)
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(urlSetting.Settings), &importSetting)
		if err != nil {
			log.Error(err)
		} else {
			err = t.SyncUrls(request.ImportRedirectUrlsRequest{
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
	formSettings, err := t.GetFormSettings()
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
			err = t.SyncForms(request.ImportFormRequest{
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

// func (t *TimeMachineSubscriber) ExecuteSyncForms2() {
// 	log.Debug("Start sync forms")
// 	type ImportSetting struct {
// 		SpreadSheetUrl string `json:"spreadsheet_url"`
// 		AutoImport     bool   `json:"auto"`
// 		Interval       uint64 `json:"interval"`
// 	}

// 	log.Debug("TimeMachineSubscriber: ExecuteSyncForms")
// 	formSettings, err := t.GetFormSettings2()
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	} else {
// 		log.Debug(formSettings)
// 		var importSetting ImportSetting
// 		err = json.Unmarshal([]byte(formSettings.Settings), &importSetting)
// 		if err != nil {
// 			log.Error(err)
// 		} else {
// 			err = t.SyncForms(request.ImportFormRequest{
// 				SpreadsheetUrl: importSetting.SpreadSheetUrl,
// 				AutoImport:     importSetting.AutoImport,
// 				Interval:       importSetting.Interval,
// 			})
// 			if err != nil {
// 				log.Error(err)
// 			}
// 		}
// 	}
// }

func (t *TimeMachineSubscriber) ExecuteSyncForms3() {
	log.Debug("Start sync forms")
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint64 `json:"interval"`
	}

	log.Debug("TimeMachineSubscriber: ExecuteSyncForms")
	formSettings, err := t.GetFormSettings3()
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
			err = t.SyncForms(request.ImportFormRequest{
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
	formSettings, err := t.GetFormSettings4()
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
			err = t.SyncForms(request.ImportFormRequest{
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

	deviceSetting, err := t.GetSyncToDosSettings()

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
				err = t.ImportToDoList(request.ImportFormRequest{
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
	monitor.ResetGoogleAPIRequestMonitor()
}
