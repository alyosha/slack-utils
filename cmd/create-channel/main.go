package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/alyosha/slack-utils/helpers"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type config struct {
	UserToken string `envconfig:"USER_TOKEN" required:"true"`
	FileName  string `envconfig:"FILE_NAME" required:"true"`
}

func main() {
	os.Exit(_main())
}

func _main() int {
	log.Println("starting up")

	var env config
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("error processing environment variables: %s", err)
		return 1
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter channel name: ")
	channelName, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("failed to read string for channel name: %s", err)
		return 1
	}

	if channelName == "\n" {
		log.Print("need a valid channel name")
		return 1
	}

	if len(channelName) > utils.ChannelNameMaxLen {
		log.Print("channel name is too long")
		return 1
	}

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Enter init message, leave blank if unnecessary: ")
	initMsg, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("failed to read string for init msg: %s", err)
		return 1
	}

	client := slack.New(env.UserToken)
	channelHandler := &utils.Channel{
		Client:      client,
		ChannelName: channelName,
	}

	userHandler := &utils.User{
		Client: client,
	}

	userEmails, err := utils.UnpackSingleColCSV(env.FileName)
	if err != nil {
		log.Printf("failed to unpack single column csv: %s", err)
		return 1
	}
	userIDs := userHandler.EmailsToSlackIDs(userEmails)

	_, err = channelHandler.CreateChannel(userIDs, initMsg)
	if err != nil {
		log.Printf("failed to open channel: %s", err)
		return 1
	}

	log.Printf("finished with no issues")

	return 0
}
