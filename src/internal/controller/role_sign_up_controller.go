package controller

import (
	"net/http"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type RoleOrgSignUpController struct {
	UseCase *usecase.RoleOrgSignUpUseCase
}

func NewRoleOrgSignUpController(uc *usecase.RoleOrgSignUpUseCase) *RoleOrgSignUpController {
	return &RoleOrgSignUpController{UseCase: uc}
}

// POST /api/rolesignup
func (ctrl *RoleOrgSignUpController) CreateOrUpdate(c *gin.Context) {
	var role entity.SRoleOrgSignUp
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := ctrl.UseCase.UpdateOrCreateExecute(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role created or updated successfully"})
}

// GET /api/rolesignup
func (ctrl *RoleOrgSignUpController) GetAll(c *gin.Context) {
	roles, err := ctrl.UseCase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// GET /api/rolesignup/:name
func (ctrl *RoleOrgSignUpController) GetByRoleName(c *gin.Context) {
	roleName := c.Param("name")
	if roleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing role name"})
		return
	}

	role, err := ctrl.UseCase.GetByRoleName(roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (ctrl *RoleOrgSignUpController) Get4AdminWeb(c *gin.Context) {
	role, err := ctrl.UseCase.GetByRoleName("Child")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Role 'Child' not found"})
		return
	}
	res := make([]entity.SRoleOrgSignUp, 0)
	res = append(res, entity.SRoleOrgSignUp{
		ID:        role.ID,
		RoleName:  role.RoleName,
		OrgCode:   role.OrgCode,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	})
	c.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}
