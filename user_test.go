package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestEmailsToSlackIDs(t *testing.T) {
	testCases := []struct {
		description   string
		respUsersList []byte
		wantErr       string
		emails        []string
		wantIDs       []string
	}{
		{
			description:   "successful retrieval of member emails",
			respUsersList: []byte(usersListResp),
			emails:        []string{"spengler@ghostbusters.example.com", "glenda@south.oz.coven"},
			wantIDs:       []string{"U0G9QF9C6", "W07QCRPA4"},
		},
		{
			description:   "failure to retrieve users list",
			respUsersList: []byte(usersListErrResp),
			wantErr:       "invalid_cursor",
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

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

			ids, err := EmailsToSlackIDs(client, tc.emails)

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
				return
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
					return
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
					return
				}
			}

			if len(ids) != len(tc.wantIDs) {
				t.Fatalf("expected to receive %v ids, got %v instead", len(tc.wantIDs), len(ids))
				return
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
	testCases := []struct {
		description   string
		respUsersList []byte
		wantErr       string
		emails        []string
		wantIDs       []string
	}{
		{
			description:   "successful retrieval of member emails",
			respUsersList: []byte(usersListResp),
			emails:        []string{"spengler@ghostbusters.example.com", "glenda@south.oz.coven"},
			wantIDs:       []string{"U0G9QF9C6", "W07QCRPA4"},
		},
		{
			description:   "failure to retrieve users list",
			respUsersList: []byte(usersListErrResp),
			wantErr:       "invalid_cursor",
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

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

			idEmailPairs, err := EmailsToSlackIDsInclusive(client, tc.emails)

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
				return
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
					return
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
					return
				}
				return
			}

			if len(idEmailPairs) != len(tc.wantIDs) {
				t.Fatalf("expected to receive %v ids, got %v instead", len(tc.wantIDs), len(idEmailPairs))
				return
			}

			for i, idEmailPair := range idEmailPairs {
				if tc.emails[i] != idEmailPair[0] {
					t.Fatalf("expected email: %v, got: %v", tc.emails[1], idEmailPair[0])
					return
				}
				if tc.wantIDs[i] != idEmailPair[1] {
					t.Fatalf("expected id: %v, got: %v", tc.wantIDs[1], idEmailPair[1])
				}
			}
		})
	}
}
