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
