package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestPostMsg(t *testing.T) {
	testCases := []struct {
		description   string
		msg           Msg
		respPostMsg   []byte
		wantTS        string
		wantChannelID string
		wantErr       string
	}{
		{
			description:   "successfully posted message",
			msg:           Msg{Body: "Hey!"},
			respPostMsg:   []byte(mockPostMsgResp),
			wantTS:        "1503435956.000247",
			wantChannelID: "C1H9RESGL",
		},
		{
			description: "failure to post message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgErrResp),
			wantErr:     "invalid_blocks",
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

			Client := &Client{
				Client: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			ts, err := Client.PostMsg(Msg{}, "C1H9RESGL")

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

			if ts != tc.wantTS {
				t.Fatalf("expected timestamp: %s, got: %s", tc.wantTS, ts)
			}
		})
	}
}

func TestPostThreadMsg(t *testing.T) {
	testCases := []struct {
		description string
		msg         Msg
		respPostMsg []byte
		wantErr     string
	}{
		{
			description: "successfully posted message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgResp),
		},
		{
			description: "failure to post message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgErrResp),
			wantErr:     "invalid_blocks",
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

			Client := &Client{
				Client: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			err := Client.PostThreadMsg(Msg{}, "C1H9RESGL", "1503435956.000247")

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

func TestPostEphemeralMsg(t *testing.T) {
	testCases := []struct {
		description   string
		msg           Msg
		respPostMsg   []byte
		wantTS        string
		wantChannelID string
		wantErr       string
	}{
		{
			description:   "successfully posted ephemeral message",
			msg:           Msg{Body: "Hey!"},
			respPostMsg:   []byte(mockPostMsgResp),
			wantTS:        "1503435956.000247",
			wantChannelID: "C1H9RESGL",
		},
		{
			description: "failure to post ephemeral message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgErrResp),
			wantErr:     "invalid_blocks",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postEphemeral", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respPostMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			Client := &Client{
				Client: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			err := Client.PostEphemeralMsg(Msg{}, "C1H9RESGL", "U12345")

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

func TestUpdateMsg(t *testing.T) {
	testCases := []struct {
		description   string
		msg           Msg
		respUpdateMsg []byte
		wantErr       string
	}{
		{
			description:   "successfully posted message",
			msg:           Msg{Body: "Hey!"},
			respUpdateMsg: []byte(mockUpdateMsgResp),
		},
		{
			description:   "failure to post message",
			msg:           Msg{Body: "Hey!"},
			respUpdateMsg: []byte(mockPostMsgErrResp),
			wantErr:       "invalid_blocks",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/chat.update", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respUpdateMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			Client := &Client{
				Client: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			err := Client.UpdateMsg(Msg{}, "C1H9RESGL", "1503435957.000237")

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

func TestDeleteMsg(t *testing.T) {
	testCases := []struct {
		description   string
		respDeleteMsg []byte
		wantErr       string
	}{
		{
			description:   "successfully posted message",
			respDeleteMsg: []byte(mockSuccessResp),
		},
		{
			description:   "failure to post message",
			respDeleteMsg: []byte(mockPostMsgErrResp),
			wantErr:       "invalid_blocks",
		},
	}

	responseURLPath := "/commands/XXXXXXXX/00000000/YYYYYYYY"
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(responseURLPath, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(tc.respDeleteMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			Client := &Client{
				Client: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			err := Client.DeleteMsg("C1H9RESGL", "1503435957.000237", fmt.Sprintf("%s%s", testServ.URL, responseURLPath))

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
