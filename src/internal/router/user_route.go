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

func setupUserRoutes(engine *gin.Engine, dbConn *gorm.DB, config config.AppConfig) {
	sessionRepository := repository.SessionRepository{
		OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		AuthorizeEncryptKey:    config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}
	secureMiddleware := middleware.SecuredMiddleware{SessionRepository: sessionRepository}

	userEntityController := &controller.UserEntityController{
		GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
		},
		CreateUserEntityUseCase: &usecase.CreateUserEntityUseCase{
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
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
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
			MenuRepository:         &repository.MenuRepository{DBConn: dbConn},
			UserEntityRepository:   &repository.UserEntityRepository{DBConn: dbConn},
			OrganizationRepository: &repository.OrganizationRepository{DBConn: dbConn},
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
	}

	componentController := &controller.ComponentController{
		GetComponentUseCase: &usecase.GetComponentUseCase{
			ComponentRepository: &repository.ComponentRepository{DBConn: dbConn},
		},
	}

	userAccess := engine.Group("v1/")
	{
		loginController := &controller.LoginController{DBConn: dbConn,
			AuthorizeUseCase: usecase.AuthorizeUseCase{
				UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
				DeviceRepository:     &repository.DeviceRepository{DBConn: dbConn},
				SessionRepository:    sessionRepository,
			},
		}
		userAccess.POST("/login", loginController.UserLogin)
	}

	user := engine.Group("v1/user")
	{
		user.GET("/current-user", secureMiddleware.Secured(), userEntityController.GetCurrentUser)
		user.GET("/all", secureMiddleware.Secured(), userEntityController.GetAllUserEntity)
		user.GET("/:id", secureMiddleware.Secured(), userEntityController.GetUserEntityById)
		user.GET("/name/:username", secureMiddleware.Secured(), userEntityController.GetUserEntityByName)
		user.GET("/:id/children", secureMiddleware.Secured(), userEntityController.GetChildrenOfGuardian)

		user.POST("/init", userEntityController.CreateUserEntity)
		user.POST("/update", secureMiddleware.Secured(), userEntityController.UpdateUserEntity)
		user.POST("/block/:id", secureMiddleware.Secured(), userEntityController.BlockUser)
		user.POST("/role/update", secureMiddleware.Secured(), userEntityController.UpdateUserRole)

		user.GET("/org/:organization_id/:user_id", secureMiddleware.Secured(), userEntityController.GetUserOrgInfo)
		user.GET("/org/:organization_id/manager", secureMiddleware.Secured(), userEntityController.GetAllOrgManagerInfo)
		user.POST("/org/update", secureMiddleware.Secured(), userEntityController.UpdateUserOrgInfo)

		user.GET("/:id/func", secureMiddleware.Secured(), userEntityController.GetAllUserAuthorize)
		user.POST("/func", secureMiddleware.Secured(), userEntityController.UpdateUserAuthorize)
		user.DELETE("/func", secureMiddleware.Secured(), userEntityController.DeleteUserAuthorize)

		user.GET("/pre-register/", secureMiddleware.Secured(), userEntityController.GetAllPreRegisterUser)
		user.POST("/pre-register/", secureMiddleware.Secured(), userEntityController.CreatePreRegister)
	}

	teacherApplication := engine.Group("/v1/user/teacher/application")
	{
		teacherApplication.GET("/", secureMiddleware.Secured(), userEntityController.GetAllTeacherFormApplication)
		teacherApplication.GET("/:id", secureMiddleware.Secured(), userEntityController.GetTeacherFormApplicationByID)

		teacherApplication.POST("/", secureMiddleware.Secured(), userEntityController.CreateTeacherFormApplication)
		teacherApplication.POST("/:id/approve", secureMiddleware.Secured(), userEntityController.ApproveTeacherFormApplication)
		teacherApplication.POST("/:id/block", secureMiddleware.Secured(), userEntityController.BlockTeacherFormApplication)
	}

	staffApplication := engine.Group("/v1/user/staff/application")
	{
		staffApplication.GET("/", secureMiddleware.Secured(), userEntityController.GetAllStaffFormApplication)
		staffApplication.GET("/:id", secureMiddleware.Secured(), userEntityController.GetStaffFormApplicationByID)

		staffApplication.POST("/", secureMiddleware.Secured(), userEntityController.CreateStaffFormApplication)
		staffApplication.POST("/:id/approve", secureMiddleware.Secured(), userEntityController.ApproveStaffFormApplication)
		staffApplication.POST("/:id/block", secureMiddleware.Secured(), userEntityController.BlockStaffFormApplication)
	}

	studentApplication := engine.Group("/v1/user/student/application")
	{
		studentApplication.GET("/", secureMiddleware.Secured(), userEntityController.GetAllStudentFormApplication)
		studentApplication.GET("/:id", secureMiddleware.Secured(), userEntityController.GetStudentFormApplicationByID)

		studentApplication.POST("/", secureMiddleware.Secured(), userEntityController.CreateStudentFormApplication)
		studentApplication.POST("/:id/approve", secureMiddleware.Secured(), userEntityController.ApproveStudentFormApplication)
		studentApplication.POST("/:id/block", secureMiddleware.Secured(), userEntityController.BlockStudentFormApplication)
	}

	userRole := engine.Group("v1/user-role", secureMiddleware.Secured())
	{
		userRole.GET("/all", userRoleController.GetAllRole)
		userRole.GET("/:id", userRoleController.GetRoleById)
		userRole.GET("/name/:role_name", userRoleController.GetRoleByName)

		userRole.POST("/init", userRoleController.CreateRole)
		userRole.POST("/", userRoleController.UpdateRole)

		userRole.DELETE("/:id", userRoleController.DeleteRole)
	}

	functionClaim := engine.Group("v1/function-claim", secureMiddleware.Secured())
	{
		functionClaim.GET("/all", functionClaimController.GetAllFunctionClaim)
		functionClaim.GET("/:id", functionClaimController.GetFunctionClaimById)
		functionClaim.GET("/name/:function_name", functionClaimController.GetFunctionClaimByName)

		functionClaim.POST("/init", functionClaimController.CreateFunctionClaim)
		functionClaim.POST("/", functionClaimController.UpdateFunctionClaim)

		functionClaim.DELETE("/:id", functionClaimController.DeleteFunctionClaim)
	}

	functionClaimPermission := engine.Group("v1/function-claim-permission", secureMiddleware.Secured())
	{
		functionClaimPermission.GET("/function/:function_claim_id/all", functionClaimPermissionController.GetAllFunctionClaimPermission)
		functionClaimPermission.GET("/:id", functionClaimPermissionController.GetFunctionClaimPermissionById)
		functionClaimPermission.GET("/name/:permission_name", functionClaimPermissionController.GetFunctionClaimPermissionByName)

		functionClaimPermission.POST("/init", functionClaimPermissionController.CreateFunctionClaimPermission)
		functionClaimPermission.POST("/", functionClaimPermissionController.UpdateRoleClaimPermission)

		functionClaimPermission.DELETE("/:id", functionClaimPermissionController.DeleteRoleClaimPermission)
	}

	userMenu := engine.Group("v1/user-menu", secureMiddleware.Secured())
	{
		userMenu.GET("/super-admin", menuController.GetSuperAdminMenu)
		userMenu.GET("/org/:id", menuController.GetOrgMenu)
		userMenu.GET("/student/:id", menuController.GetStudentMenu)
		userMenu.GET("/teacher/:id", menuController.GetTeacherMenu)
		userMenu.GET("/user/:id", menuController.GetUserMenu)

		userMenu.POST("/super-admin", secureMiddleware.ValidateSuperAdminRole(), menuController.UploadSuperAdminMenu)
		userMenu.POST("/org", menuController.UploadOrgMenu)
		userMenu.POST("/user", menuController.UploadUserMenu)
	}

	component := engine.Group("v1/component", secureMiddleware.Secured())
	{
		component.GET("/keys", componentController.GetAllComponentKey)
	}
}
