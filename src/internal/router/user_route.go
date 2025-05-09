package router

import (
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupUserRoutes(engine *gin.Engine, dbConn *gorm.DB, config config.AppConfig) {
	sessionRepository := repository.SessionRepository{
		AuthorizeEncryptKey: config.AuthorizeEncryptKey,

		TokenExpireTimeInHour: time.Duration(config.TokenExpireDurationInHour),
	}

	userEntityController := &controller.UserEntityController{
		GetUserEntityUseCase: &usecase.GetUserEntityUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
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
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		UpdateUserAuthorizeUseCase: &usecase.UpdateUserAuthorizeUseCase{
			UserEntityRepository: &repository.UserEntityRepository{DBConn: dbConn},
		},
		DeleteUserAuthorizeUseCase: &usecase.DeleteUserAuthorizeUseCase{
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
		user.GET("/all", userEntityController.GetAllUserEntity)
		user.GET("/:id", userEntityController.GetUserEntityById)
		user.GET("/name/:username", userEntityController.GetUserEntityByName)
		user.GET("/:id/children", userEntityController.GetChildrenOfGuardian)

		user.POST("/init", userEntityController.CreateUserEntity)
		user.POST("/update", userEntityController.UpdateUserEntity)
		user.POST("/role/update", userEntityController.UpdateUserRole)

		user.GET("/org/:organization_id/:user_id", userEntityController.GetUserOrgInfo)
		user.GET("/org/:organization_id/manager", userEntityController.GetAllOrgManagerInfo)
		user.POST("/org/update", userEntityController.UpdateUserOrgInfo)

		user.GET("/:id/func", userEntityController.GetAllUserAuthorize)
		user.POST("/func", userEntityController.UpdateUserAuthorize)
		user.DELETE("/func", userEntityController.DeleteUserAuthorize)
	}

	userRole := engine.Group("v1/user-role")
	{
		userRole.GET("/all", userRoleController.GetAllRole)
		userRole.GET("/:id", userRoleController.GetRoleById)
		userRole.GET("/name/:role_name", userRoleController.GetRoleByName)

		userRole.POST("/init", userRoleController.CreateRole)
		userRole.POST("/", userRoleController.UpdateRole)

		userRole.DELETE("/:id", userRoleController.DeleteRole)
	}

	functionClaim := engine.Group("v1/function-claim")
	{
		functionClaim.GET("/all", functionClaimController.GetAllFunctionClaim)
		functionClaim.GET("/:id", functionClaimController.GetFunctionClaimById)
		functionClaim.GET("/name/:function_name", functionClaimController.GetFunctionClaimByName)

		functionClaim.POST("/init", functionClaimController.CreateFunctionClaim)
		functionClaim.POST("/", functionClaimController.UpdateFunctionClaim)

		functionClaim.DELETE("/:id", functionClaimController.DeleteFunctionClaim)
	}

	functionClaimPermission := engine.Group("v1/function-claim-permission")
	{
		functionClaimPermission.GET("/function/:function_claim_id/all", functionClaimPermissionController.GetAllFunctionClaimPermission)
		functionClaimPermission.GET("/:id", functionClaimPermissionController.GetFunctionClaimPermissionById)
		functionClaimPermission.GET("/name/:permission_name", functionClaimPermissionController.GetFunctionClaimPermissionByName)

		functionClaimPermission.POST("/init", functionClaimPermissionController.CreateFunctionClaimPermission)
		functionClaimPermission.POST("/", functionClaimPermissionController.UpdateRoleClaimPermission)

		functionClaimPermission.DELETE("/:id", functionClaimPermissionController.DeleteRoleClaimPermission)
	}
}
