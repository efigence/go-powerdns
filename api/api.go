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
//	List QueryListCB
}


func (api Api)ParseRequest(raw string) (Query,error) {
	var objmap map[string]json.RawMessage
	var err error
	err = json.Unmarshal([]byte(raw),&objmap)
	if string(objmap[`method`]) == "\"lookup\"" {
		var part QueryLookup
		err := json.Unmarshal(objmap[`parameters`],&part)
		return part,err
	}
	return nil,err
}

func (api Api)Parse(raw string) (QueryResponse, error) {
	var objmap map[string]json.RawMessage
	var err error
	err = json.Unmarshal([]byte(raw),&objmap)
	if string(objmap[`method`]) == "\"lookup\"" {
		var query QueryLookup
		err := json.Unmarshal(objmap[`parameters`],&query)
		if err != nil {
			var n QueryResponse
			return n, err
		}
		return api.callbacks.Lookup.Query(query)
	}
	var n QueryResponse
	return n, err
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
	DomainId string `json:"domain_id"` // optional
}

func (data QueryLookup) Dump() string {
	return fmt.Sprintf("%+v", data)
}
