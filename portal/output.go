package portal

import (
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

const (
	outputThreshold float64 = 1.0
	outputExponent  int     = 0
	outputContact   uint64  = 1
)

func (p *Portal) processOutputs(device *model.Device, r *request.Outputs) ([]model.Probe, string, error) {
	if r == nil || r.Entry == nil {
		level.Debug(p.logger).Log("msg", "client request contains no output information", "device", device.ID())
		return emptyProbes, "", nil
	}

	level.Debug(p.logger).Log("msg", "processing output information", "device", device.ID(), "entries", len(r.Entry))

	probes := make([]model.Probe, 0, len(r.Entry))

	for _, entry := range r.Entry {
		index := uint64(entry.ID)

		probe, _ := device.FindProbeById(index, model.OutputProbeGroup)
		probe.Name = entry.Name

		level.Debug(p.logger).Log("msg", "processing output entry information", "device", device.ID(), "entry", entry.ID)
		p.processOutputEntry(&probe, entry)

		probes = append(probes, probe)
	}

	return probes, "", nil
}

func (p *Portal) processOutputEntry(probe *model.Probe, entry request.OutputEntry) {
	contact, _ := probe.FindContactById(outputContact)
	changes := contact.Changes
	value := contact.Value

	if entry.ValueLog == nil || entry.ValueLog.Value == nil {
		level.Debug(p.logger).Log("msg", "value log information is missing for output", "device", probe.DeviceID(), "entry", entry.ID)
		contact.Value = float64(entry.Value)
	} else {
		v, e := calculateValueLog(entry.ValueLog.Value, outputExponent, outputThreshold, value)
		contact.Value = v
		contact.Changes = changes + e
	}

	contact.Name = entry.Name

	probe.SetContact(contact)
}
