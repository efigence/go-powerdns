package schema

import "fmt"

// API calls
// https://doc.powerdns.com/md/authoritative/backend-remote/ for full docs

// Lookup call. Required for any plugin
type QueryLookup struct {
	// Query type. pdns will always ask for SOA records for domain, and often just query for ANY record even if source query was of type A
	QType string `json:"qtype"`
	// if pdns doesn't find domain via direct query it will ask for `*.domain` so implementing search for * is not neccesary
	QName      string `json:"qname"`
	Remote     string `json:"remote"`      // optional
	Local      string `json:"local"`       // optional
	RealRemote string `json:"real-remote"` // optional
	// -1 if unknown
	ZoneId int `json:"zone-id"`
}

func (data QueryLookup) Dump() string {
	return fmt.Sprintf("%+v", data)
}

type QueryList struct {
	ZoneName string `json:"zonename"`
	DomainId int    `json:"domain_id"` // optional
}
