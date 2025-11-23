package handler

import (
	"avito-intern-test/internal/core"
	"avito-intern-test/internal/handler/common"
	teamerr "avito-intern-test/internal/service/team"
	"encoding/json"
	"errors"
	"net/http"
)

type TeamHandler struct {
	service teamService
}

func NewTeamHandler(service teamService) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateWithMembersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "invalid json body")
	}
	team, err := h.service.CreateWithMembers(ctx, req.TeamName, req.Members)
	if err != nil {
		if errors.Is(err, teamerr.ErrTeamAlreadyExists) {
			common.RespondAPIError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
		}
		if code, msg, ok := common.ParseCodeMessage(err); ok {
			if code == core.ErrorUserExists {
				common.RespondAPIError(w, http.StatusConflict, code, msg)
			}
		}
		common.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	members := make([]TeamMemberDTO, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, TeamMemberDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}
	resp := CreateTeamResponse{
		Team: TeamDTO{
			TeamName: team.Name,
			Members:  members,
		},
	}
	common.RespondWithJSON(w, http.StatusCreated, resp)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		common.RespondWithError(w, http.StatusBadRequest, "team_name is required")
	}
	members, err := h.service.GetTeamMembers(ctx, teamName)

	if errors.Is(err, teamerr.ErrTeamNotFound) {
		common.RespondAPIError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
	}

	items := make([]TeamMemberDTO, 0, len(members))
	for _, m := range members {
		items = append(items, TeamMemberDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}
	resp := TeamDTO{
		TeamName: teamName,
		Members:  items,
	}
	common.RespondWithJSON(w, http.StatusOK, resp)
}
