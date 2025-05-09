package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DBConn *gorm.DB = nil

func GetCodeCounterList(context *gin.Context) {
	var rq request.GetCodeCountingsRequest

	err := context.ShouldBindQuery(&rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid request",
		})
		return
	}

	r, err := usecase.GetCodeCountings(DBConn, rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.GetCodeCountingsResponse{
			Codes:  r.Codes,
			Paging: r.Paging,
		},
	})
}

func UpdateCodeCounter(context *gin.Context) {
	var rq request.UpdateCodeCountingRequest
	err := context.ShouldBindJSON(&rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err = usecase.UpdateCodeCounting(DBConn, rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Message: "success",
	})
}
