//bot for slack with marketing client
package slackApi

import (
	"encoding/json"
_	"fmt"
	"github.com/radario/marketingstatbot/mbot/webHookHandler"
	"github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/radario/marketingstatbot/mbot/db"
	"github.com/radario/marketingstatbot/mbot/marketingClient"
	"golang.org/x/net/context"
	"log"
	"strings"
)

type SlackBot struct {
	botToken string
	database db.Store
	client   *marketingClient.MarketingClient
}

func NewBot(botToken string, store *db.Store, client *marketingClient.MarketingClient) *SlackBot {
	return &SlackBot{botToken, *store, client}
}

func (b *SlackBot) SetToken(token string) {
	b.botToken = token
}

func (b *SlackBot) Start() {

	bot := slackbot.New(b.botToken)

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear("(?i)(.get trans count).*").MessageHandler(b.getTransactionCountHandler)
	toMe.Hear("(?i)(.get user count).*").MessageHandler(b.getUserCountHandler)
	toMe.Hear("(?i)(.show).*").MessageHandler(b.showHandler)
	toMe.Hear("(?i)(.del).*").MessageHandler(b.delDbHandler)
	bot.Run()
}

func (b *SlackBot) showHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {

	b.database.GetAll()
}

//return count of transaction of client
func (b *SlackBot) getTransactionCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {

	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m["provider"] = args[len(args) - 1]
	m["user_id"] = args[len(args) - 2]

	okAction := slack.AttachmentAction{
		Text:  "yes",
		Type:  "button",
		Name:  "submit",
		Value: "yes",
	}
	cancelAction := slack.AttachmentAction{
		Text:  "no",
		Type:  "button",
		Name:  "cancel",
		Value: "no",
	}
	str := "Do you want to get user count of " +m["user_id"]+ m["provide"] + "?"
	attach := slack.Attachment{
		Title:      str,
		Actions:    []slack.AttachmentAction{okAction, cancelAction},
		CallbackID: "get_transaction_count",
	}

	attachments := []slack.Attachment{attach}
	bot.ReplyWithAttachments(evt, attachments, slackbot.WithoutTyping)

	var res chan string


	go webHookHandler.WebHookHandler(res)


	select {
	case tmp := <-res:
		if tmp == "yes"{
			response, err := b.client.GetTransactionCount(m["user_id"], m["provider"])
			if err != nil{
				bot.Reply(evt,err.Error(),slackbot.WithoutTyping)
			}
			bot.Reply(evt,response,slackbot.WithoutTyping)

		}else {

		}
	}
	enc, err := json.Marshal(m)
	if err != nil {
		log.Println(enc)
		bot.Reply(evt, err.Error(), slackbot.WithTyping)
		return
	}
	b.database.Save(enc)
	//bot.Reply(evt, response,slackbot.WithTyping)

}

//return count of user of client
func (b *SlackBot) getUserCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {

	args := evt.Text[15:]
	params := strings.Fields(args)
	if len(params) != 2 {
		bot.Reply(evt, "wrong arguments", slackbot.WithTyping)
		bot.Reply(evt, "arg1 host_id, arg2=provider", slackbot.WithTyping)
		return
	}
	response, err := b.client.GetUserCount(params[0], params[1])
	if err != nil {
		log.Println(err)
		bot.Reply(evt, err.Error(), slackbot.WithoutTyping)
		return
	}

	m := make(map[string]string)
	m["host_id"] = params[0]
	m["provider"] = params[1]
	m["response"] = response
	enc, _ := json.Marshal(m)
	if err != nil {
		log.Println(enc)
		bot.Reply(evt, err.Error(), slackbot.WithTyping)
		return
	}
	b.database.Save(enc)
	bot.Reply(evt, response, slackbot.WithTyping)
}

//delete all data from database
func (b *SlackBot) delDbHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	b.database.DeleteAll()
	bot.Reply(evt, "db has been deleted", slackbot.WithoutTyping)
}
