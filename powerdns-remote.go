package main

import (
	"embed"
	"github.com/efigence/go-powerdns/backend/ipredir"
	"github.com/efigence/go-powerdns/backend/yamldb"
	"github.com/efigence/go-powerdns/webapi"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
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
	YAMLDir    string
}

func main() {
	app := cli.NewApp()
	app.Name = "powerdns-remote"
	app.Description = "powerdns file backend"
	app.Version = version
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "help, h", Usage: "show help"},
		cli.StringFlag{
			Name:  "listen-addr",
			Value: "127.0.0.1:63636",
			Usage: "HTTP API listen address",
		},
		cli.StringFlag{
			Name:     "yaml-dir",
			Usage:    "path to dir of DNS yamls",
			Required: true,
		},
	}
	app.Action = func(c *cli.Context) error {
		cfg := Config{
			ListenAddr: c.String("listen-addr"),
			YAMLDir:    c.String("yaml-dir"),
		}
		if c.Bool("help") {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		log.Infof("Starting %s version: %s", app.Name, version)
		m, _ := yamldb.New()
		err := m.LoadDir(cfg.YAMLDir)
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
		log.Info("loaded YAML from %s", cfg.YAMLDir)
		r, _ := ipredir.New(m)

		signalChannel := make(chan os.Signal, 2)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1)
		go func() {
			for sig := range signalChannel {
				switch sig {
				case os.Interrupt:
					log.Infof("got interrupt, exiting")
					os.Exit(0)
				case syscall.SIGTERM:
					log.Infof("got SIGTERM, exiting")
					os.Exit(0)
				case syscall.SIGUSR1:
					log.Info("reloading records from file")
					err := m.UpdateDir(cfg.YAMLDir)
					if err != nil {
						log.Errorf("error reloading %s: %s", cfg.YAMLDir, err)
					}
				}
			}
		}()

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
		return nil
	}
	app.Run(os.Args)
}
