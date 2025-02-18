package controller

import (
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SyncController struct {
	SynchronizeUseCase *usecase.SynchronizeUseCase
	GetSettingsUseCase *usecase.GetSettingsUseCase
}

// Start sync operation
// @Summary Start sync operation
// @Description Start sync operation
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} response.UserListResponse
// @Failure 400 {object} response.FailedResponse
// @Router /v1/admin/sync/start [post]
func (receiver SyncController) StartSync(context *gin.Context) {
	appSetting, _ := receiver.GetSettingsUseCase.GetSettings()

	receiver.SynchronizeUseCase.StartSync(*usecase.TheTimeMachine, appSetting)

	log.Debug("Restart sync operation")

	context.JSON(200, response.SucceedResponse{
		Data: "Time machine has been started",
	})
}

// Stop sync operation
// @Summary Stop sync operation
// @Description Stop sync operation
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} response.UserListResponse
// @Failure 400 {object} response.FailedResponse
// @Router /v1/admin/sync/stop [post]
func (receiver SyncController) StopSync(context *gin.Context) {
	receiver.SynchronizeUseCase.StopSync(*usecase.TheTimeMachine)

	log.Debug("Stop sync operation")

	context.JSON(200, response.SucceedResponse{
		Data: "Time machine has been stopped",
	})
}
