package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func (h Handler) GetRegisters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "application/json")
		// TODO: Replace this h.reg
		err := json.NewEncoder(w).Encode(h.registers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h Handler) GetRegister(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Path
	addressStr := strings.TrimPrefix(request, "/register/")
	address, err := strconv.Atoi(addressStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	response, err := h.db.GetRowByAddress(address)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		slog.Debug("GET request for /register/<ADDRESS>", "address", address, "response", response)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "PUT":
		query := r.URL.Query()
		value := query.Get("value")

		if value == "" {
			value = query.Get("val")
		}
		if value == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		slog.Debug("PUT request for /register/<ADDRESS>", "address", address, "value", fValue)
		err = h.db.SetAddressValue(address, fValue)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
