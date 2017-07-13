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
	Callback chan bool
}

func NewWebHookHandler()(*WebHook){
	return &WebHook{make(chan bool)}
}

func (web WebHook) Start() {
	r := chi.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("catch")
		w.Write([]byte("hello"))
		res, _ := ioutil.ReadAll(r.Body)
		if len(res) < 8 {
			//web.Callback <- false
		} else {
			jsonStr, _ := url.QueryUnescape(string(res)[8:])
			rand.Seed(time.Now().UTC().UnixNano())
			var s slack.AttachmentActionCallback
			json.Unmarshal([]byte(jsonStr), &s)
			web.Callback <- true

		}
	})
	http.ListenAndServe(":1112", r)
}
