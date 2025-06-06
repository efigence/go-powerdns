package ipredir

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
	assert.NoError(t, backend.AddRecord(testRecords["wildcard"]))
	assert.NoError(t, backend.AddRecord(testRecords["www"]))
	assert.NoError(t, backend.AddRecord(testRecords["zone"]))
	q := api.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	assert.NoError(t, err)
	assert.Equal(t, res, api.DNSRecordList{})
}

func TestRecordLookup(t *testing.T) {
	backend, err := New("t-data/dns.yaml")
	require.NoError(t, err)
	assert.NoError(t, backend.AddRecord(testRecords["wildcard"]))
	assert.NoError(t, backend.AddRecord(testRecords["www"]))
	assert.NoError(t, backend.AddRecord(testRecords["www2"]))
	assert.NoError(t, backend.AddRecord(testRecords["www3"]))
	assert.NoError(t, backend.AddRecord(testRecords["zone"]))

	q := api.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	require.NoError(t, err)
	// ShouldContain craps itself on structs, work around it
	correctOutput := api.DNSRecordList{}

	sort.Sort(res)
	sort.Sort(correctOutput)

	assert.Equal(t, correctOutput, res)
}

func BenchmarkRecordLookup(b *testing.B) {
	backend, _ := New("t-data/dns.yaml")
	backend.AddRecord(testRecords["wildcard"])
	backend.AddRecord(testRecords["www"])
	backend.AddRecord(testRecords["www2"])
	backend.AddRecord(testRecords["www3"])
	backend.AddRecord(testRecords["zone"])
	q := api.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backend.Lookup(q)
	}

}

func TestRedir(t *testing.T) {
	backend, _ := New("")
	backend.AddRedirIp("127.0.0.1", "127.0.0.2")

	t.Run("Adding IP", func(t *testing.T) {
		q := api.QueryLookup{
			QType:  "A",
			QName:  "www.example.com",
			Remote: "127.0.0.1",
		}
		res, err := backend.Lookup(q)
		require.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "127.0.0.2", res[0].Content)

		q.QType = "SOA"
		res, err = backend.Lookup(q)
		require.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Contains(t, res[0].Content, "example.com 1 10 10 10 10")

		q.Remote = "127.0.1.1"
		res, err = backend.Lookup(q)
		require.NoError(t, err)
		assert.Len(t, res, 0)
	})
	t.Run("Removing IP", func(t *testing.T) {
		backend.AddRedirIp("127.1.1.1", "127.0.0.2")
		backend.AddRedirIp("127.1.2.1", "127.0.0.2")
		backend.AddRedirIp("127.1.3.1", "127.0.0.2")
		backend.DeleteRedirIp("127.1.2.1")
		q := api.QueryLookup{
			QType:  "A",
			QName:  "www.example.com",
			Remote: "127.1.2.1",
		}
		res, err := backend.Lookup(q)
		require.NoError(t, err)
		assert.Len(t, res, 0)

		q.Remote = "127.1.1.1"
		res, err = backend.Lookup(q)
		require.NoError(t, err)
		assert.Len(t, res, 1)
	})

}
