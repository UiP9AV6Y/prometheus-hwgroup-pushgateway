package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
)

const (
	sensor_contact = "sensor_contact"
)

var (
	sensorContactInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, sensor_contact, "info"),
		"Generic sensor contact information.",
		labels("sensor", "name", "unit"), nil,
	)
	sensorContactStatusDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, sensor_contact, "status"),
		"Operational status of the sensor contact.",
		labels("sensor"), nil,
	)
	sensorContactReadingsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, sensor_contact, "readings"),
		"Measured value of the sensor contact.",
		labels("sensor"), nil,
	)
)

func collectSensor(p model.Probe, address, device string, ch chan<- prometheus.Metric) error {
	s := strconv.FormatUint(p.ID(), 10)

	for _, c := range p.Contacts {
		info, err := prometheus.NewConstMetric(
			sensorContactInfoDesc, prometheus.GaugeValue, 1.0, address, device, s, c.Name, c.Unit,
		)
		if err != nil {
			return err
		}
		ch <- info

		status, err := prometheus.NewConstMetric(
			sensorContactStatusDesc, prometheus.GaugeValue, float64(c.Status), address, device, s,
		)
		if err != nil {
			return err
		}
		ch <- status

		readings, err := prometheus.NewConstMetric(
			sensorContactReadingsDesc, prometheus.GaugeValue, c.Value, address, device, s,
		)
		if err != nil {
			return err
		}
		ch <- readings

		// sensor probes generally have only a single contact. if the slice happens
		// to have more than one contact, we would not be able to track its metrics
		// anyway, as we intentionally omit the contact label from the description,
		// which would cause the collector to complain about duplicate metrics
		break
	}

	return nil
}
