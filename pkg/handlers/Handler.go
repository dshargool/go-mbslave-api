package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dshargool/go-mbslave-api.git/pkg/types"
	"github.com/simonvetter/modbus"
)

type Handler struct {
	registers          map[types.InstrumentTag]types.ModbusTag
	db                 *types.SqlDb
	MbSlave            *modbus.ModbusServer
	AllowNullRegisters bool
}

func New(config types.Configuration, db *types.SqlDb) Handler {
	return Handler{
		registers:          config.Registers,
		db:                 db,
		MbSlave:            nil,
		AllowNullRegisters: config.AllowNullRegister,
	}
}

func (h Handler) HandleRequests(port int) {
	http.HandleFunc("/all_registers", h.GetRegisters)
	http.HandleFunc("/tag/", h.GetTag)
	http.HandleFunc("/register/", h.GetRegister)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}
