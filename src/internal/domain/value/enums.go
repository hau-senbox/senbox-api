package value

import (
	"errors"
	"strings"
	"time"
)

type QuestionType int

const (
	QuestionDate QuestionType = iota + 1
	QuestionTime
	QuestionDateTime
	QuestionDurationForward
	QuestionDurationBackward
	QuestionScale
	QuestionQRCode
	QuestionSelection
	QuestionInText
	QuestionCount
	QuestionNumber
	QuestionPhoto
	QuestionMultipleChoice
	QuestionButtonCount
	QuestionSingleChoice

	QuestionButtonList
	QuestionMessageBox
	QuestionShowPic
	QuestionButton
	QuestionPlayVideo
	QuestionQRCodeFront
	QuestionChoiceToggle
	QuestionSection

	QuestionDateUser
	QuestionTimeUser
	QuestionDateTimeUser
	QuestionDurationForwardUser
	QuestionDurationBackwardUser
	QuestionScaleUser
	QuestionQRCodeUser
	QuestionSelectionUser
	QuestionTextUser
	QuestionCountUser
	QuestionNumberUser
	QuestionPhotoUser
	QuestionMultipleChoiceUser
	QuestionButtonCountUser
	QuestionSingleChoiceUser
	QuestionButtonListUser
	QuestionMessageBoxUser
	QuestionShowPicUser
	QuestionButtonUser
	QuestionPlayVideoUser
	QuestionQRCodeFrontUser
	QuestionChoiceToggleUser
	QuestionFormSection
	QuestionFormSendImmediately
	QuestionSignature
	QuestionWeb
	QuestionSignUpPreSetValue1
	QuestionSignUpPreSetValue3
	QuestionDraggableList
	QuestionSendMessage
	QuestionSendNotification
	QuestionCodeCounting
	QuestionRandomizer
	QuestionDocument
	QuestionQRCodeGenerator
	QuestionSignUpPreSetValue2
	QuestionWebUser

	QuestionPresetNickname
	QuestionPresetEmail
	QuestionPresetDob
	QuestionPresetPassword
	QuestionPresetRole
	QuestionPresetConditionAccept
	QuestionPresetRoleSelectWorkingAddress

	SignUpButtonConfiguration1
	SignUpButtonConfiguration2
	SignUpButtonConfiguration3
	SignUpButtonConfiguration4
	SignUpButtonConfiguration5
	SignUpButtonConfiguration6
	SignUpButtonConfiguration7
	SignUpButtonConfiguration8
	SignUpButtonConfiguration9
	SignUpButtonConfiguration10

	UserInformationValue1
	UserInformationValue2
	UserInformationValue3
	UserInformationValue4
	UserInformationValue5
	UserInformationValue6
	UserInformationValue7

	CameraSquareLens
	SubmitText

	MessageText1
	MessageText2
	ResponseText1
	ResponseText2

	PdfViewer
	PdfPicker

	OrganizationName
	ApplicationContent
	WaterCup

	OutNrTotal
)

type DeviceStatus int

const (
	DeviceStatus_Suspend DeviceStatus = iota + 1
	DeviceStatus_ModeT
	DeviceStatus_ModeP
	DeviceStatus_ModeS
	DeviceStatus_Deactive
	DeviceStatus_ModeL
)

func GetDeviceStatusFromString(status string) (DeviceStatus, error) {
	switch strings.ToLower(status) {
	case "suspended":
		return DeviceStatus_Suspend, nil
	case "mode t":
		return DeviceStatus_ModeT, nil
	case "mode p":
		return DeviceStatus_ModeP, nil
	case "mode s":
		return DeviceStatus_ModeS, nil
	case "deactivated":
		return DeviceStatus_Deactive, nil
	case "mode l":
		return DeviceStatus_ModeL, nil
	default:
		return 0, errors.New("invalid device status")
	}
}

func GetDeviceModeFromString(status string) (DeviceMode, error) {
	switch strings.ToLower(status) {
	case "suspended":
		return DeviceModeSuspended, nil
	case "mode t":
		return DeviceModeT, nil
	case "mode p":
		return DeviceModeP, nil
	case "mode s":
		return DeviceModeS, nil
	case "deactivated":
		return DeviceModeDeactivated, nil
	case "mode l":
		return DeviceModeL, nil
	default:
		return "", errors.New("invalid device status")
	}
}

