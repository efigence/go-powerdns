package api

// Single DNS record structure
type DNSRecord struct{
	QType string `json:"qtype"`
	QName string `json:"qname"`
	Content string `json:"content"`
	Ttl int32 `json:"ttl"`
	DomainId int `json:"domain_id,omitempty"`
	ScopeMask string `json:"scopeMask,omitempty"`
	AuthString string`json:"auth,omitempty"`
}


// Domain + SOA data

type DNSDomain struct {
	Name string
	Serial uint32
	PrimaryNs string
	Owner string
	// Assuming they are same data type as TTL, RFC doesnt say if those are unsigned or not
	Expiry int32
	Refresh int32
	Retry int32
	Nxdomain int32
}


// interface for backend

type DomainBackend interface {
	AddDomain(domain DNSDomain) error
	AddRecord(domain string, record DNSRecord) error
	Search(q QueryLookup) ([]DNSRecord, error)
	List(q QueryList) ([]DNSRecord, error)
}
