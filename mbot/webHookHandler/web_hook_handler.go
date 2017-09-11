package webHookHandler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/nlopes/slack"
	"github.com/radario/MarketingSlackBot/mbot/db"
	"github.com/radario/MarketingSlackBot/mbot/entities"
	"github.com/radario/MarketingSlackBot/mbot/marketingClient"
	"github.com/radario/MarketingSlackBot/mbot/textConstants"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type WebHook struct {
	client   *marketingClient.MarketingClient
	database db.Store
}

const answerToUserTemplate = "<@%s> %s"

func NewWebHookHandler(client *marketingClient.MarketingClient, database db.Store) *WebHook {
	return &WebHook{client, database}
}

type j struct {
	Value []slack.Attachment `json:"attachments"`
}

func (web WebHook) Start() {
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		res, _ := ioutil.ReadAll(r.Body)

		if len(res) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jsonStr, _ := url.QueryUnescape(string(res)[8:])
		var attachmentsCallback slack.AttachmentActionCallback
		json.Unmarshal([]byte(jsonStr), &attachmentsCallback)
		switch attachmentsCallback.CallbackID {
		case textConstants.AddUserLetterCountMethod:
			{
				addLettersMethodAnswer(&w, &attachmentsCallback, &web)
			}
		case textConstants.UpdateSendgridEmailMethod:
			{
				updateSendgridAnswer(&w, &attachmentsCallback, &web)
			}
		}
	})
	http.ListenAndServe(":1113", r)
}

func (web WebHook) userLettersCount(value string) int {
	var valueJson entities.UserLettersCount
	json.Unmarshal([]byte(value), &valueJson)
	statusCode, err := web.client.AddLettersToHost(valueJson.HostId, valueJson.Provider, valueJson.LettersCount)
	if err != nil {
		return 0
	}
	if statusCode == http.StatusOK {
		m := make(map[string]string)
		m["method"] = textConstants.AddUserLetterCountMethod
		m["user"] = valueJson.UserId
		m[textConstants.ProviderKey] = valueJson.Provider
		m[textConstants.HostIdKey] = valueJson.HostId
		m[textConstants.LettersCountKey] = valueJson.LettersCount
		err = web.database.Save(m)
		if err != nil {
			log.Print(err)
		}

	}
	return statusCode
}

func (web WebHook) updateSendgridEmail(value string) int {
	var valueJson entities.UserSendGrid
	json.Unmarshal([]byte(value), &valueJson)
	statusCode, err := web.client.UpdateSendgridEmail(valueJson.HostId, valueJson.Provider, valueJson.Email)
	if err != nil {
		return 0
	}
	if statusCode == http.StatusOK {
		m := make(map[string]string)
		m["method"] = textConstants.UpdateSendgridEmailMethod
		m[textConstants.ProviderKey] = valueJson.Provider
		m[textConstants.HostIdKey] = valueJson.HostId
		m[textConstants.EmailKey] = valueJson.Email
		m["user"] = valueJson.UserId
		err = web.database.Save(m)
		if err != nil {
			log.Print(err)
		}
	}
	return statusCode
}

func addLettersMethodAnswer(w *http.ResponseWriter, callback *slack.AttachmentActionCallback, web *WebHook) {
	user := callback.User.ID
	if callback.Actions[0].Value == "no" {
		response := fmt.Sprintf(answerToUserTemplate, user, textConstants.CanceledEventText)
		(*w).Write([]byte(response))
		return
	}
	httpCode := web.userLettersCount(callback.Actions[0].Value)

	switch httpCode {
	case http.StatusOK:
		{
			response := fmt.Sprintf(answerToUserTemplate, user, textConstants.ApproveEventText)
			(*w).Write([]byte(response))
		}
	case http.StatusNotFound:
		response := fmt.Sprintf(answerToUserTemplate, user, textConstants.UserDoesNotExistText)
		(*w).Write([]byte(response))
	case http.StatusInternalServerError:
		response := fmt.Sprintf(answerToUserTemplate, user, textConstants.ServerErrorText)
		(*w).Write([]byte(response))
	default:
		response := fmt.Sprintf(answerToUserTemplate, user, textConstants.RequestErrorText)
		(*w).Write([]byte(response))
	}
}

func updateSendgridAnswer(w *http.ResponseWriter, callback *slack.AttachmentActionCallback, web *WebHook) {
	user := callback.User.ID
	if callback.Actions[0].Value == "no" {
		response := fmt.Sprintf(answerToUserTemplate, user, textConstants.CanceledEventText)
		(*w).Write([]byte(response))
		return
	}
	httpCode := web.updateSendgridEmail(callback.Actions[0].Value)

	switch httpCode {
	case http.StatusCreated:
		{
			response := fmt.Sprintf(answerToUserTemplate, user, textConstants.ApproveEventText)
			(*w).Write([]byte(response))
		}
	case http.StatusNotFound:
		{
			response := fmt.Sprintf(answerToUserTemplate, user, textConstants.UserDoesNotExistText)
			(*w).Write([]byte(response))
		}
	case http.StatusInternalServerError:
		{
			response := fmt.Sprintf(answerToUserTemplate, user, textConstants.ServerErrorText)
			(*w).Write([]byte(response))
		}
	default:
		{
			response := fmt.Sprintf(answerToUserTemplate, user, textConstants.RequestErrorText)
			(*w).Write([]byte(response))

		}
	}
}
