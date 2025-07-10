package api

import (
	"encoding/json"
	"github.com/efigence/go-powerdns/schema"
	"go.uber.org/zap"
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

func (api Api) Parse(raw string) (schema.QueryResponse, error) {
	var err error
	// parse "first level" of json to get type of query
	var objmap map[string]json.RawMessage
	err = json.Unmarshal([]byte(raw), &objmap)
	if err != nil {
		var n schema.QueryResponse
		return n, err
	}
	switch string(objmap[`method`]) {
	case `"lookup"`:
		var query schema.QueryLookup
		err := json.Unmarshal(objmap[`parameters`], &query)
		if err != nil {
			var n schema.QueryResponse
			return n, err
		}
		resp, err := api.dns.Lookup(query)
		return schema.QueryResponse{Result: resp}, err
	case `"list"`:
		var query schema.QueryList
		err := json.Unmarshal(objmap[`parameters`], &query)
		if err != nil {
			var n schema.QueryResponse
			return n, err
		}
		resp, err := api.dns.List(query)
		return recordToResponse(resp, err)
	case `"initialize"`:
		return schema.ResponseOk(), err
	case `"getAllDomains"`:
		return schema.ResponseFailed(), err
	default:
		var v interface{}
		err := json.Unmarshal(objmap[`parameters`], v)
		if err != nil {
			api.l.Error("could not parse [%s]", raw)
		} else {
			api.l.Errorf("unsupported cmd %+v", v)
		}
		return schema.ResponseFailed(), err
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
