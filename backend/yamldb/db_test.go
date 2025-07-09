package yamldb

import (
	"github.com/efigence/go-powerdns/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var testStrings []string

var testRecords = map[string]schema.DNSRecord{
	"www1": {
		QType:   "A",
		QName:   "www.example.com",
		Content: "3.4.5.6",
		Ttl:     3000,
	},
	"www2": {
		QType:   "A",
		QName:   "www.example.com",
		Content: "3.4.5.7",
		Ttl:     3000,
	},
	"zone": {
		QType:   "A",
		QName:   "zone.example.com",
		Content: "1.2.3.5",
		Ttl:     60,
	},
	"wildcard": {
		QType:   "A",
		QName:   "*.example.com",
		Content: "1.2.3.6",
		Ttl:     60,
	},
}

func TestRecordList(t *testing.T) {
	backend, _ := New()
	require.NoError(t, backend.LoadFile("t-data/dns.yaml"))

	q := schema.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	assert.NoError(t, err, "should lookup records")
	assert.Len(t, res, 2)
	assert.NotContains(t, res, testRecords["wildcard"])
	assert.Contains(t, res, testRecords["www1"])
	assert.Contains(t, res, testRecords["www2"])
	list, err := backend.List(schema.QueryList{
		ZoneName: "example.com",
	})
	assert.Len(t, list, 4)
}
