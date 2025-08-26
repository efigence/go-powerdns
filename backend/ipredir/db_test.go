package ipredir

import (
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"slices"
	"testing"
)

var testStrings []string

var testRecords = map[string]schema.DNSRecord{
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
var cmpFunc = func(a, b schema.DNSRecord) int {
	x := a.QName + a.Content + a.QType
	y := a.QName + a.Content + a.QType
	if x < y {
		return -1
	}
	if x > y {
		return 1
	}
	return 0
}

func TestRecordInsert(t *testing.T) {
	mdb := memdb.New(zaptest.NewLogger(t).Sugar())
	backend, err := New(mdb)
	require.NoError(t, err)
	mdb.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	})
	assert.NoError(t, backend.AddRecord(testRecords["wildcard"]))
	assert.NoError(t, backend.AddRecord(testRecords["www"]))
	assert.NoError(t, backend.AddRecord(testRecords["zone"]))
	q := schema.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)

	assert.NoError(t, err)
	assert.Equal(t, res, []schema.DNSRecord{})
	dom, err := backend.GetRootDomainFor("zone.example.com")
	assert.Equal(t, "example.com", dom)
}

func TestRecordLookup(t *testing.T) {
	backend, err := New(memdb.New(zaptest.NewLogger(t).Sugar()))
	require.NoError(t, err)
	assert.NoError(t, backend.AddRecord(testRecords["wildcard"]))
	assert.NoError(t, backend.AddRecord(testRecords["www"]))
	assert.NoError(t, backend.AddRecord(testRecords["www2"]))
	assert.NoError(t, backend.AddRecord(testRecords["www3"]))
	assert.NoError(t, backend.AddRecord(testRecords["zone"]))

	q := schema.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	require.NoError(t, err)
	// ShouldContain craps itself on structs, work around it
	correctOutput := []schema.DNSRecord{}

	slices.SortFunc(res, cmpFunc)
	slices.SortFunc(correctOutput, cmpFunc)
	assert.Equal(t, correctOutput, res)
}

func BenchmarkRecordLookup(b *testing.B) {
	backend, _ := New(memdb.New(zaptest.NewLogger(b).Sugar()))
	backend.AddRecord(testRecords["wildcard"])
	backend.AddRecord(testRecords["www"])
	backend.AddRecord(testRecords["www2"])
	backend.AddRecord(testRecords["www3"])
	backend.AddRecord(testRecords["zone"])
	q := schema.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		backend.Lookup(q)
	}

}

func TestRedir(t *testing.T) {
	backend, _ := New(memdb.New(zaptest.NewLogger(t).Sugar()))
	backend.AddRedirIp("127.0.0.1", "127.0.0.2")

	t.Run("Adding IP", func(t *testing.T) {
		q := schema.QueryLookup{
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
		q := schema.QueryLookup{
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
