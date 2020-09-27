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
	slashRespond    func(ctx context.Context, cmd *slack.SlashCommand) error
	callbackRespond func(ctx context.Context, callback *slack.InteractionCallback) error

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
// the end-user, but using this method bypasses that behavior. Depending
// on the interaction, such a delay might lead an end user to believe the
// request did not go through and cause them to retry, so use with caution.
func (c *Client) RespondSlash(r *http.Request, respond slashRespond, cmd *slack.SlashCommand) {
	endpoint := r.URL.Path
	timeout := c.getTimeout(endpoint, c.slashResponseConfig)
	newCtx, cancel := context.WithTimeout(copyContext(r.Context()), timeout)

	go func() {
		defer cancel()
		respondCh := make(chan error)

		go func() {
			select {
			case err := <-respondCh:
				if err != nil {
					c.notifyResponseFailure(endpoint, err)
				}
			case <-newCtx.Done():
			}
		}()

		select {
		case respondCh <- respond(newCtx, cmd):
		case <-newCtx.Done():
			if c.slashResponseConfig.WarnDeadlineExceeded {
				c.notifyResponseTimeout(endpoint, timeout, newCtx.Err())
			}
			return
		}
	}()
}

// RespondCallback allows the parent goroutine to finish executing while
// the response function continues. Slack requires a 2xx response to be
// returned within three seconds of interaction or an error is shown to
// the end-user, but using this method bypasses that behavior. Depending
// on the interaction, such a delay might lead an end user to believe the
// request did not go through and cause them to retry, so use with caution.
func (c *Client) RespondCallback(r *http.Request, respond callbackRespond, callback *slack.InteractionCallback) {
	endpoint := r.URL.Path
	timeout := c.getTimeout(endpoint, c.callbackResponseConfig)
	newCtx, cancel := context.WithTimeout(copyContext(r.Context()), timeout)

	go func() {
		defer cancel()
		respondCh := make(chan error)

		go func() {
			select {
			case err := <-respondCh:
				if err != nil {
					c.notifyResponseFailure(endpoint, err)
				}
			case <-newCtx.Done():
			}
		}()

		select {
		case respondCh <- respond(newCtx, callback):
		case <-newCtx.Done():
			if c.callbackResponseConfig.WarnDeadlineExceeded {
				c.notifyResponseTimeout(endpoint, timeout, newCtx.Err())
			}
			return
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

func (c *Client) notifyResponseTimeout(endpoint string, timeout time.Duration, err error) {
	c.SendToErrChannel(
		fmt.Sprintf(
			"*response timeout*\n*endpoint*: `%s`\n*timeout duration*: `%d`\n*timestamp*: `%d`",
			endpoint,
			timeout,
			time.Now().Unix(),
		),
		err,
	)
}

func (c *Client) notifyResponseFailure(endpoint string, err error) {
	c.SendToErrChannel(
		fmt.Sprintf(
			"*response failure*\n*endpoint*: `%s`\n*timestamp*: `%d`",
			endpoint,
			time.Now().Unix(),
		),
		err,
	)
}

func copyContext(ctx context.Context) context.Context { return contextCopy{ctx} }

type contextCopy struct{ parent context.Context }

func (v contextCopy) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (v contextCopy) Done() <-chan struct{}             { return nil }
func (v contextCopy) Err() error                        { return nil }
func (v contextCopy) Value(key interface{}) interface{} { return v.parent.Value(key) }
