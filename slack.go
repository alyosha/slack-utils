package utils

import (
	"context"
	"errors"
	"math/rand"

	"github.com/slack/slack"
)

// Responder is an interface implemented when the ability to post response
// messages to the workspace is required
type Responder interface {
	GetClient() *slack.Client
}

// Listener is used for interacting with real time messaging (RTM) events
type Listener struct {
	Client *slack.Client
	BotID  string
}

// Slack is a general purpose struct used when only the client is required
type Slack struct {
	Client *slack.Client
}

// Channel is used in opening/interacting with public Slack channels
type Channel struct {
	UserClient  *slack.Client
	BotClient   *slack.Client
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

// GetClient is the method used to extract the Slack client from the request context
func GetClient(ctx context.Context) (*slack.Client, error) {
	val := ctx.Value(slackClientKey{})
	client, ok := val.(*slack.Client)
	if !ok {
		return nil, errors.New("error extracting the Slack client from context")
	}

	return client, nil
}

// WithContext embeds values into to the request context
func WithContext(ctx context.Context, signingSecret string, client *slack.Client) context.Context {
	return addClient(addSigningSecret(ctx, signingSecret), client)
}

func (l *Listener) GetClient() *slack.Client {
	return l.Client
}

func (s *Slack) GetClient() *slack.Client {
	return s.Client
}

func addSigningSecret(ctx context.Context, signingSecret string) context.Context {
	return context.WithValue(ctx, signingSecretKey{}, signingSecret)
}

func addClient(ctx context.Context, client *slack.Client) context.Context {
	return context.WithValue(ctx, slackClientKey{}, client)
}

func getSigningSecret(ctx context.Context) (string, error) {
	val := ctx.Value(signingSecretKey{})
	secret, ok := val.(string)
	if !ok {
		return "", errors.New("error extracting the signing secret from context")
	}

	return secret, nil
}
