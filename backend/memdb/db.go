package memdb

import (
	"github.com/efigence/go-powerdns/api"
	"github.com/efigence/go-powerdns/backend/schema"
	//	"gopkg.in/mem.v2"
)

func New() (*MemDomains, error) {
	var v MemDomains
	var err error
	v.DomainRecords = make(map[string]map[string]schema.DNSRecordList)
	v.Domains = make(map[string]schema.DNSDomain)
	return &v, err
}

type MemDomains struct {
	DomainRecords map[string]map[string]schema.DNSRecordList
	Domains       map[string]schema.DNSDomain
}

// add domain to DB
func (d *MemDomains) AddDomain(domain schema.DNSDomain) error {
	var err error
	// some defaults
	if domain.Owner == "" {
		domain.Owner = "hostmaster." + domain.Name
	}
	if domain.Refresh == 0 {
		domain.Refresh = 86400 * 2
	}
	if domain.Retry == 0 {
		domain.Retry = 60 * 15
	}
	if domain.Expiry == 0 {
		domain.Expiry = 86400 * 14
	}
	if domain.Nxdomain == 0 {
		domain.Nxdomain = 60 * 30
	}
	d.Domains[domain.Name] = domain
	d.AddRecord(schema.GenerateSoaFromDomain(domain))
	return err
}

// add records to DB
func (d *MemDomains) AddRecord(record schema.DNSRecord) error {
	var err error
	if d.DomainRecords[record.QName] == nil {
		d.DomainRecords[record.QName] = make(map[string]schema.DNSRecordList)
	}
	d.DomainRecords[record.QName][record.QType] = append(d.DomainRecords[record.QName][record.QType], record)
	return err
}

// return records for query
func (d *MemDomains) Lookup(query api.QueryLookup) (schema.DNSRecordList, error) {
	var err error
	if query.QType == `ANY` {
		var recordsAny schema.DNSRecordList
		for _, records := range d.DomainRecords[query.QName] {
			recordsAny = append(recordsAny, records...)
		}
		return recordsAny, err
	} else {
		return d.DomainRecords[query.QName][query.QType], err
	}
}

// return all records for domain (For AXFR-type requests)
func (d *MemDomains) List(api.QueryList) (schema.DNSRecordList, error) {
	var err error
	r := make([]schema.DNSRecord, 0)
	return r, err
}
