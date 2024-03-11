package portal

import (
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

const (
	meterThreshold float64 = 5.0
)

func (p *Portal) processMeters(device *model.Device, r *request.Meters) ([]model.Probe, string, error) {
	if r == nil || r.Entry == nil {
		level.Debug(p.logger).Log("msg", "client request contains no meter information", "device", device.ID())
		return emptyProbes, "", nil
	}

	level.Debug(p.logger).Log("msg", "processing meter information", "device", device.ID(), "entries", len(r.Entry))

	probes := make([]model.Probe, 0, len(r.Entry))
	for _, entry := range r.Entry {
		index := uint64(entry.ID)

		probe, _ := device.FindProbeById(index, model.MeterProbeGroup)
		probe.Name = entry.Name
		probe.PrimaryAddress = uint64(entry.Address)
		probe.SecondaryAddress = uint64(entry.SecAddress)

		if entries := entry.Values.Entry; entries == nil {
			level.Debug(p.logger).Log("msg", "client request contains no meter entry information", "device", device.ID())
		} else {
			level.Debug(p.logger).Log("msg", "processing meter entry information", "device", device.ID(), "entries", len(entries))
			p.processMeterEntries(&probe, entries)
		}

		probes = append(probes, probe)
	}

	return probes, "", nil
}

func (p *Portal) processMeterEntries(probe *model.Probe, entries []request.MeterValueEntry) {
	contacts := make([]model.Contact, 0, len(entries))
	for _, entry := range entries {
		contact, _ := probe.FindContactById(uint64(entry.ID))
		changes := contact.Changes
		value := calculateValue(entry.Value, entry.Exp)

		if absDiffFloat(absDiffFloat(value, contact.Value), 0) > meterThreshold {
			changes += 1
		}

		contact.Status = uint8(entry.Status)
		contact.Name = entry.Name
		contact.Unit = entry.Unit
		contact.Value = value
		contact.Changes = changes

		contacts = append(contacts, contact)
	}

	probe.Contacts = contacts
}