func GetDeviceStatusStringAtMode(status DeviceMode) string {
	switch status {
	case DeviceModeSuspended:
		return "suspended"
	case DeviceModeT:
		return "mode t"
	case DeviceModeP:
		return "mode p"
	case DeviceModeS:
		return "mode s"
	case DeviceModeDeactivated:
		return "deactivated"
	case DeviceModeL:
		return "mode l"
	default:
		return ""
	}
}

type Status int

const (
	Inactive Status = iota
	Active
	PendingDelete
)

// GetRawValue return question type in string from its enum value
func GetRawValue(questionType QuestionType) string {
	switch questionType {
	case QuestionDate:
		return "date"
	case QuestionTime:
		return "time"
	case QuestionDateTime:
		return "datetime"
	case QuestionDurationForward:
		return "duration_forward"
	case QuestionDurationBackward:
		return "duration_backward"
	case QuestionScale:
		return "scale"
	case QuestionQRCode:
		return "qr_code"
	case QuestionSelection:
		return "selection"
	case QuestionInText:
		return "in_text"
	case QuestionCount:
		return "count"
	case QuestionNumber:
		return "number"
	case QuestionPhoto:
		return "photo"
	case QuestionMultipleChoice:
		return "multiple_choice"
	case QuestionButtonCount:
		return "button_count"
	case QuestionSingleChoice:
		return "single_choice"
	case QuestionButtonList:
		return "button_list"
	case QuestionMessageBox:
		return "message_box"
	case QuestionShowPic:
		return "show_pic"
	case QuestionButton:
		return "button"
	case QuestionPlayVideo:
		return "play_video"
	case QuestionQRCodeFront:
		return "qr_code_front"
	case QuestionChoiceToggle:
		return "choice_toggle"
	case QuestionSection:
		return "section"

	case QuestionSignUpPreSetValue1:
		return "preset_value1"
	case QuestionSignUpPreSetValue3:
		return "preset_value3"
	case QuestionSignUpPreSetValue2:
		return "preset_value2"

	case QuestionPresetNickname:
		return "preset_nickname"
	case QuestionPresetEmail:
		return "preset_email"
	case QuestionPresetDob:
		return "preset_dob"
	case QuestionPresetPassword:
		return "preset_password"
	case QuestionPresetConditionAccept:
		return "preset_condition_accept"
	case QuestionPresetRole:
		return "preset_role"
	case QuestionPresetRoleSelectWorkingAddress:
		return "preset_role_select_working_address"

	case SignUpButtonConfiguration1:
		return "preset_sign_up_button_1"
	case SignUpButtonConfiguration2:
		return "preset_sign_up_button_2"
	case SignUpButtonConfiguration3:
		return "preset_sign_up_button_3"
	case SignUpButtonConfiguration4:
		return "preset_sign_up_button_4"
	case SignUpButtonConfiguration5:
		return "preset_sign_up_button_5"
	case SignUpButtonConfiguration6:
		return "preset_sign_up_button_6"
	case SignUpButtonConfiguration7:
		return "preset_sign_up_button_7"
	case SignUpButtonConfiguration8:
		return "preset_sign_up_button_8"
	case SignUpButtonConfiguration9:
		return "preset_sign_up_button_9"
	case SignUpButtonConfiguration10:
		return "preset_sign_up_button_10"

	case UserInformationValue1:
		return "user_information_value_1"
	case UserInformationValue2:
		return "user_information_value_2"
	case UserInformationValue3:
		return "user_information_value_3"
	case UserInformationValue4:
		return "user_information_value_4"
	case UserInformationValue5:
		return "user_information_value_5"
	case UserInformationValue6:
		return "user_information_value_6"
	case UserInformationValue7:
		return "user_information_value_7"

	case QuestionFormSection:
		return "form_section"
	case QuestionFormSendImmediately:
		return "send_immediately"
	case QuestionSignature:
		return "signature"
	case QuestionWeb:
		return "web"
	case QuestionWebUser:
		return "web_user"
	case QuestionDraggableList:
		return "draggable_list"
	case QuestionSendMessage:
		return "send_message"
	case QuestionCodeCounting:
		return "code_counting"
	case QuestionRandomizer:
		return "randomizer"
	case QuestionDocument:
		return "document"
	case QuestionQRCodeGenerator:
		return "qrcode_generator"
	case QuestionDateUser:
		return "date_user"
	case QuestionTimeUser:
		return "time_user"
	case QuestionDateTimeUser:
		return "datetime_user"
	case QuestionDurationForwardUser:
		return "duration_forward_user"
	case QuestionDurationBackwardUser:
		return "duration_backward_user"
	case QuestionScaleUser:
		return "scale_user"
	case QuestionQRCodeUser:
		return "qr_code_user"
	case QuestionSelectionUser:
		return "selection_user"
	case QuestionTextUser:
		return "text_user"
	case QuestionCountUser:
		return "count_user"
	case QuestionNumberUser:
		return "number_user"
	case QuestionPhotoUser:
		return "photo_user"
	case QuestionMultipleChoiceUser:
		return "multiple_choice_user"
	case QuestionButtonCountUser:
		return "button_count_user"
	case QuestionSingleChoiceUser:
		return "single_choice_user"
	case QuestionButtonListUser:
		return "button_list_user"
	case QuestionMessageBoxUser:
		return "message_box_user"
	case QuestionShowPicUser:
		return "show_pic_user"
	case QuestionButtonUser:
		return "button_user"
	case QuestionPlayVideoUser:
		return "play_video_user"
	case QuestionQRCodeFrontUser:
		return "qr_code_front_user"
	case QuestionChoiceToggleUser:
		return "choice_toggle_user"
	case QuestionSendNotification:
		return "send_notification"

	case CameraSquareLens:
		return "camera_square_lens"
	case SubmitText:
		return "submit_text"

	case MessageText1:
		return "message_text_1"
	case MessageText2:
		return "message_text_2"
	case ResponseText1:
		return "response_text_1"
	case ResponseText2:
		return "response_text_2"

	case PdfViewer:
		return "pdf_viewer"
	case PdfPicker:
		return "pdf_picker"

	case OrganizationName:
		return "organization_name"
	case ApplicationContent:
		return "application_content"
	case WaterCup:
		return "water_cup"
	case OutNrTotal:
		return "out_nr_total"
	}

	return ""
}

