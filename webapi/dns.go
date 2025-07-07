package webapi

import (
	"github.com/gin-gonic/gin"

	"bytes"
	"github.com/op/go-logging"
	//	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	"net/http"
)

var log = logging.MustGetLogger("main")

func (w *WebBackend) Dns(c *gin.Context) {
	// FIXME
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	s := buf.String()
	log.Warning(s)
	resp, err := w.dnsApi.Parse(s)
	if err != nil {
		log.Error("failure on responding query: %+v", err)
	}

	c.JSON(http.StatusOK, resp)
}
