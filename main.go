package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

type OpcTag string

type EndDevice struct {
    Name string `json:"name"`
    IpAddress net.IP `json:"IpAddress"`
    DeviceTags []ModbusTag
}

type ModbusTag struct {
    Tag OpcTag `json:"tag"`
    Description string `json:"desc"`
    Address int `json:"address"`
    Divisor int `json:"divisor"`
}

    var Tags = []ModbusTag{
        {Tag: "TestTag", Description: "Testing Tag for testing", Address: 1234, Divisor: 10},
        {Tag: "TestTag2", Description: "Testing Tag 2 for testing", Address: 4321, Divisor: 10},
    }
    var Devices = EndDevice{
        Name: "TestDevice",
        IpAddress: net.ParseIP("127.0.0.1"),
        DeviceTags: Tags,
    }

func returnAllEndDevice(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint hit: returnAllEndDevice")
    json.NewEncoder(w).Encode(Devices)
}

func main() {

	fmt.Println("Starting server at port 8081")
    http.HandleFunc("/all_devices", returnAllEndDevice)
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal(err)
    }
}
