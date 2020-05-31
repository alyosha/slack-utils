package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/kylelemons/godebug/pretty"
	"github.com/slack-go/slack"
)

const (
	testSecret1           = "e6b19c573432dcc6b075501d51b51bb8"
	testSecret2           = "f6b19c573432dcc6b075501d51b51bb8"
	testInvalidSigningSig = "v0=12345"
	testCallbackRaw       = `payload=%7B%22type%22%3A%22block_actions%22%2C%22user%22%3A%7B%22id%22%3A%22U12345678%22%2C%22username%22%3A%22fakenameyo%22%2C%22name%22%3A%22fakenameyo%22%2C%22team_id%22%3A%22T0000000%22%7D%2C%22api_app_id%22%3A%22A00000000%22%2C%22token%22%3A%22faketoken%22%2C%22container%22%3A%7B%22type%22%3A%22message%22%2C%22message_ts%22%3A%221589970639.001400%22%2C%22channel_id%22%3A%22G0000000%22%2C%22is_ephemeral%22%3Atrue%7D%2C%22trigger_id%22%3A%220000000000.1111111111.222222222222aaaaaaaaaaaaaa%22%2C%22team%22%3A%7B%22id%22%3A%22T0000000%22%2C%22domain%22%3A%22domain%22%7D%2C%22channel%22%3A%7B%22id%22%3A%22G0000000%22%2C%22name%22%3A%22privategroup%22%7D%2C%22response_url%22%3A%22https%3A%5C%2F%5C%2Fhooks.slack.com%5C%2Factions%5C%2FT0000000F%5C%2F000000000%5C%2FYYYYYYYYYYY%22%2C%22actions%22%3A%5B%7B%22action_id%22%3A%22cancel_action%22%2C%22block_id%22%3A%22channel_id_block%22%2C%22text%22%3A%7B%22type%22%3A%22plain_text%22%2C%22text%22%3A%22Done%22%2C%22emoji%22%3Atrue%7D%2C%22value%22%3A%22done%22%2C%22style%22%3A%22primary%22%2C%22type%22%3A%22button%22%2C%22action_ts%22%3A%221589971722.911477%22%7D%5D%7D`
)

var (
	testReqTsValid = fmt.Sprintf("%d", time.Now().Unix())
)

