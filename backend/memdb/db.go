package memdb

import (
	"github.com/efigence/go-powerdns/api"
	//	"gopkg.in/mem.v2"

)

func asdf() {
	_, _ = api.New(api.CallbackList{})
}

func New(file string) (api.DomainBackend, error) {
	var v MemDomains
	var err error
	v.DomainRecords = make(map[string]map[string]api.DNSRecordList)
	v.Domains = make(map[string]api.DNSDomain)
	return &v, err
}

type MemDomains struct {
	DomainRecords map[string]map[string]api.DNSRecordList
	Domains       map[string]api.DNSDomain
}

// add domain to DB
func (d *MemDomains) AddDomain(domain api.DNSDomain) error {
	var err error
	d.Domains[domain.Name] = domain
	return err
}

// add records to DB
func (d *MemDomains) AddRecord(domain string, record api.DNSRecord) error {
	var err error
	if (d.DomainRecords[domain] == nil) {
		d.DomainRecords[domain] = make(map[string]api.DNSRecordList)
	}
	d.DomainRecords[domain][record.QName] = append(d.DomainRecords[domain][record.QName], record)
	return err
}

// return records for query
func (d *MemDomains) Search(query api.QueryLookup) (api.DNSRecordList, error) {
	var err error
	return d.DomainRecords[`example.com`][query.QName],err
}

// return all records for domain (For AXFR-type requests)
func (d *MemDomains) List(api.QueryList) (api.DNSRecordList, error) {
	var err error
	r := make([]api.DNSRecord, 0)
	return r, err
}