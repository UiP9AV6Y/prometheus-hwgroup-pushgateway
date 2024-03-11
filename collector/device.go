package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
)

const (
	device = "device"
)

var (
	deviceInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, device, "info"),
		"General hardware, software, and firmware information.",
		labels("firmware", "model", "name"), nil,
	)
	deviceUptimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, device, "uptime_seconds"),
		"System uptime in seconds.",
		labels(), nil,
	)
	deviceLastPushDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, device, "last_push_timestamp_seconds"),
		"Unix timestamp of the last push received from a device.",
		labels(), nil,
	)
)

func collectDevice(d model.Device, ch chan<- prometheus.Metric) error {
	m := strconv.FormatUint(uint64(d.Model), 10)
	i := strconv.FormatUint(d.ID(), 10)

	info, err := prometheus.NewConstMetric(
		deviceInfoDesc, prometheus.GaugeValue, 1.0,
		d.Address, i, d.Firmware, m, d.Name,
	)
	if err != nil {
		return err
	}
	ch <- info

	uptime, err := prometheus.NewConstMetric(
		deviceUptimeDesc, prometheus.GaugeValue, float64(d.Uptime), d.Address, i,
	)
	if err != nil {
		return err
	}
	ch <- uptime

	lastPush, err := prometheus.NewConstMetric(
		deviceLastPushDesc, prometheus.GaugeValue, float64(d.LastSeen), d.Address, i,
	)
	if err != nil {
		return err
	}
	ch <- lastPush

	return nil
}
