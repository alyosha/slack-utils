package utils

const channelCreateResp = `{
    "ok": true,
    "channel": {
        "id": "C0DEL09A5",
        "name": "endeavor",
        "is_channel": true,
        "created": 1502833204,
        "creator": "U061F7AUR",
        "is_archived": false,
        "is_general": false,
        "name_normalized": "endeavor",
        "is_shared": false,
        "is_org_shared": false,
        "is_member": true,
        "is_private": false,
        "is_mpim": false,
        "last_read": "0000000000.000000",
        "latest": null,
        "unread_count": 0,
        "unread_count_display": 0,
        "members": [
            "U061F7AUR"
        ],
        "topic": {
            "value": "",
            "creator": "",
            "last_set": 0
        },
        "purpose": {
            "value": "",
            "creator": "",
            "last_set": 0
        },
        "previous_names": []
    }
}`

const inviteMembersResp = `{
    "ok": true,
    "channel": {
        "id": "C1H9RESGL",
        "name": "busting",
        "is_channel": true,
        "created": 1466025154,
        "creator": "U0G9QF9C6",
        "is_archived": false,
        "is_general": false,
        "name_normalized": "busting",
        "is_shared": false,
        "is_org_shared": false,
        "is_member": true,
        "is_private": false,
        "is_mpim": false,
        "last_read": "1503435963.000307",
        "latest": {
            "user": "U1QNSQB9U",
            "text": "<@U1QNSQB9U|protobot> has left the channel",
            "type": "message",
            "subtype": "channel_leave",
            "ts": "1503435963.000307"
        },
        "unread_count": 0,
        "unread_count_display": 0,
        "members": [
            "U0G9QF9C6",
            "U1QNSQB9U"
        ],
        "topic": {
            "value": "My Topic",
            "creator": "U0G9QF9C6",
            "last_set": 1503435128
        },
        "purpose": {
            "value": "My Purpose",
            "creator": "U0G9QF9C6",
            "last_set": 1503435128
        },
        "previous_names": []
    }
}`

const postMsgResp = `{
    "ok": true,
    "channel": "C1H9RESGL",
    "ts": "1503435956.000247",
    "message": {
        "text": "Here's a message for you",
        "username": "ecto1",
        "bot_id": "B19LU7CSY",
        "attachments": [
            {
                "text": "This is an attachment",
                "id": 1,
                "fallback": "This is an attachment's fallback"
            }
        ],
        "type": "message",
        "subtype": "bot_message",
        "ts": "1503435956.000247"
    }
}`

const channelInfoResp = `{
    "ok": true,
    "channel": {
        "id": "C1H9RESGL",
        "name": "busting",
        "is_channel": true,
        "created": 1466025154,
        "creator": "U0G9QF9C6",
        "is_archived": false,
        "is_general": false,
        "name_normalized": "busting",
        "is_shared": false,
        "is_org_shared": false,
        "is_member": true,
        "is_private": false,
        "is_mpim": false,
        "last_read": "1503435939.000101",
        "latest": {
            "text": "Containment unit is 98% full",
            "username": "ecto1138",
            "bot_id": "B19LU7CSY",
            "attachments": [
                {
                    "text": "Don't get too attached",
                    "id": 1,
                    "fallback": "This is an attachment fallback"
                }
            ],
            "type": "message",
            "subtype": "bot_message",
            "ts": "1503435956.000247"
        },
        "unread_count": 1,
        "unread_count_display": 1,
        "members": [
            "U0G9QF9C6",
            "U1QNSQB9U"
        ],
        "topic": {
            "value": "Spiritual containment strategies",
            "creator": "U0G9QF9C6",
            "last_set": 1503435128
        },
        "purpose": {
            "value": "Discuss busting ghosts",
            "creator": "U0G9QF9C6",
            "last_set": 1503435128
        },
        "previous_names": [
            "dusting"
        ]
    }
}`

