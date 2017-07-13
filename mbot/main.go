package main

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/radario/marketingstatbot/mbot/db"
	"github.com/radario/marketingstatbot/mbot/marketingClient"
	"github.com/radario/marketingstatbot/mbot/slackApi"
)

var config struct {
	HttpAddress string `long:"http_address" env:"HTTP_ADDRESS"`
	Token       string `long:"token" env:"TOKEN"`
	Path        string `long:"db_path" env:"DB_PATH"`
	HttpToken   string `long:"http_token" env:"HTTP_TOKEN"`
}

func main() {
	flags.Parse(&config)
	var database db.Store
	database = db.NewBoltDb(config.Path)
	client := marketingClient.NewMarketingCliet(config.HttpAddress, config.HttpToken)
	bot := slackApi.NewBot(config.Token, &database, client)
	bot.Start()
}
