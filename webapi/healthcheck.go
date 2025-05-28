package webapi

import (
	//	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	"github.com/zenazn/goji/web"
	"net/http"
)

func (w *WebApp) Healthcheck(c web.C, wr http.ResponseWriter, r *http.Request) {
	w.render.Text(wr, http.StatusOK, "I am working\n")
}
