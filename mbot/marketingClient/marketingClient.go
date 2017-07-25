package marketingClient

import (
	"io/ioutil"
	"net/http"
	"log"
)

type MarketingClient struct {
	baseApiUrl string
	httpToken  string
}

func NewMarketingClient(apiUrl, token string) *MarketingClient {
	return &MarketingClient{apiUrl, token}
}

func (client *MarketingClient) GetUserCount(userId string, provider string) (string, error) {
	const method = "customers/count"
	req, err := http.NewRequest("GET", client.baseApiUrl+method, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpToken)
	q := req.URL.Query()
	q.Add("host_id", userId)
	q.Add("provider", provider)
	req.URL.RawQuery = q.Encode()
	log.Print(req.URL)
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
	const method = "customer_transactions/count"
	req, err := http.NewRequest("GET", client.baseApiUrl + method, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpToken)
	q := req.URL.Query()
	q.Add("host_id", userId)
	q.Add("provider", provider)
	req.URL.RawQuery = q.Encode()
	log.Print(q)
	response, err := http.DefaultClient.Do(req)
	if err != nil{
		return "",err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil

}
