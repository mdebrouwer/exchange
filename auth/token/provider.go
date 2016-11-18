package token

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"

	"github.com/mdebrouwer/exchange/auth"
)

type provider struct {
	tokens  []string
	store   Store
	cookies *sessions.CookieStore
}

func NewProvider(signingKey, encryptionKey []byte, tokens []string, store Store) auth.Provider {
	return &provider{tokens: tokens, store: store, cookies: sessions.NewCookieStore(signingKey, encryptionKey)}
}

func (p *provider) has(value []byte) bool {
	str := string(value)
	for _, t := range p.tokens {
		if t == str {
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
		if p.has(token) {
			session, err := p.cookies.Get(r, "session")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			session.Values["token"] = token
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func (p *provider) getUserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := p.cookies.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	val := session.Values["token"]
	token, ok := val.([]byte)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !p.has(token) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	user, ok, err := p.store.Get(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *provider) postUserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := p.cookies.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	val := session.Values["token"]
	token, ok := val.([]byte)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !p.has(token) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ok, err = p.store.Register(token, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *provider) UserHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			p.getUserHandler(w, r)
		case "POST":
			p.postUserHandler(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func (p *provider) Filter(delegate func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := p.cookies.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		val := session.Values["token"]
		token, ok := val.([]byte)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !p.has(token) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, ok, err := p.store.Get(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		context.Set(r, "user", user)
		defer context.Clear(r)

		delegate(w, r)
	}
}
