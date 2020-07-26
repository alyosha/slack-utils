package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestSendToLogChannel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description       string
		msg               Msg
		respPostMsg       []byte
		missingLogChannel bool
		wantErr           string
	}{
		{
			description: "successfully posted message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgResp),
		},
		{
			description:       "failure because log channel not configured",
			msg:               Msg{Body: "Hey!"},
			respPostMsg:       []byte(mockPostMsgErrResp),
			missingLogChannel: true,
			wantErr:           errLogChannelNotConfigured.Error(),
		},
		{
			description: "failure to post error message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgErrResp),
			wantErr:     "c.PostMsg > invalid_blocks",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respPostMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				SlackAPI: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			if !tc.missingLogChannel {
				client.logChannel = "C1234"
			}

			err := client.SendToLogChannel(Msg{})

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
				}
			}
		})
	}
}

func TestSendToErrChannel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description       string
		msg               Msg
		respPostMsg       []byte
		respPostThreadMsg []byte
		wantErr           string
		missingErrChannel bool
		client            *Client
	}{
		{
			description: "successfully posted message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgResp),
		},
		{
			description:       "failure because err channel not confired",
			msg:               Msg{Body: "Hey!"},
			respPostMsg:       []byte(mockPostMsgErrResp),
			respPostThreadMsg: []byte(mockPostMsgResp),
			missingErrChannel: true,
			wantErr:           errLogChannelNotConfigured.Error(),
		},
		{
			description:       "failure to post message",
			msg:               Msg{Body: "Hey!"},
			respPostMsg:       []byte(mockPostMsgErrResp),
			respPostThreadMsg: []byte(mockPostMsgResp),
			wantErr:           "c.PostMsg > invalid_blocks",
		},
		{
			description:       "failure to post thread message",
			msg:               Msg{Body: "Hey!"},
			respPostMsg:       []byte(mockPostMsgResp),
			respPostThreadMsg: []byte(mockPostMsgErrResp),
			wantErr:           "c.PostThreadMsg > invalid_blocks",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			var postMsgCalls int

			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				postMsgCalls++
				if postMsgCalls > 1 {
					_, _ = w.Write(tc.respPostThreadMsg)
				}
				_, _ = w.Write(tc.respPostMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				SlackAPI: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			if !tc.missingErrChannel {
				client.errChannel = "C1234"
			}

			err := client.SendToErrChannel("", errors.New(""))

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
				}
			}
		})
	}
}
