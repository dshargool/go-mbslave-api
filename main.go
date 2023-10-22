package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dshargool/go-modbus-api.git/pkg/handlers"
	"github.com/dshargool/go-modbus-api.git/pkg/types"
)

var config = types.Configuration{}

func returnAllEndDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: returnAllEndDevice")
	json.NewEncoder(w).Encode(config.Registers)
}

func main() {

	config, _ = config.ReadConfig("config.json")

	app := handlers.New(config)
	app.HandleRequests(config.Port)
}
