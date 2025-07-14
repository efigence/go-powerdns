package memdb

import (
	"fmt"
	"github.com/efigence/go-powerdns/schema"
	"strings"
	//	"gopkg.in/mem.v2"
)

func New() *MemDomains {
	var v MemDomains
	v.DomainRecords = make(map[string]map[string][]schema.DNSRecord)
	v.Domains = make(map[string]schema.DNSDomain)
	v.PerDomainRecords = map[string][]schema.DNSRecord{}
	return &v
}

type MemDomains struct {
	DomainRecords    map[string]map[string][]schema.DNSRecord
	Domains          map[string]schema.DNSDomain
	PerDomainRecords map[string][]schema.DNSRecord
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
		d.DomainRecords[record.QName] = make(map[string][]schema.DNSRecord)
	}
	d.DomainRecords[record.QName][record.QType] = append(d.DomainRecords[record.QName][record.QType], record)
	if _, ok := d.PerDomainRecords[domName]; !ok {
		d.PerDomainRecords[domName] = []schema.DNSRecord{}
	}
	d.PerDomainRecords[domName] = append(d.PerDomainRecords[domName], record)
	return err
}

// return records for query
func (d *MemDomains) Lookup(query schema.QueryLookup) ([]schema.DNSRecord, error) {
	var err error
	if query.QType == `ANY` {
		var recordsAny []schema.DNSRecord
		if rec, ok := d.DomainRecords[query.QName]; ok {
			for _, records := range rec {
				recordsAny = append(recordsAny, records...)
			}
			return recordsAny, err
		}
		_, dom, _ := strings.Cut(query.QName, ".")
		if rec, ok := d.DomainRecords["*."+dom]; ok {
			for _, records := range rec {
				for _, r := range records {
					r.QName = query.QName
					recordsAny = append(recordsAny, r)
				}
			}
			return recordsAny, err
		}
		return recordsAny, err
	} else {
		if v, ok := d.DomainRecords[query.QName][query.QType]; ok {
			return v, err
		}
		_, dom, _ := strings.Cut(query.QName, ".")
		if v, ok := d.DomainRecords["*."+dom][query.QType]; ok {
			for idx := range v {
				v[idx].QName = query.QName
			}
			return v, err
		}
	}
	return []schema.DNSRecord{}, nil
}

// return all records for domain (For AXFR-type requests)
func (d *MemDomains) List(q schema.QueryList) (r []schema.DNSRecord, err error) {
	if v, ok := d.PerDomainRecords[q.ZoneName]; ok {
		return v, err
	}
	return r, err
}

func (d *MemDomains) ListDomains(disabled bool) ([]schema.DNSDomain, error) {
	domlist := []schema.DNSDomain{}
	for _, dom := range d.Domains {
		domlist = append(domlist, schema.DNSDomain{
			Name:     dom.Name,
			NS:       dom.NS,
			Owner:    dom.Owner,
			Serial:   dom.Serial,
			Refresh:  dom.Refresh,
			Retry:    dom.Retry,
			Expiry:   dom.Expiry,
			Nxdomain: dom.Nxdomain,
		})
	}
	return domlist, nil

}
