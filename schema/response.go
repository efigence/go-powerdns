package schema

// JSONable response to pdns

type QueryResponse struct {
	Result interface{} `json:"result"`
}

func NewResponse() QueryResponse {
	var v QueryResponse
	v.Result = make(map[string]interface{})
	return v
}

// basic OK response
func ResponseFailed() QueryResponse {
	var v QueryResponse
	v.Result = false
	return v
}

// basic "query failed" response
func ResponseOk() QueryResponse {
	var v QueryResponse
	v.Result = true
	return v
}
