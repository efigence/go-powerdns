package ipredir

import (
	"errors"
	"fmt"
	"github.com/efigence/go-powerdns/api"
	"strings"
	"sync"
	//	"gopkg.in/mem.v2"
)

func New(file string) (DomainBackend, error) {
	var v ipredirDomains
	var err error
	v.redirMap = make(map[string]string)
	return &v, err
}

type DomainBackend interface {
	api.DomainBackend
	AddRedirIp(srcIp string, target string) error
	DeleteRedirIp(string) error
	SetRedirIp(map[string]string) error
	ListRedirIp() (map[string]string, error)
}

type ipredirDomains struct {
	// map of host ip -> target IP redir
	redirMap map[string]string
	sync.RWMutex
}

// add domain to DB
func (d *ipredirDomains) AddDomain(domain api.DNSDomain) error {
	var err error
	return err
}

// add records to DB
func (d *ipredirDomains) AddRecord(record api.DNSRecord) error {
	var err error

	return err
}

// Returns nil if request should not be redirected and DNS records if it should
func (d *ipredirDomains) Lookup(query api.QueryLookup) (api.DNSRecordList, error) {
	var err error
	if val, ok := d.redirMap[query.Remote]; ok {
		if query.QType == "SOA" {
			// pretend we know the domain's root
			splitDomain, err := api.ExpandDNSName(query.QName)
			if err != nil {
				return api.DNSRecordList{}, err
			}
			var res api.DNSRecord
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
				return api.DNSRecordList{res}, err
			} else { // someone thinks we're root domain.... nope
				return api.DNSRecordList{}, errors.New(fmt.Sprintf("too short domain %+v, we're not handling root", splitDomain))
			}
		} else {
			var res api.DNSRecord
			res.QName = query.QName
			res.QType = "A"
			res.Content = val
			return api.DNSRecordList{res}, err
		}
	}
	return api.DNSRecordList{}, nil
}

// return all records for domain (For AXFR-type requests)
// Returns nil if request should not be redirected and DNS records if it should
func (d *ipredirDomains) List(api.QueryList) (api.DNSRecordList, error) {
	var err error
	return api.DNSRecordList{}, err
	return nil, err
}
