package auth

import (
	"net/http"
)

type Provider interface {
	SessionHandler() func(w http.ResponseWriter, r *http.Request)
	UserHandler() func(w http.ResponseWriter, r *http.Request)
	Filter(delegate func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)
}
