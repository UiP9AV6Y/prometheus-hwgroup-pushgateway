package response

import (
	"net/url"
)

const (
	StatusOK   int = 0
	StatusBusy int = 1
)

type Portal struct {
	Status        int
	PushPeriod    int      `xml:"PushPeriod,omitempty"`
	LogPeriod     int      `xml:"LogPeriod,omitempty"`
	APLimit       int      `xml:"APLimit,omitempty"`
	PortalMessage string   `xml:"PortalMessage,omitempty"`
	Endpoint      *url.URL `xml:"ServerAddres,omitempty"`
}

func NewPortal() *Portal {
	result := &Portal{}

	return result
}
