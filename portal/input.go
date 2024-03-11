package portal

import (
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

const (
	inputThreshold float64 = 1.0
	inputExponent  int     = 0
	inputContact   uint64  = 1
)

func (p *Portal) processInputs(device *model.Device, r *request.Inputs) ([]model.Probe, string, error) {
	if r == nil || r.Entry == nil {
		level.Debug(p.logger).Log("msg", "client request contains no input information", "device", device.ID())
		return emptyProbes, "", nil
	}

	level.Debug(p.logger).Log("msg", "processing input information", "device", device.ID(), "entries", len(r.Entry))

	probes := make([]model.Probe, 0, len(r.Entry))

	for _, entry := range r.Entry {
		index := uint64(entry.ID)

		probe, _ := device.FindProbeById(index, model.InputProbeGroup)
		probe.Name = entry.Name

		level.Debug(p.logger).Log("msg", "processing input entry information", "device", device.ID(), "entry", entry.ID)
		p.processInputEntry(&probe, entry)

		probes = append(probes, probe)
	}

	return probes, "", nil
}

func (p *Portal) processInputEntry(probe *model.Probe, entry request.InputEntry) {
	contact, _ := probe.FindContactById(inputContact)
	changes := contact.Changes
	value := contact.Value

	if entry.ValueLog == nil || entry.ValueLog.Value == nil {
		level.Debug(p.logger).Log("msg", "value log information is missing for input", "device", probe.DeviceID(), "entry", entry.ID)
		contact.Value = float64(entry.Value)
	} else {
		v, e := calculateValueLog(entry.ValueLog.Value, inputExponent, inputThreshold, value)
		contact.Value = v
		contact.Changes = changes + e
	}

	contact.Name = entry.Name

	probe.SetContact(contact)
}
