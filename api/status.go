package api

import (
	"encoding/json"
	"net/http"
)

func (a *apiServer) statusHandler() http.HandlerFunc {

	type Response struct {
		Status       string   `json:"status"`
		WatchedFiles []string `json:"watchedFiles"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		status := Response{
			Status:       "GOOD",
			WatchedFiles: []string{"/etc/fstab", "/etc/fstab"},
		}

		body, err := json.Marshal(status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a.writeResponse(w, body)
	}
}
