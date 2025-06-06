package memdb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	//	"reflect"
	"github.com/efigence/go-powerdns/api"
	"sort"
)

var testStrings []string

var testRecords = map[string]api.DNSRecord{
	"www": {
		QType:   "A",
		QName:   "www.example.com",
		Content: "1.2.3.2",
		Ttl:     60,
	},
	"www2": {
		QType:   "A",
		QName:   "www.example.com",
		Content: "1.2.3.3",
		Ttl:     60,
	},
	"www3": {
		QType:   "A",
		QName:   "www.example.com",
		Content: "1.2.3.4",
		Ttl:     60,
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

func TestRecordInsert(t *testing.T) {
	backend, err := New("t-data/dns.yaml")
	require.NoError(t, err)

	require.NoError(t, backend.AddRecord(testRecords["wildcard"]))
	require.NoError(t, backend.AddRecord(testRecords["www"]))
	require.NoError(t, backend.AddRecord(testRecords["zone"]))
	q := api.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	require.NoError(t, err)
	assert.Equal(t, api.DNSRecordList{testRecords["www"]}, res)
}

func TestRecordLookup(t *testing.T) {
	backend, _ := New("t-data/dns.yaml")
	require.NoError(t, backend.AddRecord(testRecords["wildcard"]))
	require.NoError(t, backend.AddRecord(testRecords["www"]))
	require.NoError(t, backend.AddRecord(testRecords["www2"]))
	require.NoError(t, backend.AddRecord(testRecords["www3"]))
	require.NoError(t, backend.AddRecord(testRecords["zone"]))
	q := api.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	require.NoError(t, err)
	// ShouldContain craps itself on structs, work around it
	correctOutput := api.DNSRecordList{testRecords["www"], testRecords["www2"], testRecords["www3"]}

	sort.Sort(res)
	sort.Sort(correctOutput)
	assert.Equal(t, correctOutput, res)
}

func TestRecordLookupAny(t *testing.T) {
	backend, _ := New("t-data/dns.yaml")
	backend.AddRecord(testRecords["wildcard"])
	backend.AddRecord(testRecords["www"])
	backend.AddRecord(testRecords["www2"])
	backend.AddRecord(testRecords["www3"])
	backend.AddRecord(testRecords["zone"])

	q := api.QueryLookup{
		QType: "ANY",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	require.NoError(t, err)
	correctOutput := api.DNSRecordList{testRecords["www"], testRecords["www2"], testRecords["www3"]}

	sort.Sort(res)
	sort.Sort(correctOutput)

	assert.Equal(t, correctOutput, res)

}
