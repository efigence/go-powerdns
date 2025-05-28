package ipredir

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var backend, _ = New("")

func TestAdd(t *testing.T) {
	Convey("Add IP redir", t, func() {
		err := backend.AddRedirIp("1.2.3.4", "5.6.7.8")
		So(err, ShouldEqual, nil)
		redirIp, err := backend.ListRedirIp()
		So(err, ShouldEqual, nil)
		So(redirIp["1.2.3.4"], ShouldEqual, "5.6.7.8")
	})
	Convey("Batch set IP", t, func() {
		backend.AddRedirIp("1.2.3.4", "5.6.7.8")
		err := backend.SetRedirIp(map[string]string{
			"2.2.2.2": "2.3.3.3",
			"3.2.2.2": "3.3.3.3",
			"4.2.2.2": "4.3.3.3",
		})
		So(err, ShouldEqual, nil)
		redirIp, err := backend.ListRedirIp()
		So(err, ShouldEqual, nil)
		Convey("Previous value should not exist", func() {
			So(redirIp["1.2.3.4"], ShouldNotEqual, "5.6.7.8")
		})
		Convey("New value should be set", func() {
			So(redirIp["2.2.2.2"], ShouldEqual, "2.3.3.3")
			So(redirIp["3.2.2.2"], ShouldEqual, "3.3.3.3")
			So(redirIp["4.2.2.2"], ShouldEqual, "4.3.3.3")
		})

	})
	Convey("Delete IP", t, func() {
		err := backend.SetRedirIp(map[string]string{
			"5.2.2.2": "5.3.3.3",
			"6.2.2.2": "6.3.3.3",
			"7.2.2.2": "7.3.3.3",
		})
		err = backend.DeleteRedirIp("6.2.2.2")
		So(err, ShouldEqual, nil)
		err = backend.DeleteRedirIp("99.99.99.99")
		So(err, ShouldEqual, nil)
		redirIp, err := backend.ListRedirIp()
		So(err, ShouldEqual, nil)
		Convey("Deleted IP should not exist", func() {
			So(redirIp, ShouldNotContainKey, "99.99.99.99")
			So(redirIp, ShouldNotContainKey, "6.2.2.2")
		})
		Convey("Non-deleted IPs should exist", func() {
			So(redirIp["5.2.2.2"], ShouldEqual, "5.3.3.3")
			So(redirIp["7.2.2.2"], ShouldEqual, "7.3.3.3")
		})

	})

}
