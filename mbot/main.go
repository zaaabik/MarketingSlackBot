package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/radario/MarketingSlackBot/mbot/db"
	"github.com/radario/MarketingSlackBot/mbot/marketingClient"
	"github.com/radario/MarketingSlackBot/mbot/slackApi"
	"log"
	"os"
)

var config struct {
	BaseApiUrl     string `long:"base_api_url" env:"BASE_URL"`
	BotUserToken   string `long:"bot_token" env:"BOT_TOKEN"`
	Host           string `long:"db_host" env:"DB_HOST"`
	HttpTokenValue string `long:"http_token_value" env:"HTTP_TOKEN_VALUE"`
	HttpTokenKey   string `long:"http_token_key" env:"HTTP_TOKEN_KEY"`
}

const (
	CantCreateDatabaseExitCode = 1
	WrongFlagsExitCode         = 2
)

func main() {
	flags.Parse(&config)
	log.Print("started!")
	if config.HttpTokenKey == "" || config.BotUserToken == "" ||
		config.BaseApiUrl == "" {
		os.Exit(WrongFlagsExitCode)
		log.Print("exit with 1 status!")
	}
	var database db.Store
	database, err := db.NewMongoDb(config.Host)
	if err != nil {
		log.Print("dbproblems")
		os.Exit(CantCreateDatabaseExitCode)
	}
	defer database.Close()
	client := marketingClient.NewMarketingClient(config.BaseApiUrl, config.HttpTokenValue, config.HttpTokenKey)
	bot := slackApi.NewBot(config.BotUserToken, &database, client)
	bot.Start()
}
