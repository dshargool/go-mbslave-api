package handlers

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/simonvetter/modbus"
)

func (h Handler) MbInit(port int) *modbus.ModbusServer {
	mbServer, err := modbus.NewServer(&modbus.ServerConfiguration{
		URL:        "tcp://0.0.0.0:" + strconv.Itoa(port),
		Timeout:    10 * time.Second,
		MaxClients: 5,
	}, &h)
	if err != nil {
		slog.Error("Unable to initialize modbus slave: " + err.Error())
		os.Exit(1)
	}
	return mbServer
}

func (h Handler) MbStart() {
	err := h.MbSlave.Start()
	if err != nil {
		slog.Error("Unable to start modbus slave: " + err.Error())
		os.Exit(1)
	}
}

func (h Handler) MbStop() {
	err := h.MbSlave.Stop()
	if err != nil {
		slog.Error("Unable to start modbus slave: " + err.Error())
		os.Exit(1)
	}
}

func (h *Handler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	slog.Warn("Not implemented!")
	return
}

func (h *Handler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	slog.Warn("Not implemented!")
	return
}

func (h *Handler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	// Write to DB entry with matching address.  Only update don't insert as the DbHandler should do the inserting of null values
	for i := 0; i < int(req.Quantity); i++ {
		regAddr := req.Addr + uint16(i)
		if req.IsWrite {
			slog.Debug("Updating database with holding registers", "address", regAddr, "value", req.Args[i])
			current, err := h.db.GetRowByAddress(int(regAddr))
			if err != nil {
				slog.Error("Unable to update database with holding registers", "address", regAddr, "value", req.Args[i], "err", err)
				return res, modbus.ErrProtocolError
			}
			_, err = h.db.Exec("UPDATE datapoints SET value = $1 WHERE address = $2", float64(req.Args[i])/current.Divisor, regAddr)
			if err != nil {
				slog.Error("Unable to update database with holding registers", "address", regAddr, "value", req.Args[i], "err", err)
				return res, modbus.ErrProtocolError
			}
		} else {
			slog.Debug("Reading holding registers", "address", regAddr)
			current, err := h.db.GetRowByAddress(int(regAddr))
			if err != nil {
				slog.Error("Unable to read from database", "address", regAddr, "error", err.Error())
				return res, modbus.ErrIllegalDataAddress
			}
			res = append(res, uint16(current.Value*current.Divisor))
		}
	}
	return
}

func (h *Handler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	slog.Warn("Not implemented!")
	return
}
