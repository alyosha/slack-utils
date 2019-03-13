package utils

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
	"golang.org/x/sync/errgroup"
)

// ChannelNameMaxLen is the max character length for a Slack channel name
const ChannelNameMaxLen = 21
const ErrorInviteSelf = "cant_invite_self"

// CreateChannel opens a new public channel and invites the provided list of member IDs, optionally posting an initial message
func (c *Channel) CreateChannel(userIDs []string, initMsg Message, postAsBot bool) (string, error) {
	channel, err := c.UserClient.CreateChannel(c.ChannelName)
	if err != nil {
		return "", fmt.Errorf("failed to create channel: %v", err)
	}

	if channel == nil {
		log.Print("invalid channel")
		return "", nil
	}

	for _, user := range userIDs {
		_, err = c.UserClient.InviteUserToChannel(channel.ID, user)
		if err != nil && err.Error() != ErrorInviteSelf {
			return "", fmt.Errorf("failed to invite user to channel: %v", err)
		}
	}

	client := c.UserClient
	if postAsBot && c.BotClient != nil {
		client = c.BotClient
	}

	if initMsg.Body != "\n" {
		_, ts, err := client.PostMessage(
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

// GetChannelMembers returns a list of members for a given channel
func (s *Slack) GetChannelMembers(channelID string) ([]string, error) {
	channel, err := s.Client.GetChannelInfo(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel info: %v", err)
	}

	return channel.Members, nil
}

// GetChannelMemberEmails returns a list of emails for members of a given channel
func (s *Slack) GetChannelMemberEmails(channelID string) ([]string, error) {
	var eg errgroup.Group
	var memberIDs []string
	var allUsers []slack.User

	eg.Go(func() error {
		channel, err := s.Client.GetChannelInfo(channelID)
		if err == nil {
			memberIDs = channel.Members
		}
		return err
	})

	eg.Go(func() error {
		users, err := s.GetAll()
		if err == nil {
			allUsers = users
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get channel member emails: %v", err)
	}

	emails := toEmails(allUsers, memberIDs)

	return emails, nil
}
