package collector

import (
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
)

const (
	Namespace = "hwg_pushgateway"
)

type Collector struct {
	db     *dao.DAO
	logger log.Logger
}

func New(db *dao.DAO, logger log.Logger) *Collector {
	result := &Collector{
		db:     db,
		logger: logger,
	}

	return result
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- deviceInfoDesc
	ch <- deviceUptimeDesc
	ch <- deviceLastPushDesc

	ch <- sensorContactInfoDesc
	ch <- sensorContactStatusDesc
	ch <- sensorContactReadingsDesc

	ch <- inputContactInfoDesc
	ch <- inputContactReadingsDesc

	ch <- outputContactInfoDesc
	ch <- outputContactReadingsDesc

	ch <- meterInfoDesc
	ch <- meterContactInfoDesc
	ch <- meterContactStatusDesc
	ch <- meterContactReadingsDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	devices, err := c.db.GetAllDevices()
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to get devices for metrics collection", "err", err)
		return
	}

	for _, d := range devices {
		device := strconv.FormatUint(d.ID(), 10)

		if err := collectDevice(d, ch); err != nil {
			level.Error(c.logger).Log("msg", "failed to collect device metrics", "err", err)
			return
		}

		for _, p := range d.Probes {
			var err error

			switch g := p.Group(); g {
			case model.SensorProbeGroup:
				err = collectSensor(p, d.Address, device, ch)
			case model.InputProbeGroup:
				err = collectInput(p, d.Address, device, ch)
			case model.OutputProbeGroup:
				err = collectOutput(p, d.Address, device, ch)
			case model.MeterProbeGroup:
				err = collectMeter(p, d.Address, device, ch)
			}

			if err != nil {
				level.Error(c.logger).Log("msg", "failed to collect device probe metrics", "device", d.ID(), "probe", p.ID(), "probe_type", p.Group(), "err", err)
				return
			}
		}
	}
}

func labels(l ...string) []string {
	result := append([]string{"instance", "device"}, l...)

	return result
}
