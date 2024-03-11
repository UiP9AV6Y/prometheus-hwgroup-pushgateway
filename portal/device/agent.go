package device

import (
	"fmt"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func (f *Factory) BuildAgent() (*request.Agent, error) {
	result := request.NewAgent()

	mac := make([]byte, 12)
	if _, err := f.gen.Read(mac); err != nil {
		return nil, err
	}

	// clear multicast bit (&^), ensure local bit (|)
	mac[0] = (mac[0] | 0x02) & 0xfe

	result.Version = "0.0.0"
	result.XmlVer = "1.0.0"
	result.Model = 254
	result.VendorId = f.gen.Int()
	result.MAC = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
	result.EMAI = fmt.Sprintf("%1X%1X%1X%1X%1X%1X%1X%1X%1X%1X%1X%1X",
		mac[0], mac[1], mac[2], mac[3], mac[4], mac[5], mac[6], mac[7], mac[8], mac[9], mac[10], mac[11])
	result.SerNum = f.gen.Uint64()

	return result, nil
}

func (f *Factory) UpdateAgent(m *request.Agent) error {
	return nil
}
