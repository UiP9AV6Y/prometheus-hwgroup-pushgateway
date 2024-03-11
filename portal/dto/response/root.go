package response

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/common"
)

type Root struct {
	Time   *common.Time `xml:"Time,omitempty"`
	Portal *Portal      `xml:"Portal,omitempty"`
	//Sensors *Sensors
}

func NewRoot() *Root {
	result := &Root{
		Time:   common.NewTime(),
		Portal: NewPortal(),
	}

	return result
}
