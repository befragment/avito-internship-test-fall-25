package routing

import (
	"github.com/go-chi/chi/v5"

	pr "avito-intern-test/internal/handler/pullrequest"
)

func RegisterPullRequestRoutes(r chi.Router, h *pr.PullRequestHandler) {
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.CreatePullRequest)
		r.Post("/merge", h.MergePullRequest)
		r.Post("/reassign", h.ReassignPullRequest)
	})
}
