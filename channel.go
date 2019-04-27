package utils

import (
	"github.com/nlopes/slack"
	"golang.org/x/sync/errgroup"
)

// ChannelNameMaxLen is the max character length for a Slack channel name
const ChannelNameMaxLen = 21
const ErrorInviteSelf = "cant_invite_self"
const ErrorAlreadyArchived = "already_archived"

// CreateChannel opens a new public channel and invites the provided list of member IDs, optionally posting an initial message
func (c *Channel) CreateChannel(channelName string, userIDs []string, initMsg Msg, postAsBot bool) (string, error) {
	channel, err := c.UserClient.CreateChannel(channelName)
	if err != nil {
		return "", err
	}

	for _, user := range userIDs {
		_, err = c.UserClient.InviteUserToChannel(channel.ID, user)
		if err != nil && err.Error() != ErrorInviteSelf {
			return "", err
		}
	}

	client := c.UserClient
	if postAsBot && c.BotClient != nil {
		client = c.BotClient
	}

	if initMsg.Body != "" {
		_, _, err := client.PostMessage(
			channel.ID,
			slack.MsgOptionText(initMsg.Body, false),
			slack.MsgOptionAttachments(initMsg.Attachments...),
			slack.MsgOptionBlocks(initMsg.Blocks...),
			slack.MsgOptionEnableLinkUnfurl(),
		)
		if err != nil {
			return "", err
		}
	}

	return channel.ID, nil
}

func (c *Channel) InviteUsers(userIDs []string, channelID string) error {
	for _, user := range userIDs {
		_, err := c.UserClient.InviteUserToChannel(channelID, user)
		if err != nil && err.Error() != ErrorInviteSelf {
			return err
		}
	}

	return nil
}

// GetChannelMembers returns a list of members for a given channel
func (s *Slack) GetChannelMembers(channelID string) ([]string, error) {
	channel, err := s.Client.GetChannelInfo(channelID)
	if err != nil {
		return nil, err
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
		users, err := s.getAll()
		if err == nil {
			allUsers = users
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return toEmails(allUsers, memberIDs), nil
}

// LeaveChannels allows the user whose token was used to create the API client to leave multiple channels
func (s *Slack) LeaveChannels(channelIDs []string) error {
	for _, channelID := range channelIDs {
		_, err := s.Client.LeaveChannel(channelID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ArchiveChannels allows the user whose token was used to create the API client to archive multiple channels
func (s *Slack) ArchiveChannels(channelIDs []string) error {
	for _, channelID := range channelIDs {
		err := s.Client.ArchiveChannel(channelID)
		if err != nil && err.Error() != ErrorAlreadyArchived {
			return err
		}
	}
	return nil
}
