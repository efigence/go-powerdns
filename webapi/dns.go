package webapi
import (
//	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	"net/http"
	"github.com/zenazn/goji/web"

)


func (w *WebApp) Dns(c web.C, wr http.ResponseWriter, r *http.Request) {
	w.render.Text(wr, http.StatusOK, "Plain text here")
}
