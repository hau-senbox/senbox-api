package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
)

var DBConn *gorm.DB = nil

func GetCodeCounterList(context *gin.Context) {
	var rq request.GetCodeCountingsRequest

	err := context.ShouldBindQuery(&rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid request",
			},
		})
		return
	}

	r, err := usecase.GetCodeCountings(DBConn, rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
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
		context.JSON(http.StatusBadRequest, response.SucceedResponse{
			Data: err.Error(),
		})
		return
	}

	err = usecase.UpdateCodeCounting(DBConn, rq)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.SucceedResponse{
			Data: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: "success",
	})
}
