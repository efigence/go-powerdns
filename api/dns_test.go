package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandDNSName(t *testing.T) {
	d := "very.long.domain.name.com"
	splittedDomain, err := ExpandDNSName(d)
	assert.NoError(t, err)
	assert.Len(t, splittedDomain, 5)
	assert.Equal(t, "very.long.domain.name.com", splittedDomain[0])
	assert.Equal(t, "long.domain.name.com", splittedDomain[1])
	assert.Equal(t, "domain.name.com", splittedDomain[2])
	assert.Equal(t, "name.com", splittedDomain[3])
	assert.Equal(t, "com", splittedDomain[4])
}

func TestSoaFromDomain(t *testing.T) {
	d := DNSDomain{
		Name:      "example.com",
		Serial:    12345,
		PrimaryNs: "ns1.example.com",
		Owner:     "hostmaster.example.com",
		Refresh:   86400,
		Retry:     300,
		Expiry:    864000,
		Nxdomain:  3600,
	}
	res := DNSRecord{
		QType:   "SOA",
		QName:   "example.com",
		Content: "ns1.example.com hostmaster.example.com 12345 86400 300 864000 3600",
		Ttl:     3600,
	}
	soaRecord := GenerateSoaFromDomain(d)
	assert.Equal(t, res, soaRecord)
}

func BenchmarkExpandDNSName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ExpandDNSName(`some.simple.dns.name`)
	}
}
