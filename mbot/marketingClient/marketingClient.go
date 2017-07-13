package marketingClient

import (
	"io/ioutil"
	"log"
	"net/http"
)

type MarketingClient struct {
	baseApiUrl string
	httpToken  string
}

func NewMarketingCliet(apiUrl, token string) *MarketingClient {
	return &MarketingClient{apiUrl, token}
}
func (client *MarketingClient) GetUserCount(userId string, provider string) (string, error) {
	const method = "/getUserCount"
	req, err := http.NewRequest("GET", client.baseApiUrl+method, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("SECURITY", client.httpToken)
	q := req.URL.Query()
	q.Add("host_id", userId)
	q.Add("provider", provider)
	req.URL.RawQuery = q.Encode()
	log.Println(req.URL.String())
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil

}

func (client *MarketingClient) GetTransactionCount(userId string, provider string) (string, error) {
	const method = "/getTransactionCount"
	req, err := http.NewRequest("GET", client.baseApiUrl+method, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("SECURITY", client.httpToken)
	q := req.URL.Query()
	q.Add("host_id", userId)
	q.Add("provider", provider)
	req.URL.RawQuery = q.Encode()
	log.Println(req.URL.String())
	response, err := http.DefaultClient.Do(req)

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil

}
