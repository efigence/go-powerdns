package webapi
import (
//	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	"net/http"
	"github.com/zenazn/goji/web"
	"bytes"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")


func (w *WebApp) Dns(c web.C, wr http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String()
	log.Warning(s)
	resp,err := w.dnsApi.Parse(s)
	if (err != nil) {
		log.Error("failure on responding query: %+v", err)
	}

	w.render.JSON(wr, http.StatusOK, resp)
}
