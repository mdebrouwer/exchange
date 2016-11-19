package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mdebrouwer/exchange/service"
)

type Config struct {
	Logfile               string                        `json:"Logfile"`
	ExchangeServiceConfig service.ExchangeServiceConfig `json:"ExchangeServiceConfig"`
}

func NewConfigFromFile(configFile string) (cfg Config) {
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deserializing file: %v", err)
		os.Exit(1)
	}
	return
}

func NewConfig(logfile string, exchangeConfig service.ExchangeServiceConfig) Config {
	return Config{
		Logfile:               logfile,
		ExchangeServiceConfig: exchangeConfig,
	}
}
