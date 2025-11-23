package routing

import (
	"github.com/go-chi/chi/v5"

	t "avito-intern-test/internal/handler/team"
)

func RegisterTeamRoutes(r chi.Router, h *t.TeamHandler) {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.CreateTeam)
		r.Get("/get", h.GetTeam)
	})
}
