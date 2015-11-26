package webapi
import (
	"github.com/unrolled/render" // or "gopkg.in/unrolled/render.v1")
)

type WebApp struct {
	render *render.Render
}

func New() WebApp {
	var v WebApp
	v.render = render.New(render.Options{
		IndentJSON: true, // FIXME only in debug mode ?
	})
	return v
}
