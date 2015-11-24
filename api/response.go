package api


// JSONable response to pdns
type QueryResponse struct {
	Result map[string]interface{} `json:"result"`
}


func NewResponse() (QueryResponse) {
	var v QueryResponse
	v.Result = make(map[string]interface{})
	return v
}
