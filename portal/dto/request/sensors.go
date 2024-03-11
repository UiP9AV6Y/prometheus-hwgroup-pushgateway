package request

type SensorStatus int

const (
	ValueUnknowSensorStatus SensorStatus = iota
	ValueOKSensorStatus
	SensorInvalidSensorStatus
	DeviceInvalidSensorStatus
	OutOfRangeLowSensorStatus
	OutOfRangeHiSensorStatus
)

type SensorEntry struct {
	Entry
	Status   SensorStatus
	Value    int
	Exp      int
	Units    string
	ValueLog *ValueLog `xml:"ValueLog,omitempty"`
}

func NewSensorEntry(id int, name string) *SensorEntry {
	log := NewValueLog(1)
	entry := Entry{
		ID:   id,
		Name: name,
	}
	result := &SensorEntry{
		Entry:    entry,
		ValueLog: log,
	}

	return result
}

func (e *SensorEntry) SetValueLog(l *ValueLog) {
	if l != nil && l.Value != nil && len(l.Value) > 0 {
		e.Value = l.Value[len(l.Value)-1]
	}

	e.ValueLog = l
}

type Sensors struct {
	Entry []SensorEntry
}

func NewSensors(e ...SensorEntry) *Sensors {
	result := &Sensors{
		Entry: e,
	}

	return result
}
