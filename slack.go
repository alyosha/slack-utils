package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/slack-go/slack"
	"golang.org/x/sync/errgroup"
)

var (
	errMissingBotToken             = errors.New("must provide bot token")
	errMissingAdminID              = errors.New("must provide admin ID")
	errSlashCommandNotFound        = errors.New("no slash command found in context")
	errInteractionCallbackNotFound = errors.New("no callback found in context")
)

type (
	slashCommandKey        struct{}
	interactionCallbackKey struct{}
)

// Client wraps the slack Client for additional utility
type Client struct {
	Client           *slack.Client
	adminID          string
	logChannelID     string
	errChannelID     string
	logAdminRequests bool
}

// ClientConfig is used to configure a new Client
type ClientConfig struct {
	BotToken         string
	AdminID          string
	LogChannelID     string
	ErrChannelID     string
	LogAdminRequests bool
}

// NewClient returns a new client based on provided config
func NewClient(cfg ClientConfig, opt ...slack.Option) (*Client, error) {
	if cfg.BotToken == "" {
		return nil, errMissingBotToken
	}

	if cfg.AdminID == "" {
		return nil, errMissingAdminID
	}

	c := &Client{
		Client:           slack.New(cfg.BotToken, opt...),
		adminID:          cfg.AdminID,
		logChannelID:     cfg.LogChannelID,
		errChannelID:     cfg.ErrChannelID,
		logAdminRequests: cfg.LogAdminRequests,
	}

	var eg errgroup.Group

	eg.Go(func() error {
		if _, err := c.Client.GetUserInfo(cfg.AdminID); err != nil {
			return fmt.Errorf("c.Client.GetUserInfo() > %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if cfg.LogChannelID != "" {
			if _, err := c.Client.GetConversationInfo(cfg.LogChannelID, false); err != nil {
				return fmt.Errorf("c.Client.GetConversationInfo() > %w", err)
			}
		}
		return nil
	})

	eg.Go(func() error {
		if cfg.ErrChannelID != "" {
			if _, err := c.Client.GetConversationInfo(cfg.ErrChannelID, false); err != nil {
				return fmt.Errorf("c.Client.GetConversationInfo() > %w", err)
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return c, nil
}

// SlashCommand retrieves the verified slash command from the context. To
// utilize this functionality, you must use the VerifySlashCommand middleware.
func SlashCommand(ctx context.Context) (*slack.SlashCommand, error) {
	val := ctx.Value(slashCommandKey{})
	cmd, ok := val.(*slack.SlashCommand)
	if !ok {
		return nil, errSlashCommandNotFound
	}
	return cmd, nil
}

// InteractionCallback retrieves the verified interaction callback from the context.
// To utilize this functionality, you must use the VerifyInteractionCallback middleware.
func InteractionCallback(ctx context.Context) (*slack.InteractionCallback, error) {
	val := ctx.Value(interactionCallbackKey{})
	callback, ok := val.(*slack.InteractionCallback)
	if !ok {
		return nil, errInteractionCallbackNotFound
	}
	return callback, nil
}

func withSlashCommand(ctx context.Context, cmd *slack.SlashCommand) context.Context {
	return context.WithValue(ctx, slashCommandKey{}, cmd)
}

func withInteractionCallback(ctx context.Context, cmd *slack.InteractionCallback) context.Context {
	return context.WithValue(ctx, interactionCallbackKey{}, cmd)
}
