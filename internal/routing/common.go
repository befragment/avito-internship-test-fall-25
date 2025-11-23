package routing

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterCommonRoutes(r chi.Router, h http.HandlerFunc) {
	r.Get("/healthcheck", h)
}
