package router

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sen-global-api/config"
	"sen-global-api/internal/controller"
	"sen-global-api/pkg/sheet"
)

func setupToDoRoutes(engine *gin.Engine, conn *gorm.DB, appConfig config.AppConfig, userSpreadsheet *sheet.Spreadsheet, uploaderSpreadsheet *sheet.Spreadsheet) {
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
