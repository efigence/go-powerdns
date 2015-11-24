package api

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"fmt"
)

var testStrings []string

func TestLookup(t *testing.T) {
	fmt.Printf("")
	lookup := `{"method":"lookup", "parameters":{"qtype":"ANY", "qname":"www.example.com", "remote":"192.0.2.24", "local":"192.0.2.1", "real-remote":"192.0.2.2", "zone-id":-1}}`
	var l qLookup;
	cb_list := CallbackList{
		Lookup: l,
	}
	Convey("Create new API", t, func() {
		_, err := New(CallbackList{})
		So(err,ShouldEqual,nil)
	})
	Convey("Lookup", t, func() {
		api, _ := New(cb_list)
		out, err := api.ParseRequest(lookup);
		So(err,ShouldEqual,nil)
		So(out.(QueryLookup).Remote,ShouldEqual,"192.0.2.24")
		So(out.(QueryLookup).RealRemote,ShouldEqual,"192.0.2.2")
		So(out.(QueryLookup).Local,ShouldEqual,"192.0.2.1")
		So(reflect.TypeOf(out).Name(), ShouldEqual,"QueryLookup")
		So(out.Dump(), ShouldNotEqual,nil)
	})

}

type qLookup struct {}

func (qLookup) Query(q QueryLookup) (QueryResponse, error) {
	var err error
	res := NewResponse()
	return res, err
}
