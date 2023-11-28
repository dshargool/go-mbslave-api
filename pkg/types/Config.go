package types

import (
	"encoding/json"
	"os"
)

type ConfigurationData struct {
	ApiPort     int         `json:"api_port"`
	ModbusPort  int         `json:"modbus_port"`
	DBPath      string      `json:"db"`
	Description string      `json:"description"`
	Registers   []ModbusTag `json:"registers"`
}
type Configuration struct {
	ApiPort    int
	ModbusPort int
	DBPath     string
	Registers  map[OpcTag]ModbusTag
}

func (c Configuration) ReadConfig(fileName string) (Configuration, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configData := ConfigurationData{}
	err = decoder.Decode(&configData)
	if err != nil {
		return Configuration{}, err
	}
	config := configData.dataToConfiguration()
	return config, nil
}

func (c ConfigurationData) dataToConfiguration() Configuration {
	config := Configuration{}
	config.Registers = make(map[OpcTag]ModbusTag)
	config.ApiPort = c.ApiPort
	config.ModbusPort = c.ModbusPort
	config.DBPath = c.DBPath
	for _, reg := range c.Registers {
		config.Registers[OpcTag(reg.Tag)] = reg
	}
	return config
}
