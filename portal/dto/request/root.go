package request

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/common"
)

type Root struct {
	Agent   *Agent       `xml:"Agent,omitempty"`
	Network *Network     `xml:"Network,omitempty"`
	Time    *common.Time `xml:Time,omitempty"`
	Portal  *Portal      `xml:"Portal,omitempty"`
	Sensors *Sensors     `xml:"Sensors,omitempty"`
	Meters  *Meters      `xml:"Meters,omitempty"`
	Inputs  *Inputs      `xml:"Inputs,omitempty"`
	Outputs *Outputs     `xml:"Outputs,omitempty"`
}

func NewRoot() *Root {
	result := &Root{
		Time:   common.NewTime(),
		Portal: NewPortal(),
	}

	return result
}
