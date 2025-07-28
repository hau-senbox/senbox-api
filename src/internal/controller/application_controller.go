package controller

import (
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type ApplicationController struct {
	StaffAppUsecase *usecase.StaffApplicationUseCase
}

func NewApplicationController(staffAppUsecase *usecase.StaffApplicationUseCase) *ApplicationController {
	return &ApplicationController{
		StaffAppUsecase: staffAppUsecase,
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

func (ctrl *ApplicationController) ApproveStaffApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StaffAppUsecase.ApproveStaffApplication(applicationID)
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

	err := ctrl.StaffAppUsecase.BlockStaffApplication(applicationID)
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
