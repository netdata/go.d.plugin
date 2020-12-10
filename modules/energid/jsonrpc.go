package energid

import "fmt"

// https://www.jsonrpc.org/specification#request_object
type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`
}

// http://www.jsonrpc.org/specification#response_object
type rpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *rpcError   `json:"error"`
	ID      int         `json:"id"`
}

// http://www.jsonrpc.org/specification#error_object
type rpcError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e rpcError) String() string {
	if e.Data == nil {
		return fmt.Sprintf("%s (code %d)", e.Message, e.Code)
	}
	return fmt.Sprintf("%s (code %d) (%v)", e.Message, e.Code, e.Data)
}
