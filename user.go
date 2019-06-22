package utils

import (
	"github.com/nlopes/slack"
)

// EmailsToSlackIDs takes in an array of email addresses and finds the IDs of
// any workplace members with those emails
func (s *Slack) EmailsToSlackIDs(emails []string) ([]string, error) {
	users, err := s.getAll()
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNoUsersInWorkplace
	}

	return toSlackIDs(users, emails), nil
}

// EmailToSlackIDsInclusive takes in an array of email addresses, finds the IDs
// of any workplace members with those emails, and returns both values
func (s *Slack) EmailsToSlackIDsInclusive(emails []string) ([][]string, error) {
	users, err := s.getAll()
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNoUsersInWorkplace
	}

	return toSlackIDsInclusive(users, emails), nil
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

func (s *Slack) getAll() ([]slack.User, error) {
	users, err := s.Client.GetUsers()
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNoUsersInWorkplace
	}

	return users, nil
}
