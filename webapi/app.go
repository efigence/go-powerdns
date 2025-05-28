package webapi

import (
	"github.com/efigence/go-powerdns/api"
	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
	//	"github.com/zenazn/goji/web"
	"fmt"
)

type WebApp struct {
	render     *render.Render
	dnsBackend dnsCB
	dnsApi     api.Api
}

func New() WebApp {
	var v WebApp
	v.render = render.New(render.Options{
		//		IndentJSON: true, // FIXME only in debug mode ?
	})
	b, err := newDNSBackend()
	if err != nil {
		log.Error("error creating DNS backend: %+v", err)
	}
	v.dnsBackend = b

	cbList := api.CallbackList{
		Lookup: b,
		List:   b,
	}
	dnsApi, err := api.New(cbList)
	if err != nil {
		log.Error("error creating DNS backend: %+v", err)
	}
	v.dnsApi = dnsApi

	return v
}

func respErr(err error) map[string]interface{} {
	v := make(map[string]interface{})
	v["response"] = false
	v["error"] = fmt.Sprintf("%s", err)
	return v
}
