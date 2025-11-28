package router

import (
	"sen-global-api/internal/controller"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupAppRoutes(r *gin.Engine, dbConn *gorm.DB) {
	// init repository + usecase
	appConfigRepo := &repository.AppConfigRepository{DBConn: dbConn}
	appConfigUC := &usecase.AppConfigUseCase{Repo: appConfigRepo}
	userBlockSettingRepo := &repository.UserBlockSettingRepository{DBConn: dbConn}
	userBlockSettingUC := &usecase.UserBlockSettingUsecase{Repo: userBlockSettingRepo}
	appConfigCtrl := &controller.AppConfigController{AppConfigUsecase: appConfigUC, UserBlockSettingUsecase: userBlockSettingUC}

	api := r.Group("/v1/app")
	{
		configs := api.Group("/configs")
		{
			configs.POST("", appConfigCtrl.Upload)
			configs.GET("", appConfigCtrl.GetAll)
		}
		api.POST("/on/need-update", appConfigCtrl.OnIsNeedToUpdate)
	}
}
