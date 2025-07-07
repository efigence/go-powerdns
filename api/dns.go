package api

import (
	"github.com/efigence/go-powerdns/backend/schema"
)

// Domain + SOA data

// interface for backend

type NXDomain struct{ Domain string }

func (n *NXDomain) Error() string {
	return "domain " + n.Domain + " not found"
}

type DomainReader interface {
	Lookup(q QueryLookup) (schema.DNSRecordList, error)
	List(q QueryList) (schema.DNSRecordList, error)
	// Find root domain for a given subdomain. Return NXDomain if it does not exist, anything else is db error
	GetRootDomainFor(string) (string, error)
}
type DomainWriter interface {
	// Add domain; that should also generate SOA record and AddRecord() it if backend doesn't handle that
	AddDomain(domain schema.DNSDomain) error
	// add DNS record. if backend stores data per-domain it should figure out on its own to which DNSDomain it belongs; pdns doesn't send domain in request.
	AddRecord(record schema.DNSRecord) error
}
type DomainReadWriter interface {
	DomainReader
	DomainWriter
}
