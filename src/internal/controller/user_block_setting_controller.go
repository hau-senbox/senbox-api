package controller

import (
	"errors"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBlockSettingController struct {
	Usecase *usecase.UserBlockSettingUsecase
}

// GET /user-block-setting?user_id=xxx
func (ctl *UserBlockSettingController) GetByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "user_id is required",
		})
		return
	}

	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid user_id format (must be UUID)",
		})
		return
	}

	result, err := ctl.Usecase.GetByUserID(userID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, response.SucceedResponse{
				Code:    http.StatusOK,
				Message: "not found",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, response.SucceedResponse{
			Code:    http.StatusNotFound,
			Message: "not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    result,
	})
}

// POST /user-block-setting
func (ctl *UserBlockSettingController) Upsert(c *gin.Context) {
	var req request.UserBlockSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if err := ctl.Usecase.Upsert(req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user block setting upserted successfully",
	})
}

// DELETE /user-block-setting/:id
func (ctl *UserBlockSettingController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	intID, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid id",
		})
		return
	}

	if err := ctl.Usecase.Delete(intID); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "user block setting deleted successfully",
	})
}
