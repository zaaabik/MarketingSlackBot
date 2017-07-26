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
	"github.com/radario/MarketingSlackBot/mbot/db"
)

type WebHook struct {
	client *marketingClient.MarketingClient
	database db.Store
}

func NewWebHookHandler(client *marketingClient.MarketingClient, database db.Store) *WebHook {
	return &WebHook{client,database}
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
			switch s.CallbackID {
			case "user/letters_count":
				{
					if s.Actions[0].Value == "no"{
						w.Write([]byte("canceled"))
						return
					}
					httpCode := web.userLettersCount(s.Actions[0].Value)
					if httpCode == http.StatusOK{
						w.Write([]byte("added"))
					}else if(httpCode == http.StatusBadRequest){
						w.Write([]byte("wrong data"))
					}else if(httpCode == http.StatusInternalServerError){
						w.Write([]byte("ooops! something went wrong "))
					}
				}

			}

		}
	})
	http.ListenAndServe(":1113", r)
}

func (web WebHook) userLettersCount(value string) (int) {
	var valueJson callbackValueJson.UserLettersCount
	json.Unmarshal([]byte(value), &valueJson)
	statusCode, err := web.client.AddLettersTohost(valueJson.HostId, valueJson.Provider, valueJson.LettersCount)
	if err != nil{
		return 0
	}
	return statusCode
}
