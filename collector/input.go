package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
)

const (
	input_contact = "input_contact"
)

var (
	inputContactInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, input_contact, "info"),
		"Generic input contact information.",
		labels("input", "name"), nil,
	)
	inputContactReadingsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, input_contact, "readings"),
		"Measured value of the input contact.",
		labels("input"), nil,
	)
)

func collectInput(p model.Probe, address, device string, ch chan<- prometheus.Metric) error {
	i := strconv.FormatUint(p.ID(), 10)

	for _, c := range p.Contacts {
		info, err := prometheus.NewConstMetric(
			inputContactInfoDesc, prometheus.GaugeValue, 1.0, address, device, i, c.Name,
		)
		if err != nil {
			return err
		}
		ch <- info

		readings, err := prometheus.NewConstMetric(
			inputContactReadingsDesc, prometheus.GaugeValue, c.Value, address, device, i,
		)
		if err != nil {
			return err
		}
		ch <- readings

		// input probes generally have only a single contact. if the slice happens
		// to have more than one contact, we would not be able to track its metrics
		// anyway, as we intentionally omit the contact label from the description,
		// which would cause the collector to complain about duplicate metrics
		break
	}

	return nil
}
