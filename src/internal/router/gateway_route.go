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

	cache "github.com/hung-senbox/senbox-cache-service/pkg/cache"
	cached_profile_gateway "github.com/hung-senbox/senbox-cache-service/pkg/cache/cached"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache/caching"
)

func setupGatewayRoutes(r *gin.Engine, dbConn *gorm.DB, appCfg config.AppConfig, consulClient *api.Client, cacheClientRedis *cache.RedisCache) {
	// init repository + usecase
	sessionRepository := repository.SessionRepository{
		OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		AuthorizeEncryptKey:    appCfg.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(appCfg.TokenExpireDurationInHour),
	}
	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	// gateway init
	departmentGW := gateway.NewDepartmentGateway("department-service", consulClient)
	cachedProfileGateway := cached_profile_gateway.NewCachedProfileGateway(cacheClientRedis)
	cachingMainService := caching.NewCachingMainService(cacheClientRedis, 0)
	cachingProfileService := caching.NewCachingProfileService(cacheClientRedis, 0)
	profileGw := gateway.NewProfileGateway("profile-service", consulClient, cachedProfileGateway, cachingProfileService)

	s3Provider := uploader.NewS3Provider(
		appCfg.S3.SenboxFormSubmitBucket.AccessKey,
		appCfg.S3.SenboxFormSubmitBucket.SecretKey,
		appCfg.S3.SenboxFormSubmitBucket.BucketName,
		appCfg.S3.SenboxFormSubmitBucket.Region,
		appCfg.S3.SenboxFormSubmitBucket.Domain,
		appCfg.S3.SenboxFormSubmitBucket.CloudfrontKeyGroupID,
		appCfg.S3.SenboxFormSubmitBucket.CloudfrontKeyPath,
	)

	// generateOwnerCodeUseCase
	generateOwnerCodeUseCase := usecase.NewGenerateOwnerCodeUseCase(
		&repository.UserEntityRepository{DBConn: dbConn},
		&repository.TeacherApplicationRepository{DBConn: dbConn},
		&repository.StudentApplicationRepository{DB: dbConn},
		&repository.StaffApplicationRepository{DBConn: dbConn},
		&repository.ChildRepository{DB: dbConn},
		&repository.ParentRepository{DBConn: dbConn},
		profileGw,
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
		ProfileGateway:           profileGw,
		GenerateOwnerCodeUseCase: generateOwnerCodeUseCase,
		CachingMainService:       cachingMainService,
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
		ProfileGateway:           profileGw,
		GenerateOwnerCodeUseCase: generateOwnerCodeUseCase,
		CachingMainService:       cachingMainService,
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
		ProfileGateway:           profileGw,
		GenerateOwnerCodeUseCase: generateOwnerCodeUseCase,
		CachingMainService:       cachingMainService,
	}

	// parent
	parentRepo := &repository.ParentRepository{DBConn: dbConn}
	parentUsecase := &usecase.ParentUseCase{
		DBConn:         dbConn,
		UserRepo:       userEntityRepository,
		ParentMenuRepo: &repository.ParentMenuRepository{DBConn: dbConn},
		ParentRepo:     parentRepo,
		UserImagesUsecase: &usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  s3Provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
		ProfileGateway:           profileGw,
		GenerateOwnerCodeUseCase: generateOwnerCodeUseCase,
		CachingMainService:       cachingMainService,
	}

	// child
	childUseCase := usecase.NewChildUseCase(
		dbConn,
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
		&repository.LanguageSettingRepository{DBConn: dbConn},
		&repository.ParentRepository{DBConn: dbConn},
		&repository.ParentChildsRepository{DBConn: dbConn},
		profileGw,
		generateOwnerCodeUseCase,
	)
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

	// user
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
		UserEntityUseCase: &usecase.UserEntityUseCase{
			DBConn:      dbConn,
			UserRepo:    &repository.UserEntityRepository{DBConn: dbConn},
			TeacherRepo: &repository.TeacherApplicationRepository{DBConn: dbConn},
			StaffRepo:   &repository.StaffApplicationRepository{DBConn: dbConn},
			UserBlockSettingUsecase: &usecase.UserBlockSettingUsecase{
				Repo:        &repository.UserBlockSettingRepository{DBConn: dbConn},
				TeacherRepo: &repository.TeacherApplicationRepository{DBConn: dbConn},
				StaffRepo:   &repository.StaffApplicationRepository{DBConn: dbConn},
			},
			UserImagesUsecase: &usecase.UserImagesUsecase{
				Repo:      &repository.UserImagesRepository{DBConn: dbConn},
				ImageRepo: &repository.ImageRepository{DBConn: dbConn},
				GetImageUseCase: &usecase.GetImageUseCase{
					UploadProvider:  s3Provider,
					ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				},
			},
			ProfileGateway: profileGw,
		},
		ParentUseCase: parentUsecase,
		ChildUseCase:  childUseCase,
	}
	api := r.Group("/v1/gateway", secureMiddleware.Secured())
	{

		// user
		user := api.Group("/users")
		{
			user.GET("/:user_id", userEntityCtrl.GetUser4Gateway)
			user.GET("/teacher/:teacher_id", userEntityCtrl.GetUserByTeacher)
			user.GET("/staff/:staff_id", userEntityCtrl.GetUserByStaff)
			user.POST("/code/generate", userEntityCtrl.GenerateUserCode)
		}

		// student
		student := api.Group("/students")
		{
			student.GET("/:student_id", userEntityCtrl.GetStudent4Gateway)
			student.POST("/code/generate", userEntityCtrl.GenerateStudentCode)
		}

		// teacher
		teacher := api.Group("/teachers")
		{
			teacher.GET("/:teacher_id", userEntityCtrl.GetTeacher4Gateway)
			teacher.GET("/get-by-user/:user_id", userEntityCtrl.GetTeachersByUser4Gateway)
			teacher.GET("/organization/:organization_id/user/:user_id", userEntityCtrl.GetTeacherByOrgAndUser4Gateway)
			teacher.POST("/code/generate", userEntityCtrl.GenerateTeacherCode)
		}

		// staff
		staff := api.Group("/staffs")
		{
			staff.GET("/:staff_id", userEntityCtrl.GetStaff4Gateway)
			staff.GET("/get-by-user/:user_id", userEntityCtrl.GetStaffsByUser4Gateway)
			staff.GET("/organization/:organization_id/user/:user_id", userEntityCtrl.GetStaffByOrgAndUser4Gateway)
			staff.POST("/code/generate", userEntityCtrl.GenerateStaffCode)
		}

		// prarent
		parent := api.Group("/parents")
		{
			parent.GET("/get-by-user/:user_id", userEntityCtrl.GetParentByUser4Gateway)
			parent.POST("/code/generate", userEntityCtrl.GenerateParentCode)
		}

		// child
		child := api.Group("/children")
		{
			child.POST("/code/generate", userEntityCtrl.GenerateChildCode)
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
			menu.POST("/department", menuController.UploadDepartmentMenu)
			menu.GET("/department/:department_id", menuController.GetDepartmentMenu4GW)
			// deparment org menu
			menu.POST("/department/organization", menuController.UploadDepartmentMenuOrganization)
			menu.GET("/department/:department_id/organization/:organization_id", menuController.GetDepartmentMenuOrganization4GW)
		}

		// image
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
		image := api.Group("/images")
		{
			image.POST("/get-url", imageController.GetUrlByKey)
			image.POST("/avatar/get-url", imageController.GetUrlIsMain4Owner)
			image.POST("/upload", imageController.UploadImage4GW)
			image.DELETE("/*key", imageController.DeleteImage4GW)
		}

		// video
		videoController := &controller.VideoController{
			GetVideoUseCase: &usecase.GetVideoUseCase{
				VideoRepository: &repository.VideoRepository{DBConn: dbConn},
				UploadProvider:  s3Provider,
			},
			UploadVideoUseCase: &usecase.UploadVideoUseCase{
				VideoRepository: &repository.VideoRepository{DBConn: dbConn},
				UploadProvider:  s3Provider,
			},
			DeleteVideoUseCase: &usecase.DeleteVideoUseCase{
				VideoRepository: &repository.VideoRepository{DBConn: dbConn},
				UploadProvider:  s3Provider,
			},
		}
		video := api.Group("/videos")
		{
			video.POST("/get-url", videoController.GetUrlByKey)
			video.POST("/upload", videoController.UploadVideo4GW)
			video.DELETE("/*key", videoController.DeleteVideo4GW)
		}

		// audio
		audioController := &controller.AudioController{
			GetAudioUseCase: &usecase.GetAudioUseCase{
				AudioRepository: &repository.AudioRepository{DBConn: dbConn},
				UploadProvider:  s3Provider,
			},
			UploadAudioUseCase: &usecase.UploadAudioUseCase{
				AudioRepository: &repository.AudioRepository{DBConn: dbConn},
				UploadProvider:  s3Provider,
			},
			DeleteAudioUseCase: &usecase.DeleteAudioUseCase{
				AudioRepository: &repository.AudioRepository{DBConn: dbConn},
				UploadProvider:  s3Provider,
			},
		}

		audio := api.Group("/audios")
		{
			audio.POST("/get-url", audioController.GetUrlByKey)
			audio.POST("/upload", audioController.UploadAudio4GW)
			audio.DELETE("/*key", audioController.DeleteAudio4GW)
		}

		// pdf
		pdfController := &controller.PdfController{
			GetPdfByKeyUseCase: &usecase.GetPdfByKeyUseCase{
				PdfRepository:  &repository.PdfRepository{DBConn: dbConn},
				UploadProvider: s3Provider,
			},
			UploadPDFUseCase: &usecase.UploadPDFUseCase{
				PdfRepository:  &repository.PdfRepository{DBConn: dbConn},
				UploadProvider: s3Provider,
			},
			DeletePDFUseCase: &usecase.DeletePDFUseCase{
				PdfRepository:  &repository.PdfRepository{DBConn: dbConn},
				UploadProvider: s3Provider,
			},
		}

		pdf := api.Group("/pdfs")
		{
			pdf.POST("/get-url", pdfController.GetUrlByKey4Gw)
			pdf.POST("/upload", pdfController.UpoadPDF4Gw)
			pdf.DELETE("/*key", pdfController.DeletePDF4Gw)
		}

		// message language
		messageLanguageRepo := repository.NewMessageLanguageRepository(dbConn)
		messageLanguageUsecase := usecase.NewMessageLanguageUseCase(messageLanguageRepo)
		messageLanguageController := controller.NewMessageLanguageController(messageLanguageUsecase)
		message := api.Group("/messages")
		{
			message.POST("", messageLanguageController.UploadMessageLanguages)
			message.GET("", messageLanguageController.GetMessageLanguages4GW)
			message.GET("/get-by-language", messageLanguageController.GetMessageLanguage4GW)
			message.DELETE("", messageLanguageController.DeleteMessageLanguage4GW)
		}
	}
}
