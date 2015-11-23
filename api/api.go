package api

import (
	"encoding/json"
	"fmt"
)


type rawQuery struct {
	m string
	p map[string]json.RawMessage
}

type Query interface {
	Dump() string
}


func ParseJson(raw string) (Query,error) {
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

// API calls
// https://doc.powerdns.com/md/authoritative/backend-remote/ for full docs


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
