package common

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RespondWithJSON(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
	body.Reset()
	body.Error.Code = code
	body.Error.Message = message
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *apiErrorBody) Reset() {
	b.Error.Code = ""
	b.Error.Message = ""
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
