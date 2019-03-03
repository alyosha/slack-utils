package utils

import (
	"math/rand"

	"github.com/nlopes/slack"
)

// Message is an intermediary struct used for posting messages
type Message struct {
	Body       string
	Attachment slack.Attachment
}

// Slack is a general purpose struct used when only the client is required
type Slack struct {
	Client *slack.Client
}

// Channel is used in opening/interacting with public Slack channels
type Channel struct {
	Client      *slack.Client
	ChannelName string
	InitMsg     string
}

// Shuffle is used in randmoizing a list of users and splitting them into
// groups of the designated size
type Shuffle struct {
	Client    *slack.Client
	GroupSize int
	Rand      *rand.Rand
}
