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

func setupCompanyRoutes(engine *gin.Engine, dbConn *gorm.DB, config config.AppConfig) {
	companyController := &controller.CompanyController{
		GetCompanyUseCase: &usecase.GetCompanyUseCase{
			CompanyRepository: &repository.CompanyRepository{DBConn: dbConn},
		},
	}

	sessionRepository := repository.SessionRepository{
		AuthorizeEncryptKey: config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}

	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	user := engine.Group("v1/company")
	{
		user.GET("/:id", secureMiddleware.Secured(), companyController.GetCompanyById)
	}
}
