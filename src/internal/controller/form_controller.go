package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"
)

type FormController struct {
	SaveFormUseCase    usecase.SaveFormUseCase
	DeleteFormUseCase  usecase.DeleteFormUseCase
	GetFormListUseCase usecase.GetFormListUseCase
	UpdateFormUseCase  usecase.UpdateFormUseCase
	SearchFormsUseCase usecase.SearchFormsUseCase
	ImportFormsUseCase *usecase.ImportFormsUseCase
}

// Create Form godoc
// @Summary Create Form
// @Description Create Form
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.SaveFormRequest true "Create Form Request"
// @Success 200 {object} response.SaveFormResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/form/create [post]
func (receiver *FormController) CreateForm(context *gin.Context) {
	var req request.SaveFormRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	form, err := receiver.SaveFormUseCase.SaveForm(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SaveFormResponse{Data: response.SaveFormResponseData{
		Id:          form.FormId,
		Spreadsheet: form.SpreadsheetUrl,
		Password:    form.Password,
		Note:        form.Note,
		CreatedAt:   form.CreatedAt,
		UpdatedAt:   form.UpdatedAt,
	}})
}

// Get Forms godoc
// @Summary Get Forms
// @Description Get Forms
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} response.GetFormListResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/form/list [get]
func (receiver *FormController) GetFormList(context *gin.Context) {
	var req request.GetFormListRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	forms, paging, err := receiver.GetFormListUseCase.GetFormList(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.GetFormListResponse{Data: forms, Paging: *paging})
}

// Delete Form godoc
// @Summary Delete Form
// @Description Delete Form
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Form ID"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/form/delete/:id [delete]
func (receiver *FormController) DeleteForm(context *gin.Context) {
	formIdString := context.Param("id")
	formId, err := strconv.Atoi(formIdString)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	err = receiver.DeleteFormUseCase.DeleteForm(uint64(formId))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Form deleted",
		},
	})
}

// Update Form godoc
// @Summary Update Form
// @Description Update Form
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Form ID"
// @Param request body request.UpdateFormRequest true "Update Form Request"
// @Success 200 {object} response.UpdateFormResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/form/:id [put]
func (receiver *FormController) UpdateForm(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "id is required",
			},
		})
		return
	}
	formId, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	var req request.UpdateFormRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	form, err := receiver.UpdateFormUseCase.UpdateForm(formId, req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.UpdateFormResponse{
		Data: response.GetFormListResponseData{
			Id:          form.FormId,
			Spreadsheet: form.SpreadsheetUrl,
			Password:    form.Password,
			Note:        form.Note,
			CreatedAt:   form.CreatedAt,
			UpdatedAt:   form.UpdatedAt,
		},
	})
}

// Search Forms godoc
// @Summary Search Forms
// @Description Search Forms
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param q query string true "Search Query"
// @Success 200 {object} response.GetFormListResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/forms/search [get]
func (receiver *FormController) SearchForms(context *gin.Context) {
	keyword := context.Query("keyword")
	if keyword == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "keyword is required",
			},
		})
		return
	}
	forms, err := receiver.SearchFormsUseCase.SearchForms(keyword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	if forms == nil {
		empty := make([]response.GetFormListResponseData, 0)
		context.JSON(http.StatusOK, response.GetFormListResponse{Data: empty})
		return
	}
	context.JSON(http.StatusOK, response.GetFormListResponse{Data: forms})
}

// Import Forms godoc
// @Summary Import Forms
// @Description Import Forms
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
// @Router /v1/admin/forms/import [post]
func (receiver *FormController) ImportForms(context *gin.Context) {
	receiver.importForms(context, usecase.FormsUploaderIndexFirst)
}

// Import Forms2 godoc
// @Summary Import Forms2
// @Description Import Forms2
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
// @Router /v1/admin/forms2/import [post]
func (receiver *FormController) ImportForms2(context *gin.Context) {
	receiver.importForms(context, usecase.FormsUploaderIndexSecond)
}

// Import Forms3 godoc
// @Summary Import Forms3
// @Description Import Forms3
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
// @Router /v1/admin/forms3/import [post]
func (receiver *FormController) ImportForms3(context *gin.Context) {
	receiver.importForms(context, usecase.FormsUploaderIndexThird)
}

// Import Forms4 godoc
// @Summary Import Forms4
// @Description Import Forms4
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
// @Router /v1/admin/forms4/import [post]
func (receiver *FormController) ImportForms4(context *gin.Context) {
	receiver.importForms(context, usecase.FormsUploaderIndexFourth)
}

func (receiver *FormController) importForms(context *gin.Context, index usecase.FormsUploaderIndex) {
	var req request.ImportFormRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}
	err := receiver.ImportFormsUseCase.ImportForms(req, index)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Form imported",
		},
	})
}

type importPartiallyFormRequest struct {
	SpreadsheetURL string `json:"spreadsheet_url" binding:"required"`
	TabName        string `json:"tab_name" binding:"required"`
}

// Import Forms Partially godoc
// @Summary Import Forms Partially
// @Description Import Forms Partially
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body importPartiallyFormRequest true "Import Forms At Tab Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/forms/import/partially [post]
func (receiver *FormController) ImportFormsPartially(context *gin.Context) {
	var req importPartiallyFormRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		})
		return
	}

	err := receiver.ImportFormsUseCase.ImportFormsPartially(req.SpreadsheetURL, req.TabName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: response.Cause{
			Code:    http.StatusOK,
			Message: "Form imported",
		},
	})
}

// Import Sign Up Forms godoc
// @Summary Import Sign Up Forms
// @Description Import Sign Up Forms
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
// @Router /v1/admin/forms/sign [post]
func (receiver *FormController) ImportSignUpForms(context *gin.Context) {
	receiver.importForms(context, usecase.FormsUploaderIndexFifth)
}
