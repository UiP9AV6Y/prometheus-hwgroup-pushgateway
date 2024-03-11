package model

type Contact struct {
	deviceID uint64
	probeID  uint64
	id       uint64

	Changes uint64
	Status  uint8
	Value   float64
	Unit    string
	Name    string
}

func NewContact(deviceID uint64, probeID uint64, id uint64) *Contact {
	result := &Contact{
		deviceID: deviceID,
		probeID:  probeID,
		id:       id,
	}

	return result
}

func (c *Contact) DeviceID() uint64 {
	return c.deviceID
}

func (c *Contact) ProbeID() uint64 {
	return c.probeID
}

func (c *Contact) ID() uint64 {
	return c.id
}
