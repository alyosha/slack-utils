package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description          string
		cfg                  ClientConfig
		respConversationInfo []byte
		respUserInfo         []byte
		wantErr              string
	}{
		{
			description: "missing bot token",
			cfg:         ClientConfig{AdminID: "U1234"},
			wantErr:     "must provide bot token",
		},
		{
			description: "missing admin ID",
			cfg:         ClientConfig{BotToken: "B1234"},
			wantErr:     "must provide admin ID",
		},
		{
			description: "failure to retrieve admin user info",
			cfg: ClientConfig{
				AdminID:  "U1234",
				BotToken: "B1234",
			},
			respUserInfo:         []byte(mockUserInfoErrResp),
			respConversationInfo: []byte(mockSuccessResp),
			wantErr:              "c.SlackAPI.GetUserInfo > user_not_found",
		},
		{
			description: "failure to retrieve channel info for log/err channel",
			cfg: ClientConfig{
				AdminID:      "U1234",
				BotToken:     "B1234",
				LogChannelID: "C1234",
				ErrChannelID: "C2345",
			},
			respConversationInfo: []byte(mockChannelMembersErrResp),
			respUserInfo:         []byte(mockSuccessResp),
			wantErr:              "c.SlackAPI.GetConversationInfo > channel_not_found",
		},
		{
			description: "success - channel info not verified if ID missing from config",
			cfg: ClientConfig{
				AdminID:  "U1234",
				BotToken: "B1234",
			},
			respConversationInfo: []byte(mockChannelMembersErrResp),
			respUserInfo:         []byte(mockSuccessResp),
		},
		{
			description: "success - fully loaded client",
			cfg: ClientConfig{
				AdminID:      "U1234",
				BotToken:     "B1234",
				LogChannelID: "C1234",
				ErrChannelID: "C2345",
			},
			respConversationInfo: []byte(mockSuccessResp),
			respUserInfo:         []byte(mockSuccessResp),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/conversations.info", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respConversationInfo)
			})
			mux.HandleFunc("/users.info", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respUserInfo)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			_, err := NewClient(tc.cfg, slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

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
func TestSendResp(t *testing.T) {
	var msg slack.Message
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := SendResp(w, slack.Message{})
		if err != nil {
			t.Fatalf("unexpected error handing request: %s", err)
		}
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got: %v", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("content type application/json, got: %v", resp.Header.Get("Content-Type"))
	}

	err := json.Unmarshal(body, &msg)
	if err != nil {
		t.Fatalf("failed to unmarshal response with error: %s", err)
	}

	if msg.ReplaceOriginal {
		t.Fatal("replace original should be false, but is true")
	}

	if msg.DeleteOriginal {
		t.Fatal("delete original should be false, but is true")
	}
}
