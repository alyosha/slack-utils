package utils

import "github.com/nlopes/slack"

// Channel is used in opening/interacting with public Slack channels
type Channel struct {
	Client      *slack.Client
	ChannelName string
	InitMsg     string
}

// User is used in interacting with user data
type User struct {
	Client *slack.Client
}
