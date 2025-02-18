package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
	"github.com/tiendc/gofn"
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
	SyncQuestionsUseCase                 usecase.SyncQuestionsUseCase
	GetButtonsQuestionDetailUseCase      usecase.GetButtonsQuestionDetailUseCase
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
	var request request.GetFormRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	form, err := receiver.GetFormByIdUseCase.GetFormByQRCode(request.QrCode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	user, err := receiver.GetUserFromTokenUseCase.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	userDevices, err := receiver.GetUserDeviceUseCase.GetUserDeviceById(request.DeviceId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	var userIds []string
	for _, userDevice := range *userDevices {
		userIds = append(userIds, userDevice.UserId.String())
	}
	isExist := gofn.ContainSlice(userIds, []string{
		user.ID.String(),
	})
	if !isExist {
		context.JSON(http.StatusUnauthorized, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized",
			},
		})
		return
	}

	device, err := receiver.GetDeviceByIdUseCase.GetDeviceById(request.DeviceId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	if form == nil {
		context.JSON(http.StatusNotFound, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusNotFound,
				Message: "Form not found",
			},
		})
		return
	}
	succeedRes, failedRes := receiver.GetQuestionByFormUseCase.GetQuestionByForm(*form, *device)

	if failedRes != nil {
		context.JSON(http.StatusBadRequest, failedRes)
		return
	}

	context.JSON(http.StatusOK, succeedRes)
}

// Get Buttons Question Detail godoc
// @Summary Get Button Question Detail
// @Description Get Button Question Detail
// @Tags Question
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param id path string true "Question ID"
// @Success 200 {object} response.GetScreenButtonsResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router /v1/question/buttons [get]
func (receiver *QuestionController) GetButtonsQuestion(context *gin.Context) {
	questionId := context.Query("id")
	buttons, err := receiver.GetButtonsQuestionDetailUseCase.Execute(questionId)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.GetScreenButtonsResponse{
		Data: response.GetScreenButtonsResponseData{
			Buttons: buttons,
		},
	})
}

// Get Show Pics Question Detail godoc
// @Summary Get Show Pics Question Detail
// @Description Get Show Pics Question Detail
// @Tags Question
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param id path string true "Question ID"
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
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	if photo == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "photo not found",
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.GetShowPicsQuestionResponse{
		Data: response.GetShowPicsQuestionResponseData{
			PhotoURL: photo,
		},
	})
}
