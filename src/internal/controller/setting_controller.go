package controller

import (
	"net/http"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type SettingController struct {
	*usecase.GetSettingsUseCase
	*usecase.UpdateOutputSubmissionSettingUseCase
	*usecase.UpdateOutputSummarySettingUseCase
	*usecase.UpdateEmailHistorySettingUseCase
	*usecase.UpdateOutputTemplateSettingUseCase
	*usecase.UpdateOutputTemplateSettingForTeacherUseCase
	*usecase.AdminSignUpUseCases
	*usecase.UpdateSettingNameUseCase
	*usecase.UpdateApiDistributorUseCase
}

// Get Settings godoc
// @Summary      Retrieve settings
// @Description  Retrieve settings
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Success      200  {object}  response.GetSettingsResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/settings/ [get]
func (receiver *SettingController) GetSettings(context *gin.Context) {
	settings, err := receiver.GetSettingsUseCase.GetSettings()
	if err != nil || settings == nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	var formSettings *response.GetSettingsResponseDataImport = nil
	var formSettings2 *response.GetSettingsResponseDataImport = nil
	var formSettings3 *response.GetSettingsResponseDataImport = nil
	var formSettings4 *response.GetSettingsResponseDataImport = nil
	var signupformsSettings *response.GetSettingsResponseDataImport = nil
	var urlSettings *response.GetSettingsResponseDataImport = nil
	var outputSettings *response.GetSettingsResponseDataSubmission = nil
	var summarySettings *response.GetSettingsResponseDataSummary = nil
	var syncDevicesSettings *response.GetSettingsResponseDataImport = nil
	var syncToDosSettings *response.GetSettingsResponseDataImport = nil
	var emailSettings *response.GetSettingsResponseDataSummary = nil
	var outputTemplate *response.GetSettingsResponseDataSummary = nil
	var outputTemplateForTeacher *response.GetSettingsResponseDataSummary = nil
	var signUpButton1 *response.GetSettingsResponseTextButton = nil
	var signUpButton2 *response.GetSettingsResponseTextButton = nil
	var signUpButton3 *response.GetSettingsResponseTextButton = nil
	var signUpButton4 *response.GetSettingsResponseTextButton = nil
	var signUpButton5 *response.GetSettingsResponseTextButton = nil
	var signUpButtonConfiguration *response.GetSettingsResponseDataSummary = nil
	var registrationForm *response.GetSettingsResponseDataSummary = nil
	var registrationSubmission *response.GetSettingsResponseDataSummary = nil
	var registrationPreset2 *response.GetSettingsResponseDataSummary = nil
	var apiDistributer *response.GetSettingsResponseAPIDistributor = nil
	var codeCountingSetting *response.GetSettingsResponseAPIDistributor = nil
	var registrationPreset1 *response.GetSettingsResponseDataSummary = nil

	if settings.Form != nil {
		formSettings = &response.GetSettingsResponseDataImport{
			SettingName:    settings.Form.SettingName,
			SpreadSheetUrl: settings.Form.SpreadSheetUrl,
			Auto:           settings.Form.AutoImport,
			Interval:       settings.Form.Interval,
		}
	}

	if settings.Form2 != nil {
		formSettings2 = &response.GetSettingsResponseDataImport{
			SettingName:    settings.Form2.SettingName,
			SpreadSheetUrl: settings.Form2.SpreadSheetUrl,
			Auto:           settings.Form2.AutoImport,
			Interval:       settings.Form2.Interval,
		}
	}

	if settings.Form3 != nil {
		formSettings3 = &response.GetSettingsResponseDataImport{
			SettingName:    settings.Form3.SettingName,
			SpreadSheetUrl: settings.Form3.SpreadSheetUrl,
			Auto:           settings.Form3.AutoImport,
			Interval:       settings.Form3.Interval,
		}
	}

	if settings.Form4 != nil {
		formSettings4 = &response.GetSettingsResponseDataImport{
			SettingName:    settings.Form4.SettingName,
			SpreadSheetUrl: settings.Form4.SpreadSheetUrl,
			Auto:           settings.Form4.AutoImport,
			Interval:       settings.Form4.Interval,
		}
	}

	if settings.Url != nil {
		urlSettings = &response.GetSettingsResponseDataImport{
			SettingName:    settings.Url.SettingName,
			SpreadSheetUrl: settings.Url.SpreadSheetUrl,
			Auto:           settings.Url.AutoImport,
			Interval:       settings.Url.Interval,
		}
	}

	if settings.Output != nil {
		outputSettings = &response.GetSettingsResponseDataSubmission{
			SettingName: settings.Output.SettingName,
			FolderUrl:   "https://drive.google.com/drive/folders/" + settings.Output.FolderId,
		}
	}

	if settings.Summary != nil {
		summarySettings = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.Summary.SettingName,
			SpreadSheetUrl: "https://docs.google.com/spreadsheets/d/" + settings.Summary.SpreadsheetId,
		}
	}

	if settings.SyncDevices != nil {
		syncDevicesSettings = &response.GetSettingsResponseDataImport{
			SettingName:    settings.SyncDevices.SettingName,
			SpreadSheetUrl: settings.SyncDevices.SpreadSheetUrl,
			Auto:           settings.SyncDevices.AutoImport,
			Interval:       settings.SyncDevices.Interval,
		}
	}

	if settings.SyncToDos != nil {
		syncToDosSettings = &response.GetSettingsResponseDataImport{
			SettingName:    settings.SyncToDos.SettingName,
			SpreadSheetUrl: settings.SyncToDos.SpreadSheetUrl,
			Auto:           settings.SyncToDos.AutoImport,
			Interval:       settings.SyncToDos.Interval,
		}
	}

	if settings.EmailSetting != nil {
		emailSettings = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.EmailSetting.SettingName,
			SpreadSheetUrl: "https://docs.google.com/spreadsheets/d/" + settings.EmailSetting.SpreadsheetId,
		}
	}

	if settings.OutputTemplate != nil {
		outputTemplate = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.OutputTemplate.SettingName,
			SpreadSheetUrl: "https://docs.google.com/spreadsheets/d/" + settings.OutputTemplate.SpreadsheetId,
		}
	}

	if settings.OutputTemplateForTeacher != nil {
		outputTemplateForTeacher = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.OutputTemplateForTeacher.SettingName,
			SpreadSheetUrl: "https://docs.google.com/spreadsheets/d/" + settings.OutputTemplateForTeacher.SpreadsheetId,
		}
	}

	if settings.SignUpButton1 != nil {
		signUpButton1 = &response.GetSettingsResponseTextButton{
			SettingName: settings.SignUpButton1.SettingName,
			Name:        settings.SignUpButton1.Name,
			Value:       settings.SignUpButton1.Value,
		}
	}

	if settings.SignUpButton2 != nil {
		signUpButton2 = &response.GetSettingsResponseTextButton{
			SettingName: settings.SignUpButton2.SettingName,
			Name:        settings.SignUpButton2.Name,
			Value:       settings.SignUpButton2.Value,
		}
	}

	if settings.SignUpButton3 != nil {
		signUpButton3 = &response.GetSettingsResponseTextButton{
			SettingName: settings.SignUpButton3.SettingName,
			Name:        settings.SignUpButton3.Name,
			Value:       settings.SignUpButton3.Value,
		}
	}

	if settings.SignUpButton4 != nil {
		signUpButton4 = &response.GetSettingsResponseTextButton{
			SettingName: settings.SignUpButton4.SettingName,
			Name:        settings.SignUpButton4.Name,
			Value:       settings.SignUpButton4.Value,
		}
	}

	if settings.SignUpButton5 != nil {
		signUpButton5 = &response.GetSettingsResponseTextButton{
			SettingName: settings.SignUpButton5.SettingName,
			Name:        settings.SignUpButton5.Name,
			Value:       settings.SignUpButton5.Value,
		}
	}

	if settings.SignUpButtonConfiguration != nil {
		signUpButtonConfiguration = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.SignUpButtonConfiguration.SettingName,
			SpreadSheetUrl: settings.SignUpButtonConfiguration.SpreadsheetId,
		}
	}

	if settings.RegistrationForm != nil {
		registrationForm = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.RegistrationForm.SettingName,
			SpreadSheetUrl: settings.RegistrationForm.SpreadsheetId,
		}
	}

	if settings.RegistrationSubmission != nil {
		registrationSubmission = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.RegistrationSubmission.SettingName,
			SpreadSheetUrl: settings.RegistrationSubmission.SpreadsheetId,
		}
	}

	if settings.RegistrationPreset2 != nil {
		registrationPreset2 = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.RegistrationPreset2.SettingName,
			SpreadSheetUrl: settings.RegistrationPreset2.SpreadsheetId,
		}
	}

	if settings.APIDistributer != nil {
		apiDistributer = &response.GetSettingsResponseAPIDistributor{
			SettingName:    settings.APIDistributer.SettingName,
			SpreadSheetUrl: settings.APIDistributer.SpreadSheetUrl,
		}
	}

	if settings.CodeCountingData != nil {
		codeCountingSetting = &response.GetSettingsResponseAPIDistributor{
			SettingName:    settings.CodeCountingData.SettingName,
			SpreadSheetUrl: settings.CodeCountingData.SpreadSheetUrl,
		}
	}

	if settings.SignUpForms != nil {
		signupformsSettings = &response.GetSettingsResponseDataImport{
			SettingName:    settings.SignUpForms.SettingName,
			SpreadSheetUrl: settings.SignUpForms.SpreadSheetUrl,
			Auto:           settings.SignUpForms.AutoImport,
			Interval:       settings.SignUpForms.Interval,
		}
	}

	if settings.RegistrationPreset1 != nil {
		registrationPreset1 = &response.GetSettingsResponseDataSummary{
			SettingName:    settings.RegistrationPreset1.SettingName,
			SpreadSheetUrl: settings.RegistrationPreset1.SpreadsheetId,
		}
	}

	context.JSON(http.StatusOK, response.GetSettingsResponse{
		Data: response.GetSettingsResponseData{
			ImportFormsSetting:        formSettings,
			ImportFormsSetting2:       formSettings2,
			ImportFormsSetting3:       formSettings3,
			ImportFormsSetting4:       formSettings4,
			ImportRedirectUrlsSetting: urlSettings,
			Output:                    outputSettings,
			Summary:                   summarySettings,
			SyncDevices:               syncDevicesSettings,
			ImportToDoListSetting:     syncToDosSettings,
			EmailHistory:              emailSettings,
			OutputTemplate:            outputTemplate,
			OutputTemplateForTeacher:  outputTemplateForTeacher,
			SignUpButton1:             signUpButton1,
			SignUpButton2:             signUpButton2,
			SignUpButton3:             signUpButton3,
			SignUpButton4:             signUpButton4,
			SignUpButton5:             signUpButton5,
			SignUpButtonConfiguration: signUpButtonConfiguration,
			RegistrationForm:          registrationForm,
			RegistrationSubmission:    registrationSubmission,
			RegistrationPreset2:       registrationPreset2,
			APIDistributer:            apiDistributer,
			CodeCountingData:          codeCountingSetting,
			SignUpFormsSetting:        signupformsSettings,
			RegistrationPreset1:       registrationPreset1,
		},
	})
}

