package main


import (
	flags "github.com/jessevdk/go-flags"
	"github.com/radario/marketingstatbot/mbot/slackApi"
	"github.com/radario/marketingstatbot/mbot/db"
)

var config struct{
	Token string `long:"token" env:"KEY"`
	Path string  `long:"db_path" env:"DB_PATH"`
}


func main()  {

	flags.Parse(&config)

	var database db.Store
	database = db.NewBoltDb(config.Path)

	bot:= slackApi.NewBot(config.Token, &database)
	bot.Start()
}
