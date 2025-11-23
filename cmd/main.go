package main

import (
	"log"

	"avito-intern-test/internal/core"
	prh "avito-intern-test/internal/handler/pullrequest"
	th "avito-intern-test/internal/handler/team"
	uh "avito-intern-test/internal/handler/user"
	prrepo "avito-intern-test/internal/repository/pullrequest"
	teamrepo "avito-intern-test/internal/repository/team"
	userrepo "avito-intern-test/internal/repository/user"
	"avito-intern-test/internal/routing"
	prsvc "avito-intern-test/internal/service/pullrequest"
	teamsvc "avito-intern-test/internal/service/team"
	usersvc "avito-intern-test/internal/service/user"
)

func main() {
	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbPool := core.MustInitPool()

	teamRepo := teamrepo.NewTeamRepository(dbPool)
	userRepo := userrepo.NewUserRepository(dbPool)
	pullRequestRepo := prrepo.NewPullRequestRepository(dbPool)

	core.StartServer(
		dbPool,
		cfg.Port,
		routing.Router(
			prh.NewPullRequestHandler(prsvc.NewPRService(
				userRepo,
				teamRepo,
				pullRequestRepo,
			)),
			th.NewTeamHandler(teamsvc.NewTeamService(
				teamRepo,
				userRepo,
			)),
			uh.NewUserHandler(usersvc.NewUserService(
				userRepo,
				pullRequestRepo,
			)),
		),
	)
}
