# slack-utils 
Collection of utility methods I frequently use in slackbot projects/other slack scripts.

## Disclaimer
At the moment this is still a personal utils library, so until major release version `1.0.0`, it is safe to expect some significant changes to existing functions.

## Highlighted functionality
### Easy verification
Easily verify incoming requests from slash commands or interactive component callbacks using one of the provided verification methods. 

Both methods expect the application's signing secret to be embedded in the request context. Set the secret as an environment variable and add it to the context in a manner similar to the following:
```
r := chi.NewRouter()
r.Use(func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := utils.WithSigningSecret(r.Context(), env.SigningSecret)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
})
```

Read more about signing secrets [here](https://api.slack.com/docs/verifying-requests-from-slack)

**Slash command verification example**

```
cmd, err := utils.VerifySlashCmd(r)
if err != nil {
  // handle error
}

if err = doSomethingWithArg(cmd.Text); err != nil {
  // handle error
}
```
**Interactive callback verification example**

```
callback, err := utils.VerifyCallbackMsg(r)
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
```

### Working with channels
**Create a new channel, invite users, and post an init message with a single command**
```
channelHandler := &utils.Channel{
	UserClient: slack.New(env.UserToken),
	BotClient: slack.New(env.BotToken),
}
err := channelHandler.CreateChannel(channelName, userIDs, utils.Msg{Body: initMsg})
```
Requires `UserClient` with `channels:write` scope. Include the `BotClient` as well if you wish to post the init message as the bot user and not as the user associated with the `UserClient` token

**Get all channel members' Slack IDs or emails**
```
client := slack.New(env.BotToken)
emails, err := utils.GetChannelMemberEmails(client, env.ChannelID)
```
Use `GetChannelMembers` for Slack IDs instead of emails

**Leave or archive multiple channels**
```
channelHandler := &utils.Channel{
	UserClient: slack.New(env.UserToken),
}
err := channel.LeaveChannels(channelIDs)
```
User `ArchiveChannels` to archive channels instead (both require `UserClient` with `channels:write` scope)

**Invite multiple users to a channel**
```
channelHandler := &utils.Channel{
	UserClient: slack.New(env.UserToken),
}
err := channelHandler.InviteUsers(userIDs)
```
Requires `UserClient` with `channels:write` scope

### Working with users
**Convert emails to Slack IDs**
```
client := slack.New(env.BotToken)
users, err := utils.EmailsToSlackIDs(client, userEmails)
```
Use `EmailsToSlackIDsInclusive` if you want to get back *both* the email and the Slack ID for each user

### Working with files
```
client := slack.New(env.BotToken)
rows, err := utils.DownloadAndReadCSV(h.client, urlPrivateDownload)
```
Per Slack API restrictions, requires the `files:read` scope on the `UserClient` and the user associated with the token must have access to the file

### Posting messages and using Blocks
Below is a pseudo-code example of how to post an interactive block message to Slack using some of the utilities offered by the library
```
client := slack.New(env.BotToken)

startDatePickerTxt := slack.NewTextBlockObject(
	slack.MarkdownType,
	"Please choose a *start date* for the new survey",
	false,
	false,
)
startDatePickerElem := utils.NewDatePickerWithOpts(startDatePickActionID, nil, time.Now())
startDatePickerSectionBlock := slack.NewSectionBlock(startDatePickerTxt, nil, nil)
startDatePickerActionBlock := slack.NewActionBlock(
    startDatePicerkBlockID, 
    startDatePickerElem, 
    utils.CancelBtn,
)

startDatePicerkMsg = utils.Msg{
	Blocks: []slack.Block{startDatePickerSectionBlock, startDatePickerActionBlock},
}
	
_, _, err := utils.PostMsg(client, startDatePickerMsg, channelID)
```
In addition to `CancelBtn`, the library also provides a number of other pre-generated block elements including the `GoBtn` and `AckBtn`. Create a new button with a pre-set style using `NewButtonWithStyle` and create a fully-loaded datepicker with `NewDatePickerWithOpts`. You can also use `DateOptToTime` to parse the selected opt back to datetime for any callbacks from a datepicker created by this method, 

## WIP
New functionality currently in the pipeline includes:
- Get all users' Slack IDs
- Get all users' emails
- Get Slack IDs for every member of a user group
- Get emails for every member of a user group
 
Suggestions or requests are always welcome
