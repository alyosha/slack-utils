package utils

import (
	"context"
	"errors"

	"github.com/slack-go/slack"
)

var (
	errSlashCommandNotFound        = errors.New("no slash command found in context")
	errInteractionCallbackNotFound = errors.New("no callback found in context")
)

type slashCommandKey struct{}
type interactionCallbackKey struct{}

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
