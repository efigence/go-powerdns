package memdb

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"reflect"
	"github.com/efigence/go-powerdns/api"
	"sort"
)

var testStrings []string


var testRecords = map[string]api.DNSRecord{
	"www": {
			QType:   "A",
			QName:   "www.example.com",
			Content: "1.2.3.4",
			Ttl:     60,
	},
	"zone":{
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
	backend,err := New("t-data/dns.yaml")
	Convey("load test data", t, func() {
		So(err,ShouldEqual,nil)
	})
	Convey("Record insert", t, func() {

		So(backend.AddRecord("example.com", testRecords["wildcard"]), ShouldEqual, nil)
		So(backend.AddRecord("example.com", testRecords["www"]), ShouldEqual, nil)
		So(backend.AddRecord("example.com", testRecords["zone"]), ShouldEqual, nil)
		q := api.QueryLookup{
			QType: "A",
			QName: "www.example.com",
		}
		res, err := backend.Search(q)
		So(err,ShouldEqual,nil)
		So(res,ShouldResemble,api.DNSRecordList{testRecords["www"]})
	})
}

func TestRecordList(t *testing.T) {
	backend,_ := New("t-data/dns.yaml")
	Convey("Record insert", t, func() {
		So(backend.AddRecord("example.com", testRecords["wildcard"]), ShouldEqual, nil)
		So(backend.AddRecord("example.com", testRecords["www"]), ShouldEqual, nil)
		So(backend.AddRecord("example.com", testRecords["zone"]), ShouldEqual, nil)
		q := api.QueryLookup{
			QType: "A",
			QName: "www.example.com",
		}
		res, err := backend.Search(q)
		So(err,ShouldEqual,nil)
		// ShouldContain craps itself on structs, work around it
		correctOutput := api.DNSRecordList{ testRecords["wildcard"],testRecords["www"],testRecords["zone"] }

		sort.Sort(res)
		sort.Sort(correctOutput)


		errmap := []bool{false, false, false}
		for idx, val := range res {
			if reflect.DeepEqual(testRecords["wildcard"],val) ||reflect.DeepEqual(testRecords["www"],val) || reflect.DeepEqual(testRecords["zone"],val)  {
					errmap[idx]=true
			}
		}

		So(errmap,ShouldResemble,[]bool{true, true, true})


		So(res,ShouldEqual,testRecords["zone"])
	})
}
