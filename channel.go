package utils

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

// ChannelNameMaxLen is the max character length for a Slack channel name
const ChannelNameMaxLen = 21
const ErrorInviteSelf = "cant_invite_self"

// CreateChannel opens a new public channel and invites the provided list of member IDs, optionally posting an initial message
func (c *Channel) CreateChannel(userIDs []string, initMsg Message) (string, error) {
	channel, err := c.Client.CreateChannel(c.ChannelName)
	if err != nil {
		return "", fmt.Errorf("failed to create channel: %v", err)
	}

	if channel == nil {
		log.Print("invalid channel")
		return "", nil
	}

	for _, user := range userIDs {
		_, err = c.Client.InviteUserToChannel(channel.ID, user)
		if err != nil && err.Error() != ErrorInviteSelf {
			return "", fmt.Errorf("failed to invite user to channel: %v", err)
		}
	}

	if initMsg.Body != "\n" {
		_, ts, err := c.Client.PostMessage(
			channel.ID,
			slack.MsgOptionText(initMsg.Body, false),
			slack.MsgOptionAttachments(initMsg.Attachments...),
			slack.MsgOptionEnableLinkUnfurl(),
		)
		if err != nil {
			return "", fmt.Errorf("failed to post message to Slack: %v", err)
		}

		log.Printf("posted message to %v at %v after successful channel open", channel.ID, ts)
	}

	return channel.ID, nil
}
