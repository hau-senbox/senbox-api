package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
)

type ComponentController struct {
	GetComponentUseCase *usecase.GetComponentUseCase
}

func (receiver *ComponentController) GetAllComponentKey(context *gin.Context) {
	keys, err := receiver.GetComponentUseCase.GetAllComponentKey()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: keys,
	})
}
