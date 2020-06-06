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
	IconURL     string // Incompatible with AsUser option
}

func getCommonOpts(msg Msg) []slack.MsgOption {
	return []slack.MsgOption{
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionBlocks(msg.Blocks...),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionAsUser(msg.AsUser),
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionIconURL(msg.IconURL),
	}
}

// PostMsg sends the provided message to the conversation designated by conversationID
func (c *Client) PostMsg(msg Msg, conversationID string) (string, error) {
	_, ts, err := c.client.PostMessage(
		conversationID,
		getCommonOpts(msg)...,
	)

	if err != nil {
		return "", err
	}

	return ts, nil
}

// PostThreadMsg posts a message response into an existing thread
func (c *Client) PostThreadMsg(msg Msg, conversationID string, threadTs string) error {
	_, _, err := c.client.PostMessage(
		conversationID,
		append(getCommonOpts(msg), slack.MsgOptionTS(threadTs))...,
	)

	return err
}

// PostEphemeralMsg sends an ephemeral message in the conversation designated by conversationID
func (c *Client) PostEphemeralMsg(msg Msg, conversationID, userID string) error {
	_, _, err := c.client.PostMessage(
		conversationID,
		append(getCommonOpts(msg), slack.MsgOptionPostEphemeral(userID))...,
	)

	return err
}

// UpdateMsg updates the provided message in the conversation designated by conversationID
func (c *Client) UpdateMsg(msg Msg, conversationID, timestamp string) error {
	_, _, _, err := c.client.UpdateMessage(
		conversationID,
		timestamp,
		getCommonOpts(msg)...,
	)

	return err
}

// DeleteMsg deletes the provided message in the conversation designated by conversationID
func (c *Client) DeleteMsg(conversationID, timestamp, responseURL string) error {
	_, _, _, err := c.client.UpdateMessage(
		conversationID,
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
