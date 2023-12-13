package types

import "database/sql"

type InstrumentTag string

type ModbusTag struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Address     int    `json:"address"`
	DataType    string `json:"datatype"`
}

type ModbusResponse struct {
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
	Address     int     `json:"address"`
	DataType    string  `json:"datatype"`
	Value       float64 `json:"value"`
	LastUpdate  string  `json:"last_update"`
}

type SqlDb struct {
	*sql.DB
}