// Create or Update Output Setting godoc
// @Summary      Create or Update Output Setting
// @Description  Create or Update Output Setting
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.UpdateOutputSubmissionSettingsRequest true "body"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/settings/output-sheet [post]
func (receiver *SettingController) UpdateOutputSubmissionSettings(context *gin.Context) {
	var req request.UpdateOutputSubmissionSettingsRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	err := receiver.UpdateSubmissionSetting(req.FolderUrl, req.SheetName)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Setting was saved successfully",
	})
}

// Create or Update Output Summary godoc
// @Summary      Create or Update Output Summary
// @Description  Create or Update Output Summary
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.UpdateOutputSummarySettingsRequest true "body"
// @Success      200  {object}  response.SucceedResponse
// @Failure      400  {object}  response.FailedResponse
// @Failure      404  {object}  response.FailedResponse
// @Failure      500  {object}  response.FailedResponse
// @Router       /v1/admin/settings/output-summary [post]
func (receiver *SettingController) UpdateOutputSummarySettings(context *gin.Context) {
	var req request.UpdateOutputSummarySettingsRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	err := receiver.UpdateOutputSummarySetting(req.SpreadsheetUrl)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Setting was saved successfully",
	})
}

// Update Email History Settings godoc
// @Summary      Update Email History Settings
// @Description  Update Email History Settings
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.UpdateEmailHistorySettingsRequest true "Update Email History Setting Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/email-history [post]
func (receiver *SettingController) UpdateEmailHistorySettings(context *gin.Context) {
	var req request.UpdateEmailHistorySettingsRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	err := receiver.UpdateEmailHistorySettingUseCase.Execute(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Email Setting was saved successfully",
	})
}

