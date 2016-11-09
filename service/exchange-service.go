package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"

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
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "bundle"}))

	r.PathPrefix("/assets/").Handler(http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "bundle"}))

	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/orders", OrderHandler).Methods("POST")

	s.logger.Printf("Listening on port: %v\n", s.cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), r)
}

func (s *ExchangeService) Stop() {}

type Message struct {
	Order string
}

func OrderHandler(resp http.ResponseWriter, req *http.Request) {
	var message Message
	err := json.NewDecoder(req.Body).Decode(&message)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	fmt.Println(message.Order)
}
