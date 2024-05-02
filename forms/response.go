package forms

type ResponseMessage struct {
	Message string       `json:"message"`
	Data    *MessageData `json:"data,omitempty"`
}

type MessageData interface{}