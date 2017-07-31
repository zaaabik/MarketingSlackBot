package webHookHandler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/nlopes/slack"
	"github.com/radario/MarketingSlackBot/mbot/db"
	"github.com/radario/MarketingSlackBot/mbot/entities"
	"github.com/radario/MarketingSlackBot/mbot/marketingClient"
	"github.com/radario/MarketingSlackBot/mbot/textConstants"
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
			case textConstants.AddUserLetterCountMethod:
				{
					user := s.User.ID
					if s.Actions[0].Value == "no" {
						response := "<@" + user + "> " + textConstants.CanceledEventText
						w.Write([]byte(response))
						return
					}
					httpCode := web.userLettersCount(s.Actions[0].Value)

					switch httpCode {
					case http.StatusOK:
						{
							response := "<@" + user + "> " + textConstants.ApproveEventText
							w.Write([]byte(response))
						}
					case http.StatusNotFound:
						response := "<@" + user + "> " + textConstants.UserDoesNotExistText
						w.Write([]byte(response))
					case http.StatusInternalServerError:
						response := "<@" + user + "> " + textConstants.ServerErrorText
						w.Write([]byte(response))
					default:
						response := "<@" + user + "> " + textConstants.RequestErrorText
						w.Write([]byte(response))
					}
				}
			case textConstants.UpdateSendgridEmail:
				{
					user := s.User.ID
					if s.Actions[0].Value == "no" {
						response := "<@" + user + "> " + textConstants.CanceledEventText
						w.Write([]byte(response))
						return
					}
					httpCode := web.userLettersCount(s.Actions[0].Value)

					switch httpCode {
					case http.StatusOK:
						{
							response := "<@" + user + "> " + textConstants.EmailChanged
							w.Write([]byte(response))
						}
					case http.StatusNotFound:
						response := "<@" + user + "> " + textConstants.UserDoesNotExistText
						w.Write([]byte(response))
					case http.StatusInternalServerError:
						response := "<@" + user + "> " + textConstants.ServerErrorText
						w.Write([]byte(response))
					default:
						response := "<@" + user + "> " + textConstants.RequestErrorText
						w.Write([]byte(response))
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
	if statusCode == http.StatusOK {
		m := make(map[string]string)
		m["method"] = textConstants.AddUserLetterCountMethod
		m[textConstants.ProviderKey] = valueJson.Provider
		m[textConstants.HostIdKey] = valueJson.HostId
		m[textConstants.LettersCountKey] = valueJson.LettersCount
		web.database.Save(m)
	}
	return statusCode
}

func (web WebHook) userLettersCount(value string) int {
	var valueJson entities.UserSendGrid
	json.Unmarshal([]byte(value), &valueJson)
	statusCode, err := web.client.UpdateSendgridEmail(valueJson.HostId,valueJson.Provider,valueJson.Email)
	if err != nil {
		return 0
	}
	if statusCode == http.StatusOK {
		m := make(map[string]string)
		m["method"] = textConstants.UpdateSendgridEmail
		m[textConstants.ProviderKey] = valueJson.Provider
		m[textConstants.HostIdKey] = valueJson.HostId
		m[textConstants.EmailKey] = valueJson.Email
		web.database.Save(m)
	}
	return statusCode
}


