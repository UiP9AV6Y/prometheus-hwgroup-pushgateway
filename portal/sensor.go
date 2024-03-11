package portal

import (
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

const (
	sensorThreshold float64 = 5.0
	sensorContact   uint64  = 1
)

func (p *Portal) processSensors(device *model.Device, r *request.Sensors) ([]model.Probe, string, error) {
	if r == nil || r.Entry == nil {
		level.Debug(p.logger).Log("msg", "client request contains no sensor information", "device", device.ID())
		return emptyProbes, "", nil
	}

	level.Debug(p.logger).Log("msg", "processing sensor information", "device", device.ID(), "entries", len(r.Entry))

	probes := make([]model.Probe, 0, len(r.Entry))

	for _, entry := range r.Entry {
		index := uint64(entry.ID)

		probe, _ := device.FindProbeById(index, model.SensorProbeGroup)
		probe.Name = entry.Name

		level.Debug(p.logger).Log("msg", "processing sensor entry information", "device", device.ID(), "entry", entry.ID)
		p.processSensorEntry(&probe, entry)

		probes = append(probes, probe)
	}

	return probes, "", nil
}

func (p *Portal) processSensorEntry(probe *model.Probe, entry request.SensorEntry) {
	contact, _ := probe.FindContactById(sensorContact)
	changes := contact.Changes
	value := contact.Value

	if entry.ValueLog == nil || entry.ValueLog.Value == nil {
		level.Debug(p.logger).Log("msg", "value log information is missing for sensor", "device", probe.DeviceID(), "entry", entry.ID)
		contact.Value = float64(entry.Value)
	} else {
		v, e := calculateValueLog(entry.ValueLog.Value, entry.Exp, sensorThreshold, value)
		contact.Value = v
		contact.Changes = changes + e
	}

	contact.Status = uint8(entry.Status)
	contact.Name = entry.Name
	contact.Unit = entry.Units

	probe.SetContact(contact)
}
