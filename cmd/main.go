package main

import (
	"golang_test_task/internal/httphandlers"

	formatter "github.com/fabienm/go-logrus-formatters"
	log "github.com/sirupsen/logrus"

	"github.com/alexflint/go-arg"
)

type AppCfg struct {
	LogLevel int  `arg:"--log-level,env:LOG_LEVEL" default:"4" help:"0-panic, 1-fatal, 2-error, 3-warn, 4-info, 5-debug, 6-trace"`
	LogGELF  bool `arg:"--log-gelf,env:LOG_GELF" default:"false" help:"Enable of disable GELF format of logs"`
}

func main() {
	var args struct {
		AppCfg
		httphandlers.RouterCfg
	}
	arg.MustParse(&args)

	if args.AppCfg.LogGELF {
		gelfFmt := formatter.NewGelf("Golang Test Task")
		log.SetFormatter(gelfFmt)
	}
	log.SetLevel(log.Level(args.LogLevel))

	log.Info("Controller is starting...")

	log.Infof("Start arguments: %+v", args)

	log.Print("HTTP server is starting...")
	router := httphandlers.NewRouter(args.RouterCfg)
	if err := router.Run(); err != nil {
		log.WithError(err).Fatal("Can`t start HTTP server")
	}
}
