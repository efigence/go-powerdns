package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	"reflect"
	"fmt"
)

var testStrings []string

var queries = map[string]string{
	"lookup":     `{"method":"lookup", "parameters":{"qtype":"ANY", "qname":"www.example.com", "remote":"192.0.2.24", "local":"192.0.2.1", "real-remote":"192.0.2.2", "zone-id":-1}}`,
	"list":       `{"method":"list", "parameters":{"zonename":"example.com","domain_id":-1}}`,
	"initialize": `{"method":"initialize", "parameters":{"command":"/path/to/something", "timeout":"2000", "something":"else"}}`,
	"badreq":     `{"asd":123}`,
}

func TestQuery(t *testing.T) {
	fmt.Printf("")
	var qLookup testQLookup
	var qList testQList
	cbList := CallbackList{
		Lookup: qLookup,
		List:   qList,
	}
	t.Run("Create new API", func(t *testing.T) {
		_, err := New(CallbackList{})
		assert.NoError(t, err)
	})
	t.Run("Init", func(t *testing.T) {
		api, _ := New(cbList)
		out, err := api.Parse(queries["initialize"])
		assert.NoError(t, err)
		assert.Equal(t, out, ResponseOk())
	})
	t.Run("Lookup", func(t *testing.T) {
		api, _ := New(cbList)
		out, err := api.Parse(queries["lookup"])
		testQueryOutput, _ := qLookup.Lookup(QueryLookup{})
		assert.NoError(t, err)
		assert.Equal(t, testQueryOutput, out)
	})
	t.Run("List", func(t *testing.T) {
		api, _ := New(cbList)
		out, err := api.Parse(queries["list"])
		testQueryOutput, _ := qList.List(QueryList{})
		assert.NoError(t, err)
		assert.Equal(t, testQueryOutput, out)
	})
	t.Run("BadReq", func(t *testing.T) {
		api, _ := New(cbList)
		out, err := api.Parse(queries["badreq"])
		assert.NoError(t, err)
		assert.Equal(t, ResponseFailed(), out)
	})
}

type testQLookup struct{}

func (testQLookup) Lookup(q QueryLookup) (QueryResponse, error) {
	var err error
	res := NewResponse()
	res.Result = []DNSRecord{
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

func (testQList) List(q QueryList) (QueryResponse, error) {
	var err error
	res := NewResponse()
	res.Result = []DNSRecord{
		{
			QType:   "A",
			QName:   "www.example.com",
			Content: "1.2.3.4",
			Ttl:     60,
		},
		{
			QType:   "MX",
			QName:   "10 example.com",
			Content: "mx1.example.com",
			Ttl:     60,
		},
		{
			QType:   "A",
			QName:   "mx.example.com",
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