func GetStringValue(questionType QuestionType) string {
	switch questionType {
	case QuestionDate:
		return "date"
	case QuestionTime:
		return "time"
	case QuestionDateTime:
		return "datetime"
	case QuestionDurationForward:
		return "duration_forward"
	case QuestionDurationBackward:
		return "duration_backward"
	case QuestionScale:
		return "scale"
	case QuestionQRCode:
		return "qr_code"
	case QuestionSelection:
		return "selection"
	case QuestionInText:
		return "in_text"
	case QuestionCount:
		return "count"
	case QuestionNumber:
		return "number"
	case QuestionPhoto:
		return "photo"
	case QuestionMultipleChoice:
		return "multiple_choice"
	case QuestionButtonCount:
		return "button_count"
	case QuestionSingleChoice:
		return "single_choice"
	case QuestionButtonList:
		return "button_list"
	case QuestionMessageBox:
		return "message_box"
	case QuestionShowPic:
		return "show_pic"
	case QuestionButton:
		return "button"
	case QuestionPlayVideo:
		return "play_video"
	case QuestionQRCodeFront:
		return "qr_code_front"
	case QuestionChoiceToggle:
		return "choice_toggle"
	case QuestionSection:
		return "section"
	case QuestionFormSection:
		return "form_section"
	case QuestionFormSendImmediately:
		return "send_immediately"
	case QuestionSignature:
		return "signature"
	case QuestionWeb:
		return "web"
	case QuestionWebUser:
		return "web_user"

	case QuestionSignUpPreSetValue1:
		return "preset_value1"
	case QuestionSignUpPreSetValue3:
		return "preset_value3"
	case QuestionSignUpPreSetValue2:
		return "preset_value2"

	case QuestionPresetNickname:
		return "preset_nickname"
	case QuestionPresetEmail:
		return "preset_email"
	case QuestionPresetDob:
		return "preset_dob"
	case QuestionPresetPassword:
		return "preset_password"
	case QuestionPresetConditionAccept:
		return "preset_condition_accept"
	case QuestionPresetRole:
		return "preset_role"
	case QuestionPresetRoleSelectWorkingAddress:
		return "preset_role_select_working_address"

	case SignUpButtonConfiguration1:
		return "preset_sign_up_button_1"
	case SignUpButtonConfiguration2:
		return "preset_sign_up_button_2"
	case SignUpButtonConfiguration3:
		return "preset_sign_up_button_3"
	case SignUpButtonConfiguration4:
		return "preset_sign_up_button_4"
	case SignUpButtonConfiguration5:
		return "preset_sign_up_button_5"
	case SignUpButtonConfiguration6:
		return "preset_sign_up_button_6"
	case SignUpButtonConfiguration7:
		return "preset_sign_up_button_7"
	case SignUpButtonConfiguration8:
		return "preset_sign_up_button_8"
	case SignUpButtonConfiguration9:
		return "preset_sign_up_button_9"
	case SignUpButtonConfiguration10:
		return "preset_sign_up_button_10"

	case UserInformationValue1:
		return "user_information_value_1"
	case UserInformationValue2:
		return "user_information_value_2"
	case UserInformationValue3:
		return "user_information_value_3"
	case UserInformationValue4:
		return "user_information_value_4"
	case UserInformationValue5:
		return "user_information_value_5"
	case UserInformationValue6:
		return "user_information_value_6"
	case UserInformationValue7:
		return "user_information_value_7"

	case QuestionDraggableList:
		return "draggable_list"
	case QuestionSendMessage:
		return "send_message"
	case QuestionSendNotification:
		return "send_notification"
	case QuestionCodeCounting:
		return "code_counting"
	case QuestionRandomizer:
		return "randomizer"
	case QuestionDocument:
		return "document"
	case QuestionQRCodeGenerator:
		return "qrcode_generator"

	case CameraSquareLens:
		return "camera_square_lens"
	case SubmitText:
		return "submit_text"

	case MessageText1:
		return "message_text_1"
	case MessageText2:
		return "message_text_2"
	case ResponseText1:
		return "response_text_1"
	case ResponseText2:
		return "response_text_2"

	case PdfViewer:
		return "pdf_viewer"
	case PdfPicker:
		return "pdf_picker"

	case OrganizationName:
		return "organization_name"
	case ApplicationContent:
		return "application_content"
	case WaterCup:
		return "water_cup"
	case OutNrTotal:
		return "out_nr_total"

	default:
		return ""
	}
}

