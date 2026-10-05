package httpx

import (
	"net/http"
)

type Routes interface {
	Register(mux *http.ServeMux)
}

func NewMux(routes []Routes, responder Responder) http.Handler {
	mux := http.NewServeMux()
	for _, r := range routes {
		r.Register(mux)
	}
	return responder.recoverPanics(responder.logRequests(mux))
}
