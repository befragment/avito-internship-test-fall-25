package routing

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"avito-intern-test/internal/handler/common"
)

func CommonRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/healthcheck", handler.Healthcheck)
	return r
}

func PublicRoutes() *chi.Mux {
	r := chi.NewRouter()
	return r
}

func AppRoutes(routers ...*chi.Mux) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	for _, router := range routers {
		r.Mount("/", router)
	}
	return r
}
