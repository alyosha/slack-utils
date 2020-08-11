package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

// Client wraps the slack client for additional utility
type Client struct {
	SlackAPI               *slack.Client
	adminID                string
	logChannel             string
	errChannel             string
	slashResponseConfig    ResponseConfig
	callbackResponseConfig ResponseConfig
}

// ClientConfig is used to configure a new Client
type ClientConfig struct {
	BotToken               string
	AdminID                string
	LogChannelID           string         // Optional
	ErrChannelID           string         // Optional
	SlashResponseConfig    ResponseConfig // Defaults to standard config
	CallbackResponseConfig ResponseConfig // Defaults to standard config
}

// NewClient returns a new client based on provided config
func NewClient(cfg ClientConfig, opts ...slack.Option) (*Client, error) {
	if cfg.BotToken == "" {
		return nil, errMissingBotToken
	}

	if cfg.AdminID == "" {
		return nil, errMissingAdminID
	}

	c := &Client{
		SlackAPI:               slack.New(cfg.BotToken, opts...),
		adminID:                cfg.AdminID,
		logChannel:             cfg.LogChannelID,
		errChannel:             cfg.ErrChannelID,
		slashResponseConfig:    cfg.SlashResponseConfig,
		callbackResponseConfig: cfg.CallbackResponseConfig,
	}

	var eg errgroup.Group

	eg.Go(func() error {
		if _, err := c.SlackAPI.GetUserInfo(cfg.AdminID); err != nil {
			return fmt.Errorf("c.SlackAPI.GetUserInfo > %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if cfg.LogChannelID != "" {
			if _, err := c.SlackAPI.GetConversationInfo(cfg.LogChannelID, false); err != nil {
				return fmt.Errorf("c.SlackAPI.GetConversationInfo > %w", err)
			}
		}
		return nil
	})

	eg.Go(func() error {
		if cfg.ErrChannelID != "" {
			if _, err := c.SlackAPI.GetConversationInfo(cfg.ErrChannelID, false); err != nil {
				return fmt.Errorf("c.SlackAPI.GetConversationInfo > %w", err)
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

// SendResp can be used to send simple responses
func SendResp(w http.ResponseWriter, resp interface{}) error {
	w.Header().Add("Content-type", "application/json")
	return json.NewEncoder(w).Encode(&resp)
}

func withSlashCommand(ctx context.Context, cmd *slack.SlashCommand) context.Context {
	return context.WithValue(ctx, slashCommandKey{}, cmd)
}

func withInteractionCallback(ctx context.Context, cmd *slack.InteractionCallback) context.Context {
	return context.WithValue(ctx, interactionCallbackKey{}, cmd)
}
