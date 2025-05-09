package controller

import (
	"net/http"
	"sen-global-api/config"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/sheet"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ToDoController struct {
	getToDoListByQRCodeUseCase *usecase.GetToDoListByQRCodeUseCase
	findDeviceFromRequestCase  *usecase.FindDeviceFromRequestCase
	markToDoAsDoneUseCase      *usecase.MarkToDoAsDoneUseCase
	updateToDoTasksUseCase     *usecase.UpdateToDoTasksUseCase
	findTodoByIdUseCase        *usecase.FindTodoByIdUseCase
	getDeviceByIdUseCase       *usecase.GetDeviceByIdUseCase
}

type getTodoByQRCodeParams struct {
	QRCode string `form:"qr_code" binding:"required"`
}

type task struct {
	Index     int    `json:"index" binding:"required"`
	Name      string `json:"name" binding:"required"`
	DueDate   string `json:"due_date" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Selection string `json:"selection" binding:"required"`
	Selected  string `json:"selected" binding:"required"`
}

type toDoResponseData struct {
	Tasks []task         `json:"tasks" binding:"required"`
	Name  string         `json:"name" binding:"required"`
	Type  value.ToDoType `json:"type" binding:"required"`
}

type toDoResponse struct {
	Data toDoResponseData `json:"data" binding:"required"`
}

// GetToDoListByQRCode godoc
// @Summary Get ToDo List By QR Code
// @Description Get ToDo List By QR Code
// @Tags ToDo
// @Accept  json
// @Produce  json
// @Param qr_code query string true "QR Code"
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} toDoResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/todo [get]
func (c *ToDoController) GetToDoListByQRCode(context *gin.Context) {
	var params getTodoByQRCodeParams
	if err := context.ShouldBindQuery(&params); err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todoList, err := c.getToDoListByQRCodeUseCase.Execute(params.QRCode)
	if err != nil {
		context.JSON(500, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	tasks := make([]task, 0)
	for _, t := range todoList.Tasks.Data.Tasks {
		tasks = append(tasks, task{
			Index:     t.Index,
			Name:      t.Name,
			DueDate:   t.DueDate,
			Value:     t.Value,
			Selection: t.Selection,
			Selected:  t.Selected,
		})
	}

	context.JSON(200, toDoResponse{
		Data: toDoResponseData{
			Tasks: tasks,
			Name:  todoList.Name,
			Type:  todoList.Type,
		},
	})
}

type markToDoAsDoneRequest struct {
	QRCode    string `json:"qr_code" binding:"required"`
	TaskIndex int    `json:"task_index"`
	Select    string `json:"select"`
	DeviceID  string `json:"device_id" binding:"required"`
}

// MarkToDoAsDone godoc
// @Summary Mark ToDo As Done
// @Description Mark ToDo As Done
// @Tags ToDo
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param req body markToDoAsDoneRequest true "Mark ToDo As Done Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/todo [post]
func (c *ToDoController) MarkToDoAsDone(context *gin.Context) {
	var req markToDoAsDoneRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	device, err := c.getDeviceByIdUseCase.Get(req.DeviceID)
	if err != nil || device == nil {
		context.JSON(http.StatusBadGateway, response.FailedResponse{
			Code:  http.StatusBadGateway,
			Error: "Cannot find device",
		})
		return
	}

	_, err = c.findTodoByIdUseCase.Execute(req.QRCode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "ToDo not found",
		})
		return
	}

	err = c.markToDoAsDoneUseCase.Execute(*device, req.QRCode, req.TaskIndex, req.Select)
	if err != nil {
		context.JSON(500, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(200, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Mark ToDo As Done successfully",
	})
}

func NewToDoController(cfg config.AppConfig, dbConn *gorm.DB, reader *sheet.Reader, writer *sheet.Writer) *ToDoController {
	return &ToDoController{
		getToDoListByQRCodeUseCase: usecase.NewGetToDoListByQRCodeUseCase(cfg, dbConn, reader),
		findDeviceFromRequestCase:  usecase.NewFindDeviceFromRequestCase(cfg, dbConn),
		markToDoAsDoneUseCase:      usecase.NewMarkToDoAsDoneUseCase(cfg, dbConn, reader, writer),
		updateToDoTasksUseCase:     usecase.NewUpdateToDoTasksUseCase(cfg, dbConn, reader, writer),
		findTodoByIdUseCase:        usecase.NewFindTodoByIdUseCase(dbConn),
		getDeviceByIdUseCase:       usecase.NewGetDeviceByIdUseCase(dbConn),
	}
}

// Update ToDo's Tasks godoc
// @Summary Update ToDo's Tasks
// @Description Update ToDo's Taks
// @Tags ToDo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param req body request.UpdateToDoTasksRequest true "Update Todo's Tasks Params"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 400 {object} response.FailedResponse
// @Router /v1/todo/tasks [put]
func (c *ToDoController) UpdateToDoTasks(context *gin.Context) {
	var req request.UpdateToDoTasksRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	_, err := c.updateToDoTasksUseCase.UpdateTask(req)
	if err != nil {
		context.JSON(500, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "ToDo Tasks successfully updated",
	})
}

// Log Task godoc
// @Summary Log Task
// @Description Log Task
// @Tags ToDo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param req body request.LogTaskRequest true "Log Task Params"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 400 {object} response.FailedResponse
// @Router /v1/todo/task/log [put]
func (c *ToDoController) LogTask(context *gin.Context) {
	var req request.LogTaskRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(400, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})

		return
	}

	device, err := c.findDeviceFromRequestCase.FindDevice(context)
	if err != nil || device == nil {
		context.JSON(http.StatusBadGateway, response.FailedResponse{
			Code:  http.StatusBadGateway,
			Error: "Cannot find device",
		})
		return
	}

	err = c.markToDoAsDoneUseCase.LogTask(req, *device)
	if err != nil {
		context.JSON(500, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(200, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Log Task successfully updated",
	})
}
