package utils

import (
	"github.com/slack-go/slack"
)

// Msg is an intermediary struct used for posting messages
type Msg struct {
	Body    string
	Blocks  []slack.Block
	AsUser  bool
	IconURL string // Incompatible with AsUser option
}

func getCommonOpts(msg Msg) []slack.MsgOption {
	return []slack.MsgOption{
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionBlocks(msg.Blocks...),
		slack.MsgOptionAsUser(msg.AsUser),
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionIconURL(msg.IconURL),
	}
}

// PostMsg sends the provided message to the conversation designated by conversationID
func (c *Client) PostMsg(msg Msg, conversationID string) (string, error) {
	_, ts, err := c.Client.PostMessage(
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
	_, _, err := c.Client.PostMessage(
		conversationID,
		append(getCommonOpts(msg), slack.MsgOptionTS(threadTs))...,
	)

	return err
}

// PostEphemeralMsg sends an ephemeral message in the conversation designated by conversationID
func (c *Client) PostEphemeralMsg(msg Msg, conversationID, userID string) error {
	_, _, err := c.Client.PostMessage(
		conversationID,
		append(getCommonOpts(msg), slack.MsgOptionPostEphemeral(userID))...,
	)

	return err
}

// UpdateMsg updates the provided message in the conversation designated by conversationID
func (c *Client) UpdateMsg(msg Msg, conversationID, timestamp string) error {
	_, _, _, err := c.Client.UpdateMessage(
		conversationID,
		timestamp,
		getCommonOpts(msg)...,
	)

	return err
}

// DeleteMsg deletes the provided message in the conversation designated by conversationID
func (c *Client) DeleteMsg(conversationID, timestamp, responseURL string) error {
	_, _, _, err := c.Client.UpdateMessage(
		conversationID,
		timestamp,
		slack.MsgOptionDeleteOriginal(responseURL),
	)

	return err
}