// Update Output Template Settings for Teacher godoc
// @Summary      Update Output Template Settings for Teacher godoc
// @Description  Update Output Template Settings for Teacher godoc
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.UpdateOutputTemplateRequest true "Update Output Template Settings godoc"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/output-template [post]
func (receiver *SettingController) UpdateOutputTemplateSettings(context *gin.Context) {
	var req request.UpdateOutputTemplateRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	err := receiver.UpdateOutputTemplateSettingUseCase.Execute(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Output Template Setting was saved successfully",
	})
}

// Update Output Template Settings godoc
// @Summary      Update Output Template Settings for Teacher godoc
// @Description  Update Output Template Settings for Teacher godoc
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.UpdateOutputTemplateRequest true "Update Output Template Settings for Teacher godoc"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/output-template-teacher [post]
func (receiver *SettingController) UpdateOutputTemplateSettingsForTeacher(context *gin.Context) {
	var req request.UpdateOutputTemplateRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	err := receiver.UpdateOutputTemplateSettingForTeacherUseCase.Execute(req)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Output Template Setting was saved successfully",
	})
}

type updateSignUpTextButtonRequest struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// Update Sign Up Button 1 godoc
// @Summary      Update Sign Up Button 1
// @Description  Update Sign Up Button 1
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateSignUpTextButtonRequest true "Update Sign Up Button 1"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/sign-up-button-1 [post]
func (receiver *SettingController) UpdateSignUpButton1(context *gin.Context) {
	var req updateSignUpTextButtonRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateSignUpButton1(req.Name, req.Value)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Sign Up Button was saved successfully",
	})
}

