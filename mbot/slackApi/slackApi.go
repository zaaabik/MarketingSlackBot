//bot for slack with marketing client
package slackApi

import (
	"github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/radario/marketingstatbot/mbot/db"
	"github.com/radario/marketingstatbot/mbot/marketingClient"
	"github.com/radario/marketingstatbot/mbot/webHookHandler"
	"golang.org/x/net/context"
	"strings"
	"log"
)

type SlackBot struct {
	server   *webHookHandler.WebHook
	botToken string
	database db.Store
	client   *marketingClient.MarketingClient
}

func NewBot(botToken string, store *db.Store, client *marketingClient.MarketingClient) *SlackBot {
	return &SlackBot{server: webHookHandler.NewWebHookHandler(), botToken: botToken, database: *store, client: client}

}

func (b *SlackBot) SetToken(token string) {
	b.botToken = token
}

func (b *SlackBot) Start() {
	bot := slackbot.New(b.botToken)
	b.server = webHookHandler.NewWebHookHandler()
	go b.server.Start()
	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear("(?i)(.get transaction count).*").MessageHandler(b.getTransactionCountHandler)
	toMe.Hear("(?i)(.get customers count).*").MessageHandler(b.getUserCountHandler)
	toMe.Hear("(?i)(.show).*").MessageHandler(b.showHandler)
	toMe.Hear("(?i)(.del).*").MessageHandler(b.delDbHandler)
	bot.Run()

}

func (b *SlackBot) showHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {

	b.database.GetAll()
}

func (b *SlackBot) getTransactionCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	log.Println(m)
	log.Println(args)
	m["provider"] = args[len(args)-1]
	m["host_id"] = args[len(args)-2]

	response, err := b.client.GetTransactionCount(m["host_id"],m["provider"])
	log.Print(response)
	if err != nil{
		return
	}


	bot.Reply(evt, response,slackbot.WithoutTyping)
}

func (b *SlackBot) getUserCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {

	args := strings.Fields(evt.Text)
	provider := args[len(args) - 1]
	hostId   := args[len(args)-2]

	response, err := b.client.GetUserCount(hostId,provider)
	if err != nil{
		return
	}
	if response == ""{
		bot.Reply(evt, "user doesnt exist" ,slackbot.WithoutTyping)
	}
	bot.Reply(evt,response,slackbot.WithoutTyping)
}

func (b *SlackBot) addLettersToUser(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent){
	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m["lettersCount"] = args[len(args) - 1]
	m["provider"] = args[len(args)-2]
	m["user_id"] = args[len(args)-3]

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
	str := "Do you want to get user count of " + m["user_id"] + m["provide"] + "?"
	attach := slack.Attachment{
		Title:      str,
		Actions:    []slack.AttachmentAction{okAction, cancelAction},
		CallbackID: "get_transaction_count",
	}
	bot.ReplyWithAttachments(evt,[]slack.Attachment{attach},slackbot.WithoutTyping)
}

func (b *SlackBot) delDbHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	b.database.DeleteAll()
	bot.Reply(evt, "db has been deleted", slackbot.WithoutTyping)
}
