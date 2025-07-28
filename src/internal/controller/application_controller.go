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
	apps, err := ctrl.StaffAppUsecase.GetAllStaffApplications()
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
