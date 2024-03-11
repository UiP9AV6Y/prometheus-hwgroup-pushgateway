package model

type Probe struct {
	deviceID uint64
	id       uint64
	group    ProbeGroup

	PrimaryAddress   uint64
	SecondaryAddress uint64
	Name             string

	Contacts []Contact
}

func NewProbe(deviceID uint64, id uint64, group ProbeGroup) *Probe {
	result := &Probe{
		deviceID: deviceID,
		id:       id,
		group:    group,
		Contacts: []Contact{},
	}

	return result
}

func (p *Probe) DeviceID() uint64 {
	return p.deviceID
}

func (p *Probe) ID() uint64 {
	return p.id
}

func (p *Probe) Group() ProbeGroup {
	return p.group
}

func (p *Probe) FindContactById(id uint64) (Contact, bool) {
	for _, c := range p.Contacts {
		if c.ID() == id {
			return c, true
		}
	}

	return *(NewContact(p.deviceID, p.id, id)), false
}

func (p *Probe) SetContact(m Contact) []Contact {
	result := p.Contacts
	p.Contacts = []Contact{m}

	return result
}
