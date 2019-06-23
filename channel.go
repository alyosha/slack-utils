package utils

import (
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// ChannelNameMaxLen is the max character length for a Slack channel name
const ChannelNameMaxLen = 21

const (
	errInviteSelfMsg      = "cant_invite_self"
	errAlreadyArchivedMsg = "already_archived"
)

var ErrNoUsersInWorkplace = errors.New("no users in workplace")

// CreateChannel opens a new public channel and invites the provided list of member IDs, optionally posting an initial message
func (c *Channel) CreateChannel(channelName string, userIDs []string, initMsg Msg, postAsBot bool) error {
	channel, err := c.UserClient.CreateChannel(channelName)
	if err != nil {
		return errors.Wrapf(err, "failed to create new channel")
	}

	for _, user := range userIDs {
		_, err = c.UserClient.InviteUserToChannel(channel.ID, user)
		if err != nil && err.Error() != errInviteSelfMsg {
			return errors.Wrapf(err, "failed to invite user to channel")
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
			return errors.Wrapf(err, "failed to post message")
		}
	}

	c.ChannelID = channel.ID

	return nil
}

func (c *Channel) InviteUsers(userIDs []string) error {
	for _, user := range userIDs {
		_, err := c.UserClient.InviteUserToChannel(c.ChannelID, user)
		if err != nil && err.Error() != errInviteSelfMsg {
			return err
		}
	}

	return nil
}

// GetChannelMembers returns a list of members for a given channel
func (c *Channel) GetChannelMembers() ([]string, error) {
	channel, err := c.UserClient.GetChannelInfo(c.ChannelID)
	if err != nil {
		return nil, err
	}

	return channel.Members, nil
}

// GetChannelMemberEmails returns a list of emails for members of a given channel
func GetChannelMemberEmails(client *slack.Client, channelID string) ([]string, error) {
	var eg errgroup.Group
	var memberIDs []string
	var allUsers []slack.User

	eg.Go(func() error {
		channel, err := client.GetChannelInfo(channelID)
		if err == nil {
			memberIDs = channel.Members
		}
		return err
	})

	eg.Go(func() error {
		users, err := getAll(client)
		if err == nil {
			allUsers = users
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	if len(allUsers) == 0 {
		return nil, ErrNoUsersInWorkplace
	}

	return toEmails(allUsers, memberIDs), nil
}

// LeaveChannels allows the user whose token was used to create the API client to leave multiple channels
func LeaveChannels(client *slack.Client, channelIDs []string) error {
	for _, channelID := range channelIDs {
		_, err := client.LeaveChannel(channelID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ArchiveChannels allows the user whose token was used to create the API client to archive multiple channels
func ArchiveChannels(client *slack.Client, channelIDs []string) error {
	for _, channelID := range channelIDs {
		err := client.ArchiveChannel(channelID)
		if err != nil && err.Error() != errAlreadyArchivedMsg {
			return err
		}
	}
	return nil
}
