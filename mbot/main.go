package main


import (
	"github.com/radario/mbot/slackApi"
)

func main()  {
	bot := slackApi.SlackBot{}
	bot.SetToken("xoxb-205632384709-FItggQTj2YZx542CmWZyqDFr")
	bot.Start()
}
