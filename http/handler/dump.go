package handler

import (
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal"
)

func Dump(p *portal.Portal, logger log.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := p.Dump(w); err != nil {
			level.Error(logger).Log("msg", "failed to encode response", "err", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})
}
