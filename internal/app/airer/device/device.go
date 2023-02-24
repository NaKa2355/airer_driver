package device

import (
	"airer_driver/internal/app/airer/driver"
	"context"

	apirem_v1 "github.com/NaKa2355/pirem_pkg/apirem.v1"
)

type Device struct {
	d    *driver.Driver
	info *apirem_v1.DeviceInfo
}

const DriverVersion = "0.1.0"

func (dev *Device) setInfo() error {
	var err error = nil
	dev.info = &apirem_v1.DeviceInfo{}
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
	dev.info.Service = apirem_v1.DeviceInfo_SERVICE_TYPE_SEND_RECEIVE
	dev.info.Type = apirem_v1.DeviceInfo_DEVICE_TYPE_WIRED
	dev.info.DriverVersion = DriverVersion
	return nil
}

func New(spiDevFile string, busyPinNum int) (*Device, error) {
	dev := &Device{}
	d, err := driver.New(spiDevFile, busyPinNum)
	if err != nil {
		return dev, err
	}
	dev.d = d
	if err := dev.setInfo(); err != nil {
		return dev, err
	}
	return dev, nil
}

func (dev *Device) GetDeviceInfo(ctx context.Context) (*apirem_v1.DeviceInfo, error) {
	return dev.info, nil
}

func (dev *Device) GetDeviceStatus(ctx context.Context) (*apirem_v1.DeviceStatus, error) {
	status := &apirem_v1.DeviceStatus{}
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

func (dev *Device) SendIr(ctx context.Context, irData *apirem_v1.RawIrData) error {
	return dev.d.SendIr(convertToDriverIrRawData(irData.OnOffPluseNs))
}

func (dev *Device) ReceiveIr(ctx context.Context) (*apirem_v1.RawIrData, error) {
	rawIrData := &apirem_v1.RawIrData{}
	irData, err := dev.d.ReceiveIr()
	if err != nil {
		return rawIrData, err
	}
	rawIrData.CarrierFreqKhz = 40
	rawIrData.OnOffPluseNs = convertToApiIrRawData(irData)
	return rawIrData, nil
}

func (dev *Device) Drop() error {
	return dev.d.Close()
}