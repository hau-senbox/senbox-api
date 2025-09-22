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

func setupUserRoutes(engine *gin.Engine, dbConn *gorm.DB, config config.AppConfig, consulClient *api.Client) {
	sessionRepository := repository.SessionRepository{
		OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		AuthorizeEncryptKey:    config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}
	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	provider := uploader.NewS3Provider(
		config.S3.SenboxFormSubmitBucket.AccessKey,
		config.S3.SenboxFormSubmitBucket.SecretKey,
		config.S3.SenboxFormSubmitBucket.BucketName,
		config.S3.SenboxFormSubmitBucket.Region,
		config.S3.SenboxFormSubmitBucket.Domain,
		config.S3.SenboxFormSubmitBucket.CloudfrontKeyGroupID,
		config.S3.SenboxFormSubmitBucket.CloudfrontKeyPath,
	)

	// department gateway init
	departmentGW := gateway.NewDepartmentGateway("department-service", consulClient)

	childUsecase := usecase.NewChildUseCase(
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
				UploadProvider:  provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
	)

	userEntityController := &controller.UserEntityController{
		GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		CreateUserEntityUseCase: &usecase.CreateUserEntityUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		CreateChildForParentUseCase: &usecase.CreateChildForParentUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		UpdateUserEntityUseCase: &usecase.UpdateUserEntityUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		UpdateUserRoleUseCase: &usecase.UpdateUserRoleUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		AuthorizeUseCase: &usecase.AuthorizeUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
			DeviceRepository:     &repository.DeviceRepository{DBConn: dbConn},
			SessionRepository:    sessionRepository,
			DBConn:               dbConn,
			ManageUserLoginUseCase: &usecase.ManageUserLoginUseCase{
				UserDevicesLoginRepository: &repository.UserDevicesLoginRepository{DBConn: dbConn},
				UserSettingRepositotry:     &repository.UserSettingRepository{DBConn: dbConn},
			},
		},
		UpdateUserOrgInfoUseCase: &usecase.UpdateUserOrgInfoUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		UpdateUserAuthorizeUseCase: &usecase.UpdateUserAuthorizeUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		DeleteUserAuthorizeUseCase: &usecase.DeleteUserAuthorizeUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		GetPreRegisterUseCase: &usecase.GetPreRegisterUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		CreatePreRegisterUseCase: &usecase.CreatePreRegisterUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		GetUserFromTokenUseCase: &usecase.GetUserFromTokenUseCase{
			UserEntityRepository: repository.UserEntityRepository{DBConn: dbConn},
			SessionRepository:    sessionRepository,
		},

		CreateUserFormApplicationUseCase: &usecase.CreateUserFormApplicationUseCase{
			UserEntityRepository:               &repository.UserEntityRepository{DBConn: dbConn},
			RoleOrgSignUpRepository:            &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			ComponentRepository:                &repository.ComponentRepository{DBConn: dbConn},
			StudentMenuRepository:              &repository.StudentMenuRepository{DBConn: dbConn},
			TeacherMenuRepository:              &repository.TeacherMenuRepository{DBConn: dbConn},
			OrganizationMenuTemplateRepository: &repository.OrganizationMenuTemplateRepository{DBConn: dbConn},
			StaffMenuRepository:                &repository.StaffMenuRepository{DBConn: dbConn},
			OrganizationRepository:             &repository.OrganizationRepository{DBConn: dbConn},
		},
		ApproveUserFormApplicationUseCase: &usecase.ApproveUserFormApplicationUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		BlockUserFormApplicationUseCase: &usecase.BlockUserFormApplicationUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		GetUserFormApplicationUseCase: &usecase.GetUserFormApplicationUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		UploadUserAvatarUseCase: &usecase.UploadUserAvatarUseCase{
			UploadImageUseCase: usecase.UploadImageUseCase{
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				UploadProvider:  provider,
			},
			DeleteImageUseCase: usecase.DeleteImageUseCase{
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				UploadProvider:  provider,
			},
			UserEntityRepository: repository.UserEntityRepository{DBConn: dbConn},
		},
		RoleOrgSignUpUseCase: &usecase.RoleOrgSignUpUseCase{
			Repo: &repository.RoleOrgSignUpRepository{DBConn: dbConn},
		},
		ChildUseCase: childUsecase,
		StudentApplicationUseCase: &usecase.StudentApplicationUseCase{
			StudentAppRepo:  &repository.StudentApplicationRepository{DB: dbConn},
			StudentMenuRepo: &repository.StudentMenuRepository{DBConn: dbConn},
			ComponentRepo:   &repository.ComponentRepository{DBConn: dbConn},
			RoleOrgRepo:     &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
				UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
				OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
			},
			DeviceRepo:       &repository.DeviceRepository{DBConn: dbConn},
			OrganizationRepo: &repository.OrganizationRepository{DBConn: dbConn},
		},
		UserBlockSettingUsecase: &usecase.UserBlockSettingUsecase{
			Repo:        &repository.UserBlockSettingRepository{DBConn: dbConn},
			TeacherRepo: &repository.TeacherApplicationRepository{DBConn: dbConn},
			StaffRepo:   &repository.StaffApplicationRepository{DBConn: dbConn},
		},
		GetUserOrganizationActiveUsecase: &usecase.GetUserOrganizationActiveUsecase{
			OrganizationRepository:       &repository.OrganizationRepository{DBConn: dbConn},
			TeacherApplicationRepository: &repository.TeacherApplicationRepository{DBConn: dbConn},
			StaffApplicationRepository:   &repository.StaffApplicationRepository{DBConn: dbConn},
		},
		UploadImageUseCase: &usecase.UploadImageUseCase{
			UploadProvider:  provider,
			ImageRepository: &repository.ImageRepository{DBConn: dbConn},
		},
		UserImagesUsecase: &usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			DeleteImageUsecase: &usecase.DeleteImageUseCase{
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				UploadProvider:  provider,
			},
			GetImageUseCase: &usecase.GetImageUseCase{
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				UploadProvider:  provider,
			},
		},
		LanguagesConfigUsecase: &usecase.LanguagesConfigUsecase{
			Repo: &repository.LanguagesConfigRepository{DBConn: dbConn},
		},
		UserSettingUseCase: &usecase.UserSettingUseCase{
			Repo: &repository.UserSettingRepository{DBConn: dbConn},
		},
	}

	userRoleController := &controller.RoleController{
		GetRoleUseCase: &usecase.GetRoleUseCase{
			RoleRepository: &repository.RoleRepository{DBConn: dbConn},
		},
		CreateRoleUseCase: &usecase.CreateRoleUseCase{
			RoleRepository: &repository.RoleRepository{DBConn: dbConn},
		},
		UpdateRoleUseCase: &usecase.UpdateRoleUseCase{
			RoleRepository: &repository.RoleRepository{DBConn: dbConn},
		},
		DeleteRoleUseCase: &usecase.DeleteRoleUseCase{
			RoleRepository: &repository.RoleRepository{DBConn: dbConn},
		},
	}

	functionClaimController := &controller.FunctionClaimController{
		GetFunctionClaimUseCase: &usecase.GetFunctionClaimUseCase{
			FunctionClaimRepository: &repository.FunctionClaimRepository{DBConn: dbConn},
		},
		CreateFunctionClaimUseCase: &usecase.CreateFunctionClaimUseCase{
			FunctionClaimRepository: &repository.FunctionClaimRepository{DBConn: dbConn},
		},
		UpdateFunctionClaimUseCase: &usecase.UpdateFunctionClaimUseCase{
			FunctionClaimRepository: &repository.FunctionClaimRepository{DBConn: dbConn},
		},
		DeleteFunctionClaimUseCase: &usecase.DeleteFunctionClaimUseCase{
			FunctionClaimRepository: &repository.FunctionClaimRepository{DBConn: dbConn},
		},
	}

	functionClaimPermissionController := &controller.FunctionClaimPermissionController{
		GetFunctionClaimPermissionUseCase: &usecase.GetFunctionClaimPermissionUseCase{
			FunctionClaimPermissionRepository: &repository.FunctionClaimPermissionRepository{DBConn: dbConn},
		},
		CreateFunctionClaimPermissionUseCase: &usecase.CreateFunctionClaimPermissionUseCase{
			FunctionClaimPermissionRepository: &repository.FunctionClaimPermissionRepository{DBConn: dbConn},
		},
		UpdateFunctionClaimPermissionUseCase: &usecase.UpdateFunctionClaimPermissionUseCase{
			FunctionClaimPermissionRepository: &repository.FunctionClaimPermissionRepository{DBConn: dbConn},
		},
		DeleteFunctionClaimPermissionUseCase: &usecase.DeleteFunctionClaimPermissionUseCase{
			FunctionClaimPermissionRepository: &repository.FunctionClaimPermissionRepository{DBConn: dbConn},
		},
	}

	menuController := &controller.MenuController{
		GetUserFromTokenUseCase: &usecase.GetUserFromTokenUseCase{
			UserEntityRepository: repository.UserEntityRepository{DBConn: dbConn},
			SessionRepository:    sessionRepository,
		},
		GetMenuUseCase: &usecase.GetMenuUseCase{
			MenuRepository:          &repository.MenuRepository{DBConn: dbConn},
			UserEntityRepository:    &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository:  &repository.OrganizationRepository{DBConn: dbConn},
			DeviceRepository:        &repository.DeviceRepository{DBConn: dbConn},
			RoleOrgSignUpRepository: &repository.RoleOrgSignUpRepository{DBConn: dbConn},
			FormRepository:          &repository.FormRepository{DBConn: dbConn},
			SubmissionRepository:    &repository.SubmissionRepository{DBConn: dbConn},
			ComponentRepository:     &repository.ComponentRepository{DBConn: dbConn},
			ChildRepository:         &repository.ChildRepository{DB: dbConn},
			StudentAppRepo:          &repository.StudentApplicationRepository{DB: dbConn},
			ChildMenuUseCase: &usecase.ChildMenuUseCase{
				Repo:          &repository.ChildMenuRepository{DBConn: dbConn},
				ComponentRepo: &repository.ComponentRepository{DBConn: dbConn},
				ChildRepo:     &repository.ChildRepository{DB: dbConn},
			},
			StudentMenuUseCase: &usecase.StudentMenuUseCase{
				StudentMenuRepo: &repository.StudentMenuRepository{DBConn: dbConn},
				ComponentRepo:   &repository.ComponentRepository{DBConn: dbConn},
				StudentAppRepo:  &repository.StudentApplicationRepository{DB: dbConn},
			},
			TeacherMenuUseCase: &usecase.TeacherMenuUseCase{
				TeacherMenuRepo:      &repository.TeacherMenuRepository{DBConn: dbConn},
				TeacherAppRepo:       &repository.TeacherApplicationRepository{DBConn: dbConn},
				ComponentRepo:        &repository.ComponentRepository{DBConn: dbConn},
				UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
			},
			StaffMenuUsecase: &usecase.StaffMenuUseCase{
				StaffMenuRepo:        &repository.StaffMenuRepository{DBConn: dbConn},
				StaffAppRepo:         &repository.StaffApplicationRepository{DBConn: dbConn},
				ComponentRepo:        &repository.ComponentRepository{DBConn: dbConn},
				UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
				GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
					UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
					OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
				},
			},
			StaffApplicationRepo: &repository.StaffApplicationRepository{DBConn: dbConn},
			TeacherRepository:    &repository.TeacherApplicationRepository{DBConn: dbConn},
			ParentMenuUsecase: &usecase.ParentMenuUseCase{
				Repo:          &repository.ParentMenuRepository{DBConn: dbConn},
				ComponentRepo: &repository.ComponentRepository{DBConn: dbConn},
				UserRepo:      &repository.UserEntityRepository{DBConn: dbConn},
			},
			UserImageUsecase: &usecase.UserImagesUsecase{
				Repo:      &repository.UserImagesRepository{DBConn: dbConn},
				ImageRepo: &repository.ImageRepository{DBConn: dbConn},
				GetImageUseCase: &usecase.GetImageUseCase{
					UploadProvider:  provider,
					ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				},
			},
			DepartmentGateway: departmentGW,
			DepartmentMenuUseCase: &usecase.DepartmentMenuUseCase{
				DepartmentMenuRepository: &repository.DepartmentMenuRepository{DBConn: dbConn},
				ComponentRepository:      &repository.ComponentRepository{DBConn: dbConn},
			},
			OrganizationEmergencyMenuRepo: &repository.OrganizationEmergencyMenuRepository{DBConn: dbConn},
		},
		UploadSuperAdminMenuUseCase: &usecase.UploadSuperAdminMenuUseCase{
			MenuRepository:      &repository.MenuRepository{DBConn: dbConn},
			ComponentRepository: &repository.ComponentRepository{DBConn: dbConn},
		},
		UploadOrgMenuUseCase: &usecase.UploadOrgMenuUseCase{
			MenuRepository:      &repository.MenuRepository{DBConn: dbConn},
			ComponentRepository: &repository.ComponentRepository{DBConn: dbConn},
		},
		UploadUserMenuUseCase: &usecase.UploadUserMenuUseCase{
			MenuRepository:      &repository.MenuRepository{DBConn: dbConn},
			ComponentRepository: &repository.ComponentRepository{DBConn: dbConn},
		},
		UploadDeviceMenuUseCase: &usecase.UploadDeviceMenuUseCase{
			MenuRepository:      &repository.MenuRepository{DBConn: dbConn},
			ComponentRepository: &repository.ComponentRepository{DBConn: dbConn},
		},
		GetOrganizationUseCase: &usecase.GetOrganizationUseCase{
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		UserImagesUsecase: &usecase.UserImagesUsecase{
			Repo:      &repository.UserImagesRepository{DBConn: dbConn},
			ImageRepo: &repository.ImageRepository{DBConn: dbConn},
			GetImageUseCase: &usecase.GetImageUseCase{
				UploadProvider:  provider,
				ImageRepository: &repository.ImageRepository{DBConn: dbConn},
			},
		},
		TeacherMenuOrganizationUseCase: &usecase.TeacherMenuOrganizationUseCase{
			TeacherMenuOrganizationRepository: &repository.TeacherMenuOrganizationRepository{DBConn: dbConn},
			ComponentRepo:                     &repository.ComponentRepository{DBConn: dbConn},
			TeacherRepo:                       &repository.TeacherApplicationRepository{DBConn: dbConn},
			DeviceRepository:                  &repository.DeviceRepository{DBConn: dbConn},
			UserImagesUsecase: &usecase.UserImagesUsecase{
				Repo:      &repository.UserImagesRepository{DBConn: dbConn},
				ImageRepo: &repository.ImageRepository{DBConn: dbConn},
				GetImageUseCase: &usecase.GetImageUseCase{
					UploadProvider:  provider,
					ImageRepository: &repository.ImageRepository{DBConn: dbConn},
				},
			},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		DepartmentMenuUseCase: &usecase.DepartmentMenuUseCase{
			DepartmentMenuRepository: &repository.DepartmentMenuRepository{DBConn: dbConn},
			ComponentRepository:      &repository.ComponentRepository{DBConn: dbConn},
		},
		DepartmentMenuOrganizationUseCase: &usecase.DepartmentMenuOrganizationUseCase{
			DepartmentMenuOrganizationRepository: &repository.DepartmentMenuOrganizationRepository{DBConn: dbConn},
			ComponentRepo:                        &repository.ComponentRepository{DBConn: dbConn},
			DeviceRepository:                     &repository.DeviceRepository{DBConn: dbConn},
			OrganizationRepository:               &repository.OrganizationRepository{DBConn: dbConn},
			DepartmentGateway:                    departmentGW,
		},
	}

	componentController := &controller.ComponentController{
		GetComponentUseCase: &usecase.GetComponentUseCase{
			ComponentRepository: &repository.ComponentRepository{DBConn: dbConn},
		},
	}

	userTokenFCMController := &controller.UserTokenFCMController{
		CreateUserTokenFCMUseCase: &usecase.CreateUserTokenFCMUseCase{
			UserTokenFCMRepository: &repository.UserTokenFCMRepository{DBConn: dbConn},
		},
		GetUserTokenFCMUseCase: &usecase.GetUserTokenFCMUseCase{
			UserTokenFCMRepository: &repository.UserTokenFCMRepository{DBConn: dbConn},
		},
	}

	userAccess := engine.Group("v1/")
	{
		loginController := &controller.LoginController{DBConn: dbConn,
			AuthorizeUseCase: usecase.AuthorizeUseCase{
				UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
				DeviceRepository:     &repository.DeviceRepository{DBConn: dbConn},
				SessionRepository:    sessionRepository,
				ManageUserLoginUseCase: &usecase.ManageUserLoginUseCase{
					UserDevicesLoginRepository: &repository.UserDevicesLoginRepository{DBConn: dbConn},
					UserSettingRepositotry:     &repository.UserSettingRepository{DBConn: dbConn},
				},
			},
		}
		logoutController := &controller.LogoutController{
			AuthorizeUseCase: &usecase.AuthorizeUseCase{
				UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
				DeviceRepository:     &repository.DeviceRepository{DBConn: dbConn},
				SessionRepository:    sessionRepository,
				ManageUserLoginUseCase: &usecase.ManageUserLoginUseCase{
					UserDevicesLoginRepository: &repository.UserDevicesLoginRepository{DBConn: dbConn},
				},
			},
		}
		userAccess.POST("/login", loginController.UserLogin)
		userAccess.POST("/logout", secureMiddleware.Secured(), logoutController.UserLogout)
	}

	// block setting
	userBlockSettingController := &controller.BlockSettingController{
		UserBlockUsecase: &usecase.UserBlockSettingUsecase{
			Repo:        &repository.UserBlockSettingRepository{DBConn: dbConn},
			TeacherRepo: &repository.TeacherApplicationRepository{DBConn: dbConn},
			StaffRepo:   &repository.StaffApplicationRepository{DBConn: dbConn},
		},
		StudentBlockUsecase: &usecase.StudentBlockSettingUsecase{
			Repo: &repository.StudentBlockSettingRepository{DBConn: dbConn},
		},
	}

	// accounts log
	accountsLogController := &controller.AccountsLogController{
		AccountsLogUseCase: &usecase.AccountsLogUseCase{
			AccountsLogRepository: &repository.AccountsLogRepository{DBConn: dbConn},
		},
	}

	user := engine.Group("v1/user")
	{
		user.GET("/current-user", secureMiddleware.Secured(), userEntityController.GetCurrentUser)
		user.GET("/all", secureMiddleware.Secured(), userEntityController.GetAllUserEntity)
		user.GET("/:id", secureMiddleware.Secured(), userEntityController.GetUserEntityByID)
		user.GET("/name/:username", secureMiddleware.Secured(), userEntityController.GetUserEntityByName)

		user.POST("/init", userEntityController.CreateUserEntity)
		user.POST("/child/init", userEntityController.CreateChildForParent)
		user.POST("/child/create", secureMiddleware.Secured(), userEntityController.CreateChild)
		user.POST("/update", secureMiddleware.Secured(), userEntityController.UpdateUserEntity)
		user.POST("/block/:id", secureMiddleware.Secured(), userEntityController.BlockUser)
		user.POST("/role/update", secureMiddleware.Secured(), userEntityController.UpdateUserRole)
		user.POST("/avatar", secureMiddleware.Secured(), userEntityController.UploadAvatar)

		user.GET("/org/:organization_id/:user_id", secureMiddleware.Secured(), userEntityController.GetUserOrgInfo)
		user.GET("/org/:organization_id/manager", secureMiddleware.Secured(), userEntityController.GetAllOrgManagerInfo)
		user.POST("/org/update", secureMiddleware.Secured(), userEntityController.UpdateUserOrgInfo)

		user.GET("/:id/func", secureMiddleware.Secured(), userEntityController.GetAllUserAuthorize)
		user.POST("/func", secureMiddleware.Secured(), userEntityController.UpdateUserAuthorize)
		user.DELETE("/func", secureMiddleware.Secured(), userEntityController.DeleteUserAuthorize)

		user.GET("/pre-register/", secureMiddleware.Secured(), userEntityController.GetAllPreRegisterUser)
		user.POST("/pre-register/", userEntityController.CreatePreRegister)
		user.GET("/role-sign-up", userEntityController.GetAllRoleOrgSignUp)
		user.GET("/child/:id", secureMiddleware.Secured(), userEntityController.GetChildByID)
		user.PUT("/child", secureMiddleware.Secured(), userEntityController.UpdateChild)

		block := user.Group("/block")
		{
			block.GET("/:user_id", userBlockSettingController.GetByUserID)
		}

		aacountsLog := user.Group("/accounts-log")
		{
			aacountsLog.POST("", secureMiddleware.Secured(), accountsLogController.CreateAccountsLog)
		}
	}

	teacherApplication := engine.Group("/v1/user/teacher/application")
	{
		// teacherApplication.GET("/", secureMiddleware.Secured(), userEntityController.GetAllTeacherFormApplication)
		// teacherApplication.GET("/:id", secureMiddleware.Secured(), userEntityController.GetTeacherFormApplicationByID)

		teacherApplication.POST("/", secureMiddleware.Secured(), userEntityController.CreateTeacherFormApplication)
		teacherApplication.POST("/:id/approve", secureMiddleware.Secured(), userEntityController.ApproveTeacherFormApplication)
		teacherApplication.POST("/:id/block", secureMiddleware.Secured(), userEntityController.BlockTeacherFormApplication)
	}

	staffApplication := engine.Group("/v1/user/staff/application")
	{
		// staffApplication.GET("/", secureMiddleware.Secured(), userEntityController.GetAllStaffFormApplication)
		// staffApplication.GET("/:id", secureMiddleware.Secured(), userEntityController.GetStaffFormApplicationByID)

		staffApplication.POST("/", secureMiddleware.Secured(), userEntityController.CreateStaffFormApplication)
		staffApplication.POST("/:id/approve", secureMiddleware.Secured(), userEntityController.ApproveStaffFormApplication)
		staffApplication.POST("/:id/block", secureMiddleware.Secured(), userEntityController.BlockStaffFormApplication)
	}

	studentApplication := engine.Group("/v1/user/student/application")
	{
		// studentApplication.GET("/", secureMiddleware.Secured(), userEntityController.GetAllStudentFormApplication)
		studentApplication.GET("/:id/:device_id", secureMiddleware.Secured(), userEntityController.GetStudent4App)
		studentApplication.POST("/", secureMiddleware.Secured(), userEntityController.CreateStudentFormApplication)
		studentApplication.POST("/:id/approve", secureMiddleware.Secured(), userEntityController.ApproveStudentFormApplication)
		studentApplication.POST("/:id/block", secureMiddleware.Secured(), userEntityController.BlockStudentFormApplication)
		studentApplication.PUT("/", secureMiddleware.Secured(), userEntityController.UpdateStudent4App)
		studentApplication.GET("/:id", secureMiddleware.Secured(), userEntityController.GetStudent4App)
		studentApplication.GET("/:id/language-config", secureMiddleware.Secured(), userEntityController.GetStudentLanguageConfig)
	}

	userRole := engine.Group("v1/user-role", secureMiddleware.Secured())
	{
		userRole.GET("/all", userRoleController.GetAllRole)
		userRole.GET("/:id", userRoleController.GetRoleByID)
		userRole.GET("/name/:role_name", userRoleController.GetRoleByName)

		userRole.POST("/init", userRoleController.CreateRole)
		userRole.POST("/", userRoleController.UpdateRole)

		userRole.DELETE("/:id", userRoleController.DeleteRole)
	}

	functionClaim := engine.Group("v1/function-claim", secureMiddleware.Secured())
	{
		functionClaim.GET("/all", functionClaimController.GetAllFunctionClaim)
		functionClaim.GET("/:id", functionClaimController.GetFunctionClaimByID)
		functionClaim.GET("/name/:function_name", functionClaimController.GetFunctionClaimByName)

		functionClaim.POST("/init", functionClaimController.CreateFunctionClaim)
		functionClaim.POST("/", functionClaimController.UpdateFunctionClaim)

		functionClaim.DELETE("/:id", functionClaimController.DeleteFunctionClaim)
	}

	functionClaimPermission := engine.Group("v1/function-claim-permission", secureMiddleware.Secured())
	{
		functionClaimPermission.GET("/function/:function_claim_id/all", functionClaimPermissionController.GetAllFunctionClaimPermission)
		functionClaimPermission.GET("/:id", functionClaimPermissionController.GetFunctionClaimPermissionByID)
		functionClaimPermission.GET("/name/:permission_name", functionClaimPermissionController.GetFunctionClaimPermissionByName)

		functionClaimPermission.POST("/init", functionClaimPermissionController.CreateFunctionClaimPermission)
		functionClaimPermission.POST("/", functionClaimPermissionController.UpdateRoleClaimPermission)

		functionClaimPermission.DELETE("/:id", functionClaimPermissionController.DeleteRoleClaimPermission)
	}

	userMenu := engine.Group("v1/user-menu", secureMiddleware.Secured())
	{
		userMenu.GET("/super-admin", menuController.GetSuperAdminMenu)
		userMenu.GET("/user/super-admin", menuController.GetSuperAdminMenu4App)
		userMenu.GET("/org/:id", menuController.GetOrgMenu)
		userMenu.GET("/user/org/:id", menuController.GetOrgMenu4App)
		userMenu.GET("/student/:id", menuController.GetStudentMenu4App)
		userMenu.GET("/teacher/:id", menuController.GetTeacherMenu4App)
		userMenu.GET("/user/:id", menuController.GetUserMenu4App)
		userMenu.GET("/device/:id", menuController.GetDeviceMenu)
		userMenu.GET("/user/device/:id", menuController.GetDeviceMenu4App)
		userMenu.GET("/device/organization/:organization_id", menuController.GetDeviceMenuByOrg)
		userMenu.GET("/section", menuController.GetSectionMenu4App)

		userMenu.POST("/super-admin", secureMiddleware.ValidateSuperAdminRole(), menuController.UploadSuperAdminMenu)
		userMenu.POST("/org", menuController.UploadOrgMenu)
		userMenu.POST("/user", menuController.UploadUserMenu)
		userMenu.POST("/device", menuController.UploadDeviceMenu)
		userMenu.GET("/common", menuController.GetCommonMenu)
		userMenu.GET("/common-by-user", menuController.GetCommonMenuByUser)
		userMenu.POST("/teacher/organization", menuController.GetTeacherMenuOrganization4App)
		// department menu org
		userMenu.GET("/department/device/:device_id/organization/:organization_id", menuController.GetDepartmentMenuOrganization4App)
		// emergency menu
		userMenu.GET("/emergency/organization/:organization_id", menuController.GetEmergencyMenu4App)
	}

	component := engine.Group("v1/component", secureMiddleware.Secured())
	{
		component.GET("/keys", componentController.GetAllComponentKey)
	}

	userTokenFCM := engine.Group("v1/user-token-fcm", secureMiddleware.Secured())
	{
		userTokenFCM.POST("/register", userTokenFCMController.CreateFCMToken)
		userTokenFCM.GET("/all/:user_id", userTokenFCMController.GetAllFCMToken)
	}

	// organization
	orgController := &controller.OrganizationController{
		OrganizationSettingUsecase: &usecase.OrganizationSettingUsecase{
			Repo:             &repository.OrganizationSettingRepository{DBConn: dbConn},
			ComponentRepo:    &repository.ComponentRepository{DBConn: dbConn},
			OrganizationRepo: &repository.OrganizationRepository{DBConn: dbConn},
		},
	}

	org := engine.Group("/v1/user/organization", secureMiddleware.Secured())
	{
		org.GET("/setting/:device_id", orgController.GetOrgSetting)
		org.GET("/:organization_id/setting/news", orgController.GetOrgSettingNews)
	}

	//device
	deviceController := &controller.DeviceController{
		GetUserFromTokenUseCase: &usecase.GetUserFromTokenUseCase{
			UserEntityRepository: repository.UserEntityRepository{DBConn: dbConn},
			//SessionRepository: repository.SessionRepository{},
		},
		DeviceUsecase: &usecase.DeviceUsecase{
			DeviceRepository: &repository.DeviceRepository{DBConn: dbConn},
		},
	}

	device := engine.Group("/v1/user/device", secureMiddleware.Secured())
	{
		device.GET("/:device_id", middleware.DataLogMiddleware(dbConn), deviceController.GetDevice4App)
	}

	// languages config
	languagesConfigController := &controller.LanguagesConfigController{
		LanguagesConfigUsecase: &usecase.LanguagesConfigUsecase{
			Repo: &repository.LanguagesConfigRepository{DBConn: dbConn},
		},
		ChildUsecase: childUsecase,
	}

	languagesConfig := engine.Group("/v1/user/languages-config", secureMiddleware.Secured())
	{
		languagesConfig.GET("/child/:child_id", languagesConfigController.GetChildLanguageConfig)
		languagesConfig.GET("/parent/:parent_id", languagesConfigController.GetParentLanguageConfig)
		languagesConfig.POST("/child", languagesConfigController.UploadLanguagesConfig)
		languagesConfig.POST("/parent", languagesConfigController.UploadLanguagesConfig)
	}
}
