package controller

import (
	"bufio"
	"net/http"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/pkg/randx"
	"sen-global-api/pkg/uploader"

	"github.com/gin-gonic/gin"
)

type PdfController struct {
	*usecase.UploadPDFUseCase
	*usecase.GetPdfByKeyUseCase
	*usecase.DeletePDFUseCase
}

type deletePDFByKeyRequest struct {
	Key string `json:"key"`
}

func (receiver *PdfController) CreatePDF(context *gin.Context) {

	fileHeader, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	folder := context.DefaultPostForm("folder", "pdf")
	fileName := context.DefaultPostForm("file_name", randx.GenString(10))
	orgIDString := context.DefaultPostForm("org_id", "")
	if orgIDString == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Missing org_id",
		})
		return
	}
	mode, err := uploader.UploadModeFromString(context.PostForm("mode"))
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	defer file.Close()

	dataBytes := make([]byte, fileHeader.Size)
	if _, err := bufio.NewReader(file).Read(dataBytes); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	url, pdf, err := receiver.UploadPDFUseCase.UploadPDF(dataBytes, folder, fileHeader.Filename, fileName, mode, orgIDString)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	if url == nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "pdf was not created",
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "pdf was create successfully",
		Data: response.PdfResponse{
			PdfName:        pdf.PdfName,
			Key:            pdf.Key,
			OrganizationID: pdf.OrganizationID,
			Extension:      pdf.Extension,
			Url:            *url,
		},
	})
}

func (recervier *PdfController) GetUrlByKey(context *gin.Context) {
	var req getUrlByKeyRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	mode, err := uploader.UploadModeFromString(req.Mode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	res, err := recervier.GetPdfByKey(req.Key, mode)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "pdf was get successfully",
		Data:    *res,
	})
}

func (recervier *PdfController) GetAllKeyByOrgID(context *gin.Context) {

	orgID := context.Query("org_id")
	if orgID == "" {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "org_id is required",
		})
		return
	}

	pdfs, err := recervier.GetPdfByKeyUseCase.GetAllKeyByOrgID(orgID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "pdfs were get successfully",
		Data:    pdfs,
	})
}

func (recervier *PdfController) DeletePDF(context *gin.Context) {
	var req deletePDFByKeyRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := recervier.DeletePDFUseCase.DeletePDF(req.Key)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "pdf deleted successfully",
	})
}