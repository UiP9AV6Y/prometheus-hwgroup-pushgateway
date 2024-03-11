package request

import (
	"net/netip"
	"net/url"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/common"
)

type Entry struct {
	ID   int
	Name string
}

type ValueLog struct {
	Value      common.IntSequence `xml:"Value,omitempty"`
	TimeOffset common.IntSequence `xml:"TimeOffset,omitempty"`
	Status     common.IntSequence `xml:"Status,omitempty"`
}

func NewValueLog(timeOffset int, v ...int) *ValueLog {
	times := make([]int, 0, len(v))
	state := make([]int, 0, len(v))

	for i := len(v) - 1; i >= 0; i-- {
		times = append(times, 0-(timeOffset*i))
		state = append(state, 1)
	}

	result := &ValueLog{
		Value:      common.IntSequence(v),
		TimeOffset: common.IntSequence(times),
		Status:     common.IntSequence(state),
	}

	return result
}

type Agent struct {
	Version  string
	XmlVer   string
	Model    int
	VendorId int `xml:"vendor_id,omitempty"`
	MAC      string
	EMAI     string
	SerNum   uint64
}

func NewAgent() *Agent {
	result := &Agent{}

	return result
}

type Network struct {
	IPAddr   *netip.Addr `xml:"IPAddr,omitempty"`
	Submask  *netip.Addr `xml:"Submask,omitempty"`
	Gateway  *netip.Addr `xml:"Gateway,omitempty"`
	HttpPort int         `xml:"HttpPort,omitempty"`
}

func NewNetwork() *Network {
	result := &Network{}

	return result
}

type Portal struct {
	DeviceName string
	PushSource int
	PushPeriod int
	LogPeriod  int
	APLimit    int
	APSourceID int
	Endpoint   *url.URL `xml:"ServerAddres,omitempty"`
	PortalPort int
	SetupCRC   int
}

func NewPortal() *Portal {
	result := &Portal{}

	return result
}
