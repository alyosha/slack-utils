package utils

import (
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/slack-go/slack"
)

var errLogChannelNotConfigured = errors.New("associated log channel not configured")

// SendToLogChannel sends the provided message to the log channel
func (c *Client) SendToLogChannel(msg Msg) error {
	if c.logChannel == "" {
		return errLogChannelNotConfigured
	}

	_, err := c.PostMsg(msg, c.logChannel)
	if err != nil {
		return fmt.Errorf("c.PostMsg > %w", err)
	}

	return nil
}

// SendToErrChannel creates an error message based on the provided
// message/error and also posts the stack in the thread
func (c *Client) SendToErrChannel(msgStr string, err error) error {
	if c.errChannel == "" {
		return errLogChannelNotConfigured
	}

	var errMsgBody string

	if msgStr == "" {
		errMsgBody = fmt.Sprintf("`%v`", err)
	} else {
		errMsgBody = fmt.Sprintf("%s\n*error*: `%v`", msgStr, err)
	}

	errMsg := Msg{
		Blocks: []slack.Block{
			NewTextBlock(errMsgBody, nil),
		},
	}

	threadMsg := Msg{
		Body: fmt.Sprintf("```\n%s\n```", string(debug.Stack())),
	}

	ts, err := c.PostMsg(errMsg, c.errChannel)
	if err != nil {
		return fmt.Errorf("c.PostMsg > %w", err)
	}

	if err := c.PostThreadMsg(threadMsg, c.errChannel, ts); err != nil {
		return fmt.Errorf("c.PostThreadMsg > %w", err)
	}

	return nil
}

func (c *Client) logRequest(cfg RequestLoggingConfig, endpoint, userID string) {
	if c.logChannel == "" || !c.shouldLog(cfg, userID) {
		return
	}

	var logMsgBody string

	if cfg.MaskUserID {
		logMsgBody = getBasicLogMsg(endpoint)
	} else {
		logMsgBody = fmt.Sprintf(
			"*endpoint:* `%s`\n*user:* <@%s>\n*timestamp:* `%d`",
			endpoint,
			userID,
			time.Now().Unix(),
		)
	}

	msg := Msg{
		Blocks: []slack.Block{
			NewTextBlock(logMsgBody, nil),
			DivBlock,
		},
	}

	if err := c.SendToLogChannel(msg); err != nil {
		_ = c.SendToErrChannel("failed to log request", err)
	}
}

func (c *Client) shouldLog(cfg RequestLoggingConfig, userID string) bool {
	if !cfg.Enabled || (cfg.ExcludeAdmin && c.adminID == userID) {
		return false
	}
	return true
}

func getBasicLogMsg(endpoint string) string {
	return fmt.Sprintf(
		"*endpoint:* `%s`\n*timestamp:* `%d`",
		endpoint,
		time.Now().Unix(),
	)
}
