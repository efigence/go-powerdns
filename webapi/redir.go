package webapi

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

var respOk = map[string]string{
	"response": "ok",
}

func respErr(e error) map[string]string {
	return map[string]string{
		"err": e.Error(),
	}
}

// func (w *WebApp) AddRedir(c web.C, wr http.ResponseWriter, r *http.Request) {
func (w *WebBackend) AddRedir(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	err := w.redirApi.AddRedirIp(from, to)
	if err == nil {
		c.JSON(http.StatusOK, respOk)
	} else {
		c.JSON(http.StatusBadRequest, respErr(err))
	}

}
func (w *WebBackend) BatchAddRedir(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, respErr(err))
		return
	}
	ipList := make(map[string]string)
	err = json.Unmarshal(raw, &ipList)
	if err != nil {
		c.JSON(http.StatusOK, respErr(err))
		return
	}
	err = w.redirApi.SetRedirIp(ipList)
	if err != nil {
		c.JSON(http.StatusOK, respErr(err))
		return
	}
	log.Notice("loaded %d IPs", len(ipList))
	c.JSON(http.StatusOK, respOk)
}
func (w *WebBackend) DeleteRedir(c *gin.Context) {
	from := c.Param("from")
	err := w.redirApi.DeleteRedirIp(from)
	if err == nil {
		c.JSON(http.StatusOK, respOk)
	} else {
		c.JSON(http.StatusOK, respErr(err))
	}
}

func (w *WebBackend) ListRedir(c *gin.Context) {
	list, _ := w.redirApi.ListRedirIp()
	c.JSON(http.StatusOK, list)
}
