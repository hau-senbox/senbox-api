package controller

import (
	"net/http"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserConfigController struct {
	*usecase.GetUserConfigUseCase
}

func (receiver *UserConfigController) GetUserConfigById(context *gin.Context) {
	userId := context.Param("id")
	if userId == "" {
		context.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: "id is required",
				},
			},
		)
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "invalid id",
			},
		})
		return
	}

	userConfig, err := receiver.GetUserConfigUseCase.GetUserConfigById(uint(id))
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
		Data: response.UserConfigResponse{
			ID:                   userConfig.ID,
			TopButtonConfig:      userConfig.TopButtonConfig,
			StudentOutputSheetId: userConfig.StudentOutputSheetId,
			TeacherOutputSheetId: userConfig.TeacherOutputSheetId,
		},
	})
}

// // Create ServiceDiscovery instance with Consul address and service name
// sd, err := consulapi.NewServiceDiscovery("inventory-service")
// if err != nil {
// 	context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 		Error: response.Cause{
// 			Code:    http.StatusInternalServerError,
// 			Message: fmt.Sprintf("error fetching service: %v", err),
// 		},
// 	})
// 	return
// }

// // Discover service
// service, err := sd.DiscoverService()
// if err != nil {
// 	context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 		Error: response.Cause{
// 			Code:    http.StatusInternalServerError,
// 			Message: fmt.Sprintf("Error discovering service: %v", err),
// 		},
// 	})
// 	return
// }

// // Define API endpoint
// apiEndpoint := "/user"

// // Example: GET request
// res, err := consulapi.CallAPI(service, apiEndpoint, "GET", nil, nil)
// if err != nil {
// 	context.JSON(http.StatusInternalServerError, response.FailedResponse{
// 		Error: response.Cause{
// 			Code:    http.StatusInternalServerError,
// 			Message: fmt.Sprintf("Error calling GET API: %v", err),
// 		},
// 	})
// 	return
// }
