package webHookHandler

import (
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/url"
	"time"
	"github.com/nlopes/slack"
	"encoding/json"
	"fmt"
	"net/http"
	"math/rand"
)

type WebHook struct {
	r *chi.Mux
}
func (w WebHook)start(){
	w.r = chi.NewRouter()
	http.ListenAndServe(":1235",w.r)
}
func (w WebHook)WebHookHandler(calback chan string) {
	w.r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		res,_ := ioutil.ReadAll(r.Body)
		jsonStr,_ := url.QueryUnescape(string(res)[8:])
		rand.Seed(time.Now().UTC().UnixNano())
		var s slack.AttachmentActionCallback
		json.Unmarshal([]byte(jsonStr),&s)
		fmt.Println(s)
		if s.Actions[0].Value == "yes"{
			calback <- "yes"
		} else {
			calback <- "no"
		}
	})
}