func TestVerifySlashCommand(t *testing.T) {
	testCases := []struct {
		description         string
		useMiddleware       bool
		secret              string
		ts                  string
		invalidHex          bool
		failFunc            VerifyFail
		succeedFunc         VerifySucceedSlash
		wantErr             error
		wantRespBody        string
		containsRespPattern string
	}{
		{
			description:   "using middleware and valid signing signature, expected command retrieved from context. no/empty success method so no extra action",
			useMiddleware: true,
			secret:        testSecret1,
			ts:            testReqTsValid,
		},
		{
			description:   "using middleware and valid signing signature, expected extra success response received",
			useMiddleware: true,
			succeedFunc: func(w http.ResponseWriter, r *http.Request, cmd *slack.SlashCommand) {
				_, _ = w.Write([]byte("OK"))
			},
			secret:       testSecret1,
			ts:           testReqTsValid,
			wantRespBody: "OK",
		},
		{
			description:   "using middleware with valid secret but timestamp is too old, verify fails and req killed, expected fail response received",
			useMiddleware: true,
			failFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				_, _ = w.Write([]byte(err.Error()))
			},
			secret:       testSecret1,
			ts:           "1531431954",
			wantRespBody: "timestamp is too old",
		},
		{
			description:   "using middleware with wrong secret, verify fails and req killed, expected fail response received",
			useMiddleware: true,
			failFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				_, _ = w.Write([]byte(err.Error()))
			},
			secret:              testSecret2,
			ts:                  testReqTsValid,
			containsRespPattern: "Expected signing signature:",
		},
		{
			description:   "using middleware and invalid signing signature, verify fails and req killed, no/empty fail method provided so no extra action",
			useMiddleware: true,
			invalidHex:    true,
		},
		{
			description:   "using middleware and invalid signing signature, verify fails and req killed. expected fail response received",
			useMiddleware: true,
			invalidHex:    true,
			ts:            testReqTsValid,
			failFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				_, _ = w.Write([]byte(err.Error()))
			},
			wantRespBody: "encoding/hex: odd length hex string",
		},
		{
			description: "not using middleware, command not found",
			wantErr:     errSlashCommandNotFound,
		},
	}

	body := url.Values{
		"command":         []string{"/command"},
		"team_domain":     []string{"team"},
		"enterprise_id":   []string{"E0001"},
		"enterprise_name": []string{"Globular%20Construct%20Inc"},
		"channel_id":      []string{"C1234ABCD"},
		"text":            []string{"text"},
		"team_id":         []string{"T1234ABCD"},
		"user_id":         []string{"U1234ABCD"},
		"user_name":       []string{"username"},
		"response_url":    []string{"https://hooks.slack.com/commands/XXXXXXXX/00000000000/YYYYYYYYYYYYYY"},
		"token":           []string{"valid"},
		"channel_name":    []string{"channel"},
		"trigger_id":      []string{"0000000000.1111111111.222222222222aaaaaaaaaaaaaa"},
	}

	encodedBody := body.Encode()

	wantCmd := slack.SlashCommand{
		Command:        "/command",
		TeamDomain:     "team",
		EnterpriseID:   "E0001",
		EnterpriseName: "Globular%20Construct%20Inc",
		ChannelID:      "C1234ABCD",
		Text:           "text",
		TeamID:         "T1234ABCD",
		UserID:         "U1234ABCD",
		UserName:       "username",
		ResponseURL:    "https://hooks.slack.com/commands/XXXXXXXX/00000000000/YYYYYYYYYYYYYY",
		Token:          "valid",
		ChannelName:    "channel",
		TriggerID:      "0000000000.1111111111.222222222222aaaaaaaaaaaaaa",
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			r := chi.NewRouter()

			if tc.useMiddleware {
				r.Use(VerifySlashCommand(testSecret1, tc.succeedFunc, tc.failFunc))
			}

			signingSig := getSigningSig(t, tc.ts, tc.secret, []byte(encodedBody))
			if tc.invalidHex {
				signingSig = testInvalidSigningSig
			}

			r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
				cmd, err := SlashCommand(r.Context())
				if err != nil {
					if err != tc.wantErr {
						t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
					}
					return
				}
				if diff := pretty.Compare(wantCmd, cmd); diff != "" {
					t.Fatalf("+got -want %s\n", diff)
				}
			})

			testServ := httptest.NewServer(r)
			defer testServ.Close()

			respBodyString := executeTestReq(t, testServ, signingSig, tc.ts, encodedBody)

			if tc.containsRespPattern != "" {
				if !strings.Contains(respBodyString, tc.containsRespPattern) {
					t.Fatalf("expected resp to contain pattern: %s, got: %s", tc.containsRespPattern, respBodyString)
				}
				return
			}

			if respBodyString != tc.wantRespBody {
				t.Fatalf("expected resp body: %s, got: %s", tc.wantRespBody, respBodyString)
			}
		})
	}
}

