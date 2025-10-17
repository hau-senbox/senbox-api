package response

type GetAllPreRegister4Web struct {
	Email      string `json:"email"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	FormQr     string `json:"form_qr"`
	FormName   string `json:"form_name"`
	FormSheet  string `json:"form_sheet"`
	CreatedAt  string `json:"created_at"`
}

type GetAllPreRegister4App struct {
	Email      string `json:"email"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	FormQr     string `json:"form_qr"`
	FormName   string `json:"form_name"`
	FormSheet  string `json:"form_sheet"`
	CreatedAt  string `json:"created_at"`
}
