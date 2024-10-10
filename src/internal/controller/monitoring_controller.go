package controller

import (
	"github.com/gin-gonic/gin"
	"sen-global-api/pkg/monitor"
)

type MonitoringController struct {
}

type responseGoogleAPIRequest struct {
	TotalRequestInitDevice      int `json:"total_request_init_device"`
	TotalRequestImportToDo      int `json:"total_request_import_to_do"`
	TotalRequestImportForm      int `json:"total_request_import_form"`
	TotalRequestGETScreenButton int `json:"total_request_get_screen_button"`
	TotalRequestGETTopButton    int `json:"total_request_get_top_button"`
}

// GetGoogleAPIMonitoring godoc
// @Summary Get Google API Monitoring
// @Description Get Google API Monitoring
// @Tags Monitoring
// @Accept  json
// @Produce  json
// @Success 200 {object} responseGoogleAPIRequest
// @Router /v1/admin/monitor/google-api [get]
func (c *MonitoringController) GetGoogleAPIMonitoring(context *gin.Context) {
	context.JSON(200, responseGoogleAPIRequest{
		TotalRequestInitDevice:      monitor.TotalRequestInitDevice,
		TotalRequestImportToDo:      monitor.TotalRequestImportToDo,
		TotalRequestImportForm:      monitor.TotalRequestImportForm,
		TotalRequestGETScreenButton: monitor.TotalRequestGETScreenButton,
		TotalRequestGETTopButton:    monitor.TotalRequestGETTopButton,
	})
}