func GetQuestionType(rawValue string) (QuestionType, error) {
	var lowerCaseValue = strings.ToLower(rawValue)
	switch lowerCaseValue {
	case "date":
		return QuestionDate, nil
	case "time":
		return QuestionTime, nil
	case "datetime":
		return QuestionDateTime, nil
	case "duration_forward":
		return QuestionDurationForward, nil
	case "duration_backward":
		return QuestionDurationBackward, nil
	case "scale":
		return QuestionScale, nil
	case "qr_code":
		return QuestionQRCode, nil
	case "selection":
		return QuestionSelection, nil
	case "in_text":
		return QuestionInText, nil
	case "count":
		return QuestionCount, nil
	case "number":
		return QuestionNumber, nil
	case "photo":
		return QuestionPhoto, nil
	case "multiple_choice":
		return QuestionMultipleChoice, nil
	case "button_count":
		return QuestionButtonCount, nil
	case "single_choice":
		return QuestionSingleChoice, nil
	case "button_list":
		return QuestionButtonList, nil
	case "message_box":
		return QuestionMessageBox, nil
	case "show_pic":
		return QuestionShowPic, nil
	case "button":
		return QuestionButton, nil
	case "play_video":
		return QuestionPlayVideo, nil
	case "qr_code_front":
		return QuestionQRCodeFront, nil
	case "choice_toggle":
		return QuestionChoiceToggle, nil
	case "section":
		return QuestionSection, nil

	case "preset_value1":
		return QuestionSignUpPreSetValue1, nil
	case "preset_value3":
		return QuestionSignUpPreSetValue3, nil
	case "preset_value2":
		return QuestionSignUpPreSetValue2, nil

	case "preset_nickname":
		return QuestionPresetNickname, nil
	case "preset_email":
		return QuestionPresetEmail, nil
	case "preset_dob":
		return QuestionPresetDob, nil
	case "preset_password":
		return QuestionPresetPassword, nil
	case "preset_condition_accept":
		return QuestionPresetConditionAccept, nil
	case "preset_role":
		return QuestionPresetRole, nil
	case "preset_role_select_working_address":
		return QuestionPresetRoleSelectWorkingAddress, nil

	case "preset_sign_up_button_1":
		return SignUpButtonConfiguration1, nil
	case "preset_sign_up_button_2":
		return SignUpButtonConfiguration2, nil
	case "preset_sign_up_button_3":
		return SignUpButtonConfiguration3, nil
	case "preset_sign_up_button_4":
		return SignUpButtonConfiguration4, nil
	case "preset_sign_up_button_5":
		return SignUpButtonConfiguration5, nil
	case "preset_sign_up_button_6":
		return SignUpButtonConfiguration6, nil
	case "preset_sign_up_button_7":
		return SignUpButtonConfiguration7, nil
	case "preset_sign_up_button_8":
		return SignUpButtonConfiguration8, nil
	case "preset_sign_up_button_9":
		return SignUpButtonConfiguration9, nil
	case "preset_sign_up_button_10":
		return SignUpButtonConfiguration10, nil

	case "user_information_value_1":
		return UserInformationValue1, nil
	case "user_information_value_2":
		return UserInformationValue2, nil
	case "user_information_value_3":
		return UserInformationValue3, nil
	case "user_information_value_4":
		return UserInformationValue4, nil
	case "user_information_value_5":
		return UserInformationValue5, nil
	case "user_information_value_6":
		return UserInformationValue6, nil
	case "user_information_value_7":
		return UserInformationValue7, nil

	case "form_section":
		return QuestionFormSection, nil
	case "send_immediately":
		return QuestionFormSendImmediately, nil
	case "signature":
		return QuestionSignature, nil
	case "web":
		return QuestionWeb, nil
	case "web_user":
		return QuestionWebUser, nil
	case "draggable_list":
		return QuestionDraggableList, nil
	case "send_message":
		return QuestionSendMessage, nil
	case "date_user":
		return QuestionDateUser, nil
	case "time_user":
		return QuestionTimeUser, nil
	case "datetime_user":
		return QuestionDateTimeUser, nil
	case "duration_forward_user":
		return QuestionDurationForwardUser, nil
	case "duration_backward_user":
		return QuestionDurationBackwardUser, nil
	case "scale_user":
		return QuestionScaleUser, nil
	case "qr_code_user":
		return QuestionQRCodeUser, nil
	case "selection_user":
		return QuestionSelectionUser, nil
	case "text_user":
		return QuestionTextUser, nil
	case "count_user":
		return QuestionCountUser, nil
	case "number_user":
		return QuestionNumberUser, nil
	case "photo_user":
		return QuestionPhotoUser, nil
	case "multiple_choice_user":
		return QuestionMultipleChoiceUser, nil
	case "button_count_user":
		return QuestionButtonCountUser, nil
	case "single_choice_user":
		return QuestionSingleChoiceUser, nil
	case "button_list_user":
		return QuestionButtonListUser, nil
	case "message_box_user":
		return QuestionMessageBoxUser, nil
	case "show_pic_user":
		return QuestionShowPicUser, nil
	case "button_user":
		return QuestionButtonUser, nil
	case "play_video_user":
		return QuestionPlayVideoUser, nil
	case "qr_code_front_user":
		return QuestionQRCodeFrontUser, nil
	case "choice_toggle_user":
		return QuestionChoiceToggleUser, nil
	case "send_notification":
		return QuestionSendNotification, nil
	case "code_counting":
		return QuestionCodeCounting, nil
	case "randomizer":
		return QuestionRandomizer, nil
	case "document":
		return QuestionDocument, nil
	case "qrcode_generator":
		return QuestionQRCodeGenerator, nil

	case "camera_square_lens":
		return CameraSquareLens, nil
	case "submit_text":
		return SubmitText, nil

	case "message_text_1":
		return MessageText1, nil
	case "message_text_2":
		return MessageText2, nil
	case "response_text_1":
		return ResponseText1, nil
	case "response_text_2":
		return ResponseText2, nil

	case "pdf_viewer":
		return PdfViewer, nil
	case "pdf_picker":
		return PdfPicker, nil

	case "organization_name":
		return OrganizationName, nil
	case "application_content":
		return ApplicationContent, nil
	case "water_cup":
		return WaterCup, nil
	case "out_nr_total":
		return OutNrTotal, nil

	default:
		return 0, errors.New("invalid raw value")
	}
}

