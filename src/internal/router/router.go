package router

import (
    "sen-global-api/config"
    "sen-global-api/pkg/sheet"

    firebase "firebase.google.com/go/v4"
    "github.com/gin-gonic/gin"
    "github.com/hashicorp/consul/api"
    "gorm.io/gorm"
)

func Route(engine *gin.Engine, dbConn *gorm.DB, userSpreadsheet *sheet.Spreadsheet, uploaderSpreadsheet *sheet.Spreadsheet, appConfig config.AppConfig, fcm *firebase.App, consulClient *api.Client) {
    setupAdminRoutes(engine, dbConn, appConfig, userSpreadsheet, uploaderSpreadsheet, fcm)
    setupDeviceRoutes(engine, dbConn, userSpreadsheet, appConfig, fcm)
    setupQuestionRoutes(engine, dbConn, appConfig)
    setupToDoRoutes(engine, dbConn, appConfig)
    setupSubmissionLogRoutes(engine, dbConn, appConfig)
    setupUserRoutes(engine, dbConn, appConfig, consulClient)
    setupOrganizationRoutes(engine, dbConn, appConfig)
    setupAppRoutes(engine, dbConn, appConfig)
    setupGatewayRoutes(engine, dbConn, appConfig, consulClient)
}