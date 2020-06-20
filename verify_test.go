package utils

import (
	"fmt"
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
)

var (
	testReqTsValid = fmt.Sprintf("%d", time.Now().Unix())
)

func TestVerifySlashCommand(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description         string
		useMiddleware       bool
		secret              string
		ts                  string
		invalidHex          bool
		logConfig           RequestLoggingConfig
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
			description:   "same success case as above with request logging",
			useMiddleware: true,
			secret:        testSecret1,
			ts:            testReqTsValid,
			logConfig:     RequestLoggingConfig{Enabled: true},
		},
		{
			description:   "using middleware and valid signing signature, expected extra success response received",
			useMiddleware: true,
			succeedFunc: func(w http.ResponseWriter, r *http.Request, cmd *slack.SlashCommand) {
				_, _ = w.Write([]byte(mockSuccessResp))
			},
			secret:       testSecret1,
			ts:           testReqTsValid,
			wantRespBody: mockSuccessResp,
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

			mux := http.NewServeMux()

			var logReqReceived bool

			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				logReqReceived = true
				_, _ = w.Write([]byte(mockSuccessResp))
			})

			testServSlack := httptest.NewServer(mux)
			defer testServSlack.Close()

			client := &Client{
				Client:     slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServSlack.URL))),
				logChannel: "C1H9RESGL",
			}

			if tc.useMiddleware {
				r.Use(client.VerifySlashCommand(testSecret1, tc.logConfig, tc.succeedFunc, tc.failFunc))
			}

			signingSig := getTestSigningSig(t, tc.ts, tc.secret, []byte(encodedBody))
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
					t.Fatalf("-got +want %s\n", diff)
				}
			})

			testServ := httptest.NewServer(r)
			defer testServ.Close()

			respBodyString := executeTestReq(t, testServ, signingSig, tc.ts, "/test", encodedBody)

			if tc.containsRespPattern != "" {
				if !strings.Contains(respBodyString, tc.containsRespPattern) {
					t.Fatalf("expected resp to contain pattern: %s, got: %s", tc.containsRespPattern, respBodyString)
				}
				return
			}

			if respBodyString != tc.wantRespBody {
				t.Fatalf("expected resp body: %s, got: %s", tc.wantRespBody, respBodyString)
			}

			if tc.logConfig.Enabled && !logReqReceived {
				t.Fatal("expected log request call did not come")
			}
		})
	}
}

func TestVerifyInteractionCallback(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description         string
		useMiddleware       bool
		secret              string
		ts                  string
		invalidHex          bool
		logConfig           RequestLoggingConfig
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
			description:   "same success case as above, with request logging",
			useMiddleware: true,
			secret:        testSecret1,
			ts:            testReqTsValid,
			logConfig:     RequestLoggingConfig{Enabled: true},
		},
		{
			description:   "using middleware and valid signing signature, expected extra success response received",
			useMiddleware: true,
			succeedFunc: func(w http.ResponseWriter, r *http.Request, cmd *slack.InteractionCallback) {
				_, _ = w.Write([]byte(mockSuccessResp))
			},
			secret:       testSecret1,
			ts:           testReqTsValid,
			wantRespBody: mockSuccessResp,
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

			var logReqReceived bool

			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				logReqReceived = true
				_, _ = w.Write([]byte(mockSuccessResp))
			})

			testServSlack := httptest.NewServer(mux)
			defer testServSlack.Close()

			client := &Client{
				Client:     slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServSlack.URL))),
				logChannel: "C1H9RESGL",
			}

			if tc.useMiddleware {
				r.Use(client.VerifyInteractionCallback(testSecret1, tc.logConfig, tc.succeedFunc, tc.failFunc))
			}

			signingSig := getTestSigningSig(t, tc.ts, tc.secret, []byte(mockCallbackRaw))
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
					t.Fatalf("-got +want %s\n", diff)
				}
			})

			testServ := httptest.NewServer(r)
			defer testServ.Close()

			respBodyString := executeTestReq(t, testServ, signingSig, tc.ts, "/test", mockCallbackRaw)

			if tc.containsRespPattern != "" {
				if !strings.Contains(respBodyString, tc.containsRespPattern) {
					t.Fatalf("expected resp to contain pattern: %s, got: %s", tc.containsRespPattern, respBodyString)
				}
				return
			}

			if respBodyString != tc.wantRespBody {
				t.Fatalf("expected resp body: %s, got: %s", tc.wantRespBody, respBodyString)
			}

			if tc.logConfig.Enabled != logReqReceived {
				t.Fatalf("log request enabled: %v, log received: %v", tc.logConfig.Enabled, logReqReceived)
			}
		})
	}
}
