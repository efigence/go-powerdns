package webapi

import (
	"net/http"
	"github.com/zenazn/goji/web"
)

func (w *WebApp) AddRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
	from := c.URLParams["from"]
	to := c.URLParams["to"]
	err := w.dnsBackend.redirBackend.AddRedirIp(from,to)
	_ = err
}
func (w *WebApp) BatchAddRedir(c web.C, wr http.ResponseWriter, r *http.Request) {}
func (w *WebApp) DeleteRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
	from := c.URLParams["from"]
	err := w.dnsBackend.redirBackend.DeleteRedirIp(from)
	_ = err

}
