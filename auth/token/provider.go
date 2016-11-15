package token

import (
	"io/ioutil"
	"net/http"

	"github.com/mdebrouwer/exchange/auth"
)

type provider struct {
	tokens []string
	store  Store
}

func NewProvider(tokens []string, store Store) auth.Provider {
	return &provider{tokens: tokens, store: store}
}

func (p *provider) has(value string) bool {
	for _, t := range p.tokens {
		if t == value {
			return true
		}
	}
	return false
}

func (p *provider) SessionHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		if p.has(string(token)) {
			w.WriteHeader(http.StatusCreated)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func (p *provider) UserHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (p *provider) Filter(delegate func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return delegate
}
