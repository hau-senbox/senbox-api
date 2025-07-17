package controller

import (
	"net/http"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type MenuController struct {
	*usecase.GetUserFromTokenUseCase
	*usecase.GetMenuUseCase
	*usecase.UploadSuperAdminMenuUseCase
	*usecase.UploadOrgMenuUseCase
	*usecase.UploadUserMenuUseCase
	*usecase.UploadDeviceMenuUseCase
}

type componentResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Order     int    `json:"order"`
	SectionID string `json:"section_id"`
}

type menuResponse struct {
	Top    []componentResponse `json:"top"`
	Bottom []componentResponse `json:"bottom"`
}

func (receiver *MenuController) GetSuperAdminMenu(context *gin.Context) {
	menus, err := receiver.GetMenuUseCase.GetSuperAdminMenu()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	topMenuResponse := make([]componentResponse, 0)
	bottomMenuResponse := make([]componentResponse, 0)
	for _, m := range menus {
		switch m.Direction {
		case menu.Top:
			topMenuResponse = append(topMenuResponse, componentResponse{
				ID:    m.Component.ID.String(),
				Name:  m.Component.Name,
				Type:  m.Component.Type.String(),
				Key:   m.Component.Key,
				Value: string(m.Component.Value),
				Order: m.Order,
			})
		case menu.Bottom:
			bottomMenuResponse = append(bottomMenuResponse, componentResponse{
				ID:    m.Component.ID.String(),
				Name:  m.Component.Name,
				Type:  m.Component.Type.String(),
				Key:   m.Component.Key,
				Value: string(m.Component.Value),
				Order: m.Order,
			})
		default:
			continue
		}
	}

	sort.Slice(topMenuResponse, func(i, j int) bool {
		return topMenuResponse[i].Order < topMenuResponse[j].Order
	})
	sort.Slice(bottomMenuResponse, func(i, j int) bool {
		return bottomMenuResponse[i].Order < bottomMenuResponse[j].Order
	})

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menuResponse{
			Top:    topMenuResponse,
			Bottom: bottomMenuResponse,
		},
	})
}

func (receiver *MenuController) GetOrgMenu(context *gin.Context) {
	organizationID := context.Param("id")
	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	//user, err := receiver.GetUserFromToken(context)
	//if err != nil {
	//	context.JSON(http.StatusInternalServerError, response.FailedResponse{
	//		Code:  http.StatusForbidden,
	//		Error: err.Error(),
	//	})
	//	return
	//}
	//
	//present := lo.ContainsBy(user.Organizations, func(org entity.SOrganization) bool {
	//	return org.ID == int64(id)
	//})
	//if !present {
	//	context.JSON(http.StatusForbidden, response.FailedResponse{
	//		Code:  http.StatusForbidden,
	//		Error: "access denied",
	//	})
	//	return
	//}

	menus, err := receiver.GetMenuUseCase.GetOrgMenu(organizationID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	topMenuResponse := make([]componentResponse, 0)
	bottomMenuResponse := make([]componentResponse, 0)
	for _, m := range menus {
		switch m.Direction {
		case menu.Top:
			topMenuResponse = append(topMenuResponse, componentResponse{
				ID:    m.Component.ID.String(),
				Name:  m.Component.Name,
				Type:  m.Component.Type.String(),
				Key:   m.Component.Key,
				Value: string(m.Component.Value),
				Order: m.Order,
			})
		case menu.Bottom:
			bottomMenuResponse = append(bottomMenuResponse, componentResponse{
				ID:    m.Component.ID.String(),
				Name:  m.Component.Name,
				Type:  m.Component.Type.String(),
				Key:   m.Component.Key,
				Value: string(m.Component.Value),
				Order: m.Order,
			})
		default:
			continue
		}
	}

	sort.Slice(topMenuResponse, func(i, j int) bool {
		return topMenuResponse[i].Order < topMenuResponse[j].Order
	})
	sort.Slice(bottomMenuResponse, func(i, j int) bool {
		return bottomMenuResponse[i].Order < bottomMenuResponse[j].Order
	})

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menuResponse{
			Top:    topMenuResponse,
			Bottom: bottomMenuResponse,
		},
	})
}

