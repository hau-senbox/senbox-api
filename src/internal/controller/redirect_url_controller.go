package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RedirectUrlController struct {
	*usecase.SaveRedirectUrlUseCase
	*usecase.GetRedirectUrlListUseCase
	*usecase.DeleteRedirectUrlUseCase
	*usecase.UpdateRedirectUrlUseCase
	*usecase.GetRedirectUrlByQRCodeUseCase
	*usecase.ImportRedirectUrlsUseCase
}

// @Summary Create redirect url
// @Description Create redirect url
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.SaveRedirectUrlRequest true "body"
// @Success 200 {object} response.SaveRedirectUrlResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/redirect-url/create [post]
func (receiver *RedirectUrlController) CreateRedirectUrl(context *gin.Context) {
	var req request.SaveRedirectUrlRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	form, err := receiver.Save(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SaveRedirectUrlResponse{Data: response.SaveRedirectUrlResponseData{
		Id:        form.ID,
		QRCode:    form.QRCode,
		TargetUrl: form.TargetUrl,
		Password:  form.Password,
		CreatedAt: form.CreatedAt,
		UpdatedAt: form.UpdatedAt,
	}})
}

// Get Get Redirect Urls godoc
// @Summary Get Redirect Url
// @Description Get Redirect Url
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.GetRedirectUrlListRequest true "Get Redirect Url List Request"
// @Success 200 {object} response.GetRedirectUrlListResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/redirect-url/list [get]
func (receiver *RedirectUrlController) GetRedirectUrlList(context *gin.Context) {
	var req request.GetRedirectUrlListRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	forms, paging, err := receiver.GetList(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.GetRedirectUrlListResponse{Data: forms, Paging: *paging})
}

// Delete Redirect Url godoc
// @Summary Delete Redirect Url
// @Description Delete Redirect Url
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Redirect Url RoleId"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/redirect-url/:id [delete]
func (receiver *RedirectUrlController) DeleteRedirectUrl(context *gin.Context) {
	formIdString := context.Param("id")
	formId, err := strconv.Atoi(formIdString)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	err = receiver.Delete(uint64(formId))
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Redirect Url deleted",
	})
}

// Update Redirect Url godoc
// @Summary Update Redirect Url
// @Description Update Redirect Url
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Redirect Url RoleId"
// @Param request body request.UpdateRedirectUrlRequest true "Update Redirect Url Request"
// @Success 200 {object} response.UpdateRedirectUrlResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/redirect-url/:id [put]
func (receiver *RedirectUrlController) UpdateRedirectUrl(context *gin.Context) {
	id := context.Param("id")
	if id == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "id is required",
		})
		return
	}
	formId, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	var req request.UpdateRedirectUrlRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	form, err := receiver.Update(formId, req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.UpdateRedirectUrlResponse{
		Data: response.GetRedirectUrlListResponseData{
			Id:           form.ID,
			QRCode:       form.QRCode,
			TargetUrl:    form.TargetUrl,
			Password:     form.Password,
			Hint:         form.Hint,
			HashPassword: form.HashPassword,
			CreatedAt:    form.CreatedAt,
			UpdatedAt:    form.UpdatedAt,
		},
	})
}

// Get Redirect Url by QR Code godoc
// @Summary Get Redirect Url by QR Code
// @Description Get Redirect Url by QR Code
// @Tags Redirect Url
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param qrcode query string true "QR Code"
// @Success 200 {object} response.GetRedirectUrlResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/redirect-url [get]
func (receiver *RedirectUrlController) GetRedirectUrlByQRCode(context *gin.Context) {
	var req request.GetRedirectUrlByQRCodeRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	form, err := receiver.GetByQRCode(req.QRCode)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.GetRedirectUrlResponse{
		Data: response.GetRedirectUrlListResponseData{
			Id:           form.ID,
			QRCode:       form.QRCode,
			TargetUrl:    form.TargetUrl,
			Password:     form.Password,
			Hint:         form.Hint,
			HashPassword: form.HashPassword,
			CreatedAt:    form.CreatedAt,
			UpdatedAt:    form.UpdatedAt,
		},
	})
}

// Import Redirect Urls godoc
// @Summary Import Redirect Urls
// @Description Import Redirect Urls
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.ImportRedirectUrlsRequest true "Import Redirect Url Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/redirect-url/import [post]
func (receiver *RedirectUrlController) ImportRedirectUrls(context *gin.Context) {
	var req request.ImportRedirectUrlsRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	err := receiver.Import(req)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Redirect Urls imported",
	})
}

type importPartiallyURLsRequest struct {
	SpreadsheetURL string `json:"spreadsheet_url" binding:"required"`
	TabName        string `json:"tab_name" binding:"required"`
}

// Import Partially Redirect Urls godoc
// @Summary Import Partially Redirect Urls
// @Description Import Partially Redirect Urls
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param request body importPartiallyURLsRequest true "Import Partially Redirect Urls Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/redirect-url/import/partially [post]
func (receiver *RedirectUrlController) ImportPartiallyRedirectUrls(context *gin.Context) {
	var req importPartiallyURLsRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	err := receiver.ImportPartially(req.SpreadsheetURL, req.TabName)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Redirect Urls imported",
	})
}
