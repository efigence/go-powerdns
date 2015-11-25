package api

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
//	"reflect"
	"fmt"
)

var testStrings []string

var	queries = map[string]string {
	"lookup": `{"method":"lookup", "parameters":{"qtype":"ANY", "qname":"www.example.com", "remote":"192.0.2.24", "local":"192.0.2.1", "real-remote":"192.0.2.2", "zone-id":-1}}`,
	"list": `{"method":"list", "parameters":{"zonename":"example.com","domain_id":-1}}`,
	"initialize": `{"method":"initialize", "parameters":{"command":"/path/to/something", "timeout":"2000", "something":"else"}}`,
	"badreq" : `{"asd":123}`,
}




func TestQuery(t *testing.T) {
	fmt.Printf("")
	var qLookup testQLookup;
	var qList testQList;
	cbList := CallbackList{
		Lookup: qLookup,
		List: qList,
	}
	Convey("Create new API", t, func() {
		_, err := New(CallbackList{})
		So(err,ShouldEqual,nil)
	})
	Convey("Init", t, func() {
		api, _ := New(cbList)
		out, err := api.Parse(queries["initialize"]);
		So(err,ShouldEqual,nil)
		So(out,ShouldResemble,ResponseOk())
	})
	Convey("Lookup", t, func() {
		api, _ := New(cbList)
		out, err := api.Parse(queries["lookup"]);
		testQueryOutput, _ := qLookup.Query(QueryLookup{})
		So(err,ShouldEqual,nil)
		So(out,ShouldResemble,testQueryOutput)
	})
	Convey("List", t, func() {
		api, _ := New(cbList)
		out, err := api.Parse(queries["list"]);
		testQueryOutput, _ := qList.Query(QueryList{})
		So(err,ShouldEqual,nil)
		So(out,ShouldResemble,testQueryOutput)
	})
	Convey("BadReq", t, func() {
		api, _ := New(cbList)
		out, err := api.Parse(queries["badreq"]);
		So(err,ShouldEqual,nil)
		So(out,ShouldResemble,ResponseFailed())
	})
}


type testQLookup struct {}

func (testQLookup) Query(q QueryLookup) (QueryResponse, error) {
	var err error
	res := NewResponse()
	res.Result = []DNSRecord{
		{
			QType: "A",
			QName: "www.example.com",
			Content: "1.2.3.4",
			Ttl: 60,
		},
	}
	return res, err
}


type testQList struct {}

func (testQList) Query(q QueryList) (QueryResponse, error) {
	var err error
	res := NewResponse()
	res.Result = []DNSRecord{
		{
			QType: "A",
			QName: "www.example.com",
			Content: "1.2.3.4",
			Ttl: 60,
		},
		{
			QType: "MX",
			QName: "10 example.com",
			Content: "mx1.example.com",
			Ttl: 60,
		},
		{
			QType: "A",
			QName: "mx.example.com",
			Content: "5.6.7.8",
			Ttl: 60,
		},
		{
			QType: "TXT",
			QName: "example.com",
			Content: "a record",
			Ttl: 60,
		},
	}
	return res, err
}
