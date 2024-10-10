package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var accessToken = os.Getenv("TELEGRAM_SENBOT_ACCESS_TOKEN")
var chatId = os.Getenv("TELEGRAM_SENBOT_CHAT_ID")

func SendMessageViaTelegram(message ...string) {
	msg := strings.Join(message, "\n")
	_, err := send(msg)
	if err != nil {
		log.Error(err)
	}
}

func getUrl() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", accessToken)
}

func send(text string) (bool, error) {
	// Global variables
	var err error
	var response *http.Response

	// Send the message
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": chatId,
		"text":    text,
	})
	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Error("Error when send message to telegram: ", err)
		return false, err
	}

	// Close the request at the end
	defer response.Body.Close()

	// Body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Error when read response body: ", err)
		return false, err
	}

	// Log
	log.Infof("Message '%s' was sent", text)
	log.Infof("Response JSON: %s", string(body))

	// Return
	return true, nil
}
