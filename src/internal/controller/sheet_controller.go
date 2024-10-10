package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/usecase"
)

// Location list from sheet
func LocationGetListController(c *gin.Context) {
	var req request.LocationSheetRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"code":    http.StatusBadRequest,
			},
		)
		return
	}
	err, d := usecase.GetLocationListDao(c, req)
	if err != nil {
		c.JSON(
			http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"code":    http.StatusBadRequest,
			},
		)
		return
	}

	c.JSON(http.StatusOK, d)

	return
}
