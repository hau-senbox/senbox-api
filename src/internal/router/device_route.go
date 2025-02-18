package router

import (
	"context"
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/middleware"
	"sen-global-api/pkg/sheet"
	"time"

	firebase "firebase.google.com/go/v4"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupDeviceRoutes(engine *gin.Engine, dbConn *gorm.DB, userSpreadsheet *sheet.Spreadsheet, config config.AppConfig, fcm *firebase.App) {
	sessionRepository := repository.SessionRepository{
		AuthorizeEncryptKey: config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}
	driveService, err := drive.NewService(context.Background(), option.WithCredentialsFile("./credentials/google_service_account.json"))
	if err != nil {
		log.Fatal("Unable to access Drive API:", err)
	}
	formRepo := &repository.FormRepository{DBConn: dbConn, DefaultRequestPageSize: config.DefaultRequestPageSize}
	deviceRepository := &repository.DeviceRepository{DBConn: dbConn, DefaultRequestPageSize: config.DefaultRequestPageSize, DefaultOutputSpreadsheetUrl: config.OutputSpreadsheetUrl}
	userEntityRepository := repository.UserEntityRepository{DBConn: dbConn}
	deviceController := &controller.DeviceController{
		DBConn: dbConn,
		UpdateDeviceSheetUseCase: &usecase.UpdateDeviceSheetUseCase{
			DeviceRepository: deviceRepository,
		},
		RegisterDeviceUseCase: &usecase.RegisterDeviceUseCase{
			UserRepository:    &repository.UserRepository{DBConn: dbConn},
			DeviceRepository:  deviceRepository,
			SessionRepository: &sessionRepository,
			SettingRepository: &repository.SettingRepository{DBConn: dbConn},
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
			DeviceRepository: deviceRepository,
		},
		FindDeviceFromRequestCase: &usecase.FindDeviceFromRequestCase{
			DeviceRepository:  deviceRepository,
			SessionRepository: &sessionRepository,
		},
		GetFormByIdUseCase: &usecase.GetFormByIdUseCase{
			FormRepository: formRepo,
		},
		TakeNoteUseCase: &usecase.TakeNoteUseCase{
			DeviceRepository: deviceRepository,
		},
		SubmitFormUseCase: &usecase.SubmitFormUseCase{
			FormRepository:         formRepo,
			QuestionRepository:     &repository.QuestionRepository{DBConn: dbConn},
			SubmissionRepository:   &repository.SubmissionRepository{DBConn: dbConn},
			SettingRepository:      &repository.SettingRepository{DBConn: dbConn},
			FormQuestionRepository: &repository.FormQuestionRepository{DBConn: dbConn},
			DeviceRepository:       deviceRepository,
			UserEntityRepository:   &userEntityRepository,
			CodeCountingRepository: repository.NewCodeCountingRepository(),
			Writer:                 userSpreadsheet.Writer,
			Reader:                 userSpreadsheet.Reader,
			OutputSpreadsheetId:    config.Google.SpreadsheetId,
			DriveService:           driveService,
			FirebaseApp:            fcm,
			DB:                     dbConn,
		},
		GetScreenButtonsByDeviceUseCase: &usecase.GetScreenButtonsByDeviceUseCase{
			Reader: userSpreadsheet.Reader,
		},
		GetTimeTableUseCase:      &usecase.GetTimeTableUseCase{DeviceRepository: deviceRepository, Reader: userSpreadsheet.Reader},
		GetSettingMessageUseCase: &usecase.GetSettingMessageUseCase{DeviceRepository: deviceRepository, Reader: userSpreadsheet.Reader},
		RefreshAccessTokenUseCase: &usecase.RefreshAccessTokenUseCase{
			SessionRepository: &sessionRepository,
			DeviceRepository:  deviceRepository,
		},
		UpdateDeviceInfoUseCase: &usecase.UpdateDeviceInfoUseCase{
			DeviceRepository:  deviceRepository,
			SettingRepository: &repository.SettingRepository{DBConn: dbConn},
			Reader:            userSpreadsheet.Reader,
			Writer:            userSpreadsheet.Writer,
		},
		SyncSubmissionUseCase: &usecase.SyncSubmissionUseCase{
			SubmissionRepository: &repository.SubmissionRepository{DBConn: dbConn},
			// DeviceRepository:      deviceRepository,
			FormRepository:        formRepo,
			QuestionRepository:    &repository.QuestionRepository{DBConn: dbConn},
			SettingRepository:     &repository.SettingRepository{DBConn: dbConn},
			UserSpreadsheetReader: userSpreadsheet.Reader,
			UserSpreadsheetWriter: userSpreadsheet.Writer,
			SendEmailUseCase: &usecase.SendEmailUseCase{
				SMTPConfig:        config.SMTP,
				SettingRepository: &repository.SettingRepository{DBConn: dbConn},
				Writer:            userSpreadsheet.Writer,
			},
			GetSettingMessageUseCase: &usecase.GetSettingMessageUseCase{
				DeviceRepository: deviceRepository,
				Reader:           userSpreadsheet.Reader,
			},
			UserEntityRepository: &userEntityRepository,
		},
		GetModeLUseCase: &usecase.GetModeLUseCase{
			Reader: userSpreadsheet.Reader,
			Writer: userSpreadsheet.Writer,
		},
		DiscoverUseCase: &usecase.DiscoverUseCase{
			DeviceRepository: deviceRepository,
		},
		DeviceSignUpUseCases: &usecase.DeviceSignUpUseCases{
			SettingRepository: &repository.SettingRepository{DBConn: dbConn},
			FormRepository:    formRepo,
			GetQuestionsByFormUseCase: &usecase.GetQuestionsByFormUseCase{
				QuestionRepository:          &repository.QuestionRepository{DBConn: dbConn},
				DeviceFormDatasetRepository: &repository.DeviceFormDatasetRepository{DBConn: dbConn},
				CodeCountingRepository:      repository.NewCodeCountingRepository(),
				DB:                          dbConn,
			},
		},
		GetBrandLogoCase: &usecase.GetBrandLogoCase{
			Reader: userSpreadsheet.Reader,
			Writer: userSpreadsheet.Writer,
		},
		GetRecentSubmissionByFormIdUseCase: usecase.NewGetRecentSubmissionByFormIdUseCase(dbConn),
		RegisterFcmDeviceUseCase:           usecase.NewRegisterFcmDeviceUseCase(dbConn, fcm),
		SendNotificationUseCase:            usecase.NewSendNotificationUseCase(dbConn, fcm),
		ResetCodeCountingUseCase:           usecase.NewResetCodeCountingUseCase(dbConn),
		GetDevicesByUserIdUseCase: &usecase.GetDevicesByUserIdUseCase{
			DeviceRepository: deviceRepository,
		},
		GetUserFromTokenUseCase: &usecase.GetUserFromTokenUseCase{
			UserEntityRepository: userEntityRepository,
			SessionRepository:    sessionRepository,
		},
		GetUserDeviceUseCase: &usecase.GetUserDeviceUseCase{
			UserEntityRepository: &userEntityRepository,
		},
	}

	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	v1 := engine.Group("v1/device")
	{
		v1.GET("/:device_id", deviceController.GetDeviceById)
		v1.GET("/user/:user_id", deviceController.GetAllDeviceByUserId)
		// Init for first setting
		v1.POST("/init", secureMiddleware.Secured(), deviceController.InitDeviceV1)
		v1.POST("/refresh-token", deviceController.RefreshAccessToken)
		v1.POST("/messaging/fcm/register", deviceController.RegisterFCM)
		v1.PUT("/note", secureMiddleware.Secured(), deviceController.TakeNote)
		smtpController := &controller.SMTPController{
			SendEmailUseCase: &usecase.SendEmailUseCase{
				SMTPConfig:        config.SMTP,
				SettingRepository: &repository.SettingRepository{DBConn: dbConn},
				Writer:            userSpreadsheet.Writer,
			},
			FindDeviceFromRequestCase: &usecase.FindDeviceFromRequestCase{
				DeviceRepository:  deviceRepository,
				SessionRepository: &sessionRepository,
			},
		}
		v1.POST("/send/email", smtpController.SendEmailFromDevice)
		v1.GET("/time-table", secureMiddleware.Secured(), deviceController.GetTimeTable)
		v1.GET("/messages", secureMiddleware.Secured(), deviceController.GetSettingMessage)

		v1.PUT("/update-info", secureMiddleware.Secured(), deviceController.UpdateDeviceInfo)

		v1.GET("/status/:device_id", secureMiddleware.Secured(), deviceController.GetDeviceStatus)

		v1.POST("/reserve", deviceController.Reserve)
		v1.GET("/mode-l", secureMiddleware.Secured(), deviceController.GetModeL)
		v1.GET("/brand-logo", secureMiddleware.Secured(), deviceController.GetBrandLogo)

		v1.POST("/discover", deviceController.Discover)

		v1.GET("/sign-up", deviceController.GetSignUp)

		v1.GET("/sign-up/form", deviceController.GetSignUpForm)

		v1.POST("/sign-up/form", deviceController.SubmitSignUpForm)

		v1.GET("/sign-up/pre-set-2", deviceController.GetPreset2)

		v1.GET("/sign-up/pre-set-1", deviceController.GetPreset1)
	}

	form := engine.Group("v1/form", secureMiddleware.Secured())
	{
		form.POST("/submit", deviceController.SubmitForm)
		form.GET("/submission/last", deviceController.GetLastSubmissionByForm)

		form.GET("/device/sign-up", deviceController.GetDeviceSignUp)
		form.PUT("/device/sign-up", deviceController.UpdateDeviceSignUp)
	}

	redirectUrl := engine.Group("v1/redirect-url")
	{
		redirectController := &controller.RedirectUrlController{
			SaveRedirectUrlUseCase:        nil,
			GetRedirectUrlListUseCase:     nil,
			DeleteRedirectUrlUseCase:      nil,
			UpdateRedirectUrlUseCase:      nil,
			GetRedirectUrlByQRCodeUseCase: &usecase.GetRedirectUrlByQRCodeUseCase{RedirectUrlRepository: &repository.RedirectUrlRepository{DBConn: dbConn}},
		}
		redirectUrl.GET("", redirectController.GetRedirectUrlByQRCode)
	}

	setting := engine.Group("v1/buttons")
	{
		setting.GET("/screen/:device_id", deviceController.GetScreenButtons)
		setting.GET("/top", deviceController.GetTopButtons)
		setting.POST("/notification", secureMiddleware.Secured(), deviceController.SenNotification)
	}

	codeCounting := engine.Group("v1/code-counting", secureMiddleware.Secured())
	{
		codeCounting.PUT("/reset", deviceController.ResetCodeCounting)
	}
}
