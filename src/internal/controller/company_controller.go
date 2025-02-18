package controller

import (
	"net/http"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CompanyController struct {
	*usecase.GetCompanyUseCase
}

func (receiver CompanyController) GetCompanyById(context *gin.Context) {
	userId := context.Param("id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "id is required",
				},
			},
		)
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid id",
			},
		})
		return
	}

	company, err := receiver.GetCompanyUseCase.GetCompanyById(uint(id))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.CompanyResponse{
			ID:          company.ID,
			CompanyName: company.CompanyName,
			Address:     company.Address,
			Description: company.Description,
		},
	})
}
