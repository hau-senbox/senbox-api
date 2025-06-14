package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionController struct {
	DBConn                               *gorm.DB
	GetUserQuestionsUseCase              usecase.GetUserQuestionsUseCase
	GetUserFromTokenUseCase              usecase.GetUserFromTokenUseCase
	GetQuestionByIdUseCase               usecase.GetQuestionByIdUseCase
	GetDeviceIdFromTokenUseCase          usecase.GetDeviceIdFromTokenUseCase
	GetQuestionByFormUseCase             usecase.GetQuestionsByFormUseCase
	GetFormByIdUseCase                   usecase.GetFormByIdUseCase
	GetAllQuestionsUseCase               usecase.GetAllQuestionsUseCase
	CreateFormUseCase                    usecase.CreateFormUseCase
	GetRawQuestionFromSpreadsheetUseCase usecase.GetRawQuestionFromSpreadsheetUseCase
	GetShowPicsQuestionDetailUseCase     usecase.GetShowPicsQuestionDetailUseCase
	FindDeviceFromRequestCase            usecase.FindDeviceFromRequestCase
	GetUserDeviceUseCase                 usecase.GetUserDeviceUseCase
	GetDeviceByIdUseCase                 usecase.GetDeviceByIdUseCase
}

// Get Form's Questions by QR Code godoc
// @Summary Get Form's Questions by QR Code
// @Description Get Form's Questions by QR Code
// @Tags Form
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param qr_code path string true "QR Code"
// @Success 200 {object} response.QuestionListResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router /v1/form [post]
func (receiver *QuestionController) GetFormQRCode(context *gin.Context) {
	var req request.GetFormRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	form, err := receiver.GetFormByIdUseCase.GetFormByQRCode(req.QrCode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	succeedRes, failedRes := receiver.GetQuestionByFormUseCase.GetQuestionByForm(*form)

	if failedRes != nil {
		context.JSON(http.StatusBadRequest, failedRes)
		return
	}

	context.JSON(http.StatusOK, succeedRes)
}

// Get Show Pics Question Detail godoc
// @Summary Get Show Pics Question Detail
// @Description Get Show Pics Question Detail
// @Tags Question
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param id path string true "Question RoleId"
// @Success 200 {object} response.GetShowPicsQuestionResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router /v1/question/show-pics [get]
func (receiver *QuestionController) GetShowPicsQuestion(context *gin.Context) {
	questionId := context.Query("id")
	photo, err := receiver.GetShowPicsQuestionDetailUseCase.Execute(questionId)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	if photo == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "photo not found",
		})
		return
	}
	context.JSON(http.StatusOK, response.GetShowPicsQuestionResponse{
		Data: response.GetShowPicsQuestionResponseData{
			PhotoURL: photo,
		},
	})
}
