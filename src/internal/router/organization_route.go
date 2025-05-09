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

func setupOrganizationRoutes(engine *gin.Engine, dbConn *gorm.DB, config config.AppConfig) {
	sessionRepository := repository.SessionRepository{
		AuthorizeEncryptKey: config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}

	userEntityRepository := repository.UserEntityRepository{DBConn: dbConn}

	organizationController := &controller.OrganizationController{
		GetOrganizationUseCase: &usecase.GetOrganizationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		CreateOrganizationUseCase: &usecase.CreateOrganizationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		UserJoinOrganizationUseCase: &usecase.UserJoinOrganizationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			SessionRepository:      sessionRepository,
		},
		GetUserFromTokenUseCase: &usecase.GetUserFromTokenUseCase{
			UserEntityRepository: userEntityRepository,
			SessionRepository:    sessionRepository,
		},
	}

	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	user := engine.Group("/v1/organization")
	{
		user.GET("/", secureMiddleware.Secured(), organizationController.GetAllOrganization)
		user.GET("/:id", secureMiddleware.Secured(), organizationController.GetOrganizationById)
		user.GET("/:id/users", secureMiddleware.Secured(), organizationController.GetAllUserByOrganization)
		user.POST("/", secureMiddleware.Secured(), secureMiddleware.ValidateSuperAdminRole(), organizationController.CreateOrganization)
		user.POST("/join", secureMiddleware.Secured(), organizationController.UserJoinOrganization)
	}
}
