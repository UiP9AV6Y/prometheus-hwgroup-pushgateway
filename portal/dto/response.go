package dto

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/response"
)

func NewResponse() *response.Root {
	return response.NewRoot()
}
