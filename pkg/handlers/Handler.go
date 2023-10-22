package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dshargool/go-modbus-api.git/pkg/types"
)

type handler struct {
    devices []types.ModbusSlave
    registers []types.ModbusTag
}

func New(config types.Configuration) handler {
    return handler{
        devices: config.Slaves,
        registers: config.Registers,
    }
}

func (h handler) HandleRequests(port int) {
    http.HandleFunc("/all_devices", h.GetDevice)
    http.HandleFunc("/all_tags", h.GetTag)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}
