package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/slack-go/slack"
)

const (
	verifyActionTypeSucceedSlash    = "succeed_slash"
	verifyActionTypeSucceedCallback = "succeed_callback"
	verifyActionTypeFail            = "fail"
	verifyActionTypeMalformed       = "malformed"
)

// The following func types are used to configure custom additional actions on
// verify middleware success/failure (e.g. logging, etc.)
type (
	VerifyAction interface {
		verifyActionType() string
	}

	VerifySucceedSlash    func(w http.ResponseWriter, r *http.Request, cmd *slack.SlashCommand)
	VerifySucceedCallback func(w http.ResponseWriter, r *http.Request, cmd *slack.InteractionCallback)
	VerifyFail            func(w http.ResponseWriter, r *http.Request, err error)
)

func (v VerifySucceedSlash) verifyActionType() string {
	if v == nil {
		return verifyActionTypeMalformed
	}
	return verifyActionTypeSucceedSlash
}

func (v VerifySucceedCallback) verifyActionType() string {
	if v == nil {
		return verifyActionTypeMalformed
	}
	return verifyActionTypeSucceedCallback
}

func (v VerifyFail) verifyActionType() string {
	if v == nil {
		return verifyActionTypeMalformed
	}
	return verifyActionTypeFail
}

// VerifySlash is a middleware that will automatically verify the
// authenticity of the incoming request and embed the unmarshalled SlashCommand
// in the context on success. Use the optional succeed/fail parameters to
// configure additional behavior on sucess/failure, or simply provide nil if
// no further action is required.
func (c *Client) VerifySlash(signingSecret string, logConfig RequestLoggingConfig, verifyActions ...VerifyAction) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cmd, err := verifySlashCommand(r, signingSecret)
			if err != nil {
				fail(w, r, err, verifyActions...)
				return
			}
			ctx := withSlashCommand(r.Context(), cmd)
			for _, action := range verifyActions {
				if action.verifyActionType() == verifyActionTypeSucceedSlash {
					action.(VerifySucceedSlash)(w, r, cmd)
				}
			}
			c.logRequest(logConfig, r.URL.Path, cmd.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// VerifyCallback is a middleware that will automatically verify
// the authenticity of the incoming request and embed the unmarshalled
// InteractionCallback in the context on success. Use the optional succeed/fail
// parameters to configure additional behavior on sucess/failure, or simply
// provide nil if no further action is required.
func (c *Client) VerifyCallback(signingSecret string, logConfig RequestLoggingConfig, verifyActions ...VerifyAction) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callback, err := verifyInteractionCallback(r, signingSecret)
			if err != nil {
				fail(w, r, err, verifyActions...)
				return
			}
			ctx := withInteractionCallback(r.Context(), callback)
			for _, action := range verifyActions {
				if action.verifyActionType() == verifyActionTypeSucceedCallback {
					action.(VerifySucceedCallback)(w, r, callback)
				}
			}
			c.logRequest(logConfig, r.URL.Path, callback.User.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func fail(w http.ResponseWriter, r *http.Request, err error, verifyActions ...VerifyAction) {
	for _, action := range verifyActions {
		if action == nil {
			continue
		}
		if action.verifyActionType() == verifyActionTypeFail {
			action.(VerifyFail)(w, r, err)
		}
	}
}

func verifyInteractionCallback(r *http.Request, signingSecret string) (verifiedBody *slack.InteractionCallback, err error) {
	if r.Method != http.MethodPost {
		return nil, err
	}

	buf, err := checkSecretAndWriteBody(r, signingSecret)
	if err != nil {
		return nil, err
	}

	jsonBody, err := url.QueryUnescape(strings.Replace(buf.String(), "payload=", "", 1))
	if err != nil {
		return
	}

	msg := &slack.InteractionCallback{}
	if err := json.Unmarshal([]byte(jsonBody), msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func verifySlashCommand(r *http.Request, signingSecret string) (verifiedBody *slack.SlashCommand, err error) {
	if r.Method != http.MethodPost {
		return nil, err
	}

	buf, err := checkSecretAndWriteBody(r, signingSecret)
	if err != nil {
		return nil, err
	}

	body, err := url.ParseQuery(string(buf.String()))
	if err != nil {
		return nil, err
	}

	msg := parseCmd(body)

	return &msg, nil
}

func checkSecretAndWriteBody(r *http.Request, signingSecret string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		return buf, err
	}

	dest := io.MultiWriter(&buf, &sv)
	if _, err := io.Copy(dest, r.Body); err != nil {
		return buf, err
	}

	if err := sv.Ensure(); err != nil {
		return buf, err
	}

	return buf, nil
}

func parseCmd(body url.Values) (s slack.SlashCommand) {
	s.Token = body.Get("token")
	s.TeamID = body.Get("team_id")
	s.TeamDomain = body.Get("team_domain")
	s.EnterpriseID = body.Get("enterprise_id")
	s.EnterpriseName = body.Get("enterprise_name")
	s.ChannelID = body.Get("channel_id")
	s.ChannelName = body.Get("channel_name")
	s.UserID = body.Get("user_id")
	s.UserName = body.Get("user_name")
	s.Command = body.Get("command")
	s.Text = body.Get("text")
	s.ResponseURL = body.Get("response_url")
	s.TriggerID = body.Get("trigger_id")

	return s
}
