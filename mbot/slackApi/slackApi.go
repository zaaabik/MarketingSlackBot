//bot for slack with marketing client
package slackApi

import (
	"encoding/json"
	"github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/radario/MarketingSlackBot/mbot/callbackValueJson"
	"github.com/radario/MarketingSlackBot/mbot/db"
	"github.com/radario/MarketingSlackBot/mbot/marketingClient"
	"github.com/radario/MarketingSlackBot/mbot/webHookHandler"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type SlackBot struct {
	server   *webHookHandler.WebHook
	botToken string
	database db.Store
	client   *marketingClient.MarketingClient
}

func NewBot(botToken string, store *db.Store, client *marketingClient.MarketingClient) *SlackBot {
	return &SlackBot{server: webHookHandler.NewWebHookHandler(client, *store), botToken: botToken, database: *store, client: client}

}

func (b *SlackBot) SetToken(token string) {
	b.botToken = token
}

func (b *SlackBot) Start() {
	bot := slackbot.New(b.botToken)
	go b.server.Start()
	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	toMe.Hear(`((<@\d+>\s*)+|(^\s*))(\.add (-)?\d* letters \w+ \w+\s*$)`).MessageHandler(b.addLettersToUser)
	toMe.Hear(`((<@\d+>\s*)+|(^\s*))(\.get transaction count \w+ \w+\s*$)`).MessageHandler(b.getTransactionCountHandler)
	toMe.Hear(`((<@\d+>\s*)+|(^\s*))(\.get customers count \w+ \w+\s*$)`).MessageHandler(b.getUserCountHandler)
	toMe.Hear(`\.show`).MessageHandler(b.showHandler)
	toMe.Hear(`\.del`).MessageHandler(b.delDbHandler)
	toMe.Hear(`\.help`).MessageHandler(b.showHelp)
	toMe.Hear(".*").MessageHandler(b.unknownCommand)
	bot.Run()
}

func (b *SlackBot) showHelp(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	addLettersHelp := "- add [letters Count] letters [host_id] [provider]"
	getTransCoutHelp := "- get transaction count [host_id] [provider]"
	getCustomersCountHelp := " -get customers count [host_id] [provider]"
	res := "command list\n" + addLettersHelp + "\n" + getTransCoutHelp + "\n" + getCustomersCountHelp
	bot.Reply(evt, res, slackbot.WithoutTyping)
}

func (b *SlackBot) showHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	b.database.GetAll()
}

func (b *SlackBot) unknownCommand(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	//checking is message written by bot
	if evt.User != "" {
		bot.Reply(evt, "unknown command\nwrite .help ", slackbot.WithoutTyping)
	}
}

func (b *SlackBot) getTransactionCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m["provider"] = args[len(args)-1]
	m["host_id"] = args[len(args)-2]
	response, err, httpCode := b.client.GetTransactionCount(m["host_id"], m["provider"])

	if err != nil {
		bot.Reply(evt, "<@"+evt.User+"> "+"program error", slackbot.WithoutTyping)
		return
	}

	if httpCode == http.StatusInternalServerError {
		bot.Reply(evt, "<@"+evt.User+"> "+"ooops! somthing went wrong. Try again", slackbot.WithoutTyping)
		return
	}
	if httpCode == http.StatusBadRequest {
		bot.Reply(evt, "<@"+evt.User+"> "+"wrong data", slackbot.WithoutTyping)
		return
	}

	m["response"] = response
	m["user"] = evt.User
	b.database.Save(m)
	bot.Reply(evt, "<@"+evt.User+"> "+response, slackbot.WithoutTyping)
}

func (b *SlackBot) getUserCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {

	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m["provider"] = args[len(args)-1]
	m["host_id"] = args[len(args)-2]
	response, err, httpCode := b.client.GetUserCount(m["host_id"], m["provider"])

	if err != nil {
		bot.Reply(evt, "<@"+evt.User+"> "+"ooops! program error", slackbot.WithoutTyping)
		return
	}

	if httpCode == http.StatusInternalServerError {
		bot.Reply(evt, "<@"+evt.User+"> "+"ooops! somthing went wrong. Try again", slackbot.WithoutTyping)
		return
	}
	if httpCode == http.StatusBadRequest {
		bot.Reply(evt, "<@"+evt.User+"> "+"wrong data", slackbot.WithoutTyping)
		return
	}

	m["response"] = response
	m["user"] = evt.User
	b.database.Save(m)
	bot.Reply(evt, "<@"+evt.User+"> "+response, slackbot.WithoutTyping)
}

func (b *SlackBot) addLettersToUser(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m["lettersCount"] = args[len(args)-4]
	m["provider"] = args[len(args)-1]
	m["host_id"] = args[len(args)-2]

	value := callbackValueJson.UserLettersCount{m["host_id"], m["provider"], m["lettersCount"]}
	jsonValue, err := json.Marshal(value)
	_ = err
	okAction := slack.AttachmentAction{
		Text:  "yes",
		Type:  "button",
		Name:  "submit",
		Value: string(jsonValue),
	}
	cancelAction := slack.AttachmentAction{
		Text:  "no",
		Type:  "button",
		Name:  "cancel",
		Value: "no",
	}
	str := "Do you want to get user count of " + m["user_id"] + m["provider"] + "?"
	attach := slack.Attachment{
		Title:      str,
		Actions:    []slack.AttachmentAction{okAction, cancelAction},
		CallbackID: "user/letters_count",
	}
	bot.ReplyWithAttachments(evt, []slack.Attachment{attach}, slackbot.WithoutTyping)
}

func (b *SlackBot) delDbHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	b.database.DeleteAll()
	bot.Reply(evt, "db has been deleted", slackbot.WithoutTyping)
}