func (receiver *MenuController) GetStudentMenu(context *gin.Context) {
	userID := context.Param("id")
	if userID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	menus, err := receiver.GetMenuUseCase.GetStudentMenu(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	res := make([]componentResponse, 0)
	for _, m := range menus {
		res = append(res, componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(m.Component.Value),
			Order: m.Order,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Order < res[j].Order
	})

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *MenuController) GetTeacherMenu(context *gin.Context) {
	userID := context.Param("id")
	if userID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	menus, err := receiver.GetMenuUseCase.GetTeacherMenu(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	res := make([]componentResponse, 0)
	for _, m := range menus {
		res = append(res, componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(m.Component.Value),
			Order: m.Order,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Order < res[j].Order
	})

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *MenuController) GetUserMenu(context *gin.Context) {
	userID := context.Param("id")
	if userID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	menus, err := receiver.GetMenuUseCase.GetUserMenu(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	res := make([]componentResponse, 0)
	for _, m := range menus {
		res = append(res, componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(m.Component.Value),
			Order: m.Order,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Order < res[j].Order
	})

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *MenuController) GetDeviceMenu(context *gin.Context) {
	deviceID := context.Param("id")
	if deviceID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	menus, err := receiver.GetMenuUseCase.GetDeviceMenu(deviceID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	type deviceComponentResponse struct {
		OrganizationID string              `json:"organization_id"`
		Components     []componentResponse `json:"components"`
	}

	resMap := make(map[string][]componentResponse)
	for _, m := range menus {
		key := m.OrganizationID.String()
		resMap[key] = append(resMap[key], componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(m.Component.Value),
			Order: m.Order,
		})
	}

	for key := range resMap {
		sort.Slice(resMap[key], func(i, j int) bool {
			return resMap[key][i].Order < resMap[key][j].Order
		})
	}

	// Convert map to slice
	res := make([]deviceComponentResponse, 0, len(resMap))
	for key, components := range resMap {
		res = append(res, deviceComponentResponse{
			OrganizationID: key,
			Components:     components,
		})
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *MenuController) GetDeviceMenuByOrg(context *gin.Context) {
	organizationID := context.Param("organization_id")
	if organizationID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	menus, err := receiver.GetMenuUseCase.GetDeviceMenuByOrg(organizationID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	res := make([]componentResponse, 0)
	for _, m := range menus {
		res = append(res, componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(m.Component.Value),
			Order: m.Order,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Order < res[j].Order
	})

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: res,
	})
}

func (receiver *MenuController) UploadSuperAdminMenu(context *gin.Context) {
	var req request.UploadSuperAdminMenuRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UploadSuperAdminMenuUseCase.Upload(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "menu was upload successfully",
	})
}

func (receiver *MenuController) UploadOrgMenu(context *gin.Context) {
	var req request.UploadOrgMenuRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	present := lo.ContainsBy(user.Organizations, func(org entity.SOrganization) bool {
		return org.ID.String() == req.OrganizationID
	})
	if !present {
		context.JSON(http.StatusForbidden, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: "access denied",
		})
		return
	}

	err = receiver.UploadOrgMenuUseCase.Upload(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "menu was upload successfully",
	})
}

func (receiver *MenuController) UploadUserMenu(context *gin.Context) {
	var req request.UploadUserMenuRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	present := lo.ContainsBy(user.Organizations, func(org entity.SOrganization) bool {
		return org.ID.String() == req.OrganizationID
	})
	if !present {
		context.JSON(http.StatusForbidden, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: "access denied",
		})
		return
	}

	err = receiver.UploadUserMenuUseCase.Upload(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "menu was upload successfully",
	})
}

func (receiver *MenuController) UploadDeviceMenu(context *gin.Context) {
	var req request.UploadDeviceMenuRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	user, err := receiver.GetUserFromToken(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: err.Error(),
		})
		return
	}

	present := lo.ContainsBy(user.Roles, func(org entity.SRole) bool {
		return org.Role == entity.SuperAdmin || org.Role == entity.Admin
	})
	if !present {
		context.JSON(http.StatusForbidden, response.FailedResponse{
			Code:  http.StatusForbidden,
			Error: "access denied",
		})
		return
	}

	err = receiver.UploadDeviceMenuUseCase.Upload(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "menu was upload successfully",
	})
}

func (receiver *MenuController) GetCommonMenu(context *gin.Context) {
	result := receiver.GetMenuUseCase.GetCommonMenu(context)
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: result,
	})
}

func (receiver *MenuController) GetCommonMenuByUser(context *gin.Context) {
	result := receiver.GetMenuUseCase.GetCommonMenuByUser(context)
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: result,
	})
}

func (receiver *MenuController) UploadSectionMenu(context *gin.Context) {
	var req request.UploadSectionMenuRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UploadUserMenuUseCase.UploadSectionMenu(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Section menu was upload successfully",
	})
}

func (receiver *MenuController) GetSectionMenu(context *gin.Context) {

	menus, err := receiver.GetMenuUseCase.GetSectionMenu()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	// res := make([]componentResponse, 0)
	// for _, m := range menus {
	// 	res = append(res, componentResponse{
	// 		ID:    m.Component.ID.String(),
	// 		Name:  m.Component.Name,
	// 		Type:  m.Component.Type.String(),
	// 		Key:   m.Component.Key,
	// 		Value: string(m.Component.Value),
	// 		Order: m.Order,
	// 	})
	// }

	// sort.Slice(res, func(i, j int) bool {
	// 	return res[i].Order < res[j].Order
	// })

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menus,
	})
}
