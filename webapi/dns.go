package webapi
import (
//	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	"net/http"
	"github.com/zenazn/goji/web"
	"github.com/efigence/go-powerdns/api"
	"bytes"
	"github.com/op/go-logging"
	"strings"
)

var log = logging.MustGetLogger("main")


func (w *WebApp) Dns(c web.C, wr http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var record api.DNSRecord
	response := api.NewResponse()
	s := buf.String()
	if strings.Contains(s,`SOA`) {
		log.Warning("SOA")
		record.QType = "SOA"
		record.QName = "example.com"
		record.Content = "ns1.example.com hostmaster.example.com 2014000000 3000 3000 3000 3000"
		record.Ttl = 61
	}	else {
		record.QType = "A"
		record.QName = "www.example.com"
		record.Content = "5.6.7.8"
		record.Ttl = 62
	}

	records := make([]api.DNSRecord,0)
	records = append(records, record)
	response.Result = records
	w.render.JSON(wr, http.StatusOK, response)
}
