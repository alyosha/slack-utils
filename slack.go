package utils

import (
	"context"
	"errors"
	"math/rand"

	"github.com/nlopes/slack"
)

type signingSecretKey struct{}

// Slack is a general purpose struct used when only the client is required
type Slack struct {
	Client *slack.Client
}

// Channel is used in opening/interacting with a single Slack channel
type Channel struct {
	UserClient *slack.Client
	BotClient  *slack.Client
	ChannelID  string
}

// Shuffle is used in randmoizing a list of users and splitting them into
// groups of the designated size
type Shuffle struct {
	Client    *slack.Client
	GroupSize int
	Rand      *rand.Rand
}

// WithSigningSecret embeds the signing secret value into to the request context
func WithSigningSecret(ctx context.Context, signingSecret string) context.Context {
	return context.WithValue(ctx, signingSecretKey{}, signingSecret)
}

func getSigningSecret(ctx context.Context) (string, error) {
	val := ctx.Value(signingSecretKey{})
	secret, ok := val.(string)
	if !ok {
		return "", errors.New("error extracting the signing secret from context")
	}

	return secret, nil
}
