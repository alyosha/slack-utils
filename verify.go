package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/nlopes/slack"
)

// VerifyCallbackMsg confirms the validity of the interaction callback via
// the signing secret embedded in the context and returns the verified message body
func VerifyCallbackMsg(r *http.Request) (verifiedBody *slack.InteractionCallback, err error) {
	if r.Method != http.MethodPost {
		return nil, err
	}

	buf, err := checkSecretAndWriteBody(r)
	if err != nil {
		return nil, err
	}

	jsonBody, err := url.QueryUnescape(buf.String()[8:])
	if err != nil {
		return
	}

	var msg *slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonBody), &msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// VerifySlashCmd confirms the validity of the slash command message via
// the signing secret embedded in the context and returns the verified message body
func VerifySlashCmd(r *http.Request) (verifiedBody *slack.SlashCommand, err error) {
	if r.Method != http.MethodPost {
		return nil, err
	}

	buf, err := checkSecretAndWriteBody(r)
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

func checkSecretAndWriteBody(r *http.Request) (bytes.Buffer, error) {
	var buf bytes.Buffer
	signingSecret, err := getSigningSecret(r.Context())
	if err != nil {
		return buf, err
	}

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
