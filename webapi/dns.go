package webapi
import (
//	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	"net/http"
	"github.com/zenazn/goji/web"
	"github.com/efigence/go-powerdns/api"
	"bytes"
	"strings"
)






func (w *WebApp) Dns(c web.C, wr http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var record api.DNSRecord
	response := api.NewResponse()
	if strings.Contains(buf.String(),`SOA`) {
		record.QType = "SOA"
		record.QName = "example.com"
		record.Content = "ns1.example.com hostmaster.example.com 3000 3000 3000 3000"
		record.Ttl = 60
	} else {
		record.QType = "A"
		record.QName = "www.example.com"
		record.Content = "1.2.3.4"
		record.Ttl = 60
	}

	records := make([]api.DNSRecord,0)
	records =append(records, record)
	response.Result = records
	w.render.JSON(wr, http.StatusOK, response)
}
