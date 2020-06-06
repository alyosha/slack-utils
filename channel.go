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
// members, optionally posting an initial message.
func (c *Client) CreateConversation(conversationName string, isPrivate bool, userIDs []string, initMsg Msg) (string, error) {
	conversation, err := c.client.CreateConversation(conversationName, isPrivate)
	if err != nil {
		return "", fmt.Errorf("c.client.CreateConversation() > %w", err)
	}

	if err = c.InviteUsers(conversation.ID, userIDs); err != nil {
		return conversation.ID, fmt.Errorf("c.InviteUsers() > %w", err)
	}

	if initMsg.Body != "" || initMsg.Blocks != nil {
		_, _, err := c.client.PostMessage(
			conversation.ID,
			getCommonOpts(initMsg)...,
		)
		if err != nil {
			return conversation.ID, fmt.Errorf("c.client.PostMessage() > %w", err)
		}
	}

	return conversation.ID, nil
}

// InviteUsers invites multiple users to a conversation
func (c *Client) InviteUsers(conversationID string, userIDs []string) error {
	for _, user := range userIDs {
		_, err := c.client.InviteUsersToConversation(conversationID, user)
		if err != nil && err.Error() != errInviteSelfMsg {
			return fmt.Errorf("c.client.InviteUsersToConversation() > %w", err)
		}
	}

	return nil
}

// ArchiveConversations allows archives multiple conversations
func (c *Client) ArchiveConversations(conversationIDs []string) error {
	for _, conversationID := range conversationIDs {
		err := c.client.ArchiveConversation(conversationID)
		if err != nil && err.Error() != errAlreadyArchivedMsg {
			return fmt.Errorf("c.client.ArchiveConversation() > %w", err)
		}
	}

	return nil
}

// GetConversationMembers returns a list of members for a given conversation
func (c *Client) GetConversationMembers(conversationID string) ([]string, error) {
	conversation, err := c.client.GetConversationInfo(conversationID, false)
	if err != nil {
		return nil, fmt.Errorf("c.client.GetConversationInfo() > %w", err)
	}

	return conversation.Members, nil
}

// GetConversationMemberEmails returns a list of emails for members of a given conversation
func (c *Client) GetConversationMemberEmails(conversationID string) ([]string, error) {
	var eg errgroup.Group
	var memberIDs []string
	var allUsers []slack.User

	eg.Go(func() error {
		conversation, err := c.client.GetConversationInfo(conversationID, false)
		if err != nil {
			return fmt.Errorf("c.client.GetConversationInfo() > %w", err)
		}
		memberIDs = conversation.Members
		return nil
	})

	eg.Go(func() error {
		users, err := c.getAll()
		if err != nil {
			return fmt.Errorf("c.getAll() > %w", err)
		}
		allUsers = users
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return toEmails(allUsers, memberIDs), nil
}
