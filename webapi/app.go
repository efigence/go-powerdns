package webapi
import (
	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	"github.com/efigence/go-powerdns/backend/ipredir"
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/api"
)

type WebApp struct {
	render *render.Render
	redirBackend ipredir.DomainBackend
	memBackend  memdb.DomainBackend
}

func New() WebApp {
	var v WebApp
	v.render = render.New(render.Options{
		IndentJSON: true, // FIXME only in debug mode ?
	})
	v.redirBackend,_ = ipredir.New("")
	v.memBackend,_ = memdb.New("")
	v.memBackend.AddDomain(api.DNSDomain{
		Name: "pdns.internal",
		PrimaryNs: "ns1.pdns.internal",
		Owner: "hostmaster.pdns.internal",
	});
	v.memBackend.AddRecord(api.DNSRecord{
		QType: "A",
		QName: "pdns.internal",
		Content: "127.0.0.1",
	})
	return v
}
