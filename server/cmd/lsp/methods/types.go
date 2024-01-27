package methods

type RequestMessage struct {
	Method string      `json:"method"`
	Id     int         `json:"id"`
	Params interface{} `json:"params"`
	Message
}

type Message struct {
	Jsonrpc string `json:"jsonrpc"`
}