func IsGeneralQuestionType(questionType QuestionType) bool {
	switch questionType {
	case QuestionDate,
		QuestionTime,
		QuestionDateTime,
		QuestionDurationForward,
		QuestionDurationBackward,
		QuestionScale,
		QuestionQRCode,
		QuestionSelection,
		QuestionInText,
		QuestionCount,
		QuestionNumber,
		QuestionPhoto,
		QuestionMultipleChoice,
		QuestionButtonCount,
		QuestionButtonCountUser,
		QuestionSingleChoice,
		QuestionButtonList,
		QuestionMessageBox,
		QuestionShowPic,
		QuestionButton,
		QuestionPlayVideo,
		QuestionQRCodeFront,
		QuestionChoiceToggle,
		QuestionSection,
		QuestionFormSection,
		QuestionFormSendImmediately,
		QuestionSignature,
		QuestionWeb,

		QuestionSignUpPreSetValue1,
		QuestionSignUpPreSetValue3,
		QuestionSignUpPreSetValue2,

		QuestionDraggableList,
		QuestionSendMessage,
		QuestionSendNotification,
		QuestionCodeCounting,
		QuestionRandomizer,
		QuestionDocument,
		QuestionQRCodeGenerator,

		QuestionPresetNickname,
		QuestionPresetEmail,
		QuestionPresetDob,
		QuestionPresetPassword,
		QuestionPresetConditionAccept,
		QuestionPresetRole,
		QuestionPresetRoleSelectWorkingAddress,

		SignUpButtonConfiguration1,
		SignUpButtonConfiguration2,
		SignUpButtonConfiguration3,
		SignUpButtonConfiguration4,
		SignUpButtonConfiguration5,
		SignUpButtonConfiguration6,
		SignUpButtonConfiguration7,
		SignUpButtonConfiguration8,
		SignUpButtonConfiguration9,
		SignUpButtonConfiguration10,

		UserInformationValue1,
		UserInformationValue2,
		UserInformationValue3,
		UserInformationValue4,
		UserInformationValue5,
		UserInformationValue6,
		UserInformationValue7,

		MessageText1,
		MessageText2,
		ResponseText1,
		ResponseText2,

		PdfViewer,
		PdfPicker,

		CameraSquareLens,
		SubmitText,

		OrganizationName,
		ApplicationContent,
		WaterCup,
		OutNrTotal:
		return true
	default:
		return false
	}
}

