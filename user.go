package utils

import (
	"github.com/nlopes/slack"
)

// GetAll returns all users for a workspace
func (s *Slack) GetAll() ([]slack.User, error) {
	users, err := s.Client.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// EmailsToSlackIDs takes in an array of email addresses and finds the IDs of
// any workplace members with those emails
func (s *Slack) EmailsToSlackIDs(emails []string) ([]string, error) {
	users, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	return toSlackIDs(users, emails), nil
}

// EmailToSlackIDsInclusive takes in an array of email addresses, finds the IDs
// of any workplace members with those emails, and returns both values
func (s *Slack) EmailsToSlackIDsInclusive(emails []string) ([][]string, error) {
	users, err := s.GetAll()
	if err != nil {
		return nil, err
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
