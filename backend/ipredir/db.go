package ipredir

import (
	"github.com/efigence/go-powerdns/api"
	//	"gopkg.in/mem.v2"

)

func asdf() {
	_, _ = api.New(api.CallbackList{})
}

func New(file string) (api.DomainBackend, error) {
	var v ipredirDomains
	var err error
	v.DomainRecords = make(map[string]map[string]api.DNSRecordList)
	v.Domains = make(map[string]api.DNSDomain)
	return &v, err
}

type ipredirDomains struct {
	DomainRecords map[string]map[string]api.DNSRecordList
	Domains       map[string]api.DNSDomain
}

// add domain to DB
func (d *ipredirDomains) AddDomain(domain api.DNSDomain) error {
	var err error
	d.Domains[domain.Name] = domain
	return err
}

// add records to DB
func (d *ipredirDomains) AddRecord(domain string, record api.DNSRecord) error {
	var err error
	if (d.DomainRecords[record.QName] == nil) {
		d.DomainRecords[record.QName] = make(map[string]api.DNSRecordList)
	}
	d.DomainRecords[record.QName][record.QType] = append(d.DomainRecords[record.QName][record.QType], record)
	return err
}

// return records for query
func (d *ipredirDomains) Lookup(query api.QueryLookup) (api.DNSRecordList, error) {
	var err error
	return d.DomainRecords[query.QName][query.QType],err
}

// return all records for domain (For AXFR-type requests)
func (d *ipredirDomains) List(api.QueryList) (api.DNSRecordList, error) {
	var err error
	r := make([]api.DNSRecord, 0)
	return r, err
}
