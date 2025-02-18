package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
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
	d, err := usecase.GetLocationListDao(c, req)
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
}
