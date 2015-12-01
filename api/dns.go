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

// sortable list of records, used usually as response
type DNSRecordList []DNSRecord


// sort helper
func (slice DNSRecordList) Len() int {
	return len(slice)
}

// sort helper
func (slice DNSRecordList) Less(a, b int) bool {
	return ( slice[a].QName +slice[a].Content + slice[a].QType ) < ( slice[b].QName +slice[b].Content + slice[b].QType );
}

// sort helper
func (slice DNSRecordList) Swap(a, b int) {
	slice[a], slice[b] = slice[b], slice[a]
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
func (d *DNSDomain)UpdateSerial() {
	d.Serial += 1 // Yes, overflow is completely fine here,
}


// interface for backend

type DomainBackend interface {
	AddDomain(domain DNSDomain) error
	AddRecord(domain string, record DNSRecord) error
	Search(q QueryLookup) (DNSRecordList, error)
	List(q QueryList) (DNSRecordList, error)
}
