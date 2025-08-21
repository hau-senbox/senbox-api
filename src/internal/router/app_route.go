package router

import (
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupAppRoutes(r *gin.Engine, dbConn *gorm.DB, appCfg config.AppConfig) {
	// init repository + usecase
	appConfigRepo := &repository.AppConfigRepository{DBConn: dbConn}
	appConfigUC := &usecase.AppConfigUseCase{Repo: appConfigRepo}
	appConfigCtrl := &controller.AppConfigController{AppConfigUsecase: appConfigUC}

	api := r.Group("/v1/app")
	{
		configs := api.Group("/configs")
		{
			configs.POST("", appConfigCtrl.Upload)
			configs.GET("", appConfigCtrl.GetAll)
		}
	}
}
