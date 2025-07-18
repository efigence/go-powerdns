package api

import (
	"encoding/json"
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/schema"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"regexp"
	"testing"
	//	"reflect"
	"fmt"
)

var testStrings []string

var queries = map[string]string{
	"lookup":               `{"method":"lookup", "parameters":{"qtype":"ANY", "qname":"www.example.com", "remote":"192.0.2.24", "local":"192.0.2.1", "real-remote":"192.0.2.2", "zone-id":-1}}`,
	"list":                 `{"method":"list", "parameters":{"zonename":"example.com","domain_id":-1}}`,
	"initialize":           `{"method":"initialize", "parameters":{"command":"/path/to/something", "timeout":"2000", "something":"else"}}`,
	"getAllDomains":        `{"method": "getAllDomains", "parameters": {"include_disabled": true}}`,
	"getAllDomainMetadata": `{"method":"getalldomainmetadata", "parameters":{"name":"example.com"}}`,
	"badreq":               `{"asd":123}`,
	"getUpdatedMaster": `{"method": "getUpdatedMasters", "parameters": {}}
`,
}

func testQBackend() schema.DomainReader {
	m := memdb.New()
	m.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	})
	m.AddRecord(schema.DNSRecord{
		QType:   "A",
		QName:   "www.example.com",
		Content: "1.2.3.4",
		Ttl:     60,
	})
	m.AddRecord(schema.DNSRecord{
		QType:   "MX",
		QName:   "example.com",
		Content: "10 mx1.example.com",
		Ttl:     60,
	})
	m.AddRecord(schema.DNSRecord{
		QType:   "A",
		QName:   "mx1.example.com",
		Content: "5.6.7.8",
		Ttl:     60,
	})
	m.AddRecord(schema.DNSRecord{
		QType:   "TXT",
		QName:   "example.com",
		Content: "a record",
		Ttl:     60,
	})

	return m
}

var qLookup = testQLookup{}
var qList = testQList{}
var qDomain = testQDomain{}

func TestQuery(t *testing.T) {
	fmt.Printf("")
	t.Run("Init", func(t *testing.T) {
		api, _ := New(testQBackend(), zaptest.NewLogger(t).Sugar())
		out, err := api.Parse(queries["initialize"])
		assert.NoError(t, err)
		assert.Equal(t, out, schema.ResponseOk())
	})
	t.Run("Lookup", func(t *testing.T) {
		api, _ := New(testQBackend(), zaptest.NewLogger(t).Sugar())
		out, err := api.Parse(queries["lookup"])
		testQueryOutput, _ := qLookup.Lookup(schema.QueryLookup{})
		assert.NoError(t, err)
		outj, _ := json.MarshalIndent(out, "", " ")
		testj, _ := json.MarshalIndent(testQueryOutput, "", " ")
		assert.Equal(t, string(testj), string(outj))
	})
	t.Run("List", func(t *testing.T) {
		api, _ := New(testQBackend(), zaptest.NewLogger(t).Sugar())
		out, err := api.Parse(queries["list"])
		testQueryOutput, _ := qList.List(schema.QueryList{})
		assert.NoError(t, err)
		outj, _ := json.MarshalIndent(out, "", " ")
		testj, _ := json.MarshalIndent(testQueryOutput, "", " ")
		assert.Equal(t, string(testj), string(outj))

	})
	t.Run("domainList", func(t *testing.T) {
		api, _ := New(testQBackend(), zaptest.NewLogger(t).Sugar())
		out, err := api.Parse(queries["getAllDomains"])
		testQueryOutput, _ := qDomain.ListDomains(schema.QueryLookup{})
		assert.NoError(t, err)
		outj, _ := json.MarshalIndent(out, "", " ")
		testj, _ := json.MarshalIndent(testQueryOutput, "", " ")
		re := regexp.MustCompile(`(?m)^.*last_check.*\n?`)
		outjs := re.ReplaceAllString(string(outj), "")
		testjs := re.ReplaceAllString(string(testj), "")
		assert.Equal(t, testjs, outjs)
	})
	t.Run("domainMetadata", func(t *testing.T) {
		api, _ := New(testQBackend(), zaptest.NewLogger(t).Sugar())
		out, err := api.Parse(queries["getAllDomainMetadata"])
		testQueryOutput := schema.QueryResponse{
			Result: map[string]string{},
		}
		assert.NoError(t, err)
		outj, _ := json.MarshalIndent(out, "", " ")
		testj, _ := json.MarshalIndent(testQueryOutput, "", " ")
		assert.Equal(t, string(testj), string(outj))
	})
	t.Run("BadReq", func(t *testing.T) {
		api, _ := New(testQBackend(), zaptest.NewLogger(t).Sugar())
		out, err := api.Parse(queries["badreq"])
		assert.Error(t, err)
		assert.Equal(t, schema.ResponseFailed(), out)
	})
	t.Run("GetUpdatedMaster", func(t *testing.T) {
		api, _ := New(testQBackend(), zaptest.NewLogger(t).Sugar())
		out, err := api.Parse(queries["getUpdatedMaster"])
		assert.Error(t, err)
		assert.Equal(t, schema.ResponseFailed(), out)
	})
}

type testQLookup struct{}

func (testQLookup) Lookup(q schema.QueryLookup) (schema.QueryResponse, error) {
	var err error
	res := schema.NewResponse()
	res.Result = []schema.DNSRecord{
		{
			QType:   "A",
			QName:   "www.example.com",
			Content: "1.2.3.4",
			Ttl:     60,
		},
	}
	return res, err
}

type testQList struct{}

func (testQList) List(q schema.QueryList) (schema.QueryResponse, error) {
	var err error
	res := schema.NewResponse()
	res.Result = []schema.DNSRecord{
		{
			QType:   "SOA",
			QName:   "example.com",
			Content: "ns1.example.com hostmaster.example.com 0 172800 900 1209600 1800",
			Ttl:     1800,
		},
		{
			QType:   "A",
			QName:   "www.example.com",
			Content: "1.2.3.4",
			Ttl:     60,
		},
		{
			QType:   "MX",
			QName:   "example.com",
			Content: "10 mx1.example.com",
			Ttl:     60,
		},
		{
			QType:   "A",
			QName:   "mx1.example.com",
			Content: "5.6.7.8",
			Ttl:     60,
		},
		{
			QType:   "TXT",
			QName:   "example.com",
			Content: "a record",
			Ttl:     60,
		},
	}
	return res, err
}

type testQDomain struct{}

func (testQDomain) ListDomains(q schema.QueryLookup) (schema.QueryResponse, error) {
	var err error
	res := schema.NewResponse()
	res.Result = []schema.PDNSDomain{{
		ID:   0,
		Zone: "example.com.",
		Kind: "native",
	}}
	return res, err
}
func (testQDomain) ListDomainMetadata(q schema.QueryLookup) (schema.QueryResponse, error) {
	var err error
	res := schema.NewResponse()
	res.Result = map[string]string{}
	return res, err
}
