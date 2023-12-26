package handlers

import (
	"log/slog"
	"net/http"
)

func (h Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		slog.Error("Healthcheck received non-get request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var row string
	err := h.db.QueryRow("SELECT tag FROM datapoints;").Scan(&row)
	if err != nil {
		slog.Error("Unable to open database table", "error", err.Error())
		w.WriteHeader(http.StatusFailedDependency)
		return
	}
	w.Header().Add("Content-Type", "application/plain-text")
	w.Header().Add("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    slog.Info("Healthcheck")
}
