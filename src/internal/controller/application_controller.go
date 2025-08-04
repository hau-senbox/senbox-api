package controller

import (
	"net/http"
	"sen-global-api/helper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"strings"

	"github.com/gin-gonic/gin"
)

type ApplicationController struct {
	StaffAppUsecase   *usecase.StaffApplicationUseCase
	StudentAppUsecase *usecase.StudentApplicationUseCase
	TeacherAppUsecase *usecase.TeacherApplicationUseCase
	SyncDataUsecase   *usecase.SyncDataUsecase
}

func NewApplicationController(
	staffAppUsecase *usecase.StaffApplicationUseCase,
	studentAppUsecase *usecase.StudentApplicationUseCase,
	teacherAppUsecase *usecase.TeacherApplicationUseCase) *ApplicationController {
	return &ApplicationController{
		StaffAppUsecase:   staffAppUsecase,
		StudentAppUsecase: studentAppUsecase,
		TeacherAppUsecase: teacherAppUsecase,
	}
}

// GetAllStaffApplications retrieves all staff applications
func (ctrl *ApplicationController) GetAllStaffApplications(ctx *gin.Context) {
	apps, err := ctrl.StaffAppUsecase.GetAllStaffApplications(ctx)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: apps,
	})
}

// GetAllStudentApplications retrieves all staff applications
func (ctrl *ApplicationController) GetAllStudentApplications(ctx *gin.Context) {
	apps, err := ctrl.StudentAppUsecase.GetAllStudentApplications(ctx)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: apps,
	})
}

func (ctrl *ApplicationController) GetAllTeacherApplications(ctx *gin.Context) {
	apps, err := ctrl.TeacherAppUsecase.GetAllTeacherApplications(ctx)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: apps,
	})
}

func (ctrl *ApplicationController) GetDetailStudentApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	app, err := ctrl.StudentAppUsecase.GetDetailStudentApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: app,
	})
}

func (ctrl *ApplicationController) GetDetailTeacherApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	app, err := ctrl.TeacherAppUsecase.GetDetailTeacherApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: app,
	})
}

func (ctrl *ApplicationController) GetDetailStaffApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	app, err := ctrl.StaffAppUsecase.GetDetailStaffApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: app,
	})
}

func (ctrl *ApplicationController) ApproveStaffApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StaffAppUsecase.ApproveStaffApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application approved successfully",
	})
}

func (ctrl *ApplicationController) BlockStaffApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StaffAppUsecase.BlockStaffApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application blocked successfully",
	})
}

func (ctrl *ApplicationController) ApproveStudentApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StudentAppUsecase.ApproveStudentApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application approved successfully",
	})
}

func (ctrl *ApplicationController) BlockStudentApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StudentAppUsecase.BlockStudentApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application blocked successfully",
	})
}

func (ctrl *ApplicationController) ApproveTeacherApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.TeacherAppUsecase.ApproveTeacherApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application approved successfully",
	})
}

func (ctrl *ApplicationController) BlockTeacherApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.TeacherAppUsecase.BlockTeacherApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application blocked successfully",
	})
}

func (ctrl *ApplicationController) SyncDataDemoV3(ctx *gin.Context) {
	var req request.SyncDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, response.FailedResponse{
			Code:    400,
			Message: "invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Tách chuỗi form_qrs thành slice FormNotes
	if req.FormNotesStr != "" {
		// Tách theo dấu phẩy và loại bỏ khoảng trắng thừa
		splitted := strings.Split(req.FormNotesStr, ",")
		req.FormNotes = make([]string, 0, len(splitted))
		for _, s := range splitted {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				req.FormNotes = append(req.FormNotes, trimmed)
			}
		}
	}

	lastSubmitedTimeStr, err := ctrl.SyncDataUsecase.ExcuteCreateAndSyncFormAnswer(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Error:   err.Error(),
		})
		return
	}

	// Parse lại chuỗi thời gian gốc
	t, err := helper.ParseFlexibleTime(lastSubmitedTimeStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Invalid SubmittedAt format",
			Error:   err.Error(),
		})
		return
	}

	lastSubmitTimeISO := t.Format("2006-01-02T15:04:05.000Z")

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Waiting sync data",
		Data: map[string]interface{}{
			"last_submit_time": lastSubmitTimeISO,
		},
	})
}

func (ctrl *ApplicationController) CheckStatusSyncQueue(ctx *gin.Context) {
	hasPending, err := ctrl.SyncDataUsecase.HasPendingSyncQueue()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check sync queue",
			Error:   err.Error(),
		})
		return
	}

	if !hasPending {
		ctx.JSON(http.StatusConflict, response.FailedResponse{
			Code:    http.StatusConflict,
			Message: "A sync is already in progress",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "No sync in progress. Ready to start sync.",
		Data: map[string]interface{}{
			"status": value.SyncQueueStatusDone,
		},
	})
}

func (ctrl *ApplicationController) GetAllSycnQueue(ctx *gin.Context) {

	queues, err := ctrl.SyncDataUsecase.GetAllSyncQueue()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch sync queues",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Get all sync success",
		Data:    queues,
	})
}
