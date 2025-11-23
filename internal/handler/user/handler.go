package handler

import (
	"encoding/json"
	"net/http"

	"avito-intern-test/internal/core"
	"avito-intern-test/internal/handler/common"
)

type UserHandler struct {
	service userService
}

func NewUserHandler(service userService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "invalid json body")
	} else {
		user, err := h.service.SetIsActive(ctx, req.UserID, req.IsActive)
		if err != nil {
			common.RespondAPIError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		} else {
			common.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"user": userToDTO(user),
			})
		}
	}
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "user_id is required")
	} else {
		prs, err := h.service.GetReviewerPRs(ctx, userID)
		if err != nil {
			if code, msg, ok := common.ParseCodeMessage(err); ok && code == core.ErrorNotFound {
				common.RespondAPIError(w, http.StatusNotFound, code, msg)
			} else {
				common.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			items := make([]PullRequestShortDTO, 0, len(prs))
			for _, p := range prs {
				items = append(items, PullRequestShortDTO{
					PullRequestID:   p.PullRequestID,
					PullRequestName: p.PullRequestName,
					AuthorID:        p.AuthorID,
					Status:          string(p.Status),
				})
			}
			resp := GetReviewResponse{
				UserID:       userID,
				PullRequests: items,
			}
			common.RespondWithJSON(w, http.StatusOK, resp)
		}
	}
}
