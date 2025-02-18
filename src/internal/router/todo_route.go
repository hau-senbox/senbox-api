package router

import (
	"context"
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/pkg/sheet"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func setupToDoRoutes(engine *gin.Engine, conn *gorm.DB, appConfig config.AppConfig) {
	v1 := engine.Group("/v1")
	{
		ctx := context.Background()
		spreadSheet, err := sheet.NewUserSpreadsheet(appConfig, ctx)
		if err != nil {
			log.Fatal(err)
		}
		todoController := controller.NewToDoController(appConfig, conn, spreadSheet.Reader, spreadSheet.Writer)
		v1.GET("/todo", todoController.GetToDoListByQRCode)
		v1.POST("/todo", todoController.MarkToDoAsDone)

		v1.PUT("/todo/tasks", todoController.UpdateToDoTasks)

		v1.POST("/todo/task/log", todoController.LogTask)
	}
}
