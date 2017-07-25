package webHookHandler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/nlopes/slack"
	"github.com/radario/MarketingSlackBot/mbot/callbackValueJson"
	"github.com/radario/MarketingSlackBot/mbot/marketingClient"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"log"
	"fmt"
)

type WebHook struct {
	client *marketingClient.MarketingClient
}

func NewWebHookHandler(client *marketingClient.MarketingClient) *WebHook {
	return &WebHook{client}
}

type j struct {
	Value []slack.Attachment `json:"attachments"`
}

func (web WebHook) Start() {
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		res, _ := ioutil.ReadAll(r.Body)
		if len(res) < 8 {

		} else {
			jsonStr, _ := url.QueryUnescape(string(res)[8:])
			var s slack.AttachmentActionCallback
			json.Unmarshal([]byte(jsonStr), &s)
			fmt.Println(s.CallbackID)
			switch s.CallbackID {
			case "user/letters_count":
				{
					log.Print("user/letters_count")
					web.userLettersCount(s.Actions[0].Value)
				}

			}

		}
	})
	http.ListenAndServe(":1113", r)
}

func (web WebHook) userLettersCount(value string) (string, error) {
	var valueJson callbackValueJson.UserLettersCount
	json.Unmarshal([]byte(value), &valueJson)
	lettersCountInt, err := strconv.Atoi(valueJson.LettersCount)
	response, err := web.client.AddLettersTohost(valueJson.HostId, valueJson.Provider, lettersCountInt)
	return response, err
}
