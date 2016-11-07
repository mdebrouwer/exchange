package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mdebrouwer/exchange/service"
)

type Config struct {
	Logfile               string                         `json:"Logfile"`
	ExchangeServiceConfig *service.ExchangeServiceConfig `json:"ExchangeServiceConfig"`
}

func NewConfigFromFile(configFile *string) *Config {
	raw, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v", err)
		os.Exit(1)
	}

	cfg := new(Config)
	json.Unmarshal(raw, cfg)
	return cfg
}

func NewConfig(logfile string, exchangeConfig *service.ExchangeServiceConfig) *Config {
	cfg := new(Config)
	cfg.Logfile = logfile
	cfg.ExchangeServiceConfig = exchangeConfig
	return cfg
}
