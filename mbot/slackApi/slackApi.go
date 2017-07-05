package slackApi

import (
	"github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/radario/MarketingSlackBot/mbot/request"
	"golang.org/x/net/context"
	"github.com/radario/MarketingSlackBot/mbot/db"
	"log"
)



type SlackBot struct {
	botToken	string
	database	db.Store
}

func NewBot(token string,store *db.Store) (*SlackBot){
	return &SlackBot{token,*store}
}

func (b *SlackBot)SetToken(token string){
	b.botToken = token
}

func (b *SlackBot)Start()  {
	bot := slackbot.New(b.botToken)
	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear("(?i)(.get).*").MessageHandler(b.getHandler)
	toMe.Hear("(?i)(.show).*").MessageHandler(b.showHandler)
	toMe.Hear("(?i)(.set).*").MessageHandler(b.setHandler)
	toMe.Hear("(?i)(.del).*").MessageHandler(b.delDbHandler)
	bot.Run()
}

func (b *SlackBot)showHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	b.database.GetAll()
}

func (b *SlackBot)getHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {
	request:= request.Request{RequestBody:evt.Text,User:evt.User}
	err := request.Send()
	if err != nil{
		log.Println(err)
		bot.Reply(evt,err.Error(),slackbot.WithTyping)
		return
	}
	enc, err := request.Encode()
	if err != nil {
		log.Println(enc)
		bot.Reply(evt, err.Error(), slackbot.WithTyping)
		return
	}
	b.database.Save(enc)
	bot.Reply(evt,request.Response,slackbot.WithTyping)
}

func (b *SlackBot)setHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {
	request:= request.Request{User:evt.User,RequestBody:evt.Text}
	err := request.Send()
	if err != nil{
		log.Println(err)
		bot.Reply(evt,err.Error(),slackbot.WithoutTyping)
		return
	}
	enc, err := request.Encode()
	if err != nil{
		log.Println(enc)
		bot.Reply(evt,err.Error(),slackbot.WithTyping)
		return
	}
	b.database.Save(enc)
	bot.Reply(evt,request.Response,slackbot.WithTyping)
}

func (b *SlackBot)delDbHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {
	b.database.DeleteAll()
	bot.Reply(evt,"db has been deleted",slackbot.WithoutTyping)
}
