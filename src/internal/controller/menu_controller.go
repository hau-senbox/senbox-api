package controller

import (
	"encoding/json"
	"net/http"
	"sen-global-api/helper"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"sort"

	log "github.com/sirupsen/logrus"

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
	*usecase.UploadSectionMenuUseCase
	*usecase.ChildMenuUseCase
	*usecase.StudentMenuUseCase
	*usecase.StudentApplicationUseCase
	*usecase.TeacherMenuUseCase
	*usecase.DeviceMenuUseCase
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
	MenuIconKey string              `json:"menu_icon_key,omitempty"`
	Top         []componentResponse `json:"top"`
	Bottom      []componentResponse `json:"bottom"`
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
		normalizedValue, err := helper.NormalizeComponentValue(m.Component.Value)
		if err != nil {
			log.Println("Normalize error:", err)
		}
		switch m.Direction {
		case menu.Top:
			topMenuResponse = append(topMenuResponse, componentResponse{
				ID:    m.Component.ID.String(),
				Name:  m.Component.Name,
				Type:  m.Component.Type.String(),
				Key:   m.Component.Key,
				Value: string(normalizedValue),
				Order: m.Order,
			})
		case menu.Bottom:
			bottomMenuResponse = append(bottomMenuResponse, componentResponse{
				ID:    m.Component.ID.String(),
				Name:  m.Component.Name,
				Type:  m.Component.Type.String(),
				Key:   m.Component.Key,
				Value: string(normalizedValue),
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

func (receiver *MenuController) GetSuperAdminMenu4App(context *gin.Context) {
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
		normalizedValue, err := helper.NormalizeComponentValue(m.Component.Value)
		if err != nil {
			log.Println("Normalize error:", err)
		}

		var compVal components.ComponentFullValue
		if err := json.Unmarshal(normalizedValue, &compVal); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		// Bỏ qua nếu không visible
		if !compVal.Visible {
			continue
		}

		switch m.Direction {
		case menu.Top:
			topMenuResponse = append(topMenuResponse, componentResponse{
				ID:    m.Component.ID.String(),
				Name:  m.Component.Name,
				Type:  m.Component.Type.String(),
				Key:   m.Component.Key,
				Value: string(normalizedValue),
				Order: m.Order,
			})
		case menu.Bottom:
			bottomMenuResponse = append(bottomMenuResponse, componentResponse{
				ID:    m.Component.ID.String(),
				Name:  m.Component.Name,
				Type:  m.Component.Type.String(),
				Key:   m.Component.Key,
				Value: string(normalizedValue),
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

func (receiver *MenuController) GetOrgMenu4App(context *gin.Context) {
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

		normalizedValue, err := helper.NormalizeComponentValue(m.Component.Value)
		if err != nil {
			log.Println("Normalize error:", err)
		}

		var compVal components.ComponentFullValue
		if err := json.Unmarshal(normalizedValue, &compVal); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		// Bỏ qua nếu không visible
		if !compVal.Visible {
			continue
		}

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

	// get menu icon key

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menuResponse{
			Top:    topMenuResponse,
			Bottom: bottomMenuResponse,
		},
	})
}

func (receiver *MenuController) GetStudentMenu4App(context *gin.Context) {
	studentID := context.Param("id")
	if studentID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "id is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	menus, err := receiver.GetMenuUseCase.GetStudentMenu4App(studentID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menus,
	})
}

func (receiver *MenuController) GetTeacherMenu4App(context *gin.Context) {
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

	menus, err := receiver.GetMenuUseCase.GetTeacherMenu4App(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menus,
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
		normalizedValue, err := helper.NormalizeComponentValue(m.Component.Value)
		if err != nil {
			log.Println("Normalize error:", err)
		}
		res = append(res, componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(normalizedValue),
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

func (receiver *MenuController) GetUserMenu4App(context *gin.Context) {
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
		normalizedValue, err := helper.NormalizeComponentValue(m.Component.Value)
		if err != nil {
			log.Println("Normalize error:", err)
		}

		var compVal components.ComponentFullValue
		if err := json.Unmarshal(normalizedValue, &compVal); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		// Bỏ qua nếu không visible
		if !compVal.Visible {
			continue
		}
		res = append(res, componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(normalizedValue),
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
		normalizedValue, err := helper.NormalizeComponentValue(m.Component.Value)
		if err != nil {
			log.Println("Normalize error:", err)
		}
		key := m.OrganizationID.String()
		resMap[key] = append(resMap[key], componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(normalizedValue),
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

func (receiver *MenuController) GetDeviceMenu4App(context *gin.Context) {
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
		normalizedValue, err := helper.NormalizeComponentValue(m.Component.Value)
		if err != nil {
			log.Println("Normalize error:", err)
		}
		var compVal components.ComponentFullValue
		if err := json.Unmarshal(normalizedValue, &compVal); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		// Bỏ qua nếu không visible
		if !compVal.Visible {
			continue
		}

		key := m.OrganizationID.String()
		resMap[key] = append(resMap[key], componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(normalizedValue),
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

func (receiver *MenuController) GetDeviceMenu4Admin(context *gin.Context) {
	deviceID := context.Param("id")
	if deviceID == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: "deviceID is required",
				Code:  http.StatusBadRequest,
			},
		)
		return
	}

	menus, err := receiver.DeviceMenuUseCase.GetByDeviceID(deviceID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menus,
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
		normalizedValue, err := helper.NormalizeComponentValue(m.Component.Value)
		if err != nil {
			log.Println("Normalize error:", err)
		}
		res = append(res, componentResponse{
			ID:    m.Component.ID.String(),
			Name:  m.Component.Name,
			Type:  m.Component.Type.String(),
			Key:   m.Component.Key,
			Value: string(normalizedValue),
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

	// Validate DeleteComponentIDs for each item
	for _, item := range req {
		if err := item.Validate(); err != nil {
			context.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			})
			return
		}
	}

	err := receiver.UploadSectionMenuUseCase.UploadSectionMenuV2(context, req)
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

func (receiver *MenuController) UploadStudentMenu(context *gin.Context) {
	var req request.UploadSectionMenuStudentRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	log.WithFields(log.Fields{
		"request": req,
	}).Info("Received UploadStudentMenu request")

	err := receiver.UploadSectionMenuUseCase.UploadStudentMenu(context, req)
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

func (receiver *MenuController) UploadTeacherMenu(context *gin.Context) {
	var req request.UploadSectionMenuTeacherRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UploadSectionMenuUseCase.UploadTeacherMenu(context, req)
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

func (receiver *MenuController) UploadStaffMenu(context *gin.Context) {
	var req request.UploadSectionMenuStaffRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UploadSectionMenuUseCase.UploadStaffMenu(context, req)
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

func (receiver *MenuController) UploadChildMenu(context *gin.Context) {
	var req request.UploadSectionMenuChildRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UploadSectionMenuUseCase.UploadChildMenu(context, req)
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

func (receiver *MenuController) UploadDeviceSectionMenu(context *gin.Context) {
	var req request.UploadSectionMenuDeviceRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UploadSectionMenuUseCase.UploadDeviceMenu(context, req)
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

func (receiver *MenuController) UploadParentMenu(context *gin.Context) {
	var req request.UploadSectionMenuParentRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UploadSectionMenuUseCase.UploadParentMenu(context, req)
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

	menus, err := receiver.GetMenuUseCase.GetSectionMenu(context)
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

func (receiver *MenuController) GetSectionMenu4WebAdmin(context *gin.Context) {

	menus, err := receiver.GetMenuUseCase.GetSectionMenu4WebAdmin(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menus,
	})
}

func (receiver *MenuController) GetSectionMenu4App(context *gin.Context) {

	menus, err := receiver.GetMenuUseCase.GetSectionMenu4App(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menus,
	})
}

func (receiver *MenuController) GetChildMenuByChildID(context *gin.Context) {
	childID := context.Param("id")

	if childID == "" {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: "Id is required",
		})
		return
	}

	menus, err := receiver.ChildMenuUseCase.GetByChildID(childID, false)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: menus,
	})
}

func (receiver *MenuController) UpdateIsShowChildMenu(context *gin.Context) {
	var req request.UpdateChildMenuRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.ChildMenuUseCase.UpdateIsShowByChildAndComponentID(req)

	if err != nil {
		context.JSON(http.StatusOK, response.SucceedResponse{
			Code: http.StatusOK,
			Data: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: "Updated",
	})
}

func (receiver *MenuController) UpdateIsShowStudentMenu(context *gin.Context) {
	var req request.UpdateStudentMenuRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.StudentMenuUseCase.UpdateIsShowByStudentAndComponentID(context, req)

	if err != nil {
		context.JSON(http.StatusBadRequest, response.SucceedResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: "Updated",
	})
}

func (receiver *MenuController) UpdateIsShowTeacherMenu(context *gin.Context) {
	var req request.UpdateTeacherMenuRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.TeacherMenuUseCase.UpdateIsShow(context, req)

	if err != nil {
		context.JSON(http.StatusBadRequest, response.SucceedResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: "Updated",
	})
}

func (receiver *MenuController) DeleteSectionMenu(context *gin.Context) {
	componentID := context.Param("id")

	if componentID == "" {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: "Id is required",
		})
		return
	}

	err := receiver.UploadSectionMenuUseCase.DeleteSectionMenu(componentID)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.SucceedResponse{
			Code: http.StatusBadRequest,
			Data: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code: http.StatusOK,
		Data: "Deleted",
	})
}
