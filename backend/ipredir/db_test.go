package ipredir

import (
	. "github.com/smartystreets/goconvey/convey"
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
	Convey("load test data", t, func() {
		So(err, ShouldEqual, nil)
	})
	Convey("Record insert", t, func() {

		So(backend.AddRecord(testRecords["wildcard"]), ShouldEqual, nil)
		So(backend.AddRecord(testRecords["www"]), ShouldEqual, nil)
		So(backend.AddRecord(testRecords["zone"]), ShouldEqual, nil)
		q := api.QueryLookup{
			QType: "A",
			QName: "www.example.com",
		}
		res, err := backend.Lookup(q)
		So(err, ShouldEqual, nil)
		So(res, ShouldResemble, api.DNSRecordList{})
	})
}

func TestRecordLookup(t *testing.T) {
	backend, _ := New("t-data/dns.yaml")
	backend.AddRecord(testRecords["wildcard"])
	backend.AddRecord(testRecords["www"])
	backend.AddRecord(testRecords["www2"])
	backend.AddRecord(testRecords["www3"])
	backend.AddRecord(testRecords["zone"])

	Convey("Lookup", t, func() {
		q := api.QueryLookup{
			QType: "A",
			QName: "www.example.com",
		}
		res, err := backend.Lookup(q)
		So(err, ShouldEqual, nil)
		// ShouldContain craps itself on structs, work around it
		correctOutput := api.DNSRecordList{}

		sort.Sort(res)
		sort.Sort(correctOutput)

		So(res, ShouldResemble, correctOutput)
	})
}

func TestRedir(t *testing.T) {
	backend, _ := New("")
	backend.AddRedirIp("127.0.0.1", "127.0.0.2")

	Convey("Adding IP", t, func() {
		q := api.QueryLookup{
			QType:  "A",
			QName:  "www.example.com",
			Remote: "127.0.0.1",
		}
		res, err := backend.Lookup(q)
		So(err, ShouldEqual, nil)
		So(len(res), ShouldBeGreaterThan, 0)
		So(res[0].Content, ShouldEqual, "127.0.0.2")

		q.QType = "SOA"
		res, err = backend.Lookup(q)
		So(err, ShouldEqual, nil)
		So(len(res), ShouldBeGreaterThan, 0)
		So(res[0].Content, ShouldContainSubstring, "example.com 1 10 10 10 10")

		q.Remote = "127.0.1.1"
		res, err = backend.Lookup(q)
		So(err, ShouldEqual, nil)
		So(len(res), ShouldEqual, 0)

	})
	Convey("Removing IP", t, func() {
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
		So(err, ShouldEqual, nil)
		So(len(res), ShouldEqual, 0)

		q.Remote = "127.1.1.1"
		res, err = backend.Lookup(q)
		So(err, ShouldEqual, nil)
		So(len(res), ShouldBeGreaterThan, 0)
	})

}