type UserInfoInputType int

const (
	UserInfoInputTypeKeyboard   UserInfoInputType = iota
	UserInfoInputTypeBarcode    UserInfoInputType = UserInfoInputTypeKeyboard + 1
	UserInfoInputTypeBackOffice UserInfoInputType = UserInfoInputTypeBarcode + 2
)

func GetRawValueOfUserInfoInputType(userInfoInputType UserInfoInputType) string {
	switch userInfoInputType {
	case UserInfoInputTypeKeyboard:
		return "keyboard"
	case UserInfoInputTypeBarcode:
		return "scanned"
	case UserInfoInputTypeBackOffice:
		return "back_office"
	default:
		return "scanned"
	}
}

func GetUserInfoInputTypeFromString(userInfoInputType string) (UserInfoInputType, error) {
	switch strings.ToLower(userInfoInputType) {
	case "keyboard":
		return UserInfoInputTypeKeyboard, nil
	case "scanned":
		return UserInfoInputTypeBarcode, nil
	case "back_office":
		return UserInfoInputTypeBackOffice, nil
	default:
		return UserInfoInputTypeBarcode, nil
	}
}

func GetRawStatusValue(status Status) string {
	switch status {
	case Active:
		return "active"
	case Inactive:
		return "inactive"
	case PendingDelete:
		return "pending_delete"
	default:
		return "unknown"
	}
}