func TestVerifyInteractionCallback(t *testing.T) {
	testCases := []struct {
		description         string
		useMiddleware       bool
		secret              string
		ts                  string
		invalidHex          bool
		failFunc            VerifyFail
		succeedFunc         VerifySucceedCallback
		wantErr             error
		wantRespBody        string
		containsRespPattern string
	}{
		{
			description:   "using middleware and valid signing signature, expected callback retrieved from context. no/empty success method so no extra action",
			useMiddleware: true,
			secret:        testSecret1,
			ts:            testReqTsValid,
		},
		{
			description:   "using middleware and valid signing signature, expected extra success response received",
			useMiddleware: true,
			succeedFunc: func(w http.ResponseWriter, r *http.Request, cmd *slack.InteractionCallback) {
				_, _ = w.Write([]byte("OK"))
			},
			secret:       testSecret1,
			ts:           testReqTsValid,
			wantRespBody: "OK",
		},
		{
			description:   "using middleware with valid secret but timestamp is too old, verify fails and req killed, expected fail response received",
			useMiddleware: true,
			failFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				_, _ = w.Write([]byte(err.Error()))
			},
			secret:       testSecret1,
			ts:           "1531431954",
			wantRespBody: "timestamp is too old",
		},
		{
			description:   "using middleware with wrong secret, verify fails and req killed, expected fail response received",
			useMiddleware: true,
			failFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				_, _ = w.Write([]byte(err.Error()))
			},
			secret:              testSecret2,
			ts:                  testReqTsValid,
			containsRespPattern: "Expected signing signature:",
		},
		{
			description:   "using middleware and invalid signing signature, verify fails and req killed, no/empty fail method provided so no extra action",
			useMiddleware: true,
			invalidHex:    true,
		},
		{
			description:   "using middleware and invalid signing signature, verify fails and req killed. expected fail response received",
			useMiddleware: true,
			invalidHex:    true,
			ts:            testReqTsValid,
			failFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				_, _ = w.Write([]byte(err.Error()))
			},
			wantRespBody: "encoding/hex: odd length hex string",
		},
		{
			description: "not using middleware, command not found",
			wantErr:     errInteractionCallbackNotFound,
		},
	}

	wantCallback := &slack.InteractionCallback{
		Type:        "block_actions",
		Token:       "faketoken",
		ResponseURL: "https://hooks.slack.com/actions/T0000000F/000000000/YYYYYYYYYYY",
		TriggerID:   "0000000000.1111111111.222222222222aaaaaaaaaaaaaa",
		Team: slack.Team{
			ID:     "T0000000",
			Domain: "domain",
		},
		Channel: slack.Channel{
			GroupConversation: slack.GroupConversation{
				Conversation: slack.Conversation{
					ID: "G0000000",
				},
				Name: "privategroup",
			},
		},
		User: slack.User{
			ID:     "U12345678",
			TeamID: "T0000000",
			Name:   "fakenameyo",
		},
		ActionCallback: slack.ActionCallbacks{
			BlockActions: []*slack.BlockAction{
				{
					ActionID: "cancel_action",
					BlockID:  "channel_id_block",
					Type:     "button",
					Text: slack.TextBlockObject{
						Type:     "plain_text",
						Text:     "Done",
						Emoji:    true,
						Verbatim: false,
					},
					Value:    "done",
					ActionTs: "1589971722.911477",
				},
			},
		},
		APIAppID: "A00000000",
		Container: slack.Container{
			Type: "message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			r := chi.NewRouter()

			if tc.useMiddleware {
				r.Use(VerifyInteractionCallback(testSecret1, tc.succeedFunc, tc.failFunc))
			}

			signingSig := getSigningSig(t, tc.ts, tc.secret, []byte(testCallbackRaw))
			if tc.invalidHex {
				signingSig = testInvalidSigningSig
			}

			r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
				callback, err := InteractionCallback(r.Context())
				if err != nil {
					if err != tc.wantErr {
						t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
					}
					return
				}
				if diff := pretty.Compare(wantCallback, callback); diff != "" {
					t.Fatalf("+got -want %s\n", diff)
				}
			})

			testServ := httptest.NewServer(r)
			defer testServ.Close()

			respBodyString := executeTestReq(t, testServ, signingSig, tc.ts, testCallbackRaw)

			if tc.containsRespPattern != "" {
				if !strings.Contains(respBodyString, tc.containsRespPattern) {
					t.Fatalf("expected resp to contain pattern: %s, got: %s", tc.containsRespPattern, respBodyString)
				}
				return
			}

			if respBodyString != tc.wantRespBody {
				t.Fatalf("expected resp body: %s, got: %s", tc.wantRespBody, respBodyString)
			}
		})
	}
}

func executeTestReq(t *testing.T, testServ *httptest.Server, signingSig, ts string, encodedBody string) string {
	req, err := http.NewRequest(http.MethodPost, testServ.URL+"/test", strings.NewReader(encodedBody))
	if err != nil {
		t.Fatal("failed to create new http request", err)
	}

	req.Header.Set("X-Slack-Signature", signingSig)
	req.Header.Set("X-Slack-Request-Timestamp", ts)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("failed to execute http request", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("failed to read http response body", err)
	}

	return string(respBody)
}

func getSigningSig(t *testing.T, timestamp, secret string, reqBody []byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	if _, err := hash.Write([]byte(fmt.Sprintf("v0:%s:", timestamp))); err != nil {
		t.Fatal("failed writing test hash", err)
	}

	if _, err := hash.Write(reqBody); err != nil {
		t.Fatal("failed writing test hash", err)
	}

	return fmt.Sprintf("v0=%s", hex.EncodeToString(hash.Sum(nil)))
}
