package marketingClient

import (
	"bytes"
	"fmt"
	"github.com/radario/MarketingSlackBot/mbot/textConstants"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const methodTemplate = "%s%s"

type MarketingClient struct {
	baseApiUrl                     string
	httpTokenValue                 string
	httpTokenKey                   string
	getCustomersCountMethod        string
	getTransactionCountMethod      string
	addLettersToHostMethod         string
	upgradeSendgridMethod          string
	createScenarioByCampaignMethod string
	unlockUserMethod               string
	lockUserMethod                 string
}

func NewMarketingClient(apiUrl, tokenValue, tokenKey string) *MarketingClient {
	client := &MarketingClient{baseApiUrl: apiUrl, httpTokenValue: tokenValue, httpTokenKey: tokenKey}
	client.getCustomersCountMethod = fmt.Sprintf(methodTemplate, apiUrl, textConstants.GetCustomersCountMethod)
	client.getTransactionCountMethod = fmt.Sprintf(methodTemplate, apiUrl, textConstants.GetCustomersTransactionMethod)
	client.addLettersToHostMethod = fmt.Sprintf(methodTemplate, apiUrl, textConstants.AddUserLetterCountMethod)
	client.upgradeSendgridMethod = fmt.Sprintf(methodTemplate, apiUrl, textConstants.UpdateSendgridEmailMethod)
	client.createScenarioByCampaignMethod = fmt.Sprintf(methodTemplate, apiUrl, textConstants.CreateScenarioByCampaignMethod)
	client.unlockUserMethod = fmt.Sprintf(methodTemplate, apiUrl, textConstants.UnlockUserMethod)
	client.lockUserMethod = fmt.Sprintf(methodTemplate, apiUrl, textConstants.LockUserMethod)
	return client
}

func (client *MarketingClient) GetUserCount(userId string, provider string) (string, int, error) {
	req, err := http.NewRequest("GET", client.getCustomersCountMethod, nil)
	if err != nil {
		return "", 0, err
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpTokenValue)
	q := req.URL.Query()
	q.Add(textConstants.HostIdKey, userId)
	q.Add(textConstants.ProviderKey, provider)
	req.URL.RawQuery = q.Encode()
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", 0, err
	}

	return string(body), response.StatusCode, nil
}

func (client *MarketingClient) GetTransactionCount(userId string, provider string) (string, int, error) {
	req, err := http.NewRequest("GET", client.getTransactionCountMethod, nil)
	if err != nil {
		return "", 0, err
	}

	req.Header.Add("X-MARKETING-SECURITY", client.httpTokenValue)
	q := req.URL.Query()
	q.Add(textConstants.HostIdKey, userId)
	q.Add(textConstants.ProviderKey, provider)
	req.URL.RawQuery = q.Encode()
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", 0, err
	}

	return string(body), response.StatusCode, nil
}

func (client *MarketingClient) AddLettersToHost(userId string, provider string, lettersCount string) (int, error) {
	form := url.Values{}
	form.Set(textConstants.HostIdKey, userId)
	form.Set(textConstants.ProviderKey, provider)
	form.Set(textConstants.LettersCountKey, lettersCount)

	buffer := new(bytes.Buffer)
	buffer.WriteString(form.Encode())

	req, err := http.NewRequest("PUT", client.addLettersToHostMethod, buffer)
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

	req, err := http.NewRequest("PUT", client.upgradeSendgridMethod, buffer)
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

func (client *MarketingClient) LockUser(userId string, provider string, lock bool) (int, error) {
	form := url.Values{}
	form.Set(textConstants.HostIdKey, userId)
	form.Set(textConstants.ProviderKey, provider)
	form.Set(textConstants.Lock, "true")

	buffer := new(bytes.Buffer)
	buffer.WriteString(form.Encode())

	req, err := http.NewRequest("POST", client.lockUserMethod, buffer)
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

func (client *MarketingClient) UnlockUser(userId string, provider string, lock bool) (int, error) {
	form := url.Values{}
	form.Set(textConstants.HostIdKey, userId)
	form.Set(textConstants.ProviderKey, provider)
	form.Set(textConstants.Lock, "false")

	buffer := new(bytes.Buffer)
	buffer.WriteString(form.Encode())

	req, err := http.NewRequest("POST", client.unlockUserMethod, buffer)
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

	req, err := http.NewRequest("PUT", client.createScenarioByCampaignMethod, buffer)
	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError, err
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
