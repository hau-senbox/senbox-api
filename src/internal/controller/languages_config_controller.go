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
	LanguagesConfigUsecase *usecase.LanguagesConfigUsecase
	ChildUsecase           *usecase.ChildUseCase
	GetUserFromToken       *usecase.GetUserFromTokenUseCase
}

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

	list, err := c.LanguagesConfigUsecase.GetLanguagesConfigByOwner(ctx, ownerID, ownerRole)
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

	err := c.LanguagesConfigUsecase.UploadLanguagesConfig(
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

func (c *LanguagesConfigController) GetChildLanguageConfig(ctx *gin.Context) {

	childID := ctx.Param("child_id")
	if childID == "" {
		ctx.JSON(http.StatusBadGateway, response.FailedResponse{
			Code:    http.StatusBadGateway,
			Message: "child id is required",
		})
	}

	//check child belong to user access
	userID, ok := getUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.FailedResponse{
			Code:  http.StatusUnauthorized,
			Error: "Unauthorized: invalid user_id",
		})
		return
	}
	isParent, _ := c.ChildUsecase.IsParentOfChild(userID, childID)
	if !isParent {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Access Denied",
		})
		return
	}
	list, err := c.LanguagesConfigUsecase.GetLanguagesConfigByOwner(ctx, childID, value.OwnerRoleLangChild)
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

func (c *LanguagesConfigController) GetParentLanguageConfig(ctx *gin.Context) {

	parentID := ctx.Param("parent_id")
	if parentID == "" {
		ctx.JSON(http.StatusBadGateway, response.FailedResponse{
			Code:    http.StatusBadGateway,
			Message: "child id is required",
		})
	}

	list, err := c.LanguagesConfigUsecase.GetLanguagesConfigByOwner(ctx, parentID, value.OwnerRoleLangParent)
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

func getUserID(ctx *gin.Context) (string, bool) {
	val, exists := ctx.Get("user_id")
	if !exists {
		return "", false
	}
	userID, ok := val.(string)
	return userID, ok
}

func (c *LanguagesConfigController) GetStudyLanguage4OrganizationAssign4Web(ctx *gin.Context) {

	currentUser, err := c.GetUserFromToken.GetUserFromToken(ctx)
	res, err := c.LanguagesConfigUsecase.GetStudyLanguage4OrganizationAssign4Web(ctx, currentUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get study language",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    res,
	})

}
