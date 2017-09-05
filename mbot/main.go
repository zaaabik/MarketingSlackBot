package main

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/radario/MarketingSlackBot/mbot/db"
	"github.com/radario/MarketingSlackBot/mbot/marketingClient"
	"github.com/radario/MarketingSlackBot/mbot/slackApi"
	"os"
)

var config struct {
	BaseApiUrl     string `long:"base_url" env:"BASE_URL"`
	BotUserToken   string `long:"bot_token" env:"BOT_TOKEN"`
	DatabasePath   string `long:"db_path" env:"DB_PATH"`
	HttpTokenValue string `long:"http_token_value" env:"HTTP_TOKEN_VALUE"`
	HttpTokenKey   string `long:"http_token_key" env:"HTTP_TOKEN_KEY"`
}

const (
	CantCreateDatabaseExitCode = 1
	WrongFlagsExitCode         = 2
)

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		os.Exit(CantCreateDatabaseExitCode)
	}
	var database db.Store
	database, err = db.NewBoltDb(config.DatabasePath)
	defer database.Close()
	if err != nil {
		os.Exit(WrongFlagsExitCode)
	}
	client := marketingClient.NewMarketingClient(config.BaseApiUrl, config.HttpTokenValue, config.HttpTokenKey)
	bot := slackApi.NewBot(config.BotUserToken, &database, client)
	bot.Start()
}
