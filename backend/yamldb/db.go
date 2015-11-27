package yamldb

import (
	"github.com/efigence/go-powerdns/api"
)

func asdf() {
	_, _ = api.New(api.CallbackList{})
}

func New() api.DomainBackend {
	var v YAMLDomains
	v.DomainRecords = make(map[string]map[string][]api.DNSRecord)
	v.Domains = make(map[string]api.DNSDomain)
	return &v
}

type YAMLDomains struct {
	DomainRecords map[string]map[string][]api.DNSRecord
	Domains       map[string]api.DNSDomain
}


// add domain to DB
func (d *YAMLDomains) AddDomain(api.DNSDomain) error {
	var err error
	return err
}

// add records to DB
func (d *YAMLDomains) AddRecord(domain string, record api.DNSRecord) error {
	var err error
	return err
}

// return records for query
func (d *YAMLDomains) Search(api.QueryLookup) ([]api.DNSRecord, error) {
	var err error
	r := make([]api.DNSRecord, 0)
	return r, err
}

// return all records for domain (For AXFR-type requests)
func (d *YAMLDomains) List(api.QueryList) ([]api.DNSRecord, error) {
	var err error
	r := make([]api.DNSRecord, 0)
	return r, err
}
