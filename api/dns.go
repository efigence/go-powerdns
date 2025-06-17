package api

import (
	"github.com/efigence/go-powerdns/backend/schema"
)

// Domain + SOA data

// interface for backend

type DomainReader interface {
	Lookup(q QueryLookup) (schema.DNSRecordList, error)
	List(q QueryList) (schema.DNSRecordList, error)
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
