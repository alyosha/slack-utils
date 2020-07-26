package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestCreateConversation(t *testing.T) {
	testCases := []struct {
		description       string
		inviteMembers     []string
		initMsg           Msg
		respChannelCreate []byte
		respInviteMembers []byte
		respPostMsg       []byte
		wantID            string
		wantErr           string
	}{
		{
			description:       "successful conversation creation, no additional invites",
			inviteMembers:     []string{},
			initMsg:           Msg{},
			respChannelCreate: []byte(mockChannelCreateResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "err creating conversation",
			inviteMembers:     []string{},
			initMsg:           Msg{},
			respChannelCreate: []byte(mockChannelCreateErrResp),
			wantErr:           "c.SlackAPI.CreateConversation > invalid_name_specials",
		},
		{
			description:       "successful conversation creation including additional invites",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{},
			respChannelCreate: []byte(mockChannelCreateResp),
			respInviteMembers: []byte(mockInviteMembersResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "successful conversation creation, inviting members including self but no error returned",
			inviteMembers:     []string{"U0G9QF9C6"},
			initMsg:           Msg{},
			respChannelCreate: []byte(mockChannelCreateResp),
			respInviteMembers: []byte(mockCantInviteSelfErrResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "err inviting members",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{},
			respChannelCreate: []byte(mockChannelCreateResp),
			respInviteMembers: []byte(mockInviteMembersErrResp),
			wantErr:           "c.InviteUsers > c.SlackAPI.InviteUsersToConversation > cant_invite",
		},
		{
			description:       "successful conversation creation including additional invites, successful message post",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{Body: "Hey!"},
			respChannelCreate: []byte(mockChannelCreateResp),
			respInviteMembers: []byte(mockInviteMembersResp),
			respPostMsg:       []byte(mockPostMsgResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "successful conversation creation including additional invites, failure to post init message",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{Body: "Hey!"},
			respChannelCreate: []byte(mockChannelCreateResp),
			respInviteMembers: []byte(mockInviteMembersResp),
			respPostMsg:       []byte(mockPostMsgErrResp),
			wantErr:           "c.SlackAPI.PostMessage > invalid_blocks",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/conversations.create", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respChannelCreate)
			})
			mux.HandleFunc("/conversations.invite", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respInviteMembers)
			})
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respPostMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				SlackAPI: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			conversationID, err := client.CreateConversation("general", false, tc.inviteMembers, tc.initMsg)

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

			if conversationID != tc.wantID {
				t.Fatalf("expected conversation id: %s, got: %s", tc.wantID, conversationID)
			}
		})
	}
}

func TestInviteUsers(t *testing.T) {
	testCases := []struct {
		description       string
		inviteMembers     []string
		respInviteMembers []byte
		wantErr           string
	}{
		{
			description:       "successful invite of users",
			inviteMembers:     []string{"UABC123EFG"},
			respInviteMembers: []byte(mockInviteMembersResp),
		},
		{
			description:       "successful invite, no error returned for invite members resp",
			inviteMembers:     []string{"UABC123EFG", "U0G9QF9C6"},
			respInviteMembers: []byte(mockInviteMembersResp),
		},
		{
			description:       "expect error",
			inviteMembers:     []string{"UABC123EFG"},
			respInviteMembers: []byte(mockInviteMembersErrResp),
			wantErr:           "c.SlackAPI.InviteUsersToConversation > cant_invite",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/conversations.invite", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respInviteMembers)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				SlackAPI: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			err := client.InviteUsers("C1H9RESGL", tc.inviteMembers)

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

func TestGetConversationMembers(t *testing.T) {
	testCases := []struct {
		description          string
		respConversationInfo []byte
		wantErr              string
		wantIDs              []string
	}{
		{
			description:          "successful retrieval of member IDs",
			respConversationInfo: []byte(mockChannelInfoResp),
			wantIDs:              []string{"U0G9QF9C6", "U1QNSQB9U"},
		},
		{
			description:          "failure to retrieve member IDs",
			respConversationInfo: []byte(mockChannelInfoErrResp),
			wantErr:              "c.SlackAPI.GetConversationInfo > channel_not_found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/conversations.info", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respConversationInfo)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				SlackAPI: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			members, err := client.GetConversationMembers("C1H9RESGL")

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

			if len(members) != len(tc.wantIDs) {
				t.Fatalf("expected to receive %v ids, got %v instead", len(tc.wantIDs), len(members))
			}

			for i, member := range members {
				if tc.wantIDs[i] != member {
					t.Fatalf("expected to receive id: %v, got: %v", tc.wantIDs[i], member)
				}
			}
		})
	}
}

func TestGetConversationMemberEmails(t *testing.T) {
	testCases := []struct {
		description          string
		respConversationInfo []byte
		respUsersList        []byte
		wantErr              string
		wantEmails           []string
	}{
		{
			description:          "successful retrieval of member emails",
			respConversationInfo: []byte(mockChannelInfoResp),
			respUsersList:        []byte(mockUsersListResp),
			wantEmails:           []string{"spengler@ghostbusters.example.com"},
		},
		{
			description:          "failure to retrieve conversation info",
			respConversationInfo: []byte(mockChannelInfoErrResp),
			respUsersList:        []byte(mockUsersListResp),
			wantErr:              "c.SlackAPI.GetConversationInfo > channel_not_found",
		},
		{
			description:          "failure to retrieve user list",
			respConversationInfo: []byte(mockChannelInfoResp),
			respUsersList:        []byte(mockUsersListErrResp),
			wantErr:              "c.getAll > c.SlackAPI.GetUsers > invalid_cursor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/conversations.info", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respConversationInfo)
			})
			mux.HandleFunc("/users.list", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respUsersList)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				SlackAPI: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			emails, err := client.GetConversationMemberEmails("C1H9RESGL")

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

			if len(emails) != len(tc.wantEmails) {
				t.Fatalf("expected to receive %v emails, got %v instead", len(tc.wantEmails), len(emails))
			}

			for i, email := range emails {
				if tc.wantEmails[i] != email {
					t.Fatalf("expected to receive email: %v, got: %v", tc.wantEmails[i], email)
				}
			}
		})
	}
}

func TestArchiveConversations(t *testing.T) {
	testCases := []struct {
		description         string
		respArchiveChannels []byte
		wantErr             string
	}{
		{
			description:         "successfully archived conversations",
			respArchiveChannels: []byte(mockSuccessResp),
		},
		{
			description:         "no error returned for already archived conversation",
			respArchiveChannels: []byte(mockChannelAlreadyArchivedErrResp),
		},
		{
			description:         "failure to archive conversations",
			respArchiveChannels: []byte(mockChannelsArchiveErrResp),
			wantErr:             "c.SlackAPI.ArchiveConversation > invalid_auth",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/conversations.archive", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respArchiveChannels)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := &Client{
				SlackAPI: slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL))),
			}

			err := client.ArchiveConversations([]string{"C1H9RESGL"})

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