const usersListResp = `{
    "ok": true,
    "members": [
        {
            "id": "U0G9QF9C6",
            "team_id": "T012AB3C4",
            "name": "spengler",
            "deleted": false,
            "color": "9f69e7",
            "real_name": "spengler",
            "tz": "America/Los_Angeles",
            "tz_label": "Pacific Daylight Time",
            "tz_offset": -25200,
            "profile": {
                "avatar_hash": "ge3b51ca72de",
                "status_text": "Print is dead",
                "status_emoji": ":books:",
                "real_name": "Egon Spengler",
                "display_name": "spengler",
                "real_name_normalized": "Egon Spengler",
                "display_name_normalized": "spengler",
                "email": "spengler@ghostbusters.example.com",
                "image_24": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
                "image_32": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
                "image_48": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
                "image_72": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
                "image_192": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
                "image_512": "https://.../avatar/e3b51ca72dee4ef87916ae2b9240df50.jpg",
                "team": "T012AB3C4"
            },
            "is_admin": true,
            "is_owner": false,
            "is_primary_owner": false,
            "is_restricted": false,
            "is_ultra_restricted": false,
            "is_bot": false,
            "updated": 1502138686,
            "is_app_user": false,
            "has_2fa": false
        },
        {
            "id": "W07QCRPA4",
            "team_id": "T0G9PQBBK",
            "name": "glinda",
            "deleted": false,
            "color": "9f69e7",
            "real_name": "Glinda Southgood",
            "tz": "America/Los_Angeles",
            "tz_label": "Pacific Daylight Time",
            "tz_offset": -25200,
            "profile": {
                "avatar_hash": "8fbdd10b41c6",
                "image_24": "https://a.slack-edge.com...png",
                "image_32": "https://a.slack-edge.com...png",
                "image_48": "https://a.slack-edge.com...png",
                "image_72": "https://a.slack-edge.com...png",
                "image_192": "https://a.slack-edge.com...png",
                "image_512": "https://a.slack-edge.com...png",
                "image_1024": "https://a.slack-edge.com...png",
                "image_original": "https://a.slack-edge.com...png",
                "first_name": "Glinda",
                "last_name": "Southgood",
                "title": "Glinda the Good",
                "phone": "",
                "skype": "",
                "real_name": "Glinda Southgood",
                "real_name_normalized": "Glinda Southgood",
                "display_name": "Glinda the Fairly Good",
                "display_name_normalized": "Glinda the Fairly Good",
                "email": "glenda@south.oz.coven"
            },
            "is_admin": true,
            "is_owner": false,
            "is_primary_owner": false,
            "is_restricted": false,
            "is_ultra_restricted": false,
            "is_bot": false,
            "updated": 1480527098,
            "has_2fa": false
        }
    ],
    "cache_ts": 1498777272,
		"response_metadata": {
			"next_cursor": ""
    }
}`

const channelsLeaveResp = `{
    "ok": true
}`

const channelsArchiveResp = `{
    "ok": true
}`

const channelCreateErrResp = `{
    "ok": false,
    "error": "invalid_name_specials",
    "detail": "Value passed for 'name' contained unallowed special characters."
}`

const inviteMembersErrResp = `{
    "ok": false,
    "error": "cant_invite"
}`

const cantInviteSelfErrResp = `{
    "ok": false,
    "error": "cant_invite_self"
}`

const postMsgErrResp = `{
    "ok": false,
    "error": "too_many_attachments"
}`

const channelInfoErrResp = `{
    "ok": false,
    "error": "channel_not_found"
}`

const usersListErrResp = `{
    "ok": false,
    "error": "invalid_cursor"
}`

const channelsLeaveErrResp = `{
    "ok": false,
    "error": "invalid_auth"
}`

const channelsArchiveErrResp = `{
    "ok": false,
    "error": "invalid_auth"
}`

const channelAlreadyArchivedErrResp = `{
    "ok": false,
    "error": "already_archived"
}`
