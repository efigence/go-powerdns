package webapi

import (
	"encoding/json"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"net/http"
)

var respOk = map[string]string{
	"response": "ok",
}

func (w *WebApp) AddRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
	from := c.URLParams["from"]
	to := c.URLParams["to"]
	err := w.dnsBackend.redirBackend.AddRedirIp(from, to)
	if err == nil {
		w.render.JSON(wr, http.StatusOK, respOk)
	} else {
		w.render.JSON(wr, http.StatusOK, respErr(err))
	}

}
func (w *WebApp) BatchAddRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.render.JSON(wr, http.StatusOK, respErr(err))
		return
	}
	ipList := make(map[string]string)
	err = json.Unmarshal(raw, &ipList)
	if err != nil {
		w.render.JSON(wr, http.StatusOK, respErr(err))
		return
	}
	err = w.dnsBackend.redirBackend.SetRedirIp(ipList)
	if err != nil {
		w.render.JSON(wr, http.StatusOK, respErr(err))
		return
	}
	log.Notice("loaded %d IPs", len(ipList))
	w.render.JSON(wr, http.StatusOK, respOk)
}
func (w *WebApp) DeleteRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
	from := c.URLParams["from"]
	err := w.dnsBackend.redirBackend.DeleteRedirIp(from)
	if err == nil {
		w.render.JSON(wr, http.StatusOK, respOk)
	} else {
		w.render.JSON(wr, http.StatusOK, respErr(err))
	}
}

func (w *WebApp) ListRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
	list, _ := w.dnsBackend.redirBackend.ListRedirIp()

	w.render.JSON(wr, http.StatusOK, list)
}
