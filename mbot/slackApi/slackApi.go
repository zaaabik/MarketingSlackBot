//bot for slack with marketing client
package slackApi

import (
	"encoding/json"
	"github.com/adampointer/go-slackbot"
	"github.com/nlopes/slack"
	"github.com/radario/MarketingSlackBot/mbot/db"
	"github.com/radario/MarketingSlackBot/mbot/entities"
	"github.com/radario/MarketingSlackBot/mbot/marketingClient"
	"github.com/radario/MarketingSlackBot/mbot/messagesRegExp"
	"github.com/radario/MarketingSlackBot/mbot/textConstants"
	"github.com/radario/MarketingSlackBot/mbot/webHookHandler"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"regexp"
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
	toMe.Hear(messagesRegExp.AddLettersToUserRegExp).MessageHandler(b.addLettersToUser)
	toMe.Hear(messagesRegExp.GetTransactionCountRegExp).MessageHandler(b.getTransactionCountHandler)
	toMe.Hear(messagesRegExp.GetCustomersCountRegExp).MessageHandler(b.getCustomersCountHandler)
	toMe.Hear(messagesRegExp.UpdateSendgridEmailRegExp).MessageHandler(b.updateSendgridEmail)
	toMe.Hear(messagesRegExp.CreateScenarioByCompainRegExp).MessageHandler(b.createScenarioByCampaign)
	toMe.Hear(messagesRegExp.ShowDbRegExp).MessageHandler(b.showHandler)
	toMe.Hear(messagesRegExp.DeleteDbRegExp).MessageHandler(b.delDbHandler)
	toMe.Hear(messagesRegExp.HelpRegExp).MessageHandler(b.showHelp)
	toMe.Hear(messagesRegExp.AllRegExp).MessageHandler(b.unknownCommand)
	bot.Run()
}

func (b *SlackBot) showHelp(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {

	res := textConstants.Help
	bot.Reply(evt, res, slackbot.WithoutTyping)
}

func (b *SlackBot) showHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	b.database.GetAll()
}

func (b *SlackBot) unknownCommand(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	log.Println(evt.Msg)
	//checking is message written by bot
	if evt.User != "" {
		bot.Reply(evt, textConstants.UnknownCommand, slackbot.WithoutTyping)
	}
}

func (b *SlackBot) getTransactionCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m["method"] = textConstants.GetCustomersTransactionMethod
	m[textConstants.ProviderKey] = args[len(args)-1]
	m[textConstants.HostIdKey] = args[len(args)-2]
	response, err, httpCode := b.client.GetTransactionCount(m[textConstants.HostIdKey], m[textConstants.ProviderKey])

	if err != nil {
		bot.Reply(evt, "<@"+evt.User+"> "+textConstants.RequestErrorText, slackbot.WithoutTyping)
		return
	}

	switch httpCode {
	case http.StatusInternalServerError:
		{
			bot.Reply(evt, "<@"+evt.User+"> "+textConstants.ServerErrorText, slackbot.WithoutTyping)
			return
		}
	case http.StatusNotFound:
		{
			bot.Reply(evt, "<@"+evt.User+"> "+textConstants.UserDoesNotExistText, slackbot.WithoutTyping)
			return
		}
	}

	m["response"] = response
	m["user"] = evt.User
	b.database.Save(m)
	bot.Reply(evt, "<@"+evt.User+"> "+response, slackbot.WithoutTyping)
}

func (b *SlackBot) getCustomersCountHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m["method"] = textConstants.GetCustomersCountMethod
	m[textConstants.ProviderKey] = args[len(args)-1]
	m[textConstants.HostIdKey] = args[len(args)-2]
	response, err, httpCode := b.client.GetUserCount(m[textConstants.HostIdKey], m[textConstants.ProviderKey])

	if err != nil {
		bot.Reply(evt, "<@"+evt.User+"> "+textConstants.RequestErrorText, slackbot.WithoutTyping)
		return
	}

	switch httpCode {
	case http.StatusInternalServerError:
		{
			bot.Reply(evt, "<@"+evt.User+"> "+textConstants.RequestErrorText, slackbot.WithoutTyping)
			return
		}
	case http.StatusNotFound:
		{
			bot.Reply(evt, "<@"+evt.User+"> "+textConstants.UserDoesNotExistText, slackbot.WithoutTyping)
			return
		}
	}

	m["response"] = response
	m["user"] = evt.User
	b.database.Save(m)
	bot.Reply(evt, "<@"+evt.User+"> "+response, slackbot.WithoutTyping)
}

func (b *SlackBot) addLettersToUser(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m[textConstants.LettersCountKey] = args[len(args)-4]
	m[textConstants.ProviderKey] = args[len(args)-1]
	m[textConstants.HostIdKey] = args[len(args)-2]
	value := entities.UserLettersCount{m[textConstants.HostIdKey], m[textConstants.ProviderKey], m[textConstants.LettersCountKey]}
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
	str := "Do you want to add " + m[textConstants.LettersCountKey] + " letters to " + m[textConstants.HostIdKey] + " " + m[textConstants.ProviderKey] + "?"
	attach := slack.Attachment{
		Title:      str,
		Actions:    []slack.AttachmentAction{okAction, cancelAction},
		CallbackID: textConstants.AddUserLetterCountMethod,
	}
	bot.ReplyWithAttachments(evt, []slack.Attachment{attach}, slackbot.WithoutTyping)
}

func (b *SlackBot) updateSendgridEmail(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	args := strings.Fields(evt.Text)
	m := make(map[string]string)
	m[textConstants.ProviderKey] = args[len(args)-1]
	m[textConstants.HostIdKey] = args[len(args)-2]
	re := regexp.MustCompile(`(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`)
	email := re.FindStringSubmatch(evt.Msg.Text)
	if len(email) > 0 {
		m[textConstants.EmailKey] = email[0]
	} else {
		return
	}

	value := entities.UserSendGrid{m[textConstants.HostIdKey], m[textConstants.ProviderKey], m[textConstants.EmailKey]}
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
	str := "Do you want to set email:" + m[textConstants.EmailKey] + " to " + m[textConstants.HostIdKey] + " " + m[textConstants.ProviderKey] + "?"
	attach := slack.Attachment{
		Title:      str,
		Actions:    []slack.AttachmentAction{okAction, cancelAction},
		CallbackID: textConstants.UpdateSendgridEmailMethod,
	}
	bot.ReplyWithAttachments(evt, []slack.Attachment{attach}, slackbot.WithoutTyping)
}

func (b *SlackBot) createScenarioByCampaign(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	args := strings.Fields(evt.Text)
	id := strings.Split(args[len(args)-1], "/")
	campaignId :=(id[len(id) - 1])
	result := strings.Replace(campaignId,">","",-1)
	log.Print(result)
	m := make(map[string]string)
	m[textConstants.ScenarioName] = args[len(args)-2]
	m[textConstants.CampaignId] = result

	go func() {
		httpCode, err := b.client.CreateScenarioByCampaign(m[textConstants.CampaignId], m[textConstants.ScenarioName])
		if err != nil {
			log.Print(err)
			return
		}
		switch httpCode {
		case http.StatusCreated:
			{
				bot.Reply(evt, "Created!", slackbot.WithoutTyping)
			}
		case http.StatusNotFound:
			{
				bot.Reply(evt, "Campaign doesn`t exist", slackbot.WithoutTyping)
			}
		default:
			bot.Reply(evt, "fail", slackbot.WithoutTyping)
		}
	}()
}

func (b *SlackBot) delDbHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	b.database.DeleteAll()
	bot.Reply(evt, "db has been deleted", slackbot.WithoutTyping)
}
