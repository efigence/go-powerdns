package ipredir

import (
	"errors"
	"fmt"
	"github.com/efigence/go-powerdns/schema"
	"strings"
	"sync"
	//	"gopkg.in/mem.v2"
)

func New(backend schema.DomainReader) (*ipredirDomains, error) {
	var v ipredirDomains
	v.backend = backend
	var err error
	v.redirMap = make(map[string]string)
	return &v, err
}

type ipredirDomains struct {
	// map of host ip -> target IP redir
	backend  schema.DomainReader
	redirMap map[string]string
	sync.RWMutex
}

// add domain to DB
func (d *ipredirDomains) AddDomain(domain schema.DNSDomain) error {
	var err error
	return err
}

// add records to DB
func (d *ipredirDomains) AddRecord(record schema.DNSRecord) error {
	var err error

	return err
}
func (d *ipredirDomains) GetRootDomainFor(s string) (string, error) {
	return d.backend.GetRootDomainFor(s)
}

// Returns nil if request should not be redirected and DNS records if it should
func (d *ipredirDomains) Lookup(query schema.QueryLookup) ([]schema.DNSRecord, error) {
	var err error
	if val, ok := d.redirMap[query.Remote]; ok {
		if query.QType == "SOA" {
			// pretend we know the domain's root
			splitDomain, err := schema.ExpandDNSName(query.QName)
			if err != nil {
				return []schema.DNSRecord{}, err
			}
			var res schema.DNSRecord
			res.QType = "SOA"
			if len(splitDomain) > 1 {
				res.QName = splitDomain[len(splitDomain)-2]
				content := []string{
					"ns1.",
					res.QName,
					" hostmaster.",
					res.QName,
					" 1",
					" 10 10 10 10", // TTL 10 on everything
				}
				res.Content = strings.Join(content, "")
				res.Ttl = 10
				return []schema.DNSRecord{res}, err
			} else { // someone thinks we're root domain.... nope
				return []schema.DNSRecord{}, errors.New(fmt.Sprintf("too short domain %+v, we're not handling root", splitDomain))
			}
		} else {
			var res schema.DNSRecord
			res.QName = query.QName
			res.QType = "A"
			res.Content = val
			return []schema.DNSRecord{res}, err
		}
	}
	return []schema.DNSRecord{}, nil
}

// return all records for domain (For AXFR-type requests)
// Returns nil if request should not be redirected and DNS records if it should
func (d *ipredirDomains) List(schema.QueryList) ([]schema.DNSRecord, error) {
	var err error
	return []schema.DNSRecord{}, err
}
func (db *ipredirDomains) ListDomains(disabled bool) ([]schema.DNSDomain, error) {
	return db.backend.ListDomains(disabled)
}
