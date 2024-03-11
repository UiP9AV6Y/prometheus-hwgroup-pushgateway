package portal

import (
	"io"
	"net/url"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/common"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/response"
)

var emptyProbes = []model.Probe{}

type Portal struct {
	config *Config
	logger log.Logger
	auth   *url.Userinfo
	db     *dao.DAO
}

func New(c *Config, db *dao.DAO, l log.Logger) (*Portal, error) {
	a, err := c.ParseCredentials()
	if err != nil {
		return nil, err
	}

	result := &Portal{
		config: c,
		logger: l,
		auth:   a,
		db:     db,
	}

	return result, nil
}

func (p *Portal) TrustLevel() TrustLevel {
	return p.config.TrustLevel
}

func (p *Portal) Dump(w io.Writer) error {
	return p.db.ExportWriter(w)
}

func (p *Portal) Process(r *Request) (*response.Root, error) {
	addr := r.Client()
	if addr == nil {
		level.Debug(p.logger).Log("msg", "client address information is missing")
		return p.Error("Insufficient meta data"), nil
	}

	if !r.CompareCredentials(p.auth) {
		if r.Credentials != nil {
			level.Info(p.logger).Log("msg", "authentication failed", "user", r.Credentials.Username(), "client", addr)
		} else {
			level.Info(p.logger).Log("msg", "authentication credentials missing", "client", addr)
		}
		return p.Error("Authentication failed"), nil
	}

	if r.Payload == nil {
		level.Info(p.logger).Log("msg", "request payload is missing", "client", addr)
		return p.Error("No push data provided"), nil
	}

	device, dErr, err := p.processAgent(addr, r.Payload.Agent)
	if err != nil {
		return nil, err
	} else if dErr != "" {
		return p.Error(dErr), nil
	}

	if tErr, err := p.processTime(device, r.Payload.Time); err != nil {
		return nil, err
	} else if tErr != "" {
		return p.Error(tErr), nil
	}

	if pErr, err := p.processPortal(device, r.Payload.Portal); err != nil {
		return nil, err
	} else if pErr != "" {
		return p.Error(pErr), nil
	}

	sProbes, sErr, err := p.processSensors(device, r.Payload.Sensors)
	if err != nil {
		return nil, err
	} else if sErr != "" {
		return p.Error(sErr), nil
	}

	iProbes, iErr, err := p.processInputs(device, r.Payload.Inputs)
	if err != nil {
		return nil, err
	} else if iErr != "" {
		return p.Error(iErr), nil
	}

	oProbes, oErr, err := p.processOutputs(device, r.Payload.Outputs)
	if err != nil {
		return nil, err
	} else if oErr != "" {
		return p.Error(oErr), nil
	}

	mProbes, mErr, err := p.processMeters(device, r.Payload.Meters)
	if err != nil {
		return nil, err
	} else if mErr != "" {
		return p.Error(mErr), nil
	}

	device.Probes = make([]model.Probe, 0, len(sProbes)+len(iProbes)+len(oProbes)+len(mProbes))
	device.Probes = append(device.Probes, sProbes...)
	device.Probes = append(device.Probes, iProbes...)
	device.Probes = append(device.Probes, oProbes...)
	device.Probes = append(device.Probes, mProbes...)

	if err := p.db.UpsertDevice(*device); err != nil {
		return nil, err
	}

	return p.Success(), nil
}

func (p *Portal) Success() *response.Root {
	rsp := response.NewRoot()

	rsp.Time.Now()
	rsp.Portal.Status = response.StatusOK
	rsp.Portal.PushPeriod = int(p.config.PushInterval.Seconds())
	rsp.Portal.LogPeriod = int(p.config.LogInterval.Seconds())
	rsp.Portal.APLimit = int(p.config.AutoPushDelay.Seconds())

	return rsp
}

func (p *Portal) Error(msg string) *response.Root {
	rsp := response.NewRoot()

	rsp.Time.Now()
	rsp.Portal.Status = response.StatusBusy
	rsp.Portal.PortalMessage = msg

	return rsp
}

func absDiffFloat(x, y float64) float64 {
	if x < y {
		return y - x
	}
	return x - y
}

func calculateValue(value, exponent int) float64 {
	result := float64(value)
	exp := float64(exponent)

	if exp < 0.0 {
		return result / (10.0 * (0.0 - exp))
	} else if exp > 0.0 {
		return result * (10.0 * exp)
	}

	return result
}

func calculateValueLog(values common.IntSequence, exponent int, threshold, current float64) (float64, uint64) {
	var changes uint64
	for _, value := range values {
		v := calculateValue(value, exponent)

		if absDiffFloat(absDiffFloat(v, current), 0) > threshold {
			current = v
			changes += 1
		}
	}

	return current, changes
}
