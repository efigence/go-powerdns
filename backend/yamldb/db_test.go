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
	"potato1": {
		QType:   "A",
		QName:   "www.potato.com",
		Content: "5.4.5.6",
		Ttl:     61,
	},
}

func TestRecordList(t *testing.T) {
	backend, _ := New()
	require.NoError(t, backend.LoadFile("../../t-data/dns.yaml"))

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

func TestYAMLDB_LoadDir(t *testing.T) {
	t.Run("No dup domains", func(t *testing.T) {
		backend, _ := New()
		require.Error(t, backend.LoadDir("../../t-data/dupdomain"))

	})
	t.Run("No yaml dir fail", func(t *testing.T) {
		backend, _ := New()
		require.Error(t, backend.LoadDir("."))

	})
	t.Run("load subdirectory", func(t *testing.T) {

		backend, _ := New()

		require.NoError(t, backend.LoadDir("../../t-data/dns"))
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
		q = schema.QueryLookup{
			QType: "A",
			QName: "www.potato.com",
		}
		res, err = backend.Lookup(q)
		assert.NoError(t, err, "should lookup records")
		assert.Len(t, res, 2)
		assert.NotContains(t, res, testRecords["wildcard"])
		assert.Contains(t, res, testRecords["potato1"])
		assert.NotContains(t, res, testRecords["www2"])
	})
}
