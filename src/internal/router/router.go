package router

import (
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sen-global-api/config"
	"sen-global-api/pkg/sheet"
)

func Route(engine *gin.Engine, dbConn *gorm.DB, userSpreadsheet *sheet.Spreadsheet, uploaderSpreadsheet *sheet.Spreadsheet, appConfig config.AppConfig, fcm *firebase.App) {
	setupAdminRoutes(engine, dbConn, appConfig, userSpreadsheet, uploaderSpreadsheet, fcm)
	setupDeviceRoutes(engine, dbConn, userSpreadsheet, uploaderSpreadsheet, appConfig, fcm)
	setupQuestionRoutes(engine, dbConn, appConfig, userSpreadsheet, uploaderSpreadsheet)
	setupToDoRoutes(engine, dbConn, appConfig, userSpreadsheet, uploaderSpreadsheet)
}
