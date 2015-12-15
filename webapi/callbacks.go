package webapi

import (
	"github.com/efigence/go-powerdns/backend/ipredir"
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/api"
)

type dnsCB struct {
	redirBackend ipredir.DomainBackend
	memBackend  memdb.DomainBackend

}

func newDNSBackend() (dnsCB, error) {
	var v dnsCB
	var err error
	v.redirBackend,_ = ipredir.New("")
	v.memBackend,_ = memdb.New("")
	v.memBackend.AddDomain(api.DNSDomain{
		Name: "pdns.internal",
		PrimaryNs: "ns1.pdns.internal",
		Owner: "hostmaster.pdns.internal",
	});
	v.memBackend.AddRecord(api.DNSRecord{
		QType: "A",
		QName: "pdns.internal",
		Content: "127.0.0.1",
	})
	return v, err
}

func (b dnsCB) Lookup(q api.QueryLookup) (api.QueryResponse, error) {
	response, err := b.memBackend.Lookup(q)
	if (len(response) < 1) {
		response, err = b.redirBackend.Lookup(q)
	}
	if (len(response) == 0 ) {
		return api.ResponseOk(), err
	} else {

		return api.QueryResponse {
			Result: response,
		}, err
	}
}

func (b dnsCB) List(q api.QueryList) (api.QueryResponse, error) {
	response, err := b.memBackend.List(q)
	if (len(response) < 1) {
		response, err = b.redirBackend.List(q)
	}
	if (len(response) == 0 ) {
		return api.ResponseOk(), err
	} else {

		return api.QueryResponse {
			Result: response,
		}, err
	}
}
