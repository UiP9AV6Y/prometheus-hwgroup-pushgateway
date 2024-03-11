package device

import (
	"net/netip"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func (f *Factory) BuildNetwork() (*request.Network, error) {
	result := request.NewNetwork()

	ip := netip.AddrFrom4([4]byte{
		127,
		byte(f.gen.Intn(254)),
		byte(f.gen.Intn(254)),
		byte(2 + f.gen.Intn(252)),
	})
	sm := netip.AddrFrom4([4]byte{
		255,
		0,
		0,
		0,
	})
	gw := netip.AddrFrom4([4]byte{
		127,
		0,
		0,
		1,
	})

	result.IPAddr = &ip
	result.Submask = &sm
	result.Gateway = &gw
	result.HttpPort = 80

	return result, nil
}

func (f *Factory) UpdateNetwork(m *request.Network) error {
	return nil
}
