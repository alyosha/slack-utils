package utils

import (
	"fmt"

	"github.com/slack-go/slack"
	"golang.org/x/sync/errgroup"
)

// ChannelNameMaxLen is the max character length for a Slack channel name
const ChannelNameMaxLen = 21

const (
	errInviteSelfMsg      = "cant_invite_self"
	errAlreadyArchivedMsg = "already_archived"
)

// CreateConversation opens a new public/private channel and invites the provided
// members, optionally posting an initial message
func (c *Client) CreateConversation(conversationName string, isPrivate bool, userIDs []string, initMsg Msg) (string, error) {
	conversation, err := c.SlackAPI.CreateConversation(conversationName, isPrivate)
	if err != nil {
		return "", fmt.Errorf("c.SlackAPI.CreateConversation > %w", err)
	}

	if err = c.InviteUsers(conversation.ID, userIDs); err != nil {
		return conversation.ID, fmt.Errorf("c.InviteUsers > %w", err)
	}

	if initMsg.Body != "" || initMsg.Blocks != nil {
		_, _, err := c.SlackAPI.PostMessage(
			conversation.ID,
			initMsg.getCommonOpts()...,
		)
		if err != nil {
			return conversation.ID, fmt.Errorf("c.SlackAPI.PostMessage > %w", err)
		}
	}

	return conversation.ID, nil
}

// InviteUsers invites multiple users to a conversation
func (c *Client) InviteUsers(conversationID string, userIDs []string) error {
	for _, user := range userIDs {
		_, err := c.SlackAPI.InviteUsersToConversation(conversationID, user)
		if err != nil && err.Error() != errInviteSelfMsg {
			return fmt.Errorf("c.SlackAPI.InviteUsersToConversation > %w", err)
		}
	}

	return nil
}

// ArchiveConversations archives multiple conversations
func (c *Client) ArchiveConversations(conversationIDs []string) error {
	for _, conversationID := range conversationIDs {
		err := c.SlackAPI.ArchiveConversation(conversationID)
		if err != nil && err.Error() != errAlreadyArchivedMsg {
			return fmt.Errorf("c.SlackAPI.ArchiveConversation > %w", err)
		}
	}

	return nil
}

// GetConversationMembers returns a list of members for a given conversation
func (c *Client) GetConversationMembers(conversationID string) ([]string, error) {
	members, _, err := c.SlackAPI.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: conversationID})
	if err != nil {
		return nil, fmt.Errorf("c.SlackAPI.GetUsersInConversation > %w", err)
	}

	return members, nil
}

// GetConversationMemberEmails returns a list of emails for members of a given conversation
func (c *Client) GetConversationMemberEmails(conversationID string) ([]string, error) {
	var eg errgroup.Group
	var memberIDs []string
	var allUsers []slack.User

	eg.Go(func() error {
		members, _, err := c.SlackAPI.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: conversationID})
		if err != nil {
			return fmt.Errorf("c.SlackAPI.GetConversationMembers > %w", err)
		}
		memberIDs = members
		return nil
	})

	eg.Go(func() error {
		users, err := c.getAll()
		if err != nil {
			return fmt.Errorf("c.getAll > %w", err)
		}
		allUsers = users
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return toEmails(allUsers, memberIDs), nil
}
