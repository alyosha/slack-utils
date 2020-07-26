package utils

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type userInfoMapKind string

const (
	mapKindIDEmail userInfoMapKind = "id_email"
	mapKindEmailID userInfoMapKind = "email_id"
)

// ErrNoUsersInWorkplace is returned when there was no error calling Slack, but
// the requested method cannot continue given there are no returned users.
var ErrNoUsersInWorkplace = errors.New("no users in workplace")

// EmailsToSlackIDs takes in an array of email addresses and finds the IDs of
// any workplace members with those emails
func (c *Client) EmailsToSlackIDs(emails []string) ([]string, error) {
	users, err := c.getAll()
	if err != nil {
		return nil, fmt.Errorf("c.getAll > %w", err)
	}

	return toSlackIDs(users, emails), nil
}

// EmailToSlackIDsInclusive takes in an array of email addresses, finds the IDs
// of any workplace members with those emails, and returns both values
func (c *Client) EmailsToSlackIDsInclusive(emails []string) ([][]string, error) {
	users, err := c.getAll()
	if err != nil {
		return nil, fmt.Errorf("c.getAll > %w", err)
	}

	return toSlackIDsInclusive(users, emails), nil
}

func (c *Client) getAll() ([]slack.User, error) {
	users, err := c.SlackAPI.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("c.SlackAPI.GetUsers > %w", err)
	}

	if len(users) == 0 {
		return nil, ErrNoUsersInWorkplace
	}

	return users, nil
}

func toSlackIDs(users []slack.User, emails []string) []string {
	var ids []string

	userEmailMap := getUserInfoMap(users, mapKindEmailID)

	for _, email := range emails {
		if userID, ok := userEmailMap[email]; ok {
			ids = append(ids, userID)
		}
	}

	return ids
}

func toSlackIDsInclusive(users []slack.User, emails []string) [][]string {
	var emailIDPairs [][]string

	userEmailMap := getUserInfoMap(users, mapKindEmailID)

	for _, email := range emails {
		if userID, ok := userEmailMap[email]; ok {
			emailIDPairs = append(emailIDPairs, []string{email, userID})
		}
	}

	return emailIDPairs
}

func toEmails(users []slack.User, userIDs []string) []string {
	var emails []string

	userEmailMap := getUserInfoMap(users, mapKindIDEmail)

	for _, id := range userIDs {
		if email, ok := userEmailMap[id]; ok {
			if email != "" {
				emails = append(emails, email)
			}
		}
	}

	return emails
}

func getUserInfoMap(users []slack.User, mapKind userInfoMapKind) map[string]string {
	userInfoMap := make(map[string]string)

	for _, user := range users {
		switch mapKind {
		case mapKindIDEmail:
			userInfoMap[user.ID] = user.Profile.Email
		case mapKindEmailID:
			userInfoMap[user.Profile.Email] = user.ID
		}
	}

	return userInfoMap
}
