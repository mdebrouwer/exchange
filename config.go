package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	svc "github.com/mdebrouwer/exchange/service"
)

type Config struct {
	Logfile               string                     `json:"Logfile"`
	ExchangeServiceConfig *svc.ExchangeServiceConfig `json:"ExchangeServiceConfig"`
}

func NewConfigFromFile(configFile *string) *Config {
	raw, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cfg := new(Config)
	json.Unmarshal(raw, cfg)
	return cfg
}

func NewConfig(logfile string, exchangeConfig *svc.ExchangeServiceConfig) *Config {
	cfg := new(Config)
	cfg.Logfile = logfile
	cfg.ExchangeServiceConfig = exchangeConfig
	return cfg
}
