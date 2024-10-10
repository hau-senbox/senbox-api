package response

type GetScreenButtonsItem struct {
	ButtonTitle     string  `json:"title"`
	ButtonValue     string  `json:"value"`
	BackgroundColor *string `json:"background_color"`
}

type GetScreenButtonsResponseData struct {
	Buttons []GetScreenButtonsItem `json:"buttons"`
}

type GetScreenButtonsResponse struct {
	Data GetScreenButtonsResponseData `json:"data"`
}
