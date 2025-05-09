package controller

import (
	"net/http"
	"sen-global-api/config"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/pkg/job"
	"sen-global-api/pkg/sheet"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ImportToDoController struct {
	*usecase.ImportToDoListUseCase
}

// Import ToDos godoc
// @Summary Import ToDos
// @Description Import ToDos
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.ImportFormRequest true "Import Form Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/todo/import [post]
func (c *ImportToDoController) ImportTodos(context *gin.Context) {
	var request request.ImportFormRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := c.ImportToDoList(request)
	if err != nil {
		context.JSON(500, response.SucceedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Import ToDos failed",
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Import ToDos successfully",
	})
}

// Import Partially ToDos godoc
// @Summary Import Partially ToDos
// @Description Import Partially ToDos
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.ImportPartiallyTodoRequest true "Import Partially ToDo Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/todo/import/partially [post]
func (c *ImportToDoController) ImportPartiallyTodos(context *gin.Context) {
	var req request.ImportPartiallyTodoRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := c.ImportPartiallyToDos(req.SpreadsheetURL, req.TabName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "TODO imported successfully",
	})
}

func NewImportToDoListController(cfg config.AppConfig, dbConn *gorm.DB, reader *sheet.Reader, writer *sheet.Writer, machine *job.TimeMachine) *ImportToDoController {
	return &ImportToDoController{
		ImportToDoListUseCase: usecase.NewImportToDoListUseCase(cfg, dbConn, reader, writer, machine),
	}
}
