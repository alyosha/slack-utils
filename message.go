package utils

import (
	"encoding/json"
	"net/http"

	"github.com/nlopes/slack"
)

// Msg is an intermediary struct used for posting messages
type Msg struct {
	Body        string
	Attachments []slack.Attachment
	AsUser      bool
}

// PostMsg sends the provided message to the channel designated by channelID
func PostMsg(client *slack.Client, msg Msg, channelID string) (string, error) {
	_, ts, err := client.PostMessage(
		channelID,
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionAsUser(msg.AsUser),
	)

	if err != nil {
		return "", err
	}

	return ts, nil
}

// PostThreadMsg posts a message response into an existing thread
func PostThreadMsg(client *slack.Client, msg Msg, channelID string, threadTs string) error {
	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionTS(threadTs),
	)

	if err != nil {
		return err
	}

	return nil
}

// SendResp can be used to send simple callback responses
func SendResp(w http.ResponseWriter, msg slack.Message) {
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(&msg)
	return
}

// SendOKAndDeleteOriginal responds with status 200 and deletes the original message
func SendOKAndDeleteOriginal(w http.ResponseWriter) {
	var msg slack.Message
	msg.DeleteOriginal = true
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&msg)
	return
}

// SendEmptyOK responds with status 200
func SendEmptyOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	return
}
