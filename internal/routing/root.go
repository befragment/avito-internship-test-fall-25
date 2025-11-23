package routing

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	common "avito-intern-test/internal/handler/common"
	prh "avito-intern-test/internal/handler/pullrequest"
	th "avito-intern-test/internal/handler/team"
	uh "avito-intern-test/internal/handler/user"
)

func Router(
	prHandler *prh.PullRequestHandler,
	teamHandler *th.TeamHandler,
	userHandler *uh.UserHandler,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	RegisterCommonRoutes(r, common.Healthcheck)
	RegisterPullRequestRoutes(r, prHandler)
	RegisterTeamRoutes(r, teamHandler)
	RegisterUserRoutes(r, userHandler)
	return r
}
