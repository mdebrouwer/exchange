package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/mdebrouwer/exchange/log"
	svc "github.com/mdebrouwer/exchange/service"
)

func main() {
	var cfg = getConfig()
	var logger = log.NewLogger(cfg.Logfile)

	logger.Printf("Configuration: %+v\n", cfg)

	logger.Println("Service starting...")
	var s = svc.NewExchangeService(logger, cfg.ExchangeServiceConfig)
	s.Start()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Printf("Signal received [%v].\n", sig)
		s.Stop()
		done <- true
	}()

	<-done

	logger.Println("Service exiting.")
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