func GetStatusFromString(status string) (Status, error) {
	switch strings.ToLower(status) {
	case "active":
		return Active, nil
	case "inactive":
		return Inactive, nil
	case "pending_delete":
		return PendingDelete, nil
	default:
		return 0, errors.New("invalid status")
	}
}

type FromApplicationStatus string

const (
	Approved FromApplicationStatus = "approved"
	Blocked  FromApplicationStatus = "blocked"
	Pending  FromApplicationStatus = "pending"
)

func (s FromApplicationStatus) String() string {
	return string(s)
}

type ImportSpreadsheetStatus int

const (
	ImportSpreadsheetStatusPending    ImportSpreadsheetStatus = iota
	ImportSpreadsheetStatusDeleted    ImportSpreadsheetStatus = ImportSpreadsheetStatusPending + 1
	ImportSpreadsheetStatusSkip       ImportSpreadsheetStatus = ImportSpreadsheetStatusPending + 2
	ImportSpreadsheetStatusNew        ImportSpreadsheetStatus = ImportSpreadsheetStatusPending + 3
	ImportSpreadsheetStatusDeactivate ImportSpreadsheetStatus = ImportSpreadsheetStatusPending + 4
)

func GetRawValueOfImportSpreadsheetStatus(importSpreadsheetStatus ImportSpreadsheetStatus) string {
	switch importSpreadsheetStatus {
	case ImportSpreadsheetStatusPending:
		return "pending"
	case ImportSpreadsheetStatusDeleted:
		return "delete"
	case ImportSpreadsheetStatusSkip:
		return "uploaded"
	case ImportSpreadsheetStatusNew:
		return "upload"
	case ImportSpreadsheetStatusDeactivate:
		return "deactivate"
	default:
		return "unknown"
	}
}

func GetImportSpreadsheetStatusFromString(importSpreadsheetStatus string) (ImportSpreadsheetStatus, error) {
	switch strings.ToLower(importSpreadsheetStatus) {
	case "pending":
		return ImportSpreadsheetStatusPending, nil
	case "delete":
		return ImportSpreadsheetStatusDeleted, nil
	case "uploaded":
		return ImportSpreadsheetStatusSkip, nil
	case "upload":
		return ImportSpreadsheetStatusNew, nil
	case "deactivate":
		return ImportSpreadsheetStatusDeactivate, nil
	default:
		return ImportSpreadsheetStatusPending, nil
	}
}

type SettingType int

const (
	SettingTypeSubmission                SettingType = iota + 1
	SettingTypeImportForms                           = SettingTypeSubmission + 1
	SettingTypeImportUrls                            = SettingTypeSubmission + 2
	SettingTypeSummary                               = SettingTypeSubmission + 3
	SettingTypeSyncDevices                           = SettingTypeSubmission + 4
	SettingTypeSyncToDos                             = SettingTypeSubmission + 5
	SettingTypeEmailHistory                          = SettingTypeSubmission + 6
	SettingTypeOutputTemplate                        = SettingTypeSubmission + 7
	SettingTypeOutputTemplateTeacher                 = SettingTypeSubmission + 8
	SettingTypeImportForms2                          = SettingTypeSubmission + 9
	SettingTypeImportForms3                          = SettingTypeSubmission + 10
	SettingTypeImportForms4                          = SettingTypeSubmission + 11
	SettingTypeSignUpButton1                         = SettingTypeSubmission + 12
	SettingTypeSignUpButton2                         = SettingTypeSubmission + 13
	SettingTypeSignUpButton3                         = SettingTypeSubmission + 14
	SettingTypeSignUpButton4                         = SettingTypeSubmission + 15
	SettingTypeSignUpForm                            = SettingTypeSubmission + 16
	SettingTypeSignUpOutput                          = SettingTypeSubmission + 17
	SettingTypeSignUpPresetValue2                    = SettingTypeSubmission + 18
	SettingTypeAPIDistributer                        = SettingTypeSubmission + 19
	SettingTypeCodeCountingData                      = SettingTypeSubmission + 20
	SettingTypeLogoRefreshInterval                   = SettingTypeSubmission + 21
	SettingTypeImportSignUpForms                     = SettingTypeSubmission + 22
	SettingTypeSignUpPresetValue1                    = SettingTypeSubmission + 23
	SettingTypeSignUpButton5                         = SettingTypeSubmission + 24
	SettingTypeSignUpButtonConfiguration             = SettingTypeSubmission + 25
)

