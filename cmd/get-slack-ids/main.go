package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/alyosha93/slack-utils/helpers"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type config struct {
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`
	FileName string `envconfig:"FILE_NAME" required:"true"`
}

func main() {
	os.Exit(_main())
}

func _main() int {
	log.Println("starting up")

	var env config
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("failed to process environment variables: %s", err)
		return 1
	}

	file, err := os.Create("users.csv")
	if err != nil {
		log.Printf("failed to create new CSV file")
	}
	defer file.Close()

	client := slack.New(env.BotToken)
	userHandler := &utils.User{
		Client: client,
	}

	userEmails, err := utils.UnpackSingleColCSV(env.FileName)
	if err != nil {
		log.Printf("failed to unpack single column csv: %s", err)
		return 1
	}
	users := userHandler.EmailsToSlackIDsInclusive(userEmails)

	w := csv.NewWriter(file)
	w.WriteAll(users)

	if err := w.Error(); err != nil {
		log.Fatalf("failed to write csv:", err)
	}

	log.Print("finished with no issues")
	return 0
}
