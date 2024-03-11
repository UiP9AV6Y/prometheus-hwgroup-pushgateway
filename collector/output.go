package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
)

const (
	output_contact = "output_contact"
)

var (
	outputContactInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, output_contact, "info"),
		"Generic output contact information.",
		labels("output", "name"), nil,
	)
	outputContactReadingsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, output_contact, "readings"),
		"Measured value of the output contact.",
		labels("output"), nil,
	)
)

func collectOutput(p model.Probe, address, device string, ch chan<- prometheus.Metric) error {
	o := strconv.FormatUint(p.ID(), 10)

	for _, c := range p.Contacts {
		info, err := prometheus.NewConstMetric(
			outputContactInfoDesc, prometheus.GaugeValue, 1.0, address, device, o, c.Name,
		)
		if err != nil {
			return err
		}
		ch <- info

		readings, err := prometheus.NewConstMetric(
			outputContactReadingsDesc, prometheus.GaugeValue, c.Value, address, device, o,
		)
		if err != nil {
			return err
		}
		ch <- readings

		// output probes generally have only a single contact. if the slice happens
		// to have more than one contact, we would not be able to track its metrics
		// anyway, as we intentionally omit the contact label from the description,
		// which would cause the collector to complain about duplicate metrics
		break
	}

	return nil
}