// Update Sign Up Button 2 godoc
// @Summary      Update Sign Up Button 2
// @Description  Update Sign Up Button 2
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateSignUpTextButtonRequest true "Update Sign Up Button 2"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/sign-up-button-2 [post]
func (receiver *SettingController) UpdateSignUpButton2(context *gin.Context) {
	var req updateSignUpTextButtonRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateSignUpButton2(req.Name, req.Value)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Sign Up Button was saved successfully",
	})
}

// Update Sign Up Button 3 godoc
// @Summary      Update Sign Up Button 3
// @Description  Update Sign Up Button 3
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateSignUpTextButtonRequest true "Update Sign Up Button 3"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/sign-up-button-3 [post]
func (receiver *SettingController) UpdateSignUpButton3(context *gin.Context) {
	var req updateSignUpTextButtonRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateSignUpButton3(req.Name, req.Value)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Sign Up Button was saved successfully",
	})
}

// Update Sign Up Button 4 godoc
// @Summary      Update Sign Up Button 4
// @Description  Update Sign Up Button 4
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateSignUpTextButtonRequest true "Update Sign Up Button 4"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/sign-up-button-4 [post]
func (receiver *SettingController) UpdateSignUpButton4(context *gin.Context) {
	var req updateSignUpTextButtonRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateSignUpButton4(req.Name, req.Value)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Sign Up Button was saved successfully",
	})
}

// Update Sign Up Button 5 godoc
// @Summary      Update Sign Up Button 5
// @Description  Update Sign Up Button 5
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateSignUpTextButtonRequest true "Update Sign Up Button 5"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/sign-up-button-5 [post]
func (receiver *SettingController) UpdateSignUpButton5(context *gin.Context) {
	var req updateSignUpTextButtonRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateSignUpButton5(req.Name, req.Value)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Sign Up Button was saved successfully",
	})
}

type updateRegistrationSpreadsheetRequest struct {
	SpreadSheetUrl string `json:"spreadsheet_url" binding:"required"`
}

// Update Registration Form godoc
// @Summary      Update Registration Form
// @Description  Update Registration Form
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateRegistrationSpreadsheetRequest true "Update Registration Form"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/registration-form [post]
func (receiver *SettingController) UpdateRegistrationForm(context *gin.Context) {
	var req updateRegistrationSpreadsheetRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateRegistrationForm(req.SpreadSheetUrl)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Registration Form was saved successfully",
	})
}

type updateSignUpButtonConfigurationRequest struct {
	SpreadSheetUrl string `json:"spreadsheet_url" binding:"required"`
}

// Update SignUpButtonConfiguration Form godoc
// @Summary      Update SignUpButtonConfiguration Form
// @Description  Update SignUpButtonConfiguration Form
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateSignUpButtonConfigurationRequest true "Update SignUpButtonConfiguration Form"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/sign-up-button-configuration [post]
func (receiver *SettingController) UpdateSignUpButtonConfiguration(context *gin.Context) {
	var req updateSignUpButtonConfigurationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateSignUpButtonConfiguration(req.SpreadSheetUrl)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "sign up buttons configuration was saved successfully",
	})
}

