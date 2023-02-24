package main

// build with this command
// $ go build -buildmode=plugin

//change plugin name below command.
//if you don't change, daemon might not able to load your plugin.
// $ cd ..; mv plugin new_name

import (
	"airer_driver/internal/app/airer/device"
	"encoding/json"

	dev_ctrler "github.com/NaKa2355/pirem_pkg/device_controller"
)

type Config struct {
	SpiDevFile string `json:"spi_dev_file"`
	BusyPin    int    `json:"busy_pin"`
}

func GetController(jsonConfig json.RawMessage) (dev_ctrler.DeviceController, error) {
	config := Config{}
	json.Unmarshal(jsonConfig, &config)
	return device.New(config.SpiDevFile, config.BusyPin)
}

func main() {
}
