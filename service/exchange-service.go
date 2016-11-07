package service

import (
	"fmt"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"

	"github.com/mdebrouwer/exchange/log"
)

type ExchangeService struct {
	logger log.Logger
	cfg    *ExchangeServiceConfig
}

func NewExchangeService(logger log.Logger, cfg *ExchangeServiceConfig) *ExchangeService {
	s := new(ExchangeService)
	s.logger = logger
	s.cfg = cfg
	return s
}

func (s *ExchangeService) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/order/", s.orderHandler)
	mux.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "static"}))

	s.logger.Printf("Listening on port: %v\n", s.cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), mux)
}

func (s *ExchangeService) orderHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Println("Handling Order Request")
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (s *ExchangeService) Stop() {

}
