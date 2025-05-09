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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupQuestionRoutes(engine *gin.Engine, conn *gorm.DB, config config.AppConfig) {
	sessionRepository := repository.SessionRepository{
		AuthorizeEncryptKey: config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}
	userEntityRepository := repository.UserEntityRepository{DBConn: conn}
	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}
	questionRepository := repository.QuestionRepository{DBConn: conn}
	formRepo := &repository.FormRepository{DBConn: conn, DefaultRequestPageSize: config.DefaultRequestPageSize}
	ctx := context.Background()
	spreadSheet, _ := sheet.NewUserSpreadsheet(config, ctx)
	questionController := controller.QuestionController{
		DBConn: conn,
		GetUserQuestionsUseCase: usecase.GetUserQuestionsUseCase{
			DeviceQuestionRepository: repository.DeviceQuestionRepository{
				DBConn: conn,
			},
		},
		GetUserFromTokenUseCase: usecase.GetUserFromTokenUseCase{
			UserEntityRepository: userEntityRepository,
			SessionRepository:    sessionRepository,
		},
		GetQuestionByIdUseCase: usecase.GetQuestionByIdUseCase{
			QuestionRepository: questionRepository,
		},
		GetDeviceIdFromTokenUseCase: usecase.GetDeviceIdFromTokenUseCase{
			SessionRepository: &sessionRepository,
			DeviceRepository:  &repository.DeviceRepository{DBConn: conn, DefaultRequestPageSize: config.DefaultRequestPageSize, DefaultOutputSpreadsheetUrl: config.OutputSpreadsheetUrl},
		},
		GetQuestionByFormUseCase: usecase.GetQuestionsByFormUseCase{
			QuestionRepository:     &questionRepository,
			CodeCountingRepository: repository.NewCodeCountingRepository(),
			DB:                     conn,
		},
		GetFormByIdUseCase: usecase.GetFormByIdUseCase{
			FormRepository: formRepo,
		},
		GetAllQuestionsUseCase: usecase.GetAllQuestionsUseCase{
			QuestionRepository: &questionRepository,
		},
		CreateFormUseCase: usecase.CreateFormUseCase{
			FormRepository:         formRepo,
			FormQuestionRepository: &repository.FormQuestionRepository{DBConn: conn},
		},
		GetRawQuestionFromSpreadsheetUseCase: usecase.GetRawQuestionFromSpreadsheetUseCase{
			SpreadsheetId:     config.Google.SpreadsheetId,
			SpreadsheetReader: spreadSheet.Reader,
		},
		GetShowPicsQuestionDetailUseCase: usecase.GetShowPicsQuestionDetailUseCase{
			QuestionRepository: &questionRepository,
		},
		FindDeviceFromRequestCase: usecase.FindDeviceFromRequestCase{
			DeviceRepository:  &repository.DeviceRepository{DBConn: conn, DefaultRequestPageSize: config.DefaultRequestPageSize, DefaultOutputSpreadsheetUrl: config.OutputSpreadsheetUrl},
			SessionRepository: &sessionRepository,
		},
		GetUserDeviceUseCase: usecase.GetUserDeviceUseCase{
			UserEntityRepository: &userEntityRepository,
		},
		GetDeviceByIdUseCase: usecase.GetDeviceByIdUseCase{
			DeviceRepository: &repository.DeviceRepository{DBConn: conn, DefaultRequestPageSize: config.DefaultRequestPageSize, DefaultOutputSpreadsheetUrl: config.OutputSpreadsheetUrl},
		},
	}

	form := engine.Group("v1/form")
	{
		form.POST("", questionController.GetFormQRCode)
	}

	question := engine.Group("v1/question", secureMiddleware.Secured())
	{
		question.GET("/show-pics", questionController.GetShowPicsQuestion)
	}
}
