package router

import (
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/middleware"
	"sen-global-api/pkg/consulapi/gateway"
	"sen-global-api/pkg/uploader"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

func setupGatewayRoutes(r *gin.Engine, dbConn *gorm.DB, appCfg config.AppConfig, consulClient *api.Client) {
	// init repository + usecase
	sessionRepository := repository.SessionRepository{
		OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		AuthorizeEncryptKey:    appCfg.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(appCfg.TokenExpireDurationInHour),
	}
	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	// department gateway init
	departmentGW := gateway.NewDepartmentGateway("department-service", consulClient)

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

	// menu ctl
	menuController := &controller.MenuController{
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
		},
		DepartmentMenuUseCase: &usecase.DepartmentMenuUseCase{
			DepartmentMenuRepository: &repository.DepartmentMenuRepository{DBConn: dbConn},
			ComponentRepository:      &repository.ComponentRepository{DBConn: dbConn},
			LanguageSettingRepo:      &repository.LanguageSettingRepository{DBConn: dbConn},
		},
		DepartmentMenuOrganizationUseCase: &usecase.DepartmentMenuOrganizationUseCase{
			DepartmentMenuOrganizationRepository: &repository.DepartmentMenuOrganizationRepository{DBConn: dbConn},
			ComponentRepo:                        &repository.ComponentRepository{DBConn: dbConn},
			OrganizationRepository:               &repository.OrganizationRepository{DBConn: dbConn},
			DeviceRepository:                     &repository.DeviceRepository{DBConn: dbConn},
			DepartmentGateway:                    departmentGW,
			LanguageSettingRepo:                  &repository.LanguageSettingRepository{DBConn: dbConn},
		},
	}

	// image ctl
	imageController := &controller.ImageController{
		GetImageUseCase: &usecase.GetImageUseCase{
			ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			UploadProvider:  s3Provider,
		},
		UploadImageUseCase: &usecase.UploadImageUseCase{
			ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			UploadProvider:  s3Provider,
		},
		DeleteImageUseCase: &usecase.DeleteImageUseCase{
			ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			UploadProvider:  s3Provider,
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
			teacher.GET("/get-by-user/:user_id", userEntityCtrl.GetTeachersByUser4Gateway)
			teacher.GET("/organization/:organization_id/user/:user_id", userEntityCtrl.GetTeacherByOrgAndUser4Gateway)
		}

		// staff
		staff := api.Group("/staffs")
		{
			staff.GET("/:staff_id", userEntityCtrl.GetStaff4Gateway)
			staff.GET("/get-by-user/:user_id", userEntityCtrl.GetStaffsByUser4Gateway)
			staff.GET("/organization/:organization_id/user/:user_id", userEntityCtrl.GetStaffByOrgAndUser4Gateway)
		}

		// organization
		org := api.Group("/organizations")
		{
			org.GET("", orgCtrl.GetAllOrganizations4Gateway)
		}

		// menu
		menu := api.Group("/menus")
		{
			// department menu
			menu.POST("/department", middleware.GeneralLoggerMiddleware(dbConn), menuController.UploadDepartmentMenu)
			menu.GET("/department/:department_id", menuController.GetDepartmentMenu4GW)
			// deparment org menu
			menu.POST("/department/organization", middleware.GeneralLoggerMiddleware(dbConn), menuController.UploadDepartmentMenuOrganization)
			menu.GET("/department/:department_id/organization/:organization_id", menuController.GetDepartmentMenuOrganization4GW)
		}

		// image
		image := api.Group("/images")
		{
			image.POST("/get-url", imageController.GetUrlByKey)
			image.POST("/avatar/get-url", imageController.GetUrlIsMain4Owner)
		}

		// message language
		messageLanguageRepo := repository.NewMessageLanguageRepository(dbConn)
		messageLanguageUsecase := usecase.NewMessageLanguageUseCase(messageLanguageRepo)
		messageLanguageController := controller.NewMessageLanguageController(messageLanguageUsecase)
		message := api.Group("/messages")
		{
			message.POST("", messageLanguageController.UploadMessageLanguage)
		}
	}
}
