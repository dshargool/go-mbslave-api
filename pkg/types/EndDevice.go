package types

import ( "net")

type OpcTag string

type ModbusSlave struct {
    Name string `json:"name"`
    Ip net.IP `json:"ip"`
}


type ModbusTag struct {
    Tag string `json:"tag"`
    Description string `json:"desc"`
    Address int `json:"address"`
    Divisor int `json:"divisor"`
    Slave int `json:"slave"`
}
