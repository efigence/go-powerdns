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
	response := api.ResponseOk()
	s := buf.String()
	log.Warning(s)
	if strings.Contains(s,`SOA`) && strings.Contains(s,`example.com`) {
		record.QType = "SOA"
		record.QName = "example.com"
		record.Content = "ns1.example.com hostmaster.example.com 2014000000 3000 3000 3000 3000"
		record.Ttl = 61
		records := make([]api.DNSRecord,0)
		records = append(records, record)
		response.Result = records
		log.Warning("responding with SOA")
	}	else {
		if strings.Contains(s,`www.example.com`) {
			record.QType = "A"
			record.QName = "www.example.com"
			record.Content = "5.6.7.8"
			record.Ttl = 5
			records := make([]api.DNSRecord,0)
			records = append(records, record)
			response.Result = records
			log.Warning("responding with A")

		} else {
			response = api.ResponseOk()
		}

	}


	w.render.JSON(wr, http.StatusOK, response)
}
