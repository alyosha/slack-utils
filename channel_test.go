package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestCreateChannel(t *testing.T) {
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
			description:       "successful channel creation, no additional invites",
			inviteMembers:     []string{},
			initMsg:           Msg{},
			respChannelCreate: []byte(mockChannelCreateResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "err creating channel",
			inviteMembers:     []string{},
			initMsg:           Msg{},
			respChannelCreate: []byte(mockChannelCreateErrResp),
			wantErr:           "failed to create new channel: invalid_name_specials",
		},
		{
			description:       "successful channel creation including additional invites",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{},
			respChannelCreate: []byte(mockChannelCreateResp),
			respInviteMembers: []byte(mockInviteMembersResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "successful channel creation, inviting members including self but no error returned",
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
			wantErr:           "failed to invite user to channel: cant_invite",
		},
		{
			description:       "successful channel creation including additional invites, successful message post",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{Body: "Hey!"},
			respChannelCreate: []byte(mockChannelCreateResp),
			respInviteMembers: []byte(mockInviteMembersResp),
			respPostMsg:       []byte(mockPostMsgResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "successful channel creation including additional invites, failure to post init message",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{Body: "Hey!"},
			respChannelCreate: []byte(mockChannelCreateResp),
			respInviteMembers: []byte(mockInviteMembersResp),
			respPostMsg:       []byte(mockPostMsgErrResp),
			wantErr:           "failed to post message: too_many_attachments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/channels.create", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respChannelCreate)
			})
			mux.HandleFunc("/channels.invite", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respInviteMembers)
			})
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respPostMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))
			channel := Channel{
				UserClient: client,
			}

			err := channel.CreateChannel("general", tc.inviteMembers, tc.initMsg)

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

			if channel.ChannelID != tc.wantID {
				t.Fatalf("expected channel id: %s, got: %s", tc.wantID, channel.ChannelID)
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
			wantErr:           "cant_invite",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/channels.invite", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respInviteMembers)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))
			channel := Channel{
				UserClient: client,
				ChannelID:  "C1H9RESGL",
			}

			err := channel.InviteUsers(tc.inviteMembers)

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

func TestGetChannelMembers(t *testing.T) {
	testCases := []struct {
		description     string
		respChannelInfo []byte
		wantErr         string
		wantIDs         []string
	}{
		{
			description:     "successful retrieval of member IDs",
			respChannelInfo: []byte(mockChannelInfoResp),
			wantIDs:         []string{"U0G9QF9C6", "U1QNSQB9U"},
		},
		{
			description:     "failure to retrieve member IDs",
			respChannelInfo: []byte(mockChannelInfoErrResp),
			wantErr:         "channel_not_found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/channels.info", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respChannelInfo)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

			members, err := GetChannelMembers(client, "C1H9RESGL")

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

func TestGetChannelMemberEmails(t *testing.T) {
	testCases := []struct {
		description     string
		respChannelInfo []byte
		respUsersList   []byte
		wantErr         string
		wantEmails      []string
	}{
		{
			description:     "successful retrieval of member emails",
			respChannelInfo: []byte(mockChannelInfoResp),
			respUsersList:   []byte(mockUsersListResp),
			wantEmails:      []string{"spengler@ghostbusters.example.com"},
		},
		{
			description:     "failure to retrieve channel info",
			respChannelInfo: []byte(mockChannelInfoErrResp),
			respUsersList:   []byte(mockUsersListResp),
			wantErr:         "channel_not_found",
		},
		{
			description:     "failure to retrieve user list",
			respChannelInfo: []byte(mockChannelInfoResp),
			respUsersList:   []byte(mockUsersListErrResp),
			wantErr:         "invalid_cursor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/channels.info", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respChannelInfo)
			})
			mux.HandleFunc("/users.list", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respUsersList)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

			emails, err := GetChannelMemberEmails(client, "C1H9RESGL")

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

func TestLeaveChannels(t *testing.T) {
	testCases := []struct {
		description       string
		respLeaveChannels []byte
		wantErr           string
	}{
		{
			description:       "successfully left channels",
			respLeaveChannels: []byte(mockSuccessResp),
		},
		{
			description:       "failure to leave channels",
			respLeaveChannels: []byte(mockChannelsLeaveErrResp),
			wantErr:           "invalid_auth",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/channels.leave", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respLeaveChannels)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))
			channel := Channel{
				UserClient: client,
			}

			err := channel.LeaveChannels([]string{"C1H9RESGL"})

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

func TestArchiveChannels(t *testing.T) {
	testCases := []struct {
		description         string
		respArchiveChannels []byte
		wantErr             string
	}{
		{
			description:         "successfully archived channels",
			respArchiveChannels: []byte(mockSuccessResp),
		},
		{
			description:         "no error returned for already archived channel",
			respArchiveChannels: []byte(mockChannelAlreadyArchivedErrResp),
		},
		{
			description:         "failure to archive channels",
			respArchiveChannels: []byte(mockChannelsArchiveErrResp),
			wantErr:             "invalid_auth",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/channels.archive", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respArchiveChannels)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))
			channel := Channel{
				UserClient: client,
			}

			err := channel.ArchiveChannels([]string{"C1H9RESGL"})

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
