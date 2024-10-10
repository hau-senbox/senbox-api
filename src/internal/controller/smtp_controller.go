package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"
)

type SMTPController struct {
	*usecase.SendEmailUseCase
	*usecase.FindDeviceFromRequestCase
}

// Send An Email godoc
// @Summary Send An Email godoc
// @Description Send An Email godoc
// @Tags Device
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer {token}"
// @Param body body request.SendEmailRequest true "Send Email request"
// @Success 200 {object} response.UserListResponse
// @Failure 400 {object} response.FailedResponse
// @Router /v1/device/send/email [post]
func (receiver *SMTPController) SendEmailFromDevice(c *gin.Context) {
	var req request.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			},
		)
		return
	}
	device, err := receiver.FindDeviceFromRequestCase.FindDevice(c)
	if err != nil || device == nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				},
			},
		)
		return
	}
	err = receiver.SendEmailUseCase.SendEmail(req.To, req.Subject, req.Body, *device)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, response.FailedResponse{
				Error: response.Cause{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				},
			},
		)
		return
	}

	c.JSON(
		http.StatusOK, response.SucceedResponse{
			Data: response.Cause{
				Code:    http.StatusOK,
				Message: "Email sent successfully",
			},
		})
}
