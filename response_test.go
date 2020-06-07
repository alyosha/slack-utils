package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/slack-go/slack"
)

func TestRespondSlash(t *testing.T) {
	t.Parallel()

	body := url.Values{
		"command":         []string{"/send-message"},
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

	testCases := []struct {
		description string
		responseCfg ResponseConfig
		executeTime time.Duration
		wantTimeout bool
	}{
		{
			description: "default config allows response to continue as expected",
			executeTime: 10 * time.Millisecond,
		},
		{
			description: "response execution takes longer than configured global timeout, no log message",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 5 * time.Millisecond,
			},
			wantTimeout: true,
			executeTime: 10 * time.Millisecond,
		},
		{
			description: "response execution takes longer than configured global timeout, log message sent",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 5 * time.Millisecond,
				WarnDeadlineExceeded:  true,
			},
			wantTimeout: true,
			executeTime: 10 * time.Millisecond,
		},
		{
			description: "response execution takes longer than configured global timeout, but overwritten by timeout map so no timeout",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 5 * time.Millisecond,
				ResponseTimeoutMap: map[string]time.Duration{
					"/send_message": 100 * time.Millisecond,
				},
			},
			executeTime: 10 * time.Millisecond,
		},
		{
			description: "global timeout would cover execution time, but overwritten by timeout map so timeout occurs",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 30 * time.Millisecond,
				ResponseTimeoutMap: map[string]time.Duration{
					"/send_message": 5 * time.Millisecond,
				},
			},
			wantTimeout: true,
			executeTime: 10 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc := tc
			t.Parallel()

			logReceivedCh := make(chan struct{})

			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				logReceivedCh <- struct{}{}
				_, _ = w.Write([]byte("ok"))
			})

			testServSlack := httptest.NewServer(mux)
			defer testServSlack.Close()

			client := &Client{
				Client:              slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServSlack.URL))),
				slashResponseConfig: tc.responseCfg,
				errChannel:          "C1H9RESGL",
			}

			r := chi.NewRouter()

			r.Use(client.VerifySlashCommand(testSecret1))

			signingSig := getTestSigningSig(t, testReqTsValid, testSecret1, []byte(encodedBody))

			originalReqDoneCh := make(chan struct{})
			respDoneCh := make(chan struct{})

			ctxCancelledCh := make(chan struct{})

			var sendMsg = func(ctx context.Context, cmd *slack.SlashCommand) {
				_ = <-originalReqDoneCh

				select {
				case <-time.After(tc.executeTime):
					respDoneCh <- struct{}{}
				case <-ctx.Done():
					ctxCancelledCh <- struct{}{}
				}
			}

			r.Post("/send_message", func(w http.ResponseWriter, r *http.Request) {
				cmd, err := SlashCommand(r.Context())
				if err != nil {
					t.Fatal("unexpected error")
				}
				client.RespondSlash(r, sendMsg, cmd)
				originalReqDoneCh <- struct{}{}
			})

			testServ := httptest.NewServer(r)
			defer testServ.Close()

			_ = executeTestReq(t, testServ, signingSig, testReqTsValid, "/send_message", encodedBody)

			select {
			case <-respDoneCh:
				if tc.wantTimeout {
					t.Fatal("expected timeout")
				}
			case <-ctxCancelledCh:
				if !tc.wantTimeout {
					t.Fatal("unexpected timeout")
				}
				select {
				case <-logReceivedCh:
					if !tc.responseCfg.WarnDeadlineExceeded {
						t.Fatal("unexpected warning log message")
					}
				case <-time.After(1500 * time.Millisecond): // buffer for log message to go through
					if tc.responseCfg.WarnDeadlineExceeded {
						t.Fatal("expected warning log message, did not receive")
					}
				}
			}
		})
	}
}

func TestRespondCallback(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		responseCfg ResponseConfig
		executeTime time.Duration
		wantTimeout bool
	}{
		{
			description: "default config allows response to continue as expected",
			executeTime: 10 * time.Millisecond,
		},
		{
			description: "response execution takes longer than configured global timeout, no log message",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 5 * time.Millisecond,
			},
			wantTimeout: true,
			executeTime: 10 * time.Millisecond,
		},
		{
			description: "response execution takes longer than configured global timeout, log message sent",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 5 * time.Millisecond,
				WarnDeadlineExceeded:  true,
			},
			wantTimeout: true,
			executeTime: 10 * time.Millisecond,
		},
		{
			description: "response execution takes longer than configured global timeout, but overwritten by timeout map so no timeout",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 5 * time.Millisecond,
				ResponseTimeoutMap: map[string]time.Duration{
					"/callback": 100 * time.Millisecond,
				},
			},
			executeTime: 10 * time.Millisecond,
		},
		{
			description: "global timeout would cover execution time, but overwritten by timeout map so timeout occurs",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 30 * time.Millisecond,
				ResponseTimeoutMap: map[string]time.Duration{
					"/callback": 5 * time.Millisecond,
				},
			},
			wantTimeout: true,
			executeTime: 10 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc := tc
			t.Parallel()

			logReceivedCh := make(chan struct{})

			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				logReceivedCh <- struct{}{}
				_, _ = w.Write([]byte("ok"))
			})

			testServSlack := httptest.NewServer(mux)
			defer testServSlack.Close()

			client := &Client{
				Client:                 slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServSlack.URL))),
				callbackResponseConfig: tc.responseCfg,
				errChannel:             "C1H9RESGL",
			}

			r := chi.NewRouter()

			r.Use(client.VerifyInteractionCallback(testSecret1))

			signingSig := getTestSigningSig(t, testReqTsValid, testSecret1, []byte(mockCallbackRaw))

			originalReqDoneCh := make(chan struct{})
			respDoneCh := make(chan struct{})

			ctxCancelledCh := make(chan struct{})

			var sendMsg = func(ctx context.Context, cmd *slack.InteractionCallback) {
				_ = <-originalReqDoneCh

				select {
				case <-time.After(tc.executeTime):
					respDoneCh <- struct{}{}
				case <-ctx.Done():
					ctxCancelledCh <- struct{}{}
				}
			}

			r.Post("/callback", func(w http.ResponseWriter, r *http.Request) {
				callback, err := InteractionCallback(r.Context())
				if err != nil {
					t.Fatal("unexpected error", err)
				}
				client.RespondCallback(r, sendMsg, callback)
				originalReqDoneCh <- struct{}{}
			})

			testServ := httptest.NewServer(r)
			defer testServ.Close()

			_ = executeTestReq(t, testServ, signingSig, testReqTsValid, "/callback", mockCallbackRaw)

			select {
			case <-respDoneCh:
				if tc.wantTimeout {
					t.Fatal("expected timeout")
				}
			case <-ctxCancelledCh:
				if !tc.wantTimeout {
					t.Fatal("unexpected timeout")
				}
				select {
				case <-logReceivedCh:
					if !tc.responseCfg.WarnDeadlineExceeded {
						t.Fatal("unexpected warning log message")
					}
				case <-time.After(1500 * time.Millisecond): // buffer for log message to go through
					if tc.responseCfg.WarnDeadlineExceeded {
						t.Fatal("expected warning log message, did not receive")
					}
				}
			}
		})
	}
}
