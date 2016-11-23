package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"

	"github.com/mdebrouwer/exchange/auth/token"
	"github.com/mdebrouwer/exchange/command"
	"github.com/mdebrouwer/exchange/log"
	"github.com/mdebrouwer/exchange/orderbook"
)

type ExchangeService struct {
	logger         log.Logger
	cfg            ExchangeServiceConfig
	commandHandler command.Handler
}

func NewExchangeService(logger log.Logger, cfg ExchangeServiceConfig) *ExchangeService {
	i := orderbook.NewInstrument("MATDEB_500", orderbook.TickSize(5))
	ob := orderbook.NewOrderbook(logger, i)
	s := new(ExchangeService)
	s.logger = logger
	s.cfg = cfg
	s.commandHandler = command.NewHandler(ob)
	return s
}

func (s *ExchangeService) Start() error {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "bundle"}))

	r.PathPrefix("/assets/").Handler(http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "bundle"}))

	tokens, err := s.cfg.AuthConfig.Tokens()
	if err != nil {
		return err
	}

	authStore, err := token.NewBoltBackedStore(s.cfg.AuthConfig.BoltPath)
	if err != nil {
		return err
	}

	auth := token.NewProvider(
		s.cfg.AuthConfig.CookieSigningKey,
		s.cfg.AuthConfig.CookieEncryptionKey,
		tokens,
		authStore,
	)

	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/orders", auth.Filter(OrderHandler(s.commandHandler))).Methods("POST")
	api.HandleFunc("/user", auth.UserHandler()).Methods("GET", "POST")
	api.HandleFunc("/sessions", auth.SessionHandler()).Methods("POST")

	bindAddress := fmt.Sprintf("%s:%v", s.cfg.Host, s.cfg.Port)
	s.logger.Printf("Listening on: %s\n", bindAddress)
	return http.ListenAndServe(bindAddress, r)
}

func (s *ExchangeService) Stop() {
	s.commandHandler.Stop()
}

type Message struct {
	Order string
}

func OrderHandler(h command.Handler) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		var message Message
		err := json.NewDecoder(req.Body).Decode(&message)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		h.Handle(command.NewOrderCommand{
			CounterParty: message.Order,
			Side:         orderbook.BUY,
			Price:        10,
			Volume:       10,
		})
	}
}
