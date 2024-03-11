package request

type OutputEntry struct {
	Entry
	Value    int
	ValueLog *ValueLog `xml:"ValueLog,omitempty"`
}

func NewOutputEntry(id int, name string) *OutputEntry {
	log := NewValueLog(1)
	entry := Entry{
		ID:   id,
		Name: name,
	}
	result := &OutputEntry{
		Entry:    entry,
		ValueLog: log,
	}

	return result
}

func (e *OutputEntry) SetValueLog(l *ValueLog) {
	if l != nil && l.Value != nil && len(l.Value) > 0 {
		e.Value = l.Value[len(l.Value)-1]
	}

	e.ValueLog = l
}

type Outputs struct {
	Entry []OutputEntry
}

func NewOutputs(e ...OutputEntry) *Outputs {
	result := &Outputs{
		Entry: e,
	}

	return result
}
