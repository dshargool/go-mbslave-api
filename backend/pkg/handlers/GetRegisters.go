package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/dshargool/go-mbslave-api.git/pkg/types"
)

func (h Handler) GetRegisters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var registers []types.ModbusResponse
		for addr := range h.registers {
			val, err := h.db.GetRowByTag(string(addr))
			if err != nil {
				slog.Error(err.Error())
				reg := h.registers[addr]
				empty_val := types.ModbusResponse{
					Tag:         string(reg.Tag),
					Description: reg.Description,
					Address:     reg.Address,
					DataType:    reg.DataType,
					Value:       -1.0,
					LastUpdate:  "",
				}
				registers = append(registers, empty_val)
			} else {
				registers = append(registers, val)
			}
		}
        jRegister, err := json.Marshal(registers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/plain-text")
        w.Header().Add("Access-Control-Allow-Origin","*")
        w.Write(jRegister)
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
		slog.Warn("Could not get row by address; row not found", err)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		slog.Warn("Could not get row by address", err)
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
