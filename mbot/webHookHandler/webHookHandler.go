package webHookHandler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/nlopes/slack"
	"github.com/radario/MarketingSlackBot/mbot/db"
	"github.com/radario/MarketingSlackBot/mbot/entities"
	"github.com/radario/MarketingSlackBot/mbot/errorsText"
	"github.com/radario/MarketingSlackBot/mbot/marketingClient"
	"io/ioutil"
	"net/http"
	"net/url"
)

type WebHook struct {
	client   *marketingClient.MarketingClient
	database db.Store
}

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
			//wrong webhook from slakc
		} else {
			jsonStr, _ := url.QueryUnescape(string(res)[8:])
			var s slack.AttachmentActionCallback
			json.Unmarshal([]byte(jsonStr), &s)
			switch s.CallbackID {
			case "user/letters_count":
				{
					if s.Actions[0].Value == "no" {
						w.Write([]byte("canceled"))
						return
					}
					httpCode := web.userLettersCount(s.Actions[0].Value)
					switch httpCode {
					case http.StatusOK:
						w.Write([]byte("added"))
					case http.StatusNotFound:
						w.Write([]byte(errorsText.UserDoesNotExistText))
					case http.StatusInternalServerError:
						w.Write([]byte(errorsText.ServerErrorText))
					default:
						w.Write([]byte(errorsText.RequestErrorText))
					}
				}

			}

		}
	})
	http.ListenAndServe(":1113", r)
}

func (web WebHook) userLettersCount(value string) int {
	var valueJson entities.UserLettersCount
	json.Unmarshal([]byte(value), &valueJson)
	statusCode, err := web.client.AddLettersTohost(valueJson.HostId, valueJson.Provider, valueJson.LettersCount)
	if err != nil {
		return 0
	}
	return statusCode
}
