package device

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func (f *Factory) BuildPortal() (*request.Portal, error) {
	result := request.NewPortal()
	result.Endpoint = f.addr
	result.DeviceName = "HWgroup Device Simulator"
	result.SetupCRC = f.gen.Int()

	return result, nil
}

func (f *Factory) UpdatePortal(m *request.Portal) error {
	m.SetupCRC = f.gen.Int()

	return nil
}
