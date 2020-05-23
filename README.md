# slack-utils 
Collection of utility methods/middleware I frequently use in slackbot projects.

## Disclaimer
Until major release version `1.0.0`, it is safe to expect some significant changes to existing functions.

## Highlighted functionality
### Easy verification
Easily verify incoming requests from slash commands/interactive callbacks using the provided verification middleware. 

The simplest way to enable request verification is as follows:
```go
r := chi.NewRouter()
r.Use(utils.VerifySlashCommand(env.SigningSecret, nil, nil))
r.Use(utils.VerifyInteractionCallback(env.SigningSecret, nil, nil))
```

The above will verify the authenticity of all incoming requests using your
signing secret and embed the verified/unmarshalled request object into the
context on success. Optionally configure additional actions to be taken on
success/failure (e.g. logging) by passing in corresponding callback methods.

Retrieve the request from the context and use it in the following manner:

**Slash command example**
```go
func Foo(w http.ResponseWriter, r *http.Request) {
  cmd, err := utils.SlashCommand(r.Context())
  if err != nil {
    // handle error
  }

  if err = doSomething(cmd.Text); err != nil {
    // handle error
  }
}
```

**Interactive callback example**
```go
func Callback(w http.ResponseWriter, r *http.Request) {
  callback, err := utils.InteractionCallback(r.Context())
  if err != nil {
   // handle error
  }

  switch callback.Type {
  case slack.InteractionTypeBlockActions:
    // handle block action callback
  case slack.InteractionTypeMessageAction:
    // handle message action callback
  case slack.InteractionTypeDialogSubmission:
    // handle dialog submission callback
  }
}
```

### Posting messages and using Blocks
Below is a pseudo-code example of how to post an interactive block message to Slack using some of the utilities offered by the library
```go
client := slack.New(env.BotToken)

startDatePickerSectionBlock := utils.NewTextBlock("Please choose a *start date* for the new survey", nil)

startDatePickerElem := utils.NewDatePickerWithOpts(startDatePickerActionID, nil, time.Now())

startDatePickerActionBlock := slack.NewActionBlock(
    startDatePickerBlockID, 
    startDatePickerElem, 
    utils.CancelBtn,
)

startDatePickerMsg = utils.Msg{
	Blocks: []slack.Block{startDatePickerSectionBlock, startDatePickerActionBlock},
}
	
_, err := utils.PostMsg(client, startDatePickerMsg, channelID)
```

To post ephemerally, use `PostEphemeralMsg` and include the target user's ID.
 
Delete normal/ephemeral messages alike in the following manner:
```go
utils.DeleteMsg(client, channelID, ts, responseURL)
``` 

### Working with channels
**Create a new channel, invite users, and post an init message with a single command**
```go
channelHandler := &utils.Channel{
	UserClient: slack.New(env.UserToken),
	BotClient: slack.New(env.BotToken),
}
err := channelHandler.CreateChannel(channelName, userIDs, utils.Msg{Body: initMsg})
```
Requires `UserClient` with `channels:write` scope. Include the `BotClient` as well if you wish to post the init message as the bot user and not as the user associated with the `UserClient` token

**Get all channel members' Slack IDs or emails**
```go
client := slack.New(env.BotToken)
emails, err := utils.GetChannelMemberEmails(client, env.ChannelID)
```
Use `GetChannelMembers` for Slack IDs instead of emails

**Leave or archive multiple channels**
```go
channelHandler := &utils.Channel{
	UserClient: slack.New(env.UserToken),
}
err := channelHandler.LeaveChannels(channelIDs)
```
User `ArchiveChannels` to archive channels instead (both methods require `UserClient` with `channels:write` scope)

**Invite multiple users to a channel**
```go
channelHandler := &utils.Channel{
	UserClient: slack.New(env.UserToken),
}
err := channelHandler.InviteUsers(userIDs)
```
Requires `UserClient` with `channels:write` scope

### Working with users
**Convert emails to Slack IDs**
```go
client := slack.New(env.BotToken)
users, err := utils.EmailsToSlackIDs(client, userEmails)
```
Use `EmailsToSlackIDsInclusive` if you want to get back *both* the email and the Slack ID for each user

### Working with files
**Read and download CSV files shared in Slack**
```go
client := slack.New(env.BotToken)
rows, err := utils.DownloadAndReadCSV(h.client, urlPrivateDownload)
```
Per Slack API restrictions, requires the `files:read` scope on the `UserClient` and the user associated with the token must have access to the file


---
Suggestions/requests for new functionality are always welcome
