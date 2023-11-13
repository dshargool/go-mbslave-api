package types

import "database/sql"

type OpcTag string

type ModbusTag struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Address     int    `json:"address"`
	Divisor     int    `json:"divisor"`
}

type ModbusResponse struct {
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
	Address     int     `json:"address"`
	Divisor     float64 `json:"divisor"`
	Value       float64 `json:"value"`
	LastUpdate  string  `json:"last_update"`
}

type SqlDb struct {
	*sql.DB
}
