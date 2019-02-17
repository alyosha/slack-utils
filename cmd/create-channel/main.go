package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/alyosha93/slack-scripts/cmd/utils"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type config struct {
	UserToken   string `envconfig:"USER_TOKEN" required:"true"`
	ChannelName string `envconfig:"CHANNEL_NAME" required:"true"`
	FileName    string `envconfig:"FILE_NAME" required:"true"`
}

func main() {
	os.Exit(_main())
}

func _main() int {
	log.Print("starting up")
	var env config
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("error processing environment variables: %s", err)
		return 1
	}

	if len(env.ChannelName) > utils.ChannelNameMaxLen {
		log.Print("channel name is too long")
		return 1
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter init message, leave blank if unnecessary: ")
	text, _ := reader.ReadString('\n')

	client := slack.New(env.UserToken)
	channelHandler := &utils.Channel{
		Client:      client,
		ChannelName: env.ChannelName,
	}

	userHandler := &utils.User{
		Client: client,
	}

	userEmails := utils.UnpackSingleColCSV(env.FileName)
	userIDs := userHandler.EmailsToSlackIDs(userEmails)

	_, err := channelHandler.CreateChannel(userIDs, text)
	if err != nil {
		log.Printf("received the following error: %s", err)
		return 1
	}

	log.Printf("finished with no issues")

	return 0
}
