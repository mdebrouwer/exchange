package service

import (
	"fmt"
	"net/http"

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
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	s.logger.Printf("Listening on port: %v\n", s.cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", s.cfg.Port), mux)
}

func (s *ExchangeService) orderHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Println("Handling Order Request")
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (s *ExchangeService) Stop() {

}
