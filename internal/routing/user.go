package routing

import (
	"github.com/go-chi/chi/v5"

	u "avito-intern-test/internal/handler/user"
)

func RegisterUserRoutes(r chi.Router, h *u.UserHandler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.SetIsActive)
		r.Get("/getReview", h.GetReview)
	})
}
