package router

import (
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupGatewayRoutes(r *gin.Engine, dbConn *gorm.DB, appCfg config.AppConfig) {
	// init repository + usecase
	sessionRepository := repository.SessionRepository{
		OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		AuthorizeEncryptKey:    appCfg.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(appCfg.TokenExpireDurationInHour),
	}
	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	studentRepo := &repository.StudentApplicationRepository{DB: dbConn}
	studentUsecase := &usecase.StudentApplicationUseCase{StudentAppRepo: studentRepo}
	userEntityCtrl := &controller.UserEntityController{StudentApplicationUseCase: studentUsecase}

	api := r.Group("/v1/gateway", secureMiddleware.Secured())
	{
		configs := api.Group("/students")
		{
			configs.GET("/:id", userEntityCtrl.GetStudent4Gateway)
		}
	}
}
