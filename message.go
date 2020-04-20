package utils

import (
	"encoding/json"
	"net/http"

	"github.com/slack-go/slack"
)

// Msg is an intermediary struct used for posting messages
type Msg struct {
	Body        string
	Blocks      []slack.Block
	Attachments []slack.Attachment
	AsUser      bool
}

func getCommonOpts(msg Msg) []slack.MsgOption {
	return []slack.MsgOption{
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionBlocks(msg.Blocks...),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionAsUser(msg.AsUser),
		slack.MsgOptionEnableLinkUnfurl(),
	}
}

// PostMsg sends the provided message to the channel designated by channelID
func PostMsg(client *slack.Client, msg Msg, channelID string) (string, error) {
	_, ts, err := client.PostMessage(
		channelID,
		getCommonOpts(msg)...,
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

	return err
}

// PostEphemeralMsg sends an ephemeral message in the channel designated by channelID
func PostEphemeralMsg(client *slack.Client, msg Msg, channelID, userID string) error {
	_, _, err := client.PostMessage(
		channelID,
		append(getCommonOpts(msg), slack.MsgOptionPostEphemeral(userID))...,
	)

	return err
}

// UpdateMsg updates the provided message in the channel designated by channelID
func UpdateMsg(client *slack.Client, msg Msg, channelID, timestamp string) error {
	_, _, _, err := client.UpdateMessage(
		channelID,
		timestamp,
		getCommonOpts(msg)...,
	)

	return err
}

// DeleteMsg deletes the provided message in the channel designated by channelID
func DeleteMsg(client *slack.Client, channelID, timestamp, responseURL string) error {
	_, _, _, err := client.UpdateMessage(
		channelID,
		timestamp,
		slack.MsgOptionDeleteOriginal(responseURL),
	)

	return err
}

// SendEmptyOK responds with status 200
func SendEmptyOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	return
}

// SendResp can be used to send simple callback responses
// NOTE: cannot be used in callback from block messages
func SendResp(w http.ResponseWriter, msg slack.Message) error {
	w.Header().Add("Content-type", "application/json")
	return json.NewEncoder(w).Encode(&msg)
}

// ReplaceOriginal replaces the original message with the newly encoded one
// NOTE: cannot be used in callback from block messages
func ReplaceOriginal(w http.ResponseWriter, msg slack.Message) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	msg.ReplaceOriginal = true
	return json.NewEncoder(w).Encode(&msg)
}

// SendOKAndDeleteOriginal responds with status 200 and deletes the original message
// NOTE: cannot be used in callback from block messages
func SendOKAndDeleteOriginal(w http.ResponseWriter) error {
	var msg slack.Message
	msg.DeleteOriginal = true
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(&msg)
}
