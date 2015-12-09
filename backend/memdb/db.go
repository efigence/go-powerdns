package memdb

import (
	"github.com/efigence/go-powerdns/api"
	//	"gopkg.in/mem.v2"

)


func New(file string) (DomainBackend, error) {
	var v MemDomains
	var err error
	v.DomainRecords = make(map[string]map[string]api.DNSRecordList)
	v.Domains = make(map[string]api.DNSDomain)
	return &v, err
}

type DomainBackend interface {
	api.DomainBackend
}

type MemDomains struct {
	DomainRecords map[string]map[string]api.DNSRecordList
	Domains       map[string]api.DNSDomain
}

// add domain to DB
func (d *MemDomains) AddDomain(domain api.DNSDomain) error {
	var err error
	// some defaults
	if (domain.Owner == "") {domain.Owner = "hostmaster." + domain.Name}
	if (domain.Refresh == 0) {domain.Refresh = 86400 * 2 }
	if (domain.Retry == 0) {domain.Retry = 60 * 15}
	if (domain.Expiry == 0) {domain.Expiry = 86400 * 14 }
	if (domain.Nxdomain == 0) {domain.Nxdomain = 60 * 30}
	d.AddRecord(api.GenerateSoaFromDomain(domain))
	return err
}

// add records to DB
func (d *MemDomains) AddRecord(record api.DNSRecord) error {
	var err error
	if (d.DomainRecords[record.QName] == nil) {
		d.DomainRecords[record.QName] = make(map[string]api.DNSRecordList)
	}
	d.DomainRecords[record.QName][record.QType] = append(d.DomainRecords[record.QName][record.QType], record)
	return err
}

// return records for query
func (d *MemDomains) Lookup(query api.QueryLookup) (api.DNSRecordList, error) {
	var err error
	return d.DomainRecords[query.QName][query.QType],err
}

// return all records for domain (For AXFR-type requests)
func (d *MemDomains) List(api.QueryList) (api.DNSRecordList, error) {
	var err error
	r := make([]api.DNSRecord, 0)
	return r, err
}
