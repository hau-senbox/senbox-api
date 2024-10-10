package monitor

// register, import 1 todo, import 1 form, screen button, top button
var (
	TotalRequestInitDevice      = 0
	TotalRequestImportToDo      = 0
	TotalRequestImportForm      = 0
	TotalRequestGETScreenButton = 0
	TotalRequestGETTopButton    = 0
)

func LogGoogleAPIRequestInitDevice() {
	TotalRequestInitDevice++
}

func LogGoogleAPIRequestImportTodo() {
	TotalRequestImportToDo++
}

func LogGoogleAPIRequestImportForm() {
	TotalRequestImportForm++
}

func LogGoogleAPIRequestGETScreenButton() {
	TotalRequestGETScreenButton++
}

func LogGoogleAPIRequestGETTopButton() {
	TotalRequestGETTopButton++
}

func ResetGoogleAPIRequestMonitor() {
	TotalRequestInitDevice = 0
	TotalRequestImportToDo = 0
	TotalRequestImportForm = 0
	TotalRequestGETScreenButton = 0
	TotalRequestGETTopButton = 0
}
