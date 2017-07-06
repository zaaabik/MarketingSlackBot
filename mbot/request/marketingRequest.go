package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"io/ioutil"
	"github.com/radario/marketingstatbot/mbot/requestTypes"
)


type MarketingRequest struct {
	User			string
	RequestType		int
	RequestParams 	[]string
	Response		string
}

func (request *MarketingRequest)Send(httpToken string,httpAdress string)(error){
	//test data
	req,err := http.NewRequest("GET",httpAdress + requestTypes.RequestUrl[request.RequestType],nil)
	if err != nil{
		return err
	}
	req.URL.Query().Add("host_id",request.RequestParams[0])
	req.URL.Query().Add("provider",request.RequestParams[1])
	req.Header.Set("SECURITY",httpToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil{
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		return err
	}
	request.Response = string(body)
	return nil
}

func (r *MarketingRequest)Encode() ([]byte,error)  {
	enc, err := json.Marshal(r)
	if err != nil{
		return nil,errors.New("cant encode data to bytes")
	}
	return enc, nil
}








