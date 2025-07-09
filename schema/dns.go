package schema

import (
	"fmt"
	"strings"
)

type DNSDomain struct {
	Name  string
	NS    []string
	Owner string
	// Assuming they are same data type as TTL, RFC doesnt say if those are unsigned or not
	Serial   uint32
	Refresh  int32
	Retry    int32
	Expiry   int32
	Nxdomain int32
}

// Single DNS record structure
type DNSRecord struct {
	QType      string `json:"qtype"`
	QName      string `json:"qname"`
	Content    string `json:"content"`
	Ttl        int32  `json:"ttl"`
	DomainId   int    `json:"domain_id,omitempty"`
	ScopeMask  string `json:"scopeMask,omitempty"`
	AuthString string `json:"auth,omitempty"`
}

// sortable list of records, used usually as response
// empty list should be treated as "no records exist" and return accordingly
type DNSRecordList []DNSRecord

func (d *DNSDomain) UpdateSerial() {
	d.Serial += 1 // Yes, overflow is completely fine here,
}
func GenerateSoaFromDomain(d DNSDomain) DNSRecord {
	var rec DNSRecord
	rec.QType = "SOA"
	rec.QName = d.Name
	content := []string{
		d.NS[0],
		" ",
		d.Owner,
		" ",
		fmt.Sprintf("%d %d %d %d %d", d.Serial, d.Refresh, d.Retry, d.Expiry, d.Nxdomain),
	}
	rec.Content = strings.Join(content, "")
	rec.Ttl = d.Nxdomain
	return rec
}

// sort helper
func (slice DNSRecordList) Len() int {
	return len(slice)
}

// sort helper
func (slice DNSRecordList) Less(a, b int) bool {
	return (slice[a].QName + slice[a].Content + slice[a].QType) < (slice[b].QName + slice[b].Content + slice[b].QType)
}

// sort helper
func (slice DNSRecordList) Swap(a, b int) {
	slice[a], slice[b] = slice[b], slice[a]
}

// generate array of domains from subdomain, specific -> generic
func ExpandDNSName(name string) ([]string, error) {
	var s []string
	var err error

	parts := strings.Split(name, `.`)
	for i := 0; i < len(parts); i++ {
		s = append(s, strings.Join(parts[i:], `.`))
	}
	return s, err
}

type NXDomain struct{ Domain string }

func (n *NXDomain) Error() string {
	return "domain " + n.Domain + " not found"
}

type DomainReader interface {
	Lookup(q QueryLookup) (DNSRecordList, error)
	List(q QueryList) (DNSRecordList, error)
	// Find root domain for a given subdomain. Return NXDomain if it does not exist, anything else is db error
	GetRootDomainFor(string) (string, error)
}
type DomainWriter interface {
	// Add domain; that should also generate SOA record and AddRecord() it if backend doesn't handle that
	AddDomain(domain DNSDomain) error
	// add DNS record. if backend stores data per-domain it should figure out on its own to which DNSDomain it belongs; pdns doesn't send domain in request.
	AddRecord(record DNSRecord) error
}
type DomainReadWriter interface {
	DomainReader
	DomainWriter
}

type IPRedir interface {
	AddRedirIp(srcIp string, target string) error
	DeleteRedirIp(string) error
	SetRedirIp(map[string]string) error
	ListRedirIp() (map[string]string, error)
}
type DomainIPRedirReader interface {
	IPRedir
	DomainReader
}
type DomainIPRedirReadWriter interface {
	IPRedir
	DomainReader
	DomainWriter
}
