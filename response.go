package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/slack-go/slack"
)

const defaultResponseTimeout = 10 * time.Second

type (
	slashRespond    func(ctx context.Context, cmd *slack.SlashCommand)
	callbackRespond func(ctx context.Context, callback *slack.InteractionCallback)

	// ResponseConfig is used to configure various behavioral characteristics of
	// the async responder methods, including max timeout, logging, etc.
	ResponseConfig struct {
		GlobalResponseTimeout time.Duration            // Defaults to 10 seconds if not set
		ResponseTimeoutMap    map[string]time.Duration // Overrides default timeout if key is set for URL path
		WarnDeadlineExceeded  bool                     // Whether to send a warning message to errChannel for timeouts
	}
)

// RespondSlash allows the parent goroutine to finish executing while
// the response function continues. Slack requires a 2xx response to be
// returned within three seconds of interaction or an error is shown to
// the end-user, but using this method bypasses that behavior
func (c *Client) RespondSlash(r *http.Request, respond slashRespond, cmd *slack.SlashCommand) {
	endpoint := r.URL.Path
	timeout := c.getTimeout(endpoint, c.slashResponseConfig)
	newCtx, cancel := context.WithTimeout(context.Background(), timeout)

	go func() {
		defer cancel()
		doneCh := make(chan struct{}, 1)

		go func() {
			respond(newCtx, cmd)
			doneCh <- struct{}{}
		}()

		select {
		case <-newCtx.Done():
			if c.slashResponseConfig.WarnDeadlineExceeded {
				c.warnResponseTimeout(endpoint, timeout, newCtx.Err())
			}
		case <-doneCh:
		}
	}()
}

// RespondCallback allows the parent goroutine to finish executing while
// the response function continues. Slack requires a 2xx response to be
// returned within three seconds of interaction or an error is shown to
// the end-user, but using this method bypasses that behavior
func (c *Client) RespondCallback(r *http.Request, respond callbackRespond, callback *slack.InteractionCallback) {
	endpoint := r.URL.Path
	timeout := c.getTimeout(endpoint, c.callbackResponseConfig)
	newCtx, cancel := context.WithTimeout(context.Background(), timeout)

	go func() {
		defer cancel()
		doneCh := make(chan struct{}, 1)

		go func() {
			respond(newCtx, callback)
			doneCh <- struct{}{}
		}()

		select {
		case <-newCtx.Done():
			if c.callbackResponseConfig.WarnDeadlineExceeded {
				c.warnResponseTimeout(endpoint, timeout, newCtx.Err())
			}
		case <-doneCh:
		}
	}()
}

func (c *Client) getTimeout(endpoint string, responseCfg ResponseConfig) time.Duration {
	if overrideTimeout, ok := responseCfg.ResponseTimeoutMap[endpoint]; ok {
		return overrideTimeout
	}

	if responseCfg.GlobalResponseTimeout == 0 {
		return defaultResponseTimeout
	}

	return responseCfg.GlobalResponseTimeout
}

func (c *Client) warnResponseTimeout(endpoint string, timeout time.Duration, err error) {
	c.SendToErrChannel(
		fmt.Sprintf(
			"response timeout/failure\n*endpoint*: `%s`\n*timeout duration*: `%d`\n*timestamp*: `%d`",
			endpoint,
			timeout,
			time.Now().Unix(),
		),
		err,
	)
}
