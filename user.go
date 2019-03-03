package utils

import (
	"log"

	"github.com/nlopes/slack"
)

// GetAll returns all users for a workspace
func (u *User) GetAll() []slack.User {
	users, err := u.Client.GetUsers()
	if err != nil {
		log.Printf("Error getting user profiles: %v", err)
		return nil
	}
	return users
}

// EmailsToSlackIDs takes in an array of email addresses and finds the IDs of
// any workplace members with those emails
func (u *User) EmailsToSlackIDs(emails []string) []string {
	users := u.GetAll()
	return toSlackIDs(users, emails)
}

// EmailToSlackIDsInclusive takes in an array of email addresses, finds the IDs
// of any workplace members with those emails, and returns both values
func (u *User) EmailsToSlackIDsInclusive(emails []string) [][]string {
	users := u.GetAll()
	return toSlackIDsInclusive(users, emails)
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
