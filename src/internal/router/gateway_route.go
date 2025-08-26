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

func setupGatewayRoutes(r *gin.Engine, dbConn *gorm.DB, appCfg config.AppConfig) {
	// init repository + usecase
	sessionRepository := repository.SessionRepository{
		OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		AuthorizeEncryptKey:    appCfg.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(appCfg.TokenExpireDurationInHour),
	}
	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	s3Provider := uploader.NewS3Provider(
		appCfg.S3.SenboxFormSubmitBucket.AccessKey,
		appCfg.S3.SenboxFormSubmitBucket.SecretKey,
		appCfg.S3.SenboxFormSubmitBucket.BucketName,
		appCfg.S3.SenboxFormSubmitBucket.Region,
		appCfg.S3.SenboxFormSubmitBucket.Domain,
		appCfg.S3.SenboxFormSubmitBucket.CloudfrontKeyGroupID,
		appCfg.S3.SenboxFormSubmitBucket.CloudfrontKeyPath,
	)

	userEntityRepository := &repository.UserEntityRepository{DBConn: dbConn}

	// student
	studentRepo := &repository.StudentApplicationRepository{DB: dbConn}
	studentUsecase := &usecase.StudentApplicationUseCase{
		StudentAppRepo: studentRepo,
		UserImagesUsecase: &usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	}

	// teacher
	teacherRepo := &repository.TeacherApplicationRepository{DBConn: dbConn}
	teacherUsecase := &usecase.TeacherApplicationUseCase{
		TeacherRepo:          teacherRepo,
		UserEntityRepository: userEntityRepository,
		UserImagesUsecase: &usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	}

	// staff
	staffRepo := &repository.StaffApplicationRepository{DBConn: dbConn}
	staffUsecase := &usecase.StaffApplicationUseCase{
		StaffAppRepo:         staffRepo,
		UserEntityRepository: userEntityRepository,
		UserImagesUsecase: &usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	}

	// user entity ctl
	userEntityCtrl := &controller.UserEntityController{
		StudentApplicationUseCase: studentUsecase,
		TeacherApplicationUseCase: teacherUsecase,
		StaffApplicationUseCase:   staffUsecase,
		GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		UserImagesUsecase: &usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	}

	// organization ctl
	orgCtrl := &controller.OrganizationController{
		GetOrganizationUseCase: &usecase.GetOrganizationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
	}

	api := r.Group("/v1/gateway", secureMiddleware.Secured())
	{

		// user
		user := api.Group("/users")
		{
			user.GET("/:user_id", userEntityCtrl.GetUser4Gateway)
		}

		// student
		student := api.Group("/students")
		{
			student.GET("/:student_id", userEntityCtrl.GetStudent4Gateway)
		}

		// teacher
		teacher := api.Group("/teachers")
		{
			teacher.GET("/:teacher_id", userEntityCtrl.GetTeacher4Gateway)
			teacher.GET("/get-by-user/:user_id", userEntityCtrl.GetTeacherByUser4Gateway)
		}

		// staff
		staff := api.Group("/staffs")
		{
			staff.GET("/:staff_id", userEntityCtrl.GetStaff4Gateway)
		}

		// organization
		org := api.Group("/organizations")
		{
			org.GET("", orgCtrl.GetAllOrganizations4Gateway)
		}
	}
}
