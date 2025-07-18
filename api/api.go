package api

import (
	"encoding/json"
	"fmt"
	"github.com/efigence/go-powerdns/schema"
	"go.uber.org/zap"
	"strings"
)

type rawQuery struct {
	m string
	p map[string]json.RawMessage
}

func recordToResponse(list []schema.DNSRecord, err error) (schema.QueryResponse, error) {
	if err != nil {
		return schema.QueryResponse{Result: schema.ResponseFailed()}, err
	}

	return schema.QueryResponse{Result: list}, nil
}
func domainToResponse(list []schema.PDNSDomain, err error) (schema.QueryResponse, error) {
	if err != nil {
		return schema.QueryResponse{Result: schema.ResponseFailed()}, err
	}

	return schema.QueryResponse{Result: list}, nil
}
func (api Api) Parse(raw string) (schema.QueryResponse, error) {
	var err error
	// parse "first level" of json to get type of query
	var objmap map[string]json.RawMessage
	err = json.Unmarshal([]byte(raw), &objmap)
	if err != nil {
		var n schema.QueryResponse
		return n, err
	}
	method := strings.ToLower(string(objmap[`method`]))
	switch method {
	case `"initialize"`:
		return schema.ResponseOk(), err
	case `"lookup"`:
		var query schema.QueryLookup
		err := json.Unmarshal(objmap[`parameters`], &query)
		if err != nil {
			var n schema.QueryResponse
			return n, err
		}
		query.QName = strings.TrimRight(query.QName, ".")
		resp, err := api.dns.Lookup(query)
		return schema.QueryResponse{Result: resp}, err
	// no such thing as disabled records in our backends so list and APILookup are functionally same
	case `"list"`, `"apilookup"`:
		var query schema.QueryList
		err := json.Unmarshal(objmap[`parameters`], &query)
		if err != nil {
			var n schema.QueryResponse
			return n, err
		}
		strings.TrimRight(query.ZoneName, ".")
		resp, err := api.dns.List(query)
		return recordToResponse(resp, err)
	case `"getupdatedmaster"`:
		return schema.ResponseFailed("method %s not implemented, use native zone", string(objmap["method"])), nil
	case `"getalldomainmetadata"`: // no metadata support yet so we only need to respond with empty successful request
		return schema.QueryResponse{
			Result: map[string]string{},
		}, nil
	case `"getalldomains"`:
		domains, err := api.dns.ListDomains(false)
		pdnsDomains := schema.NewPDNSDomainList(domains)
		return domainToResponse(pdnsDomains, err)
	default:
		api.l.Warnf("unimplemented %s %s", method, string(objmap[`parameters`]))
		return schema.ResponseFailed("method %s not implemented", method), fmt.Errorf("unimplemented")
	}
}

type Api struct {
	dns schema.DomainReader
	l   *zap.SugaredLogger
}

func New(c schema.DomainReader, log *zap.SugaredLogger) (*Api, error) {
	api := Api{
		dns: c,
		l:   log,
	}
	return &api, nil
}
