package request

type SendMailRequest struct {
	MacAddress string `json:"macAddress"`
	Content    string `json:"content"`
}
