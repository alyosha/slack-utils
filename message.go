package utils

import (
	"github.com/nlopes/slack"
)

// PostMessage sends the provided message to the channel designated by channelID
func (s *Slack) PostMessage(msg Message, channelID string) (string, error) {
	_, ts, err := s.Client.PostMessage(
		channelID,
		slack.MsgOptionText(msg.Body, false),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionEnableLinkUnfurl(),
	)

	if err != nil {
		return "", err
	}

	return ts, nil
}

// PostThreadMessage posts a message response into an existing thread
func (s *Slack) PostThreadMessage(msg Message, channelID string, threadTs string) error {
	_, _, err := s.Client.PostMessage(
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
