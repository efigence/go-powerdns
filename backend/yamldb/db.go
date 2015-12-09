package yamldb

import (
	"github.com/efigence/go-powerdns/api"
	"os"

)

func asdf() {
	_, _ = api.New(api.CallbackList{})
}

func New(file string) (api.DomainBackend, error) {
	var v YAMLDomains
	data, err := os.Open(file)
	if err != nil {
		return &v, err
	}
	v.ParseDNS(data)
	v.DomainRecords = make(map[string]map[string]api.DNSRecordList)
	v.Domains = make(map[string]api.DNSDomain)
	return &v, err
}

type YAMLDomains struct {
	DomainRecords map[string]map[string]api.DNSRecordList
	Domains       map[string]api.DNSDomain
}

// add domain to DB
func (d *YAMLDomains) AddDomain(api.DNSDomain) error {
	var err error
	return err
}

// add records to DB
func (d *YAMLDomains) AddRecord(record api.DNSRecord) error {
	var err error
	return err
}

// return records for query
func (d *YAMLDomains) Lookup(api.QueryLookup) (api.DNSRecordList, error) {
	var err error
	r := make([]api.DNSRecord, 0)
	return r, err
}

// return all records for domain (For AXFR-type requests)
func (d *YAMLDomains) List(api.QueryList) (api.DNSRecordList, error) {
	var err error
	r := make([]api.DNSRecord, 0)
	return r, err
}
