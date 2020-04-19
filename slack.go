package utils

import (
	"context"
	"errors"

	"github.com/slack-go/slack"
)

var ErrNoSecret = errors.New("no signing secret found in context")

type signingSecretKey struct{}

// Channel is used in opening/interacting with a single Slack channel
type Channel struct {
	UserClient *slack.Client
	BotClient  *slack.Client
	ChannelID  string
}

// WithSigningSecret embeds the signing secret value into to the request context
func WithSigningSecret(ctx context.Context, signingSecret string) context.Context {
	return context.WithValue(ctx, signingSecretKey{}, signingSecret)
}

func getSigningSecret(ctx context.Context) (string, error) {
	val := ctx.Value(signingSecretKey{})
	secret, ok := val.(string)
	if !ok {
		return "", ErrNoSecret
	}

	return secret, nil
}
