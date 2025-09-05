package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/deepissue/core/logging"
	"github.com/deepissue/core/option"
	"github.com/deepissue/core/server"
	"github.com/deepissue/fee_server/config"
	"github.com/deepissue/fee_server/services"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	opts := option.NewOptions()
	if err := opts.Parse(); err != nil {
		return
	}
	initialize(opts)
	logger, err := logging.NewLogger(opts.Application, &opts.Log)
	if err != nil {
		log.Fatal(err)
		return
	}

	srv, err := server.NewServer(opts, logger)
	if err != nil {
		log.Fatal(err)
		return
	}
	cfg := config.LoadConfig(opts.ConfigFile)
	logrus.Debugf("Loaded config: %v", cfg)
	db, err := xorm.NewEngine(cfg.Xorm.Driver, cfg.Xorm.Datasource[0])
	if err != nil {
		log.Fatal(err)
		return
	}

	feeService, err := services.NewFeeService(srv, db, cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	feeService.Start()
	srv.HandleSignal(func() {
		feeService.Stop()
	})

}

func initialize(opts *option.Options) {
	level, err := logrus.ParseLevel(opts.Log.Level)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true,
		PadLevelText:    true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "\t", fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	})
	writer, err := rotatelogs.New(
		path.Join(opts.Log.Path, "app.log"),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*90),
	)
	if err != nil {
		logrus.Fatalf("Failed to create rotatelogs: %v", err)
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, writer))
}
