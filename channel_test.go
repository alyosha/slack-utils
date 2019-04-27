package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	utils "github.com/alyosha/slack-utils"
	"github.com/nlopes/slack"
)

func TestCreateChannel(t *testing.T) {
	testCases := []struct {
		description       string
		inviteMembers     []string
		initMsg           utils.Msg
		respChannelCreate []byte
		respInviteMembers []byte
		respPostMsg       []byte
		wantID            string
		wantErr           string
	}{
		{
			description:       "successful channel creation, no additional invites",
			inviteMembers:     []string{},
			initMsg:           utils.Msg{},
			respChannelCreate: []byte(channelCreateResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "err creating channel",
			inviteMembers:     []string{},
			initMsg:           utils.Msg{},
			respChannelCreate: []byte(channelCreateErrResp),
			wantErr:           "invalid_name_specials",
		},
		{
			description:       "successful channel creation including additional invites",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           utils.Msg{},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(inviteMembersResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "successful channel creation, inviting members including self but no error returned",
			inviteMembers:     []string{"U0G9QF9C6"},
			initMsg:           utils.Msg{},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(cantInviteSelfErrResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "err inviting members",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           utils.Msg{},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(inviteMembersErrResp),
			wantErr:           "cant_invite",
		},
		{
			description:       "successful channel creation including additional invites, successful message post",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           utils.Msg{Body: "Hey!"},
			respChannelCreate: []byte(channelCreateResp),
			respInviteMembers: []byte(inviteMembersResp),
			respPostMsg:       []byte(postMsgResp),
			wantID:            "C0DEL09A5",
		},
		{
			description:       "successful channel creation including additional invites, successful message post",
			inviteMembers:     []string{"UABC123EFG"},
			initMsg:           utils.Msg{Body: "Hey!"},
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
			channel := utils.Channel{
				UserClient: client,
				BotClient:  nil,
			}

			id, err := channel.CreateChannel("general", tc.inviteMembers, tc.initMsg, false)

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

			if id != tc.wantID {
				t.Fatalf("expected channel id: %s, got: %s", tc.wantID, id)
			}
		})
	}
}
