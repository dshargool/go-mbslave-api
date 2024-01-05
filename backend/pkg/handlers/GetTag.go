package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func (h Handler) GetTag(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Path
	tag := strings.TrimPrefix(request, "/tag/")
	query := r.URL.Query()

	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "application/json")
		response, err := h.db.GetRowByTag(tag)
		if err == sql.ErrNoRows {
			slog.Warn("Could not find row in database", "error", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			slog.Warn("Database error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Debug("GET request for /tag/<TAG>", "tag", tag, "response", response)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "PUT":
		value := query.Get("value")
		if value == "" {
			value = query.Get("val")
		}
		if value != "" {
			fValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				slog.Error("Could not parse request value as float", err)
				w.WriteHeader(http.StatusBadRequest)
				return

			}

			slog.Info("Updating tag " + tag + " with value " +
				strconv.FormatFloat(fValue, 'f', -1, 64))
			err = h.db.SetTagValue(tag, fValue)
			if err != nil {
				slog.Error("Could not set tag value", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "GET, PUT")
	}
}
