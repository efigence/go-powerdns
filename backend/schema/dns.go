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
