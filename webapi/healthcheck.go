package webapi
import (
//	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	"net/http"
	"github.com/zenazn/goji/web"
)


func (w *WebApp) Healthcheck(c web.C, wr http.ResponseWriter, r *http.Request) {
     w.render.Text(wr, http.StatusOK, "I am working\n")
}
