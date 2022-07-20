package response

import (
	"fmt"
)

type (
	Message struct {
		StatusCode int `json:"statusCode"`
		Code       string `json:"code"`
		Message    string `json:"message"`
	}
	Response struct {
		StatusCode int         `json:"statusCode"` // statusCode == 0: normal, statusCode > 0
		MsgCode    string      `json:"msgCode"`
		Message    string      `json:"message"`
		Payload    interface{} `json:"payload"`
	}
	RespAffectedRow struct {
		RowAffected int64 `json:"rowAffected"`
	}
)

func NewMessage(statusCode int, code string, message string) Message {
	return Message{StatusCode: statusCode, Code: code, Message: message}
}

func NewResponse(statusCode int, msgCode string, msg string, payload interface{}) (resp Response) {
	var r Response
	r.StatusCode = statusCode
	r.MsgCode = msgCode
	r.Message = msg
	r.Payload = payload
	return r
}

func NewResponseByMessage(payload interface{}, message Message, v ...interface{}) Response {
	m := message.Message
	if v != nil {
		m = fmt.Sprintf(m, v...)
	}
	return NewResponse(message.StatusCode, message.Code, m, payload)
}

func NewResponseByMessageWithStatusCode(statusCode int, payload interface{}, message Message, v ...interface{}) Response {
	m := message.Message
	sc := message.StatusCode
	if v != nil {
		m = fmt.Sprintf(m, v...)
	}
	// If input statusCode is not null than use it
	if statusCode != 0{
		sc = statusCode
	}
	return NewResponse(sc, message.Code, m, payload)
}