// Update Registration Submission godoc
// @Summary      Update Registration Submission
// @Description  Update Registration Submission
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateRegistrationSpreadsheetRequest true "Update Registration Submission"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/registration-submission [post]
func (receiver *SettingController) UpdateRegistrationSubmission(context *gin.Context) {
	var req updateRegistrationSpreadsheetRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateRegistrationSubmission(req.SpreadSheetUrl)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Registration Submission was saved successfully",
	})
}

type updateRegistrationPresetRequest struct {
	FormNote string `json:"spreadsheet_url" binding:"required"`
}

// Update Registration Preset godoc
// @Summary      Update Registration Preset
// @Description  Update Registration Preset
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateRegistrationPresetRequest true "Update Registration Preset"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/registration-preset-2 [post]
func (receiver *SettingController) UpdateRegistrationPreset2(context *gin.Context) {
	var req updateRegistrationPresetRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateRegistrationPreset2(req.FormNote)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Registration Preset was saved successfully",
	})
}

// Update Registration Preset godoc
// @Summary      Update Registration Preset
// @Description  Update Registration Preset
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateRegistrationPresetRequest true "Update Registration Preset"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/registration-preset-1 [post]
func (receiver *SettingController) UpdateRegistrationPreset1(context *gin.Context) {
	var req updateRegistrationPresetRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})
		return
	}

	err := receiver.AdminSignUpUseCases.UpdateRegistrationPreset1(req.FormNote)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Registration Preset was saved successfully",
	})
}

type updateDistributerRequest struct {
	Spreadsheet string `json:"spreadsheet_url" binding:"required"`
}

// Update Distributer Preset godoc
// @Summary      Update Distributer Preset
// @Description  Update Distributer Preset
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body updateDistributerRequest true "Update Distributor Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/api-distributer [post]
func (receiver *SettingController) UpdateAPIDistributor(context *gin.Context) {
	var req updateDistributerRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})

		return
	}

	err := receiver.UpdateApiDistributorUseCase.Execute(req.Spreadsheet)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "APIDistributer was saved successfully",
	})
}

// Update Set Setting Label Name godoc
// @Summary      Set Setting Label Name
// @Description  Set Setting Label Name
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.UpdateSettingNameRequest true "Update Setting Name Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/label/name [post]
func (receiver *SettingController) SetSettingNames(context *gin.Context) {
	var req request.UpdateSettingNameRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})

		return
	}

	if err := receiver.UpdateSettingNameUseCase.Execute(req); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Setting name has been update successfully",
	})
}

// Update Code Counting Data godoc
// @Summary      Update Code Counting Data
// @Description  Update Code Counting Data
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.UpdateCodeCountingSettingRequest true "Update Code Counting Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/code-counting-data [post]
func (receiver *SettingController) UpdateCodeCountingData(context *gin.Context) {
	var req request.UpdateCodeCountingSettingRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})

		return
	}

	if err := usecase.UpdateCodeCountingDataUseCase(req); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Code Counting Data has been update successfully",
	})
}

// Setup Logo Refresh Interval godoc
// @Summary      Setup Logo Refresh Interval
// @Description  Setup Logo Refresh Interval
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Param request body request.SetupLogoRefreshIntervalRequest true "Setup Logo Refresh Interval Request"
// @Success 200 {object} response.SucceedResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/logo-refresh-interval [post]
func (receiver *SettingController) SetupLogoRefreshInterval(context *gin.Context) {
	var req request.SetupLogoRefreshIntervalRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, response.FailedResponse{
			Code:  http.StatusBadRequest,
			Error: "Invalid request body",
		})

		return
	}

	if err := usecase.SetupLogoRefreshIntervalUseCase(req); err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Code:    http.StatusOK,
		Message: "Logo refresh interval has been setup successfully",
	})
}

type getLogoRefreshIntervalResponse struct {
	Interval uint64 `json:"interval" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

// Get Logo Refresh Interval godoc
// @Summary      Get Logo Refresh Interval
// @Description  Get Logo Refresh Interval
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} getLogoRefreshIntervalResponse
// @Failure 400 {object} response.FailedResponse
// @Failure 401 {object} response.FailedResponse
// @Failure 403 {object} response.FailedResponse
// @Failure 500 {object} response.FailedResponse
// @Router /v1/admin/settings/logo-refresh-interval [get]
func (receiver *SettingController) GetLogoRefreshInterval(context *gin.Context) {
	interval, err := usecase.GetLogoRefreshIntervalUseCase()
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.FailedResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, response.SucceedResponse{
		Data: getLogoRefreshIntervalResponse{Interval: interval.IntegerValue, Title: interval.SettingName},
	})
}
