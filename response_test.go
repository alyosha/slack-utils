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
			executeTime: 100 * time.Millisecond,
		},
		{
			description: "response execution takes longer than configured global timeout, log message sent",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 5 * time.Millisecond,
				WarnDeadlineExceeded:  true,
			},
			wantTimeout: true,
			executeTime: 100 * time.Millisecond,
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
				GlobalResponseTimeout: 2 * time.Second,
				ResponseTimeoutMap: map[string]time.Duration{
					"/send_message": 10 * time.Millisecond,
				},
			},
			wantTimeout: true,
			executeTime: 100 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(st *testing.T) {
			st.Parallel()

			logReceivedCh := make(chan struct{}, 1)

			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				logReceivedCh <- struct{}{}
				_, _ = w.Write([]byte(mockSuccessResp))
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

			respDoneCh := make(chan struct{}, 1)

			ctxCancelledCh := make(chan struct{}, 1)

			sendMsg := func(ctx context.Context, cmd *slack.SlashCommand) {
				time.Sleep(tc.executeTime)
				if ctx.Err() != nil {
					ctxCancelledCh <- struct{}{}
					return
				}
				respDoneCh <- struct{}{}
			}

			r.Post("/send_message", func(w http.ResponseWriter, r *http.Request) {
				cmd, err := SlashCommand(r.Context())
				if err != nil {
					t.Fatal("unexpected error")
				}
				client.RespondSlash(r, sendMsg, cmd)
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
				if tc.responseCfg.WarnDeadlineExceeded {
					select {
					case <-logReceivedCh:
					case <-time.After(5 * time.Second):
						t.Fatal("expected warning log message, did not receive")
					}
					return
				}
				select {
				case <-logReceivedCh:
					t.Fatal("unexpected log message")
				case <-time.After(100 * time.Millisecond):
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
			executeTime: 100 * time.Millisecond,
		},
		{
			description: "response execution takes longer than configured global timeout, log message sent",
			responseCfg: ResponseConfig{
				GlobalResponseTimeout: 7 * time.Millisecond,
				WarnDeadlineExceeded:  true,
			},
			wantTimeout: true,
			executeTime: 100 * time.Millisecond,
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
				GlobalResponseTimeout: 2 * time.Second,
				ResponseTimeoutMap: map[string]time.Duration{
					"/callback": 10 * time.Millisecond,
				},
			},
			wantTimeout: true,
			executeTime: 100 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(st *testing.T) {
			st.Parallel()

			logReceivedCh := make(chan struct{}, 1)

			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				logReceivedCh <- struct{}{}
				_, _ = w.Write([]byte(mockSuccessResp))
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

			respDoneCh := make(chan struct{}, 1)

			ctxCancelledCh := make(chan struct{}, 1)

			sendMsg := func(ctx context.Context, cmd *slack.InteractionCallback) {
				time.Sleep(tc.executeTime)
				if ctx.Err() != nil {
					ctxCancelledCh <- struct{}{}
					return
				}
				respDoneCh <- struct{}{}
			}

			r.Post("/callback", func(w http.ResponseWriter, r *http.Request) {
				callback, err := InteractionCallback(r.Context())
				if err != nil {
					t.Fatal("unexpected error", err)
				}
				client.RespondCallback(r, sendMsg, callback)
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
				if tc.responseCfg.WarnDeadlineExceeded {
					select {
					case <-logReceivedCh:
					case <-time.After(5 * time.Second):
						t.Fatal("expected warning log message, did not receive")
					}
					return
				}
				select {
				case <-logReceivedCh:
					t.Fatal("unexpected log message")
				case <-time.After(100 * time.Millisecond):
				}
			}
		})
	}
}
