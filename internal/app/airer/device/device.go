package device

import (
	"airer_driver/internal/app/airer/driver"
	"context"
	"encoding/json"

	irdatav1 "github.com/NaKa2355/irdeck-proto/gen/go/common/irdata/v1"
	apiremv1 "github.com/NaKa2355/irdeck-proto/gen/go/pirem/api/v1"
)

type Device struct {
	d    *driver.Driver
	info *apiremv1.DeviceInfo
}

type DeviceConfig struct {
	SpiDevFile string `json:"spi_dev_file"`
	BusyPin    int    `json:"busy_pin"`
}

const DriverVersion = "0.1.0"

func (dev *Device) setInfo() error {
	var err error = nil
	dev.info = &apiremv1.DeviceInfo{}
	bufferSize, err := dev.d.GetBufSize()
	if err != nil {
		return err
	}
	firmVersion, err := dev.d.GetVersion()
	if err != nil {
		return err
	}
	dev.info.BufferSize = int32(bufferSize)
	dev.info.FirmwareVersion = firmVersion
	dev.info.Service = apiremv1.DeviceInfo_SERVICE_TYPE_SEND_RECEIVE
	dev.info.Type = apiremv1.DeviceInfo_DEVICE_TYPE_WIRED
	dev.info.DriverVersion = DriverVersion
	return nil
}

func (dev *Device) Init(ctx context.Context, jsonConf json.RawMessage) error {
	conf := DeviceConfig{}
	err := json.Unmarshal(jsonConf, &conf)
	if err != nil {
		return err
	}
	d, err := driver.New(conf.SpiDevFile, conf.BusyPin)
	if err != nil {
		return err
	}
	dev.d = d
	if err := dev.setInfo(); err != nil {
		return err
	}
	return nil
}

func (dev *Device) GetDeviceInfo(ctx context.Context) (*apiremv1.DeviceInfo, error) {
	return dev.info, nil
}

func (dev *Device) GetDeviceStatus(ctx context.Context) (*apiremv1.DeviceStatus, error) {
	status := &apiremv1.DeviceStatus{}
	bufferSize, err := dev.d.GetBufSize()
	if err != nil {
		return status, err
	}

	if dev.info.BufferSize == int32(bufferSize) {
		status.IsActive = true
	} else {
		status.IsActive = false
	}

	return status, nil
}

func (dev *Device) SendRawIr(ctx context.Context, irData *irdatav1.RawIrData) error {
	return dev.d.SendIr(ctx, convertToDriverIrRawData(irData.OnOffPluseNs))
}

func (dev *Device) ReceiveRawIr(ctx context.Context) (*irdatav1.RawIrData, error) {
	rawIrData := &irdatav1.RawIrData{}
	irData, err := dev.d.ReceiveIr(ctx)
	if err != nil {
		return rawIrData, err
	}
	rawIrData.CarrierFreqKhz = 40
	rawIrData.OnOffPluseNs = convertToApiIrRawData(irData)
	return rawIrData, nil
}

func (dev *Device) IsBusy(context.Context) (bool, error) {
	return dev.d.IsBusy(), nil
}

func (dev *Device) Drop() error {
	return dev.d.Close()
}
