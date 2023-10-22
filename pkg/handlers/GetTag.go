package handlers

import (
	"encoding/json"
	"net/http"
)

func (h handler) GetTag(w http.ResponseWriter, r *http.Request){
    json.NewEncoder(w).Encode(h.registers)
}
