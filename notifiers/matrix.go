package notifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/statping-ng/statping-ng/types/failures"
	"github.com/statping-ng/statping-ng/types/notifications"
	"github.com/statping-ng/statping-ng/types/notifier"
	"github.com/statping-ng/statping-ng/types/null"
	"github.com/statping-ng/statping-ng/types/services"
	"github.com/statping-ng/statping-ng/utils"
)

var _ notifier.Notifier = (*matrix)(nil)

type matrix struct {
	*notifications.Notification
}

var Matrix = &matrix{&notifications.Notification{
	Method:      "matrix",
	Title:       "Matrix",
	Description: "An open network for secure, decentralized communication.",
	Author:      "jojo",
	AuthorUrl:   "https://jojo.garden",
	Icon:        "fab fa-brackets-square",
	Delay:       time.Duration(5 * time.Second),
	SuccessData: null.NewNullString("Your service '{{.Service.Name}}' is currently online!"),
	FailureData: null.NewNullString("Your service '{{.Service.Name}}' is currently offline!"),
	DataType:    "text",
	Limits:      60,
	Form: []notifications.NotificationForm{
		{
			Type:        "text",
			Title:       "Homeserver URL",
			Placeholder: "https://matrix.org",
			SmallText:   "Enter the Matrix URL, including the http or https scheme",
			DbField:     "host",
			Required:    true,
		},
		{
			Type:        "text",
			Title:       "Room ID",
			Placeholder: "!MzCSpOkYukLxYisAbY:matrix.org",
			SmallText:   "Enter the room ID",
			DbField:     "var1",
			Required:    true,
		},
		{
			Type:        "text",
			Title:       "Token",
			Placeholder: "MDAxyz...",
			SmallText:   "Enter the user token",
			DbField:     "api_secret",
			Required:    true,
		},
	}},
}

func (m *matrix) Select() *notifications.Notification {
	return m.Notification
}

func (m *matrix) Valid(values notifications.Values) error {
	return nil
}

// OnFailure will trigger failing service
func (m *matrix) OnFailure(s services.Service, f failures.Failure) (string, error) {
	msg := ReplaceVars(m.FailureData.String, s, f)
	return m.sendMessage(msg)
}

// OnSuccess will trigger successful service
func (m *matrix) OnSuccess(s services.Service) (string, error) {
	msg := ReplaceVars(m.SuccessData.String, s, failures.Failure{})
	return m.sendMessage(msg)
}

// OnSave will trigger when this notifier is saved
func (m *matrix) OnSave() (string, error) {
	return "", nil
}

// OnTest will test the Twilio SMS messaging
func (m *matrix) OnTest() (string, error) {
	msg := "Testing the Matrix Notifier on your Statping server"
	return m.sendMessage(msg)
}

func (m *matrix) sendMessage(message string) (string, error) {
	type matrixRequestBody struct {
		MsgType string `json:"msgtype"`
		Body    string `json:"body"`
	}
	body := matrixRequestBody{
		MsgType: "m.text",
		Body:    message,
	}
	bodyByte, _ := json.Marshal(body)
	bodyReader := bytes.NewReader(bodyByte)

	parsedUrl, err := url.Parse(m.Host.String)
	if err != nil {
		return "", err
	}

	requestUrl := fmt.Sprintf("%s/_matrix/client/r0/rooms/%s:%s/send/m.room.message", m.Host.String, m.Var1.String, parsedUrl.Hostname())

	headers := []string{
		fmt.Sprintf("Authorization=%s", m.ApiSecret.String),
	}

	contents, _, err := utils.HttpRequest(requestUrl, "POST", "application/json", headers, bodyReader, time.Duration(10*time.Second), true, nil)
	if err != nil {
		return "", err
	}

	return string(contents), err
}
