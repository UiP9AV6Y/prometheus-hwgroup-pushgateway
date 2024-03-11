package dto

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func NewRequest() *request.Root {
	return request.NewRoot()
}
