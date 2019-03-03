package utils

import (
	"github.com/nlopes/slack"
)

// PostMessage sends the provided message to the channel designated by channelID
func (s *Slack) PostMessage(msg Message, channelID string) error {
	_, ts, err := s.botClient.PostMessage(
		channelID,
		slack.MsgOptionText(msg.body, false),
		slack.MsgOptionAttachments(msg.attachment),
		slack.MsgOptionEnableLinkUnfurl(),
	)

	if err != nil {
		return err
	}

	return ts
}

// PostThreadMessage posts a message response into an existing thread
func (s *Slack) PostThreadMessage(msg message, channelID string, threadTs string) error {
	_, _, err := s.botClient.PostMessage(
		channelID,
		slack.MsgOptionText(msg.body, false),
		slack.MsgOptionAttachments(msg.attachment),
		slack.MsgOptionEnableLinkUnfurl(),
		slack.MsgOptionTS(threadTs),
	)

	if err != nil {
		return err
	}

	return nil
}
