package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func (h Handler) GetRegisters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(h.registers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h Handler) GetRegister(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Path
	register := strings.TrimPrefix(request, "/register/")
	// query := r.URL.Query()
	switch r.Method {
	case "GET":
		w.Header().Add("Content-Type", "application/json")
		response, err := h.db.GetRowByRegister(string(register))
		if err == sql.ErrNoRows {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Debug("GET request for /register/<ADDRESS>", "address", register, "response", response)
		response.Value = response.Value / response.Divisor
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
