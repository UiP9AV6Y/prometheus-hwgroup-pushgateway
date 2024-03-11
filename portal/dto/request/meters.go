package request

type MeterValueEntry struct {
	Entry
	Unit   string
	Value  int
	Exp    int
	Status SensorStatus
}

func NewMeterValueEntry(id int, name string) *MeterValueEntry {
	entry := Entry{
		ID:   id,
		Name: name,
	}
	result := &MeterValueEntry{
		Entry: entry,
	}

	return result
}

type MeterValues struct {
	Entry []MeterValueEntry
}

func NewMeterValues(e ...MeterValueEntry) *MeterValues {
	result := &MeterValues{
		Entry: e,
	}

	return result
}

type MeterEntry struct {
	Entry
	Address    int
	SecAddress int
	Values     MeterValues
}

func NewMeterEntry(id int, name string) *MeterEntry {
	entry := Entry{
		ID:   id,
		Name: name,
	}
	result := &MeterEntry{
		Entry: entry,
	}

	return result
}

type Meters struct {
	Entry []MeterEntry
}

func NewMeters(e ...MeterEntry) *Meters {
	result := &Meters{
		Entry: e,
	}

	return result
}
