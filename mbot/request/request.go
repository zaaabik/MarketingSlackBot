package request

import (
	"encoding/json"
	"errors"
)


type Request struct {
	User			string
	RequestType		int
	RequestBody 	string
	Response		string
}




func (request *Request)Send()(error){
	//test data
	answerFromServer := request.RequestBody


	if(&(answerFromServer) == nil){
		return errors.New("server doesn't respond")
	}
	request.Response = answerFromServer
	_, err := request.Encode()
	if err != nil{
		return err
	}
	return nil
}

func (r *Request)Encode() ([]byte,error)  {
	enc, err := json.Marshal(r)
	if err != nil{
		return nil,errors.New("cant encode data to bytes")
	}
	return enc, nil
}










