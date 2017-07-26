package marketingClient

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type MarketingClient struct {
	baseApiUrl string
	httpToken  string
}

func NewMarketingClient(apiUrl, token string) *MarketingClient {
	return &MarketingClient{apiUrl, token}
}

func (client *MarketingClient) GetUserCount(userId string, provider string) (string, error,int) {
	const method = "customers/count"
	req, err := http.NewRequest("GET", client.baseApiUrl+method, nil)
	if err != nil {
		return "", err,0
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpToken)
	q := req.URL.Query()
	q.Add("host_id", userId)
	q.Add("provider", provider)
	req.URL.RawQuery = q.Encode()
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err,0
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err,0
	}

	return string(body), nil,response.StatusCode
}

func (client *MarketingClient) GetTransactionCount(userId string, provider string) (string, error,int) {
	const method = "customer_transactions/count"
	req, err := http.NewRequest("GET", client.baseApiUrl+method, nil)
	if err != nil {
		return "", err,0
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpToken)
	q := req.URL.Query()
	q.Add("host_id", userId)
	q.Add("provider", provider)
	req.URL.RawQuery = q.Encode()
	log.Print(q)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err,0
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err,0
	}

	return string(body), nil,response.StatusCode
}
func (client *MarketingClient) AddLettersTohost(userId string, provider string, lettersCount string) (int,error) {
	const method = "user/letters_count"
	//const method = ""
	form := url.Values{}
	form.Set("host_id", userId)
	form.Set("provider", provider)
	form.Set("lettersCount", lettersCount)

	buffer := new(bytes.Buffer)
	buffer.WriteString(form.Encode())
	//body := []byte("user_id=3&provider=radario&lettersCount=1")
	//	js,err := json.Marshal(form)
	log.Println(form.Encode())
	req, err := http.NewRequest("PATCH", client.baseApiUrl+method, buffer)
	if err != nil {
		log.Print(err)
		return http.StatusOK,err
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0,err
	}
	return response.StatusCode,nil


}
