package slackApi

import (
	"github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/radario/mbot/request"
	"golang.org/x/net/context"
	"github.com/radario/mbot/DB"
	"log"
)
type SlackBot struct {
	botToken	string
}

func (b *SlackBot)SetToken(token string){
	b.botToken = token
}

func (b *SlackBot)Start()  {
	bot := slackbot.New(b.botToken)
	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	bot.Messages(slackbot.MessageType("c"))
	toMe.Hear("(?i)(.get).*").MessageHandler(getHandler)
	toMe.Hear("(?i)(.show).*").MessageHandler(showHandler)
	toMe.Hear("(?i)(.set).*").MessageHandler(setHandler)
	toMe.Hear("(?i)(.del).*").MessageHandler(delDbHandler)
	bot.Run()
}

func showHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	DB.ShowDb()
}

func getHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {
	request:= request.Request{RequestBody:evt.Text,User:evt.User}
	response,err := request.Send()
	if err != nil{
		log.Println(err)
		bot.Reply(evt,err.Error(),slackbot.WithTyping)
	}
	bot.Reply(evt,response,slackbot.WithTyping)
}

func setHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {
	request:= request.Request{User:evt.User,RequestBody:evt.Text}
	response,err := request.Send()
	if err != nil{
		log.Println(err)
		bot.Reply(evt,err.Error(),slackbot.WithoutTyping)
	}
	bot.Reply(evt,response,slackbot.WithTyping)
}

func delDbHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {
	DB.DeleteDb()
	bot.Reply(evt,"db has been deleted",slackbot.WithoutTyping)
}
