package marketingClient

import (
	"bytes"
	"github.com/radario/MarketingSlackBot/mbot/textConstants"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type MarketingClient struct {
	baseApiUrl     string
	httpTokenValue string
	httpTokenKey   string
}

func NewMarketingClient(apiUrl, tokenValue, tokenKey string) *MarketingClient {
	return &MarketingClient{apiUrl, tokenValue, tokenKey}
}

func (client *MarketingClient) GetUserCount(userId string, provider string) (string, error, int) {
	req, err := http.NewRequest("GET", client.baseApiUrl+textConstants.GetCustomersCountMethod, nil)
	if err != nil {
		return "", err, 0
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpTokenValue)
	q := req.URL.Query()
	q.Add(textConstants.HostIdKey, userId)
	q.Add(textConstants.ProviderKey, provider)
	req.URL.RawQuery = q.Encode()
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err, 0
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err, 0
	}

	return string(body), nil, response.StatusCode
}

func (client *MarketingClient) GetTransactionCount(userId string, provider string) (string, error, int) {
	req, err := http.NewRequest("GET", client.baseApiUrl+textConstants.GetCustomersTransactionMethod, nil)
	if err != nil {
		return "", err, 0
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpTokenValue)
	q := req.URL.Query()
	q.Add(textConstants.HostIdKey, userId)
	q.Add(textConstants.ProviderKey, provider)
	req.URL.RawQuery = q.Encode()
	log.Print(q)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err, 0
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err, 0
	}

	return string(body), nil, response.StatusCode
}

func (client *MarketingClient) AddLettersTohost(userId string, provider string, lettersCount string) (int, error) {
	form := url.Values{}
	form.Set(textConstants.HostIdKey, userId)
	form.Set(textConstants.ProviderKey, provider)
	form.Set(textConstants.LettersCountKey, lettersCount)

	buffer := new(bytes.Buffer)
	buffer.WriteString(form.Encode())

	req, err := http.NewRequest("PUT", client.baseApiUrl+textConstants.AddUserLetterCountMethod, buffer)
	if err != nil {
		log.Print(err)
		return http.StatusOK, err
	}

	req.Header.Add(client.httpTokenKey, client.httpTokenValue)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	return response.StatusCode, nil

}

func (client *MarketingClient) UpdateSendgridEmail(userId string, provider string, email string) (int, error) {
	form := url.Values{}
	form.Set(textConstants.HostIdKey, userId)
	form.Set(textConstants.ProviderKey, provider)
	form.Set(textConstants.EmailKey, email)

	buffer := new(bytes.Buffer)
	buffer.WriteString(form.Encode())

	log.Println(form.Encode())
	req, err := http.NewRequest("PUT", client.baseApiUrl+textConstants.UpdateSendgridEmailMethod, buffer)
	if err != nil {
		log.Print(err)
		return http.StatusOK, err
	}

	req.Header.Add(client.httpTokenKey, client.httpTokenValue)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	return response.StatusCode, nil
}

func (client *MarketingClient) CreateScenarioByCampaign(campaignId string, scenarioName string) (int, error) {
	form := url.Values{}
	form.Set(textConstants.CampaignId, campaignId)
	form.Set(textConstants.ScenarioName, scenarioName)

	buffer := new(bytes.Buffer)
	buffer.WriteString(form.Encode())

	req, err := http.NewRequest("PUT", client.baseApiUrl+textConstants.CreateScenarioByCampaignMethod, buffer)
	if err != nil {
		log.Print(err)
		return http.StatusOK, err
	}

	req.Header.Add(client.httpTokenKey, client.httpTokenValue)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	return response.StatusCode, nil
}
