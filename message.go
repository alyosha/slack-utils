package utils

import (
	"errors"
	"fmt"

	"github.com/slack-go/slack"
)

var ErrInvalidForwardChannelID = errors.New("forward channel cannot be the same as original message channel")

// Msg is an intermediary struct used for posting messages
type Msg struct {
	Body    string
	Blocks  []slack.Block
	AsUser  bool
	IconURL string // Incompatible with AsUser option
}

// NewBasicMsg returns a simple text Msg struct
func NewBasicMsg(body string) Msg {
	return Msg{
		Body: body,
	}
}

// NewBasicBotDM returns a simple text Msg for direct messaging the user as bot
func NewBasicBotDM(body string) Msg {
	return Msg{
		Body:   body,
		AsUser: true,
	}
}

// PostMsg sends the provided message to the conversation designated by conversationID
func (c *Client) PostMsg(msg Msg, conversationID string) (string, error) {
	_, ts, err := c.SlackAPI.PostMessage(
		conversationID,
		msg.getCommonOpts()...,
	)

	if err != nil {
		return "", err
	}

	return ts, nil
}

// PostThreadMsg posts a message response into an existing thread
func (c *Client) PostThreadMsg(msg Msg, conversationID string, threadTs string) error {
	_, _, err := c.SlackAPI.PostMessage(
		conversationID,
		append(msg.getCommonOpts(), slack.MsgOptionTS(threadTs))...,
	)

	return err
}

// PostEphemeralMsg sends an ephemeral message in the conversation designated by conversationID
func (c *Client) PostEphemeralMsg(msg Msg, conversationID, userID string) error {
	_, _, err := c.SlackAPI.PostMessage(
		conversationID,
		append(msg.getCommonOpts(), slack.MsgOptionPostEphemeral(userID))...,
	)

	return err
}

// UpdateMsg updates the provided message in the conversation designated by conversationID
func (c *Client) UpdateMsg(msg Msg, conversationID, timestamp string) error {
	_, _, _, err := c.SlackAPI.UpdateMessage(
		conversationID,
		timestamp,
		msg.getCommonOpts()...,
	)

	return err
}

// DeleteMsg deletes the provided message in the conversation designated by conversationID
func (c *Client) DeleteMsg(conversationID, timestamp, responseURL string) error {
	_, _, _, err := c.SlackAPI.UpdateMessage(
		conversationID,
		timestamp,
		slack.MsgOptionDeleteOriginal(responseURL),
	)

	return err
}

// ForwardMsg shares a message permalink to another conversation
func (c *Client) ForwardMsg(originalChannelID, originalMsgTimestamp, forwardChannelID string) error {
	if originalChannelID == forwardChannelID {
		return ErrInvalidForwardChannelID
	}

	permalink, err := c.SlackAPI.GetPermalink(&slack.PermalinkParameters{Channel: originalChannelID, Ts: originalMsgTimestamp})
	if err != nil {
		return fmt.Errorf("c.SlackAPI.GetPermalink > %w", err)
	}

	if _, err := c.PostMsg(NewBasicMsg(permalink), forwardChannelID); err != nil {
		return fmt.Errorf("c.PostMsg > %w", err)
	}

	return nil
}

func (msg Msg) getCommonOpts() []slack.MsgOption {
	return []slack.MsgOption{
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionBlocks(msg.Blocks...),
		slack.MsgOptionAsUser(msg.AsUser),
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionIconURL(msg.IconURL),
	}
}
