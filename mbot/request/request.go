package request

import (
	"encoding/json"
	"github.com/radario/mbot/DB"
	"errors"
)


type Request struct {
	User			string
	RequestType		int
	RequestBody 	string
	Response		string
}




func (request *Request)Send()(string,error){
	//test data
	answerFromServer := "test"


	if(&(answerFromServer) == nil){
		return answerFromServer, errors.New("server doesn't respond")
	}
	request.Response = answerFromServer
	enc, err := request.Encode()
	if err != nil{
		return "",err
	}
	DB.Save(enc)

	return answerFromServer,nil
}

func (r *Request)Encode() ([]byte,error)  {
	enc, err := json.Marshal(r)
	if err != nil{
		return nil,errors.New("cant encode data to bytes")
	}
	return enc, nil
}










