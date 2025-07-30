package controller

import (
	"sen-global-api/helper"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type ApplicationController struct {
	StaffAppUsecase   *usecase.StaffApplicationUseCase
	StudentAppUsecase *usecase.StudentApplicationUseCase
	TeacherAppUsecase *usecase.TeacherApplicationUseCase
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

func (ctrl *ApplicationController) SyncDataDemo(ctx *gin.Context) {
	// Cấu hình dữ liệu và thông tin sheet
	spreadsheetID := "1YGe4AWf1qt8f5K5iJ6OGcZDGLGnGWWE0JmDZIr0jrn8"
	sheetName := "Sheet1"
	startCell := "A1"
	credentialsPath := "credentials/uploader_service_account.json"

	values := [][]interface{}{
		{"ID", "Name", "Score"},
		{"1", "Alice", 90},
		{"2", "Bob", 85},
		{"3", "Charlie", 92},
	}

	err := helper.WriteDataToSheet(spreadsheetID, sheetName, startCell, values, credentialsPath)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to sync data", "error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Data synced successfully to Google Sheet"})
}
