package main

import (
	"github.com/efigence/go-powerdns/api"
	"github.com/efigence/go-powerdns/webapi"
	//	"github.com/efigence/go-powerdns/backend/yamldb"
	"github.com/op/go-logging"
	"github.com/zenazn/goji"
	//	"github.com/zenazn/goji/web"
	"flag"
	"os"
	//	"strings"
)

var version string
var log = logging.MustGetLogger("main")
var stdout_log_format = logging.MustStringFormatter("%{color:bold}%{time:2006-01-02T15:04:05.9999Z-07:00}%{color:reset}%{color} [%{level:.1s}] %{color:reset}%{shortpkg}[%{longfunc}] %{message}")

type Config struct {
	ListenAddr string
	YAMLDB     string
}

func main() {
	var cfg Config
	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	stderrFormatter := logging.NewBackendFormatter(stderrBackend, stdout_log_format)
	logging.SetBackend(stderrFormatter)
	logging.SetFormatter(stdout_log_format)

	if cfg.ListenAddr == "" {
		cfg.ListenAddr = "127.0.0.1:63636"
	}
	flag.Set("bind", cfg.ListenAddr)

	log.Info("Starting app")
	log.Debug("version: %s", version)

	webApp := webapi.New()
	goji.Get("/dns", webApp.Dns)
	goji.Post("/dns", webApp.Dns)
	goji.Get("/isItWorking", webApp.Healthcheck)

	_, _ = api.New(api.CallbackList{})

	goji.Serve()
}
