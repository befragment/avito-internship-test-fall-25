package handler

import (
	"avito-intern-test/internal/handler/common"
	"encoding/json"
	"net/http"
)

type PullRequestHandler struct {
	service pullReqeustService
}

func NewPullRequestHandler(service pullReqeustService) *PullRequestHandler {
	return &PullRequestHandler{service: service}
}

func (h *PullRequestHandler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "invalid json body")
	} else if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "missing required fields")
	} else {
		pr, err := h.service.CreatePR(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
		if err != nil {
			if code, msg, ok := common.ParseCodeMessage(err); ok {
				handleCreatePullRequestError(w, code, msg, err)
			} else {
				common.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			common.RespondWithJSON(w, http.StatusCreated, CreatePRResponse{PR: prModelToDTO(*pr)})
		}
	}
}

func (h *PullRequestHandler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "invalid json body")
	} else if req.PullRequestID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "missing required fields")
	} else {
		pr, err := h.service.MergePR(ctx, req.PullRequestID)
		if err != nil {
			if code, msg, ok := common.ParseCodeMessage(err); ok {
				handleMergePullRequestError(w, code, msg, err)
			} else {
				common.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			common.RespondWithJSON(w, http.StatusOK, MergePRResponse{PR: prModelToDTO(*pr)})
		}
	}
}

func (h *PullRequestHandler) ReassignPullRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ReassignPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "invalid json body")
	} else if req.PullRequestID == "" || req.OldReviewerID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "missing required fields")
	} else {
		pr, replacedBy, err := h.service.ReassignReviewer(ctx, req.PullRequestID, req.OldReviewerID)
		if err != nil {
			if code, msg, ok := common.ParseCodeMessage(err); ok {
				handleReassignPullRequestError(w, code, msg, err)
			} else {
				common.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			common.RespondWithJSON(w, http.StatusOK, ReassignPRResponse{
				PR:         prModelToDTO(*pr),
				ReplacedBy: replacedBy,
			})
		}
	}
}