type ButtonType int

const (
	ButtonTypeScan ButtonType = iota
	ButtonTypeList ButtonType = ButtonTypeScan + 1
)

func GetButtonTypeFromString(typeInString string) (ButtonType, error) {
	switch strings.ToLower(typeInString) {
	case "scan":
		return ButtonTypeScan, nil
	case "list":
		return ButtonTypeList, nil
	default:
		return 0, errors.New("invalid button type")
	}
}

func GetRawButtonTypeValue(buttonType ButtonType) string {
	switch buttonType {
	case ButtonTypeScan:
		return "scan"
	case ButtonTypeList:
		return "list"
	default:
		return "scan"
	}
}

const (
	WorkingHoursStart = "08:00"
	WorkingHoursEnd   = "20:00"
)

type ToDoType string

const (
	ToDoTypeAssign  ToDoType = "assign"
	ToDoTypeCompose ToDoType = "compose"
)

type FormType string

const (
	FormType_General      FormType = "general"
	FormType_SelfRemember FormType = "self-remember"
)

type QuestionForMobile string

const (
	QuestionForMobile_Enabled  = "enabled"
	QuestionForMobile_Disabled = "disabled"
)

type DeviceType string

const (
	DeviceTypeIOS     = DeviceType("IOS")
	DeviceTypeANDROID = DeviceType("ANDROID")
)

type NotificationType string

const (
	NotificationType_TopButtons                 NotificationType = "top_buttons"
	NotificationType_NewFormSubmit              NotificationType = "new_form_submit"
	NotificationType_LogoRefreshIntervalChanged NotificationType = "logo_refresh_interval_changed"
	NotificationType_UserMessageChanged         NotificationType = "user_message_changed"
	NotificationType_NoteChanged                NotificationType = "note_changed"
	NotificationType_DeviceStatusChanged        NotificationType = "device_status_changed"
)

type FcmTopics string

const (
	FcmTopicsGeneral FcmTopics = "GENERAL"
)

type ScreenButtonType string

const (
	ScreenButtonType_Scan ScreenButtonType = "scan"
	ScreenButtonType_List ScreenButtonType = "list"
)

type InfoInputType string

const (
	InfoInputTypeKeyboard   InfoInputType = "keyboard"
	InfoInputTypeBarcode    InfoInputType = "scanned"
	InfoInputTypeBackOffice InfoInputType = "back_office"
)

func GetInfoInputTypeFromString(input string) InfoInputType {
	switch strings.ToLower(input) {
	case "keyboard":
		return InfoInputTypeKeyboard
	case "bar_code":
		return InfoInputTypeBarcode
	case "backoffice":
		return InfoInputTypeBackOffice
	default:
		return InfoInputTypeBackOffice
	}
}

type DeviceMode string

const (
	DeviceModeSuspended   DeviceMode = "suspended"
	DeviceModeDeactivated DeviceMode = "deactivated"
	DeviceModeT           DeviceMode = "mode t"
	DeviceModeS           DeviceMode = "mode s"
	DeviceModeP           DeviceMode = "mode p"
	DeviceModeL           DeviceMode = "mode l"
)

type DeviceConditionKey string

const (
	DeviceConditionKeyStudent      DeviceConditionKey = "key-student"
	DeviceConditionKeyTeacher      DeviceConditionKey = "key-teacher"
	DeviceConditionKeyOrganization DeviceConditionKey = "key-organization"
)

type TimeSort string

const (
	TimeShortLatest TimeSort = "latest"
	TimeShortOldest TimeSort = "oldest"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}
