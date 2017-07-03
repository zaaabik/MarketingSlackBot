package main


import (
	flags "github.com/jessevdk/go-flags"
	"github.com/radario/mbot/slackApi"
	"github.com/radario/mbot/db"
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
