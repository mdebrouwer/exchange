package service

type ExchangeServiceConfig struct {
	Url  string `json:"Url"`
	Port int32  `json:"Port"`
}

func NewExchangeServiceConfig(url string, port int32) *ExchangeServiceConfig {
	cfg := new(ExchangeServiceConfig)
	cfg.Url = url
	cfg.Port = port
	return cfg
}
