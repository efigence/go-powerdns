package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
		Name:     "example.com",
		Serial:   12345,
		NS:       []string{"ns1.example.com"},
		Owner:    "hostmaster.example.com",
		Refresh:  86400,
		Retry:    300,
		Expiry:   864000,
		Nxdomain: 3600,
	}
	res := DNSRecord{
		QType:   "SOA",
		QName:   "example.com",
		Content: "ns1.example.com hostmaster.example.com 12345 86400 300 864000 3600",
		Ttl:     3600,
	}
	soaRecord := GenerateSoaFromDomain(d)
	assert.Equal(t, res, soaRecord)
	s1 := d.Serial
	d.UpdateSerial()
	s2 := d.Serial
	assert.Greater(t, s2, s1)
}

func BenchmarkExpandDNSName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Millisecond * 40)
		_, _ = ExpandDNSName(`some.simple.dns.name`)
	}
}

func TestDNSDomain_Validate(t *testing.T) {
	d := DNSDomain{
		Name:     "example.com",
		NS:       []string{"ns1.example.com"},
		Owner:    "hostmaster.example.com",
		Serial:   12,
		Refresh:  23,
		Retry:    34,
		Expiry:   45,
		Nxdomain: 56,
	}
	assert.NoError(t, d.Validate())
	v := d
	v.Name = "example.com."
	assert.Error(t, v.Validate())
	v.Name = ""
	assert.Error(t, v.Validate())
	v.Name = "example.com"
	assert.NoError(t, v.Validate())
	v = d
	v.NS = []string{}
	assert.Error(t, v.Validate())
	v = d
	v.Owner = ""
	assert.Error(t, v.Validate())
	v = d
	v.Serial = 0
	assert.Error(t, v.Validate())
	v = d
	v.Refresh = 0
	assert.Error(t, v.Validate())
	v = d
	v.Retry = 0
	assert.Error(t, v.Validate())
	v = d
	v.Expiry = 0
	assert.Error(t, v.Validate())
	v = d
	v.Nxdomain = 0
	assert.Error(t, v.Validate())

}
