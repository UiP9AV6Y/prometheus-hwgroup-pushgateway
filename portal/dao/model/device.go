package model

type Device struct {
	id uint64

	LastSeen uint64
	Uptime   uint64
	Model    uint8
	Firmware string
	Address  string
	Name     string

	Probes []Probe
}

func NewDevice(id uint64) *Device {
	result := &Device{
		id:     id,
		Probes: []Probe{},
	}

	return result
}

func (d *Device) ID() uint64 {
	return d.id
}

func (d *Device) FindProbeById(id uint64, group ProbeGroup) (Probe, bool) {
	for _, p := range d.Probes {
		if p.ID() == id && p.Group() == group {
			return p, true
		}
	}

	return *(NewProbe(d.id, id, group)), false
}
