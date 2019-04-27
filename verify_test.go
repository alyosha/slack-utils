package utils

import (
	"context"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nlopes/slack"
)

func TestVerifyCallbackMsg(t *testing.T) {
	ts := time.Now().Unix()
	tsStr := strconv.FormatInt(ts, 10)
	testCases := []struct {
		description string
		timestamp   string
		wantErr     string
	}{
		{
			description: "timestamp too old error",
			timestamp:   "1531431954",
			wantErr:     "timestamp is too old",
		},
		{
			description: "signing signature does not match",
			timestamp:   tsStr,
			wantErr:     "Expected signing signature: adada4ed31709aef585c2580ca3267678c6a8eaeb7e0c1aca3ee57b656886b2c, but computed:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest("POST", "http://example.com/foo", nil)
			ctx := WithContext(context.Background(), "e6b19c573432dcc6b075501d51b51bb8", &slack.Client{})
			req = req.WithContext(ctx)
			req.Header.Set("X-Slack-Signature", "v0=adada4ed31709aef585c2580ca3267678c6a8eaeb7e0c1aca3ee57b656886b2c")
			req.Header.Set("X-Slack-Request-Timestamp", tc.timestamp)

			_, err := VerifyCallbackMsg(req)

			if tc.wantErr != "" && err == nil {
				t.Fatal("expected timestamp too old error, didn't receive one")
				return
			}

			if err.Error() != tc.wantErr && !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error: %s, got: %s", tc.wantErr, err)
			}
		})
	}
}

func TestVerifySlashCmd(t *testing.T) {
	ts := time.Now().Unix()
	tsStr := strconv.FormatInt(ts, 10)
	testCases := []struct {
		description string
		timestamp   string
		wantErr     string
	}{
		{
			description: "timestamp too old error",
			timestamp:   "1531431954",
			wantErr:     "timestamp is too old",
		},
		{
			description: "signing signature does not match",
			timestamp:   tsStr,
			wantErr:     "Expected signing signature: adada4ed31709aef585c2580ca3267678c6a8eaeb7e0c1aca3ee57b656886b2c, but computed:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest("POST", "http://example.com/foo", nil)
			ctx := WithContext(context.Background(), "e6b19c573432dcc6b075501d51b51bb8", &slack.Client{})
			req = req.WithContext(ctx)
			req.Header.Set("X-Slack-Signature", "v0=adada4ed31709aef585c2580ca3267678c6a8eaeb7e0c1aca3ee57b656886b2c")
			req.Header.Set("X-Slack-Request-Timestamp", tc.timestamp)

			_, err := VerifySlashCmd(req)

			if tc.wantErr != "" && err == nil {
				t.Fatal("expected timestamp too old error, didn't receive one")
				return
			}

			if err.Error() != tc.wantErr && !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error: %s, got: %s", tc.wantErr, err)
			}
		})
	}
}
