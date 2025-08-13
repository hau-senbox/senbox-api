package controller

import (
	"fmt"
	"net/http"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sen-global-api/internal/domain/value"

	"github.com/gin-gonic/gin"
)

type LanguagesConfigController struct {
	uc *usecase.LanguagesConfigUsecase
}

func NewLanguagesConfigController(uc *usecase.LanguagesConfigUsecase) *LanguagesConfigController {
	return &LanguagesConfigController{uc: uc}
}

// GET /languages-config?owner_id=xxx&owner_type=student
func (c *LanguagesConfigController) GetByOwner(ctx *gin.Context) {
	ownerID := ctx.Query("owner_id")
	ownerRoleStr := ctx.Query("owner_role")
	ownerRole, err := value.ParseOwnerRole4LangConfig(ownerRoleStr)

	if ownerID == "" || err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "owner id và owner role là bắt buộc và hợp lệ",
		})
		return
	}

	list, err := c.uc.GetLanguagesConfigByOwner(ctx, ownerID, ownerRole)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get languages config",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    list,
	})
}

func (c *LanguagesConfigController) UploadLanguagesConfig(ctx *gin.Context) {
	var req request.UploadLanguagesConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	ownerRoleStr := req.OwnerRole
	ownerRole, parseErr := value.ParseOwnerRole4LangConfig(ownerRoleStr)
	if parseErr != nil {
		ctx.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("invalid owner role: %s", ownerRoleStr),
		})
		return
	}

	err := c.uc.UploadLanguagesConfig(
		ctx,
		req.OwnerID,
		ownerRole,
		entity.LanguageConfigList(req.SpokenLang),
		entity.LanguageConfigList(req.StudyLang),
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to upload languages config",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.SucceedResponse{
		Code:    http.StatusCreated,
		Message: "Upload languages config successfully",
		Data:    nil,
	})
}
