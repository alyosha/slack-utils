package utils

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

var ErrNoUsersInWorkplace = errors.New("no users in workplace")

// EmailsToSlackIDs takes in an array of email addresses and finds the IDs of
// any workplace members with those emails
func (c *Client) EmailsToSlackIDs(emails []string) ([]string, error) {
	users, err := c.getAll()
	if err != nil {
		return nil, fmt.Errorf("c.getAll() > %w", err)
	}

	return toSlackIDs(users, emails), nil
}

// EmailToSlackIDsInclusive takes in an array of email addresses, finds the IDs
// of any workplace members with those emails, and returns both values
func (c *Client) EmailsToSlackIDsInclusive(emails []string) ([][]string, error) {
	users, err := c.getAll()
	if err != nil {
		return nil, fmt.Errorf("c.getAll() > %w", err)
	}

	return toSlackIDsInclusive(users, emails), nil
}

func (c *Client) getAll() ([]slack.User, error) {
	users, err := c.client.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("c.client.GetUsers() > %w", err)
	}

	if len(users) == 0 {
		return nil, ErrNoUsersInWorkplace
	}

	return users, nil
}

func toSlackIDs(users []slack.User, emails []string) []string {
	var ids []string
	for _, email := range emails {
		for _, user := range users {
			if user.Profile.Email == email {
				ids = append(ids, user.ID)
			}
		}
	}

	return ids
}

func toSlackIDsInclusive(users []slack.User, emails []string) [][]string {
	var emailIDPairs [][]string
	for _, email := range emails {
		for _, user := range users {
			if user.Profile.Email == email {
				emailIDPairs = append(emailIDPairs, []string{email, user.ID})
			}
		}
	}

	return emailIDPairs
}

func toEmails(users []slack.User, userIDs []string) []string {
	var emails []string
	for _, id := range userIDs {
		for _, user := range users {
			if user.ID == id && user.Profile.Email != "" {
				emails = append(emails, user.Profile.Email)
			}
		}
	}

	return emails
}
