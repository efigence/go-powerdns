package memdb

import (
	"github.com/efigence/go-powerdns/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestMemDomains_AddDomain(t *testing.T) {
	backend := New()
	require.NoError(t, backend.AddDomain(schema.DNSDomain{
		Name: "www.example2.com",
		NS:   []string{"ns1.example.com"},
	}))
	assert.NotEmpty(t, backend.Domains["www.example2.com"].Owner)
	assert.Greater(t, backend.Domains["www.example2.com"].Refresh, int32(0))
	assert.Greater(t, backend.Domains["www.example2.com"].Retry, int32(0))
	assert.Greater(t, backend.Domains["www.example2.com"].Expiry, int32(0))
	assert.Greater(t, backend.Domains["www.example2.com"].Nxdomain, int32(0))
}

func TestRecordInsert(t *testing.T) {
	backend := New()
	require.NoError(t, backend.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	}))
	require.NoError(t, backend.AddRecord(testRecords["wildcard"]))
	require.NoError(t, backend.AddRecord(testRecords["www"]))
	require.NoError(t, backend.AddRecord(testRecords["zone"]))
	q := schema.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	require.NoError(t, err)
	assert.Equal(t, []schema.DNSRecord{testRecords["www"]}, res)
	list, err := backend.List(schema.QueryList{
		ZoneName: "example.com",
	})
	assert.Len(t, list, 4)
}

func TestRecordLookup(t *testing.T) {
	backend := New()
	require.NoError(t, backend.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	}))
	require.NoError(t, backend.AddRecord(testRecords["wildcard"]))
	require.NoError(t, backend.AddRecord(testRecords["www"]))
	require.NoError(t, backend.AddRecord(testRecords["www2"]))
	require.NoError(t, backend.AddRecord(testRecords["www3"]))
	require.NoError(t, backend.AddRecord(testRecords["zone"]))
	q := schema.QueryLookup{
		QType: "A",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	require.NoError(t, err)
	// ShouldContain craps itself on structs, work around it
	correctOutput := []schema.DNSRecord{testRecords["www"], testRecords["www2"], testRecords["www3"]}

	slices.SortFunc(res, cmpFunc)
	slices.SortFunc(correctOutput, cmpFunc)
	assert.Equal(t, correctOutput, res)
	q = schema.QueryLookup{
		QType: "A",
		QName: "potato.example.com",
	}
	res, err = backend.Lookup(q)
	require.NoError(t, err)
	tr := testRecords["wildcard"]
	tr.QName = "potato.example.com"
	correctOutput = []schema.DNSRecord{tr}
	slices.SortFunc(res, cmpFunc)
	slices.SortFunc(correctOutput, cmpFunc)
	assert.Equal(t, correctOutput, res)
	q = schema.QueryLookup{
		QType: "ANY",
		QName: "tomato.example.com",
	}
	res, err = backend.Lookup(q)
	require.NoError(t, err)
	tr = testRecords["wildcard"]
	tr.QName = "tomato.example.com"
	correctOutput = []schema.DNSRecord{tr}
	slices.SortFunc(res, cmpFunc)
	slices.SortFunc(correctOutput, cmpFunc)
	assert.Equal(t, correctOutput, res)
}

func TestRecordLookupAny(t *testing.T) {
	backend := New()
	require.NoError(t, backend.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	}))
	backend.AddRecord(testRecords["wildcard"])
	backend.AddRecord(testRecords["www"])
	backend.AddRecord(testRecords["www2"])
	backend.AddRecord(testRecords["www3"])
	backend.AddRecord(testRecords["zone"])

	q := schema.QueryLookup{
		QType: "ANY",
		QName: "www.example.com",
	}
	res, err := backend.Lookup(q)
	require.NoError(t, err)
	correctOutput := []schema.DNSRecord{testRecords["www"], testRecords["www2"], testRecords["www3"]}

	slices.SortFunc(res, cmpFunc)
	slices.SortFunc(correctOutput, cmpFunc)
	assert.Equal(t, correctOutput, res)
}

func BenchmarkMemDomains_Lookup(b *testing.B) {
	backend := New()
	backend.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	})
	backend.AddRecord(testRecords["wildcard"])
	backend.AddRecord(testRecords["www"])
	backend.AddRecord(testRecords["www2"])
	backend.AddRecord(testRecords["www3"])
	backend.AddRecord(testRecords["zone"])
	b.Run("A", func(b *testing.B) {
		q := schema.QueryLookup{
			QType: "A",
			QName: "www.example.com",
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			backend.Lookup(q)
		}
	})
	b.Run("A wildcard", func(b *testing.B) {
		q := schema.QueryLookup{
			QType: "A",
			QName: "potato.example.com",
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			backend.Lookup(q)
		}
	})
	b.Run("ANY", func(b *testing.B) {
		q := schema.QueryLookup{
			QType: "ANY",
			QName: "www.example.com",
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			backend.Lookup(q)
		}
	})
	b.Run("ANY wildcard", func(b *testing.B) {
		q := schema.QueryLookup{
			QType: "A",
			QName: "potato.example.com",
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			backend.Lookup(q)
		}
	})
}

func TestMemDomains_GetRootDomainFor(t *testing.T) {
	backend := New()
	require.NoError(t, backend.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	}))
	dom, err := backend.GetRootDomainFor("very.long.test.www.example.com")
	assert.NoError(t, err)
	assert.Equal(t, "example.com", dom)

	require.NoError(t, backend.AddDomain(schema.DNSDomain{
		Name: "www.example.com",
		NS:   []string{"ns1.example.com"},
	}))
	dom, err = backend.GetRootDomainFor("very.long.test.www.example.com")
	assert.NoError(t, err)
	assert.Equal(t, "www.example.com", dom)

	require.NoError(t, backend.AddDomain(schema.DNSDomain{
		Name: "very.long.test.www.example.com",
		NS:   []string{"ns1.example.com"},
	}))
	dom, err = backend.GetRootDomainFor("very.long.test.www.example.com")
	assert.NoError(t, err)
	assert.Equal(t, "very.long.test.www.example.com", dom)

	dom, err = backend.GetRootDomainFor("very.long.test.www.example2.com")
	assert.Error(t, err)

}
