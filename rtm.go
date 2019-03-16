package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

// IsBotImperative checks the content of the message event to see if it is
// a command directed at the bot user
func (l *Listener) IsBotImperative(event *slack.MessageEvent) bool {
	if l.BotID == "" {
		log.Printf("received the following message: %s", event.Msg.Text)
		return false
	}

	msg := splitMsg(event.Msg.Text)
	if len(msg) == 0 || msg[0] != fmt.Sprintf("<@%s>", l.BotID) {
		return false
	}

	return true
}

func splitMsg(msg string) []string {
	return strings.Split(strings.TrimSpace(msg), " ")[0:]
}
