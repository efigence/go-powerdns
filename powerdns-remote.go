package main

import (
	"embed"
	"flag"
	"github.com/efigence/go-powerdns/backend/ipredir"
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/schema"
	"github.com/efigence/go-powerdns/webapi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var version string

var log *zap.SugaredLogger
var debug = true

//go:embed static templates
var embeddedWebContent embed.FS

func init() {
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	// naive systemd detection. Drop timestamp if running under it
	if os.Getenv("JOURNAL_STREAM") != "" {
		consoleEncoderConfig.TimeKey = ""
	}
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return (lvl < zapcore.ErrorLevel) != (lvl == zapcore.DebugLevel && !debug)
	})
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, os.Stderr, lowPriority),
		zapcore.NewCore(consoleEncoder, os.Stderr, highPriority),
	)
	logger := zap.New(core)
	if debug {
		logger = logger.WithOptions(
			zap.Development(),
			zap.AddCaller(),
			zap.AddStacktrace(highPriority),
		)
	} else {
		logger = logger.WithOptions(
			zap.AddCaller(),
		)
	}
	log = logger.Sugar()

}

type Config struct {
	ListenAddr string
	YAMLDB     string
}

func main() {
	var cfg Config

	if cfg.ListenAddr == "" {
		cfg.ListenAddr = "127.0.0.1:63636"
	}
	flag.Set("bind", cfg.ListenAddr)
	log.Info("Starting app")
	log.Debug("version: %s", version)
	m := memdb.New()
	m.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	})
	m.AddRecord(schema.DNSRecord{
		QType:   "A",
		QName:   "example.com",
		Content: "1.2.3.4",
		Ttl:     3600,
	})
	r, _ := ipredir.New(m)
	w, err := webapi.New(webapi.Config{
		Logger:       log.Named("web"),
		AccessLogger: log.Named("access"),
		ListenAddr:   cfg.ListenAddr,
		DNSBackend:   m,
		RedirBackend: r,
	}, embeddedWebContent)

	if err != nil {
		log.Fatalf("error setting up: %s", err)
	}
	log.Fatalf("error starting up: %s", w.Run())
}
