package controller

import (
	"fmt"
	"net/http"
	"sen-global-api/helper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/sheets/v4"
)

type ApplicationController struct {
	StaffAppUsecase   *usecase.StaffApplicationUseCase
	StudentAppUsecase *usecase.StudentApplicationUseCase
	TeacherAppUsecase *usecase.TeacherApplicationUseCase
	SyncDataUsecase   *usecase.SyncDataUsecae
}

func NewApplicationController(
	staffAppUsecase *usecase.StaffApplicationUseCase,
	studentAppUsecase *usecase.StudentApplicationUseCase,
	teacherAppUsecase *usecase.TeacherApplicationUseCase) *ApplicationController {
	return &ApplicationController{
		StaffAppUsecase:   staffAppUsecase,
		StudentAppUsecase: studentAppUsecase,
		TeacherAppUsecase: teacherAppUsecase,
	}
}

// GetAllStaffApplications retrieves all staff applications
func (ctrl *ApplicationController) GetAllStaffApplications(ctx *gin.Context) {
	apps, err := ctrl.StaffAppUsecase.GetAllStaffApplications(ctx)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: apps,
	})
}

// GetAllStudentApplications retrieves all staff applications
func (ctrl *ApplicationController) GetAllStudentApplications(ctx *gin.Context) {
	apps, err := ctrl.StudentAppUsecase.GetAllStudentApplications(ctx)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: apps,
	})
}

func (ctrl *ApplicationController) GetAllTeacherApplications(ctx *gin.Context) {
	apps, err := ctrl.TeacherAppUsecase.GetAllTeacherApplications(ctx)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: apps,
	})
}

func (ctrl *ApplicationController) GetDetailStudentApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	app, err := ctrl.StudentAppUsecase.GetDetailStudentApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: app,
	})
}

func (ctrl *ApplicationController) GetDetailTeacherApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	app, err := ctrl.TeacherAppUsecase.GetDetailTeacherApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: app,
	})
}

func (ctrl *ApplicationController) GetDetailStaffApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	app, err := ctrl.StaffAppUsecase.GetDetailStaffApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: app,
	})
}

func (ctrl *ApplicationController) ApproveStaffApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StaffAppUsecase.ApproveStaffApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application approved successfully",
	})
}

func (ctrl *ApplicationController) BlockStaffApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StaffAppUsecase.BlockStaffApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application blocked successfully",
	})
}

func (ctrl *ApplicationController) ApproveStudentApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StudentAppUsecase.ApproveStudentApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application approved successfully",
	})
}

func (ctrl *ApplicationController) BlockStudentApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.StudentAppUsecase.BlockStudentApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application blocked successfully",
	})
}

func (ctrl *ApplicationController) ApproveTeacherApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.TeacherAppUsecase.ApproveTeacherApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application approved successfully",
	})
}

func (ctrl *ApplicationController) BlockTeacherApplication(ctx *gin.Context) {
	applicationID := ctx.Param("id")
	if applicationID == "" {
		ctx.JSON(400, response.FailedResponse{
			Code:  400,
			Error: "Application ID is required",
		})
		return
	}

	err := ctrl.TeacherAppUsecase.BlockTeacherApplication(ctx, applicationID)
	if err != nil {
		ctx.JSON(500, response.FailedResponse{
			Code:  500,
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(200, response.SucceedResponse{
		Code: 200,
		Data: "Application blocked successfully",
	})
}

func (ctrl *ApplicationController) SyncDataDemo(ctx *gin.Context) {
	// 0. Parse dữ liệu từ request
	var submissions []struct {
		SubmittedAt string `json:"SubmittedAt"`
		StudentID   string `json:"StudentID"`
		UserID      string `json:"UserID"`
		FormCode    string `json:"FormCode"`
		FormName    string `json:"FormName"`
		Question    string `json:"Question"`
		Answer      string `json:"Answer"`
	}
	if err := ctx.ShouldBindJSON(&submissions); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request", "detail": err.Error()})
		return
	}

	spreadsheetID := "1YGe4AWf1qt8f5K5iJ6OGcZDGLGnGWWE0JmDZIr0jrn8"
	sheetName := "Sheet1"
	credentialsPath := "credentials/uploader_service_account.json"

	// 1. Lấy dữ liệu hiện tại của sheet
	srv, _ := helper.GetSheetsService(credentialsPath)
	time.Sleep(3 * time.Second)
	readRange := fmt.Sprintf("%s!A1:Z1", sheetName)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to read sheet header", "detail": err.Error()})
		return
	}

	// 2. Xử lý header
	currentHeader := []string{}
	if len(resp.Values) > 0 {
		for _, cell := range resp.Values[0] {
			currentHeader = append(currentHeader, fmt.Sprintf("%v", cell))
		}
	}

	// 3. Chuẩn bị header chuẩn (có thể mới hơn)
	fixedColumns := []string{"Submitted At", "Student ID", "User ID", "Form Code", "Form Name"}
	questionSet := map[string]bool{}
	for _, s := range submissions {
		questionSet[s.Question] = true
	}
	for _, col := range fixedColumns {
		questionSet[col] = false // đảm bảo không bị thêm nhầm
	}

	// 4. Tạo header mới nếu cần
	newHeader := append([]string{}, fixedColumns...)
	for q := range questionSet {
		if !contains(newHeader, q) {
			newHeader = append(newHeader, q)
		}
	}

	// Nếu header cũ khác header mới → ghi lại header mới
	if !equal(currentHeader, newHeader) {
		_, err = srv.Spreadsheets.Values.Update(spreadsheetID, fmt.Sprintf("%s!A1", sheetName), &sheets.ValueRange{
			Values: [][]interface{}{stringSliceToInterfaceSlice(newHeader)},
		}).ValueInputOption("RAW").Do()
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to update header", "detail": err.Error()})
			return
		}
	}

	// 5. Chuẩn bị dữ liệu mới
	values := [][]interface{}{}
	for _, s := range submissions {
		row := make([]interface{}, len(newHeader))
		for i, col := range newHeader {
			switch col {
			case "Submitted At":
				row[i] = s.SubmittedAt
			case "Student ID":
				row[i] = s.StudentID
			case "User ID":
				row[i] = s.UserID
			case "Form Code":
				row[i] = s.FormCode
			case "Form Name":
				row[i] = s.FormName
			case s.Question:
				row[i] = s.Answer
			default:
				row[i] = ""
			}
		}
		values = append(values, row)
	}

	// 6. Ghi dữ liệu
	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	time.Sleep(3 * time.Second)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to write data", "detail": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Sync success!"})
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func stringSliceToInterfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

