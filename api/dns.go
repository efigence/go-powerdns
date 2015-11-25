package api

// Single DNS record structure
type DNSRecord struct{
	QType string `json:"qtype"`
	QName string `json:"qname"`
	Content string `json:"content"`
	Ttl int `json:"ttl"`
	DomainId int `json:"domain_id,omitempty"`
	ScopeMask string `json:"scopeMask,omitempty"`
	AuthString string`json:"auth,omitempty"`
}
