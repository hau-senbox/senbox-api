package router

import (
	"context"
	"os"
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/middleware"
	"sen-global-api/pkg/sheet"
	"sen-global-api/pkg/uploader"
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

	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("error getting current directory: %w", err)
		return
	}
	driveService, err := drive.NewService(context.Background(), option.WithCredentialsFile(pwd+"/credentials/google_service_account.json"))
	if err != nil {
		log.Fatal("Unable to access Drive API:", err)
	}
	formRepo := &repository.FormRepository{DBConn: dbConn, DefaultRequestPageSize: config.DefaultRequestPageSize}
	deviceRepository := &repository.DeviceRepository{DBConn: dbConn, DefaultRequestPageSize: config.DefaultRequestPageSize, DefaultOutputSpreadsheetUrl: config.OutputSpreadsheetUrl}
	userEntityRepository := repository.UserEntityRepository{DBConn: dbConn}
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
	)
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
		SubmitFormUseCase: &usecase.SubmitFormUseCase{
			FormRepository:         formRepo,
			QuestionRepository:     &repository.QuestionRepository{DBConn: dbConn},
			SubmissionRepository:   &repository.SubmissionRepository{DBConn: dbConn},
			SettingRepository:      &repository.SettingRepository{DBConn: dbConn},
			FormQuestionRepository: &repository.FormQuestionRepository{DBConn: dbConn},
			AnswerRepository:       &repository.AnswerRepository{DBConn: dbConn},
			DeviceRepository:       deviceRepository,
			UserEntityRepository:   &userEntityRepository,
			CodeCountingRepository: repository.NewCodeCountingRepository(),
			Writer:                 userSpreadsheet.Writer,
			Reader:                 userSpreadsheet.Reader,
			OutputSpreadsheetID:    config.Google.SpreadsheetID,
			DriveService:           driveService,
			FirebaseApp:            fcm,
			DB:                     dbConn,
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
			UserEntityRepository: userEntityRepository,
			SessionRepository:    sessionRepository,
		},
		GetUserDeviceUseCase: &usecase.GetUserDeviceUseCase{
			UserEntityRepository: &userEntityRepository,
		},
		GetDeviceStatusUseCase: &usecase.GetDeviceStatusUseCase{
			DeviceRepository: deviceRepository,
		},
		OrgDeviceRegistrationUseCase: &usecase.OrgDeviceRegistrationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			DeviceRepository:       deviceRepository,
		},
		GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
			UserEntityRepository:   &userEntityRepository,
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		ChildUseCase: childUseCase,
	}

	answerRepo := repository.AnswerRepository{DBConn: dbConn}
	questionRepo := repository.QuestionRepository{DBConn: dbConn}
	answerUseCase := usecase.NewAnswerUseCase(answerRepo, userEntityRepository, questionRepo)
	answerController := controller.NewAnswerController(answerUseCase)

	provider := uploader.NewS3Provider(
		config.S3.SenboxFormSubmitBucket.AccessKey,
		config.S3.SenboxFormSubmitBucket.SecretKey,
		config.S3.SenboxFormSubmitBucket.BucketName,
		config.S3.SenboxFormSubmitBucket.Region,
		config.S3.SenboxFormSubmitBucket.Domain,
		config.S3.SenboxFormSubmitBucket.CloudfrontKeyGroupID,
		config.S3.SenboxFormSubmitBucket.CloudfrontKeyPath,
	)

	imageController := &controller.ImageController{
		GetImageUseCase: &usecase.GetImageUseCase{
			ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
		UploadImageUseCase: &usecase.UploadImageUseCase{
			ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
		DeleteImageUseCase: &usecase.DeleteImageUseCase{
			ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
	}

	videoController := &controller.VideoController{
		GetVideoUseCase: &usecase.GetVideoUseCase{
			VideoRepository: &repository.VideoRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
		UploadVideoUseCase: &usecase.UploadVideoUseCase{
			VideoRepository: &repository.VideoRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
		DeleteVideoUseCase: &usecase.DeleteVideoUseCase{
			VideoRepository: &repository.VideoRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
	}

	audioController := &controller.AudioController{
		GetAudioUseCase: &usecase.GetAudioUseCase{
			AudioRepository: &repository.AudioRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
		UploadAudioUseCase: &usecase.UploadAudioUseCase{
			AudioRepository: &repository.AudioRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
		DeleteAudioUseCase: &usecase.DeleteAudioUseCase{
			AudioRepository: &repository.AudioRepository{DBConn: dbConn},
			UploadProvider:  provider,
		},
	}

	pdfController := &controller.PdfController{
		UploadPDFUseCase: &usecase.UploadPDFUseCase{
			PdfRepository:  &repository.PdfRepository{DBConn: dbConn},
			UploadProvider: provider,
		},
		GetPdfByKeyUseCase: &usecase.GetPdfByKeyUseCase{
			PdfRepository:  &repository.PdfRepository{DBConn: dbConn},
			UploadProvider: provider,
		},
		DeletePDFUseCase: &usecase.DeletePDFUseCase{
			PdfRepository:  &repository.PdfRepository{DBConn: dbConn},
			UploadProvider: provider,
		},
	}

	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	v1 := engine.Group("v1/device")
	{
		v1.GET("/:device_id", deviceController.GetDeviceByID)
		v1.GET("/user/:user_id", deviceController.GetAllDeviceByUserID)
		v1.GET("/org/:organization_id", deviceController.GetAllDeviceByOrgID)
		v1.POST("/org", deviceController.RegisterOrgDevice)
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

		v1.GET("/status/:device_id", secureMiddleware.Secured(), deviceController.GetDeviceStatus)

		v1.POST("/reserve", deviceController.Reserve)

		v1.POST("/discover", deviceController.Discover)

		v1.GET("/sign-up", deviceController.GetSignUp)

		v1.GET("/sign-up/form", deviceController.GetSignUpForm)

		v1.GET("/sign-up/pre-set-2", deviceController.GetPreset2)

		v1.GET("/sign-up/pre-set-1", deviceController.GetPreset1)
	}

	form := engine.Group("v1/form", secureMiddleware.Secured())
	{
		form.POST("/submit", deviceController.SubmitForm)
		form.POST("/submission/last", deviceController.GetLastSubmissionByForm)
		form.POST("/get-submission-by-condition", deviceController.GetSubmissionByCondition)
		form.POST("/get-total-nr-submission-by-condition", deviceController.GetTotalNrSubmissionByCondition)
		form.POST("/get-memory-form", deviceController.GetSubmission4Memories)
		// form.GET("/submission/get-for-memories", deviceController.GetSubmissionChildProfile)
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
		setting.POST("/notification", secureMiddleware.Secured(), deviceController.SenNotification)
	}

	codeCounting := engine.Group("v1/code-counting", secureMiddleware.Secured())
	{
		codeCounting.PUT("/reset", deviceController.ResetCodeCounting)
	}

	image := engine.Group("v1/images", secureMiddleware.Secured())
	{
		image.POST("/", imageController.GetUrlByKey)
		image.GET("/icons", imageController.GetIcons)
		image.POST("/icon", imageController.GetIconByKey)
		image.POST("/upload", imageController.CreateImage)
		image.POST("/delete", imageController.DeleteImage)
	}

	video := engine.Group("v1/videos", secureMiddleware.Secured())
	{
		video.POST("/", videoController.GetUrlByKey)
		video.POST("/upload", videoController.CreateVideo)
		video.POST("/delete", videoController.DeleteVideo)
	}

	audio := engine.Group("v1/audios", secureMiddleware.Secured())
	{
		audio.POST("/", audioController.GetUrlByKey)
		audio.POST("/upload", audioController.CreateAudio)
		audio.POST("/delete", audioController.DeleteAudio)
	}

	pdf := engine.Group("v1/pdfs", secureMiddleware.Secured())
	{
		pdf.POST("upload", pdfController.CreatePDF)
		pdf.POST("", pdfController.GetUrlByKey)
		pdf.GET("", pdfController.GetAllKeyByOrgID)
		pdf.DELETE("", pdfController.DeletePDF)
	}

	answer := engine.Group("v1/answer", secureMiddleware.Secured())
	{
		answer.POST("/get-by-key-db", answerController.GetByKeyAndDB)
	}
}