func (ctrl *ApplicationController) SyncDataDemoV2(ctx *gin.Context) {
	var submissions []struct {
		SubmittedAt string `json:"SubmittedAt"`
		StudentID   string `json:"StudentID"`
		UserID      string `json:"UserID"`
		FormCode    string `json:"FormCode"`
		FormName    string `json:"FormName"`
		Question    string `json:"Question"`
		Answer      string `json:"Answer"`
	}
	if err := ctx.ShouldBindJSON(&submissions); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request", "detail": err.Error()})
		return
	}

	// Setup
	spreadsheetID := "1YGe4AWf1qt8f5K5iJ6OGcZDGLGnGWWE0JmDZIr0jrn8"
	sheetName := "Sheet1"
	credentialsPath := "credentials/uploader_service_account.json"

	srv, err := helper.GetSheetsService(credentialsPath)
	time.Sleep(3 * time.Second)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to init Sheets API", "detail": err.Error()})
		return
	}

	// Group theo mỗi lần nộp (dựa vào SubmittedAt + StudentID + FormCode)
	type keyStruct struct {
		SubmittedAt string
		StudentID   string
		UserID      string
		FormCode    string
		FormName    string
	}
	grouped := make(map[keyStruct]map[string]string)

	for _, s := range submissions {
		k := keyStruct{
			SubmittedAt: s.SubmittedAt,
			StudentID:   s.StudentID,
			UserID:      s.UserID,
			FormCode:    s.FormCode,
			FormName:    s.FormName,
		}
		if grouped[k] == nil {
			grouped[k] = make(map[string]string)
		}
		grouped[k][s.Question] = s.Answer
	}

	// Ghi từng dòng vào sheet
	for k, answers := range grouped {
		baseInfo := []interface{}{k.SubmittedAt, k.StudentID, k.UserID, k.FormCode, k.FormName}
		err := helper.AppendFormAnswersToSheet(srv, spreadsheetID, sheetName, baseInfo, answers)
		time.Sleep(3 * time.Second)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to append row", "detail": err.Error()})
			return
		}
	}

	ctx.JSON(200, gin.H{"message": "Sync success!"})
}

func (ctrl *ApplicationController) SyncDataDemoV3(ctx *gin.Context) {
	var req request.SyncDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, response.FailedResponse{
			Code:    400,
			Message: "invalid request",
			Error:   err.Error(),
		})
		return
	}

	lastSubmitedTime, err := ctrl.SyncDataUsecase.ExcuteCreateAndSyncFormAnswer(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to start sync process",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Waiting sync data",
		Data: map[string]interface{}{
			"last_submit_time": lastSubmitedTime,
		},
	})
}

func (ctrl *ApplicationController) CheckStatusSyncQueue(ctx *gin.Context) {
	hasPending, err := ctrl.SyncDataUsecase.HasPendingSyncQueue()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check sync queue",
			Error:   err.Error(),
		})
		return
	}

	if !hasPending {
		ctx.JSON(http.StatusConflict, response.FailedResponse{
			Code:    http.StatusConflict,
			Message: "A sync is already in progress",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "No sync in progress. Ready to start sync.",
		Data: map[string]interface{}{
			"status": value.SyncQueueStatusDone,
		},
	})
}
