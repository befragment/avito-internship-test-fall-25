package handler

import (
	"avito-intern-test/internal/handler/common"
	"net/http"
)

func handleReassignPullRequestError(w http.ResponseWriter, code string, msg string, err error) {
	if err != nil {
		if code, msg, ok := common.ParseCodeMessage(err); ok {
			switch code {
			case "NOT_FOUND":
				common.RespondAPIError(w, http.StatusNotFound, code, msg)
			case "PR_MERGED", "NOT_ASSIGNED", "NO_CANDIDATE":
				common.RespondAPIError(w, http.StatusConflict, code, msg)
			default:
				common.RespondAPIError(w, http.StatusInternalServerError, code, msg)
			}
		}
	}
}

func handleMergePullRequestError(w http.ResponseWriter, code string, msg string, err error) {
	if err != nil {
		if code, msg, ok := common.ParseCodeMessage(err); ok {
			switch code {
			case "NOT_FOUND":
				common.RespondAPIError(w, http.StatusNotFound, code, msg)
			case "PR_MERGED", "NOT_ASSIGNED", "NO_CANDIDATE":
				common.RespondAPIError(w, http.StatusConflict, code, msg)
			default:
				common.RespondAPIError(w, http.StatusInternalServerError, code, msg)
			}
		}
	}
}

func handleCreatePullRequestError(w http.ResponseWriter, code string, msg string, err error) {
	switch code {
	case "PR_EXISTS":
		common.RespondAPIError(w, http.StatusConflict, code, msg)
	case "NOT_FOUND":
		common.RespondAPIError(w, http.StatusNotFound, code, msg)
	default:
		common.RespondAPIError(w, http.StatusInternalServerError, code, msg)
	}
}