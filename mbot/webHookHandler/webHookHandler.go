package webHookHandler

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/nlopes/slack"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"log"
)

type WebHook struct {
	Callback chan string
}

func NewWebHookHandler()(*WebHook){
	return &WebHook{make(chan string,3)}
}

type j struct {
	Value []slack.Attachment `json:"attachments"`
}

func (web WebHook) Start() {
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		res, _ := ioutil.ReadAll(r.Body)
		log.Println(string(res))
		if len(res) < 8 {
			web.Callback <- "error"
			close(web.Callback)
		} else {
			jsonStr, _ := url.QueryUnescape(string(res)[8:])
			var s slack.AttachmentActionCallback
			json.Unmarshal([]byte(jsonStr), &s)
			web.Callback <- s.Actions[0].Value
			if s.Actions[0].Value == "no"{
				res := "CANCEL"

				w.Write([]byte(res))
			} else {
				res := "ACCEPT"
				w.Write([]byte(res))

			}

		}
	})
	http.ListenAndServe(":1113", r)
}
