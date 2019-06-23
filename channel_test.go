package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nlopes/slack"
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
			respChannelCreate: []byte(channelCreateResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "err creating channel",
			inviteMembers:     []string{},
			initMsg:           Msg{},
			respChannelCreate: []byte(channelCreateErrResp),
			wantErr:           "invalid_name_specials",
		},
		{
			description:       "successful channel creation including additional invites",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(inviteMembersResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "successful channel creation, inviting members including self but no error returned",
			inviteMembers:     []string{"U0G9QF9C6"},
			initMsg:           Msg{},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(cantInviteSelfErrResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "err inviting members",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(inviteMembersErrResp),
			wantErr:           "cant_invite",
		},
		{
			description:       "successful channel creation including additional invites, successful message post",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{Body: "Hey!"},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(inviteMembersResp),
			respPostMsg:       []byte(postMsgResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "successful channel creation including additional invites, failure to post init message",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           Msg{Body: "Hey!"},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(inviteMembersResp),
			respPostMsg:       []byte(postMsgErrResp),
			wantErr:           "too_many_attachments",
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
				BotClient:  nil,
			}

			err := channel.CreateChannel("general", tc.inviteMembers, tc.initMsg, false)

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
			respInviteMembers: []byte(inviteMembersResp),
		},
		{
			description:       "successful invite, no error returned for invite members resp",
			inviteMembers:     []string{"UABC123EFG", "U0G9QF9C6"},
			respInviteMembers: []byte(inviteMembersResp),
		},
		{
			description:       "expect error",
			inviteMembers:     []string{"UABC123EFG"},
			respInviteMembers: []byte(inviteMembersErrResp),
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
				BotClient:  nil,
				ChannelID:  "C1H9RESGL",
			}

			err := channel.InviteUsers(tc.inviteMembers)

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
			respChannelInfo: []byte(channelInfoResp),
			wantIDs:         []string{"U0G9QF9C6", "U1QNSQB9U"},
		},
		{
			description:     "failure to retrieve member IDs",
			respChannelInfo: []byte(channelInfoErrResp),
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
			channel := Channel{
				UserClient: client,
				ChannelID:  "C1H9RESGL",
			}

			members, err := channel.GetChannelMembers()

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

			if len(members) != len(tc.wantIDs) {
				t.Fatalf("expected to receive %v ids, got %v instead", len(tc.wantIDs), len(members))
				return
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
			respChannelInfo: []byte(channelInfoResp),
			respUsersList:   []byte(usersListResp),
			wantEmails:      []string{"spengler@ghostbusters.example.com"},
		},
		{
			description:     "failure to retrieve channel info",
			respChannelInfo: []byte(channelInfoErrResp),
			respUsersList:   []byte(usersListResp),
			wantErr:         "channel_not_found",
		},
		{
			description:     "failure to retrieve user list",
			respChannelInfo: []byte(channelInfoResp),
			respUsersList:   []byte(usersListErrResp),
			wantErr:         "no users in workplace",
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
			slackClient := Slack{
				Client: client,
			}

			emails, err := slackClient.GetChannelMemberEmails("C1H9RESGL")

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

			if len(emails) != len(tc.wantEmails) {
				t.Fatalf("expected to receive %v emails, got %v instead", len(tc.wantEmails), len(emails))
				return
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
			respLeaveChannels: []byte(channelsLeaveResp),
		},
		{
			description:       "failure to leave channels",
			respLeaveChannels: []byte(channelsLeaveErrResp),
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
			slackClient := Slack{
				Client: client,
			}

			err := slackClient.LeaveChannels([]string{"C1H9RESGL"})

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
			respArchiveChannels: []byte(channelsArchiveResp),
		},
		{
			description:         "no error returned for already archived channel",
			respArchiveChannels: []byte(channelAlreadyArchivedErrResp),
		},
		{
			description:         "failure to archive channels",
			respArchiveChannels: []byte(channelsArchiveErrResp),
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
			slackClient := Slack{
				Client: client,
			}

			err := slackClient.ArchiveChannels([]string{"C1H9RESGL"})

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
				}
			}
		})
	}
}
