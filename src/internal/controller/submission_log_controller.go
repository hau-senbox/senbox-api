package controller

import (
	"net/http"
	"sen-global-api/config"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubmissionLogController struct {
	SubmissionLogUseCase *usecase.SubmissionLogUseCase
}

func NewSubmissionLogController(appConfig config.AppConfig, conn *gorm.DB) *SubmissionLogController {
	return &SubmissionLogController{
		SubmissionLogUseCase: usecase.NewSubmissionLogUseCase(appConfig, conn),
	}
}

func (c *SubmissionLogController) GetSubmissionsFormLogs(ctx *gin.Context) {

	start := ctx.Query("start")
	end := ctx.Query("end")
	qrCode := ctx.Query("qr_code")
	customID := ctx.Query("custom_id")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	result, err := c.SubmissionLogUseCase.GetSubmissionsFormLogs(ctx, start, end, qrCode, customID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *SubmissionLogController) GetSubmissionsFormLogsBySubmit(ctx *gin.Context) {

	start := ctx.Query("start")
	end := ctx.Query("end")
	qrCode := ctx.Query("qr_code")
	customID := ctx.Query("custom_id")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	result, err := c.SubmissionLogUseCase.GetSubmissionsFormLogsBySubmit(ctx, start, end, qrCode, customID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)

}
