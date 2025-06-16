package yamldb

import (
	"github.com/efigence/go-powerdns/api"
	"github.com/efigence/go-powerdns/backend/schema"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

var testStrings []string

var testRecords = map[string]schema.DNSRecord{
	"www": {
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
	assert.NoError(t, err, "should load test data")

	err = backend.AddRecord(testRecords["wildcard"])
	assert.NoError(t, err, "should add wildcard record")

	err = backend.AddRecord(testRecords["www"])
	assert.NoError(t, err, "should add www record")

	err = backend.AddRecord(testRecords["zone"])
	assert.NoError(t, err, "should add zone record")

	q := api.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	assert.NoError(t, err, "should lookup www record")
	assert.Equal(t, schema.DNSRecordList{testRecords["www"]}, res, "lookup should return the www record")
}

func TestRecordList(t *testing.T) {
	backend, _ := New("t-data/dns.yaml")

	err := backend.AddRecord(testRecords["wildcard"])
	assert.NoError(t, err, "should add wildcard record")

	err = backend.AddRecord(testRecords["www"])
	assert.NoError(t, err, "should add www record")

	err = backend.AddRecord(testRecords["zone"])
	assert.NoError(t, err, "should add zone record")

	q := api.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	assert.NoError(t, err, "should lookup records")

	correctOutput := schema.DNSRecordList{testRecords["wildcard"], testRecords["www"], testRecords["zone"]}

	sort.Sort(res)
	sort.Sort(correctOutput)

	errmap := []bool{false, false, false}
	for idx, val := range res {
		if reflect.DeepEqual(testRecords["wildcard"], val) ||
			reflect.DeepEqual(testRecords["www"], val) ||
			reflect.DeepEqual(testRecords["zone"], val) {
			errmap[idx] = true
		}
	}
	assert.Equal(t, []bool{true, true, true}, errmap, "all returned records should match one of the test records")

	// This assertion may not make sense, as res is a list, not a single record, but kept for parity
	assert.Contains(t, res, testRecords["zone"], "result should contain the zone record")
}
