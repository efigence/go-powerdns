package memdb

import (
	"fmt"
	"github.com/efigence/go-powerdns/schema"
	"strings"
	//	"gopkg.in/mem.v2"
)

func New() *MemDomains {
	var v MemDomains
	v.DomainRecords = make(map[string]map[string]schema.DNSRecordList)
	v.Domains = make(map[string]schema.DNSDomain)
	v.PerDomainRecords = map[string]schema.DNSRecordList{}
	return &v
}

type MemDomains struct {
	DomainRecords    map[string]map[string]schema.DNSRecordList
	Domains          map[string]schema.DNSDomain
	PerDomainRecords map[string]schema.DNSRecordList
}

func (d *MemDomains) GetRootDomainFor(dom string) (root string, err error) {
	v := strings.Split(dom, ".")
	if len(v) < 2 {
		return "", fmt.Errorf("domain needs to contain at least one dot")
	}
	if len(v) == 2 {
		if _, ok := d.Domains[dom]; ok {
			return dom, nil
		}
	}
	for dd := dom; len(v) >= 2; dd = strings.Join(v, ".") {
		v = v[1:]
		if _, ok := d.Domains[dd]; ok {
			return dd, nil
		}
	}
	return "", &schema.NXDomain{Domain: dom}
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
	if len(domain.NS) < 1 {
		return fmt.Errorf("domain needs at least one NS")
	}
	d.Domains[domain.Name] = domain
	d.AddRecord(schema.GenerateSoaFromDomain(domain))
	return err
}

// add records to DB
func (d *MemDomains) AddRecord(record schema.DNSRecord) error {
	var err error
	domName, err := d.GetRootDomainFor(record.QName)
	if err != nil {
		return err
	}
	if d.DomainRecords[record.QName] == nil {
		d.DomainRecords[record.QName] = make(map[string]schema.DNSRecordList)
	}
	d.DomainRecords[record.QName][record.QType] = append(d.DomainRecords[record.QName][record.QType], record)
	if _, ok := d.PerDomainRecords[domName]; !ok {
		d.PerDomainRecords[domName] = schema.DNSRecordList{}
	}
	d.PerDomainRecords[domName] = append(d.PerDomainRecords[domName], record)
	return err
}

// return records for query
func (d *MemDomains) Lookup(query schema.QueryLookup) (schema.DNSRecordList, error) {
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
func (d *MemDomains) List(q schema.QueryList) (r schema.DNSRecordList, err error) {
	if v, ok := d.PerDomainRecords[q.ZoneName]; ok {
		return v, err
	}
	return r, err
}
