package handlers

import (
	"encoding/binary"
	"errors"
	"log/slog"
	"math"
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
	slog.Info("HandleHoldingRegisters - new request", "req", req)
	var dataType string = "none"
	for i := 0; i < int(req.Quantity); i++ {
        // Move our request address along to service the entire quantity
		regAddr := req.Addr + uint16(i)

        // If our dataType is uninitialized we try to do it from the database. 
        // If it fails here we don't know what type of data to expect to read and it will fail
		if dataType == "none" {
			dataType, err = h.db.GetDataTypeByAddress(int(regAddr))
			if err != nil {
				slog.Error("Unable to read row data type", "address", regAddr,
					"allow_null", h.AllowNullRegisters, "req", req, "err", err)
				return res, modbus.ErrProtocolError
			}
		}

        // If the requested quantity is an odd number of bits we have to make 
        // it an integer type as floats are 2 or 4 bytes wide
		if req.Quantity%2 != 0 {
			dataType = "uint16"
		}

        // Based on the data type get the number of registers we'll return
		num_regs, err := numRegsDataType(dataType)
		if err != nil {
			slog.Error("Unable to calculate number of required registers for datatype.",
				"datatype", dataType, "num_regs", num_regs)
			return res, modbus.ErrIllegalDataAddress
		}

		if req.IsWrite {
			slog.Info("Writing holding registers", "address", regAddr)

            // Put our arguments that we're interested in into a new data slice
			var data []uint16
			for j := 0; j < int(num_regs); j++ {
				data = append(data, req.Args[i+j])
			}
            // Convert the bytes of our slice to our data type
			conv_val, err := parseByteToDataType(dataType, req.Args)
			if err != nil {
				slog.Error("Unable to convert data type",
					"address", regAddr, "value", conv_val, "err", err)
				return res, modbus.ErrProtocolError
			}

			slog.Info("Updating database with holding registers",
				"address", regAddr, "data", data, "value", conv_val)

            // Write the value we received into the DB
			err = h.db.SetAddressValue(int(regAddr), conv_val)
			if err != nil {
				slog.Error("Unable to update database with holding registers",
					"address", regAddr, "value", conv_val, "err", err)
				return res, modbus.ErrProtocolError
			}

			i = i + int(num_regs) - 1

		} else {
			slog.Info("Reading holding registers", "address", regAddr)

			// Get the current value from the database
			current, err := h.db.GetRowByAddress(int(regAddr))
			if err != nil {
				slog.Error("Unable to read from database",
					"address", regAddr, "error", err.Error())
				// When we don't have a database value but allow null registers we return a 0
                // if we don't allow null values it's considered an illegal data address
				if h.AllowNullRegisters {
					slog.Warn("Setting Null Register to 0")
					current.Value = 0
				} else {
					return res, modbus.ErrIllegalDataAddress
				}
			}

			// Take our value and parse it into the datatype we expect to use
			conv_val, err := parseDataTypeToByte(dataType, float64(current.Value))
			if err != nil {
				slog.Error("Couldn't parse DataType to Byte",
					"DataType", dataType)
			}
			// Increment the addresses by the amount we're appending minus the regular increase
			i = i + int(num_regs) - 1
			res = append(res, conv_val...)
		}
	}
	slog.Info("HandleHoldingRegisters - Returning data",
		"length", len(res), "request_len", req.Quantity, "error", err)
	return
}

func (h *Handler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	slog.Warn("Not implemented!")
	return
}

func parseDataTypeToByte(dataType string, value float64) (res []uint16, err error) {
	switch dataType {
	case "float32":
		bits := math.Float32bits(float32(value))
		res = append(res, uint16((bits>>16)&0xffff))
		res = append(res, uint16((bits)&0xffff))
	case "float64":
		bits := math.Float64bits(value)
		res = append(res, uint16(bits>>48)&0xffff)
		res = append(res, uint16(bits>>32)&0xffff)
		res = append(res, uint16(bits>>16)&0xffff)
		res = append(res, uint16(bits)&0xffff)
	case "int16":
		res = append(res, uint16(int16(value)))
	case "uint16":
		res = append(res, uint16(value))
	default:
		return nil, errors.New("Can't parse dataType: " + dataType)
	}
	return res, nil
}

func parseByteToDataType(dataType string, bytes []uint16) (res float64, err error) {
	switch dataType {
	case "float32":
		b := make([]byte, 4)
		b[0] = byte(bytes[0] >> 8 & 0xff)
		b[1] = byte(bytes[0] & 0xff)
		b[2] = byte(bytes[1] >> 8 & 0xff)
		b[3] = byte(bytes[1] & 0xff)

		f_bits := binary.BigEndian.Uint32(b)
		res = float64(math.Float32frombits(f_bits))
	case "float64":
		b := make([]byte, 8)
		b[0] = byte(bytes[0] >> 8 & 0xff)
		b[1] = byte(bytes[0] & 0xff)
		b[2] = byte(bytes[1] >> 8 & 0xff)
		b[3] = byte(bytes[1] & 0xff)
		b[4] = byte(bytes[2] >> 8 & 0xff)
		b[5] = byte(bytes[2] & 0xff)
		b[6] = byte(bytes[3] >> 8 & 0xff)
		b[7] = byte(bytes[3] & 0xff)

		f_bits := binary.BigEndian.Uint64(b)
		res = math.Float64frombits(f_bits)
	case "int16":
		res = float64(bytes[0])
	case "uint16":
		res = float64(bytes[0])
	default:
		return 0, errors.New("Can't parse dataType")
	}
	return res, nil
}

func numRegsDataType(dataType string) (res uint16, err error) {
	switch dataType {
	case "float32":
		res = 2
	case "float64":
		res = 4
	case "int16":
		res = 1
	case "uint16":
		res = 1
	default:
		return 0, errors.New("Can't parse dataType: " + dataType)
	}
	return res, nil
}
