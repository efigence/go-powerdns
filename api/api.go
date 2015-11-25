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
	Query(request QueryLookup) (QueryResponse, error)
}

type QueryListCB interface {
	Query(request QueryList) (QueryResponse, error)
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
		return api.callbacks.Lookup.Query(query)
	case `"list"`:
		var query QueryList
		err := json.Unmarshal(objmap[`parameters`],&query)
		if err != nil {
			var n QueryResponse
			return n, err
		}
		return api.callbacks.List.Query(query)
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





// API calls
// https://doc.powerdns.com/md/authoritative/backend-remote/ for full docs

// Lookup call. Required for any plugin
type QueryLookup struct {
	QType string `json:"qtype"`
	QName string `json:"qname"`
	Remote string `json:"remote"` // optional
	Local string `json:"local"`// optional
	RealRemote string `json:"real-remote"` // optional
	ZoneId int `json:"zone-id"`
}


type QueryList struct {
	ZoneName string `json:"zonename"`
	DomainId int `json:"domain_id"` // optional
}

func (data QueryLookup) Dump() string {
	return fmt.Sprintf("%+v", data)
}
