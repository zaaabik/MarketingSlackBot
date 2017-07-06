package slackApi

import (
	"github.com/radario/marketingstatbot/mbot/requestTypes"
	"github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/radario/marketingstatbot/mbot/request"
	"golang.org/x/net/context"
	"github.com/radario/marketingstatbot/mbot/db"
	"log"
	"strings"
	_"fmt"
)



type SlackBot struct {
	botToken    string
	database    db.Store
	httpToken   string
	httpAddress string
}

func NewBot(botToken string,store *db.Store,httpAdress,httpToken string) (*SlackBot){
	return &SlackBot{botToken,*store,httpAdress,httpToken}
}

func (b *SlackBot)SetToken(token string){
	b.botToken = token
}

func (b *SlackBot)Start()  {

	bot := slackbot.New(b.botToken)
	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear("(?i)(.get trans count).*").MessageHandler(b.getTransactionCountHandler)
	toMe.Hear("(?i)(.get user count).*").MessageHandler(b.showHandler)
	toMe.Hear("(?i)(.show).*").MessageHandler(b.showHandler)
	toMe.Hear("(?i)(.del).*").MessageHandler(b.delDbHandler)
	bot.Run()
}

func (b *SlackBot)showHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {

	b.database.GetAll()
}
//return count of transaction of client
func (b *SlackBot)getTransactionCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {
	args := evt.Text[16:]
	pars := strings.Fields(args)
	if len(pars) != 2{
		bot.Reply(evt,"wrong arguments",slackbot.WithTyping)
		bot.Reply(evt,"arg1 host_id, arg2 provider",slackbot.WithTyping)
		return
	}
	currentRequest := request.MarketingRequest{RequestType: requestTypes.GetTransactionCount,User: evt.User,RequestParams: pars}

	err := currentRequest.Send(b.httpAddress,b.httpToken)
	if err != nil{
		log.Println(err)
		bot.Reply(evt,err.Error(),slackbot.WithTyping)
		return
	}
	enc, err := currentRequest.Encode()
	if err != nil {
		log.Println(enc)
		bot.Reply(evt, err.Error(), slackbot.WithTyping)
		return
	}
	b.database.Save(enc)
	bot.Reply(evt, currentRequest.Response,slackbot.WithTyping)
}

//return count of user of client
func (b *SlackBot)getUserCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {

	args := evt.Text[15:]
	params := strings.Fields(args)
	if len(params) != 2{
		bot.Reply(evt,"wrong arguments",slackbot.WithTyping)
		bot.Reply(evt,"arg1 host_id, arg2=provider",slackbot.WithTyping)
		return
	}
	currentRequest := request.MarketingRequest{RequestType: requestTypes.GetTransactionCount,User: evt.User,RequestParams: params}
	currentRequest.RequestType = requestTypes.GetUserCount
	err := currentRequest.Send(b.httpAddress,b.httpToken)
	if err != nil{
		log.Println(err)
		bot.Reply(evt,err.Error(),slackbot.WithoutTyping)
		return
	}

	enc, err := currentRequest.Encode()
	if err != nil{
		log.Println(enc)
		bot.Reply(evt,err.Error(),slackbot.WithTyping)
		return
	}
	b.database.Save(enc)
	bot.Reply(evt, currentRequest.Response,slackbot.WithTyping)
}
//delete all data from database
func (b *SlackBot)delDbHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent)  {
	b.database.DeleteAll()
	bot.Reply(evt,"db has been deleted",slackbot.WithoutTyping)
}
