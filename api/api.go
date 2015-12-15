package api

import (
	"encoding/json"
	"fmt"
)


type rawQuery struct {
	m string
	p map[string]json.RawMessage
}


type QueryLookupCB interface {
	Lookup(request QueryLookup) (QueryResponse, error)
}

type QueryListCB interface {
	List(request QueryList) (QueryResponse, error)
}


type CallbackList struct {
	Lookup QueryLookupCB
	List QueryListCB
}


func (api Api)Parse(raw string) (QueryResponse, error) {
	var err error
	// parse "first level" of json to get type of query
	var objmap map[string]json.RawMessage
	err = json.Unmarshal([]byte(raw),&objmap)
	if err != nil {
		var n QueryResponse
		return n, err
	}
	switch string(objmap[`method`]) {
	case `"lookup"`:
		var query QueryLookup
		err := json.Unmarshal(objmap[`parameters`],&query)
		if err != nil {
			var n QueryResponse
			return n, err
		}
		return api.callbacks.Lookup.Lookup(query)
	case `"list"`:
		var query QueryList
		err := json.Unmarshal(objmap[`parameters`],&query)
		if err != nil {
			var n QueryResponse
			return n, err
		}
		return api.callbacks.List.List(query)
	case `"initialize"`:
		return ResponseOk(), err
	default:
		return ResponseFailed(), err
	}
}

type Api struct {
	callbacks CallbackList
}

func New(c CallbackList) (Api, error) {
	var api Api
	api.callbacks = c
	var err error
	return api, err
}







func (data QueryLookup) Dump() string {
	return fmt.Sprintf("%+v", data)
}
