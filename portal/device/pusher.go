package device

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	pconfig "github.com/prometheus/common/config"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/response"
)

var ErrMalformedRequest = errors.New("malformed request")

type Pusher struct {
	logger    log.Logger
	connector *pconfig.HTTPClientConfig
	factory   *Factory
}

func NewPusher(l log.Logger, c *pconfig.HTTPClientConfig, f *Factory) *Pusher {
	result := &Pusher{
		logger:    l,
		connector: c,
		factory:   f,
	}

	return result
}

func (r *Pusher) Report(data *request.Root) (*response.Root, error) {
	client, err := pconfig.NewClientFromConfig(*r.connector, "hwgroup_device_simulator", pconfig.WithKeepAlivesDisabled(), pconfig.WithHTTP2Disabled())
	if err != nil {
		level.Error(r.logger).Log("msg", "Error generating HTTP client", "err", err)
		return nil, err
	}

	if data.Portal == nil || data.Portal.Endpoint == nil {
		level.Error(r.logger).Log("msg", "Request data does not contain an endpoint URL")
		return nil, ErrMalformedRequest
	}

	body := new(bytes.Buffer)
	enc := xml.NewEncoder(body)
	if err := enc.Encode(data); err != nil {
		level.Error(r.logger).Log("msg", "Failed to marshal request body", "err", err)
		return nil, err
	}
	if err := enc.Close(); err != nil {
		level.Error(r.logger).Log("msg", "Failed to close body encoder", "err", err)
		return nil, err
	}

	var req *http.Request
	req, err = http.NewRequest("POST", data.Portal.Endpoint.String(), body)
	if err != nil {
		level.Error(r.logger).Log("msg", "Failed to create request", "err", err)
		return nil, err
	}

	req.Header = map[string][]string{
		"Accept":       {"text/xml"},
		"Content-Type": {"text/xml"},
	}

	level.Debug(r.logger).Log("msg", "Begin request", "url", data.Portal.Endpoint)

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		return nil, err
	}

	defer func() {
		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			level.Error(r.logger).Log("msg", "Failed to discard body", "err", err)
		}
		resp.Body.Close()
	}()

	level.Debug(r.logger).Log("msg", "Request completed", "url", data.Portal.Endpoint, "status_code", resp.StatusCode, "duration", elapsed)

	result := response.NewRoot()
	if err := xml.NewDecoder(resp.Body).Decode(result); err != nil {
		level.Error(r.logger).Log("msg", "Failed to parse response", "err", err)
		return nil, err
	}

	return result, nil
}

func (r *Pusher) Run(i time.Duration) error {
	data, err := r.factory.Build()
	if err != nil {
		level.Error(r.logger).Log("msg", "Failed to create device request", "err", err)
		return err
	}

	ticker := time.NewTicker(i)
	defer ticker.Stop()

	level.Info(r.logger).Log("msg", "start reporting", "interval", i)

	for {
		select {
		case <-ticker.C:
			if err := r.factory.Update(data); err != nil {
				return err
			}

			res, err := r.Report(data)
			if err != nil {
				return err
			}

			if res.Portal == nil {
				level.Debug(r.logger).Log("msg", "response has no portal information")
				continue
			}

			if res.Portal.PortalMessage != "" {
				level.Info(r.logger).Log("msg", "received portal message", "data", res.Portal.PortalMessage)
			}

			if res.Portal.Endpoint != nil {
				level.Info(r.logger).Log("msg", "adjusting report endpoint", "url", res.Portal.Endpoint)
				data.Portal.Endpoint = res.Portal.Endpoint
			}

			if p := time.Duration(res.Portal.PushPeriod) * time.Second; p > 0 && p != i {
				level.Info(r.logger).Log("msg", "adjusting report interval", "interval", p)
				i = p
				ticker.Reset(i)
			}
		}
	}

	return nil
}
