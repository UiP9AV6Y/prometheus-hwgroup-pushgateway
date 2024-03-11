package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
)

const (
	meter         = "meter"
	meter_contact = "meter_contact"
)

var (
	meterInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, meter, "info"),
		"Generic meter information.",
		labels("meter", "name"), nil,
	)
	meterContactInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, meter_contact, "info"),
		"Generic meter contact information.",
		labels("meter", "contact", "name", "unit"), nil,
	)
	meterContactStatusDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, meter_contact, "status"),
		"Operational status of the meter contact.",
		labels("meter", "contact"), nil,
	)
	meterContactReadingsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, meter_contact, "readings"),
		"Measured value of the meter contact.",
		labels("meter", "contact"), nil,
	)
)

func collectMeter(p model.Probe, address, device string, ch chan<- prometheus.Metric) error {
	m := strconv.FormatUint(p.ID(), 10)

	meterInfo, err := prometheus.NewConstMetric(
		meterInfoDesc, prometheus.GaugeValue, 1.0, address, device, m, p.Name,
	)
	if err != nil {
		return err
	}
	ch <- meterInfo

	for _, c := range p.Contacts {
		i := strconv.FormatUint(c.ID(), 10)

		info, err := prometheus.NewConstMetric(
			meterContactInfoDesc, prometheus.GaugeValue, 1.0, address, device, m, i, c.Name, c.Unit,
		)
		if err != nil {
			return err
		}
		ch <- info

		status, err := prometheus.NewConstMetric(
			meterContactStatusDesc, prometheus.GaugeValue, float64(c.Status), address, device, m, i,
		)
		if err != nil {
			return err
		}
		ch <- status

		readings, err := prometheus.NewConstMetric(
			meterContactReadingsDesc, prometheus.GaugeValue, c.Value, address, device, m, i,
		)
		if err != nil {
			return err
		}
		ch <- readings
	}

	return nil
}
