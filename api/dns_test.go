package api
import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
)


func TestExpandDNSName(t *testing.T) {
	d := "very.long.domain.name.com"
	splittedDomain, err := ExpandDNSName(d)
	Convey("Domain dissection",t,func() {
		So(err,ShouldEqual,nil)
		So(fmt.Sprintf("records: %d",len(splittedDomain)),ShouldEqual,"records: 5")
		So(splittedDomain[0],ShouldEqual,"very.long.domain.name.com")
		So(splittedDomain[1],ShouldEqual,"long.domain.name.com")
		So(splittedDomain[2],ShouldEqual,"domain.name.com")
		So(splittedDomain[3],ShouldEqual,"name.com")
		So(splittedDomain[4],ShouldEqual,"com") // we TLD now
	})
}
func BenchmarkExpandDNSName(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _,_ = ExpandDNSName(`some.simple.dns.name`)
    }
}
