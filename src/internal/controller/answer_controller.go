package controller

import (
	"net/http"
	"strings"

	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type AnswerController struct {
	answerUseCase *usecase.AnswerUseCase
}

func NewAnswerController(answerUseCase *usecase.AnswerUseCase) *AnswerController {
	return &AnswerController{
		answerUseCase: answerUseCase,
	}
}

// POST /answers
func (ctrl *AnswerController) Create(c *gin.Context) {
	var answer entity.SAnswer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			},
		)
		return
	}
	if err := ctrl.answerUseCase.CreateAnswer(&answer); err != nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: answer,
	})
}

// // GET /answers/:id
// func (ctrl *AnswerController) GetByID(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, response.Failed("Invalid UUID"))
// 		return
// 	}

// 	answer, err := ctrl.UseCase.GetAnswerByID(id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, response.Failed("Answer not found"))
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.Success(answer))
// }

// // GET /answers/submission/:submission_id
// func (ctrl *AnswerController) GetBySubmissionID(c *gin.Context) {
// 	submissionID := c.Param("submission_id")
// 	answers, err := ctrl.UseCase.GetAnswersBySubmissionID(submissionID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, response.Failed(err.Error()))
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.Success(answers))
// }

// // PUT /answers/:id
// func (ctrl *AnswerController) Update(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, response.Failed("Invalid UUID"))
// 		return
// 	}

// 	var answer entity.SAnswer
// 	if err := c.ShouldBindJSON(&answer); err != nil {
// 		c.JSON(http.StatusBadRequest, response.Failed(err.Error()))
// 		return
// 	}
// 	answer.ID = id

// 	updated, err := ctrl.UseCase.UpdateAnswer(&answer)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, response.Failed(err.Error()))
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.Success(updated))
// }

// // DELETE /answers/:id
// func (ctrl *AnswerController) Delete(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, response.Failed("Invalid UUID"))
// 		return
// 	}

// 	if err := ctrl.UseCase.DeleteAnswer(id); err != nil {
// 		c.JSON(http.StatusInternalServerError, response.Failed(err.Error()))
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.Success("Deleted successfully"))
// }

// GET /answers/search?key=xxx&db=yyy
func (ctrl *AnswerController) GetByKeyAndDB(c *gin.Context) {

	var req request.GetAnswerByKeyAndDB
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	// Parse atr_value_string
	attr := helper.ParseAtrValueStringToStruct(req.AtrValueString)
	if attr.Key == nil || attr.DB == nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: "Missing key or db",
			},
		)
		return
	}

	res, err := ctrl.answerUseCase.GetAnswersByKeyAndDB(repository.GetSubmissionByConditionParam{
		Key:          attr.Key,
		DB:           attr.DB,
		DateDuration: attr.DateDuration,
		TimeSort:     attr.TimeSort,
	})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Code:  http.StatusInternalServerError,
				Error: err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (ctrl *AnswerController) GetTotalNrByKeyAndDb(c *gin.Context) {
	var req request.GetTotalNrByKeyAndDbRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
	}

	attr := helper.ParseAtrValueStringToStruct(req.AtrValueString)

	// check question key, question db NR
	if attr.Key != nil && !strings.Contains(*attr.Key, "NR") {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request condition: the question key must contain NR!",
		})
		return
	}

	if attr.DB != nil && !strings.Contains(*attr.DB, "NR") {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request condition: the question db must contain NR!",
		})
		return
	}

	res, err := ctrl.answerUseCase.GetTotalNrByKeyAndDb(repository.GetSubmissionByConditionParam{
		Key:          attr.Key,
		DB:           attr.DB,
		DateDuration: attr.DateDuration,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}
