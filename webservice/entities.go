package webservice

type ResponseMessage struct {
	IsOk bool `json:"is_ok"`
	Message string `json:"message"`
}

func NewResponseMessage(msg string,isOk bool) *ResponseMessage {
	return &ResponseMessage{
		Message: msg,
		IsOk: isOk,
	}
}

