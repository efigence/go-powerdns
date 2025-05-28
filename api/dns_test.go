package api

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestExpandDNSName(t *testing.T) {
	d := "very.long.domain.name.com"
	splittedDomain, err := ExpandDNSName(d)
	Convey("Domain dissection", t, func() {
		So(err, ShouldEqual, nil)
		So(fmt.Sprintf("records: %d", len(splittedDomain)), ShouldEqual, "records: 5")
		So(splittedDomain[0], ShouldEqual, "very.long.domain.name.com")
		So(splittedDomain[1], ShouldEqual, "long.domain.name.com")
		So(splittedDomain[2], ShouldEqual, "domain.name.com")
		So(splittedDomain[3], ShouldEqual, "name.com")
		So(splittedDomain[4], ShouldEqual, "com") // we TLD now
	})
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
	Convey("Generate SOA record", t, func() {
		soaRecord := GenerateSoaFromDomain(d)
		So(soaRecord, ShouldResemble, res)
	})
}

func BenchmarkExpandDNSName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ExpandDNSName(`some.simple.dns.name`)
	}
}
