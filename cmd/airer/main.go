package main

// build with this command
// $ go build -buildmode=plugin

//change plugin name below command.
//if you don't change, daemon might not able to load your plugin.
// $ cd ..; mv plugin new_name

import (
	"airer_driver/internal/app/airer/device"

	dev_plugin "github.com/NaKa2355/pirem/pkg/plugin"
	plugin "github.com/hashicorp/go-plugin"
)

type Config struct {
	SpiDevFile string `json:"spi_dev_file"`
	BusyPin    int    `json:"busy_pin"`
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dev_plugin.Handshake,
		Plugins: map[string]plugin.Plugin{
			"device_controller": &dev_plugin.DevicePlugin{Impl: &device.Device{}},
		},

		GRPCServer: plugin.DefaultGRPCServer,
	})
}
