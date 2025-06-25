package router

import (
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/middleware"
	"sen-global-api/pkg/uploader"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupOrganizationRoutes(engine *gin.Engine, dbConn *gorm.DB, config config.AppConfig) {
	sessionRepository := repository.SessionRepository{
		AuthorizeEncryptKey: config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}

	provider := uploader.NewS3Provider(
		config.S3.SenboxFormSubmitBucket.AccessKey,
		config.S3.SenboxFormSubmitBucket.SecretKey,
		config.S3.SenboxFormSubmitBucket.BucketName,
		config.S3.SenboxFormSubmitBucket.Region,
		config.S3.SenboxFormSubmitBucket.Domain,
		config.S3.SenboxFormSubmitBucket.CloudfrontKeyGroupID,
		config.S3.SenboxFormSubmitBucket.CloudfrontKeyPath,
	)

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
		GetOrgFormApplicationUseCase: &usecase.GetOrgFormApplicationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		ApproveOrgFormApplicationUseCase: &usecase.ApproveOrgFormApplicationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		BlockOrgFormApplicationUseCase: &usecase.BlockOrgFormApplicationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		CreateOrgFormApplicationUseCase: &usecase.CreateOrgFormApplicationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
		},
		UploadOrgAvatarUseCase: &usecase.UploadOrgAvatarUseCase{
			OrganizationRepository: repository.OrganizationRepository{DBConn: dbConn},
			UploadImageUseCase: usecase.UploadImageUseCase{
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				UploadProvider:  provider,
			},
			DeleteImageUseCase: usecase.DeleteImageUseCase{
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				UploadProvider:  provider,
			},
		},
	}

	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	org := engine.Group("/v1/organization")
	{
		org.GET("/", secureMiddleware.Secured(), organizationController.GetAllOrganization)
		org.GET("/:id", secureMiddleware.Secured(), organizationController.GetOrganizationByID)
		org.GET("/name", secureMiddleware.Secured(), organizationController.GetOrganizationByName)
		org.GET("/:id/users", secureMiddleware.Secured(), organizationController.GetAllUserByOrganization)
		org.POST("/", secureMiddleware.Secured(), secureMiddleware.ValidateSuperAdminRole(), organizationController.CreateOrganization)
		org.POST("/join", secureMiddleware.Secured(), organizationController.UserJoinOrganization)
		org.POST("/avatar", secureMiddleware.Secured(), organizationController.UploadAvatar)
	}

	application := engine.Group("/v1/organization/application")
	{
		application.GET("/", secureMiddleware.Secured(), organizationController.GetAllOrgFormApplication)
		application.GET("/:id", secureMiddleware.Secured(), organizationController.GetOrgFormApplicationByID)

		application.POST("/", secureMiddleware.Secured(), organizationController.CreateOrgFormApplication)
		application.POST("/:id/approve", secureMiddleware.Secured(), secureMiddleware.ValidateSuperAdminRole(), organizationController.ApproveOrgFormApplication)
		application.POST("/:id/block", secureMiddleware.Secured(), secureMiddleware.ValidateSuperAdminRole(), organizationController.BlockOrgFormApplication)
	}
}
