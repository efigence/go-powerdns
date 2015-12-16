package webapi

import (
	"net/http"
	"github.com/zenazn/goji/web"
)

var	respOk = map[string]string {
		"response": "ok",
	}

func (w *WebApp) AddRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
	from := c.URLParams["from"]
	to := c.URLParams["to"]
	err := w.dnsBackend.redirBackend.AddRedirIp(from,to)
	if err == nil {
		w.render.JSON(wr, http.StatusOK, respOk)
	} else {
		respErr := map[string]string {
			"response":"fail",
		}
		w.render.JSON(wr, http.StatusOK, respErr)
	}


}
func (w *WebApp) BatchAddRedir(c web.C, wr http.ResponseWriter, r *http.Request) {}
func (w *WebApp) DeleteRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
	from := c.URLParams["from"]
	err := w.dnsBackend.redirBackend.DeleteRedirIp(from)
		if err == nil {
		w.render.JSON(wr, http.StatusOK, respOk)
	} else {
		respErr := map[string]string {
			"response":"fail",
		}
		w.render.JSON(wr, http.StatusOK, respErr)
	}
}
