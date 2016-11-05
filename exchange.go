package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	svc "github.com/mdebrouwer/exchange/service"
	"github.com/mdebrouwer/glog"
)

func main() {
	var cfg = getConfig()
	var log = glog.NewLogger(cfg.Logfile)

	log.Infof("Configuration: %+v", cfg)

	log.Info("Service starting...")
	var s = svc.NewExchangeService(log, cfg.ExchangeServiceConfig)
	s.Start()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Infof("Signal received [%v].\n", sig)
		s.Stop()
		done <- true
	}()

	<-done

	log.Info("Service exiting.")
	os.Exit(0)
}

func getConfig() *Config {
	program := os.Args[0]
	_, file := filepath.Split(program)
	name := strings.TrimSuffix(file, filepath.Ext(file))
	configfile := fmt.Sprintf("%v.json", name)

	cfg := flag.String("c", configfile, "config file location")
	flag.Parse()

	return NewConfigFromFile(cfg)
}
