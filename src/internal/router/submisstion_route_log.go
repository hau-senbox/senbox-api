package router

import (
	"sen-global-api/config"
	"sen-global-api/internal/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupSubmissionLogRoutes(engine *gin.Engine, conn *gorm.DB, appConfig config.AppConfig) {
	group := engine.Group("api/v1")
	{
		submissionController := controller.NewSubmissionLogController(appConfig, conn)
		group.GET("/submissions/form", submissionController.GetSubmissionsFormLogs)
		group.GET("/submissions/form/submit", submissionController.GetSubmissionsFormLogsBySubmit)
	}
}
