package controller

import (
	"errors"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlockSettingController struct {
	UserBlockUsecase    *usecase.UserBlockSettingUsecase
	StudentBlockUsecase *usecase.StudentBlockSettingUsecase
}

// GET /user-block-setting?user_id=xxx
func (ctl *BlockSettingController) GetByUserID(c *gin.Context) {
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

	result, err := ctl.UserBlockUsecase.GetByUserID(userID)
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

func (ctl *BlockSettingController) GetByUserID4App(c *gin.Context) {
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

	result, err := ctl.UserBlockUsecase.GetByUserID4App(userID)
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
func (ctl *BlockSettingController) UpsertUserBlockSetting(c *gin.Context) {
	var req request.UserBlockSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if err := ctl.UserBlockUsecase.Upsert(req); err != nil {
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

// GET /student-block-setting/:student_id
func (ctl *BlockSettingController) GetByStudentID(c *gin.Context) {
	studentID := c.Param("student_id")

	if studentID == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "student_id is required",
		})
		return
	}

	if _, err := uuid.Parse(studentID); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid student_id format (must be UUID)",
		})
		return
	}

	result, err := ctl.StudentBlockUsecase.GetByStudentID(studentID)
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

// POST /student-block-setting
func (ctl *BlockSettingController) UpsertStudentBlockSetting(c *gin.Context) {
	var req request.StudentBlockSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if err := ctl.StudentBlockUsecase.Upsert(req); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "student block setting upserted successfully",
	})
}

func (ctl *BlockSettingController) OffIsNeedToUpdateByUser(c *gin.Context) {
	userID := c.GetString("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "user_id is required",
		})
		return
	}

	if err := ctl.UserBlockUsecase.OffIsNeedToUpdateByUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "off is need to update by user successfully",
	})
}

func (ctl *BlockSettingController) MigrateFirestore(c *gin.Context) {
	if err := ctl.UserBlockUsecase.MigrateFirestore(); err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "migrate firestore successfully",
	})
}
