package types

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
    Port int `json:"port"`
    Slaves []ModbusSlave
    Registers []ModbusTag
}

func (c Configuration) ReadConfig(fileName string) (Configuration, error) {
    file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Unable to open config file: ", fileName, err)
        return Configuration{}, err
    }
    defer file.Close()
    decoder := json.NewDecoder(file)
    config := Configuration{}
    err = decoder.Decode(&config)
    if err != nil {
        fmt.Println("Could not decode config file: ", fileName, err)
    }

    return config, nil
}
