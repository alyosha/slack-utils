package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestEmailsToSlackIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description   string
		respUsersList []byte
		wantErr       string
		emails        []string
		wantIDs       []string
	}{
		{
			description:   "successful retrieval of member emails",
			respUsersList: []byte(mockUsersListResp),
			emails:        []string{"spengler@ghostbusters.example.com", "glenda@south.oz.coven"},
			wantIDs:       []string{"U0G9QF9C6", "W07QCRPA4"},
		},
		{
			description:   "failure to retrieve users list",
			respUsersList: []byte(mockUsersListErrResp),
			wantErr:       "c.getAll() > c.Client.GetUsers() > invalid_cursor",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/users.list", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respUsersList)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				Client: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			ids, err := client.EmailsToSlackIDs(tc.emails)

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

			if len(ids) != len(tc.wantIDs) {
				t.Fatalf("expected to receive %v ids, got %v instead", len(tc.wantIDs), len(ids))
			}

			for i, id := range ids {
				if tc.wantIDs[i] != id {
					t.Fatalf("expected to receive id: %v, got: %v", tc.wantIDs[i], id)
				}
			}
		})
	}
}

func TestEmailsToSlackIDsInclusive(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description   string
		respUsersList []byte
		wantErr       string
		emails        []string
		wantIDs       []string
	}{
		{
			description:   "successful retrieval of member emails",
			respUsersList: []byte(mockUsersListResp),
			emails:        []string{"spengler@ghostbusters.example.com", "glenda@south.oz.coven"},
			wantIDs:       []string{"U0G9QF9C6", "W07QCRPA4"},
		},
		{
			description:   "failure to retrieve users list",
			respUsersList: []byte(mockUsersListErrResp),
			wantErr:       "c.getAll() > c.Client.GetUsers() > invalid_cursor",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/users.list", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respUsersList)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				Client: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			idEmailPairs, err := client.EmailsToSlackIDsInclusive(tc.emails)

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
				return
			}

			if len(idEmailPairs) != len(tc.wantIDs) {
				t.Fatalf("expected to receive %v ids, got %v instead", len(tc.wantIDs), len(idEmailPairs))
			}

			for i, idEmailPair := range idEmailPairs {
				if tc.emails[i] != idEmailPair[0] {
					t.Fatalf("expected email: %v, got: %v", tc.emails[1], idEmailPair[0])
				}
				if tc.wantIDs[i] != idEmailPair[1] {
					t.Fatalf("expected id: %v, got: %v", tc.wantIDs[1], idEmailPair[1])
				}
			}
		})
	}
}
