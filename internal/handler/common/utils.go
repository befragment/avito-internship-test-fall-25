package common

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RespondWithJSON(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, httpStatus int, message string) {
	RespondWithJSON(w, httpStatus, map[string]string{"error": message})
}

type apiErrorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func RespondAPIError(w http.ResponseWriter, httpStatus int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	var body apiErrorBody
	body.Error.Code = code
	body.Error.Message = message
	_ = json.NewEncoder(w).Encode(body)
}

func ParseCodeMessage(err error) (string, string, bool) {
	if err == nil {
		return "", "", false
	}
	msg := err.Error()
	parts := strings.SplitN(msg, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	code := strings.TrimSpace(parts[0])
	detail := strings.TrimSpace(parts[1])
	if code == "" || detail == "" {
		return "", "", false
	}
	return code, detail, true
}
