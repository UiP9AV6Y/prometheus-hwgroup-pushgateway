package request

type InputEntry struct {
	Entry
	Value    int
	ValueLog *ValueLog `xml:"ValueLog,omitempty"`
}

func NewInputEntry(id int, name string) *InputEntry {
	log := NewValueLog(1)
	entry := Entry{
		ID:   id,
		Name: name,
	}
	result := &InputEntry{
		Entry:    entry,
		ValueLog: log,
	}

	return result
}

func (e *InputEntry) SetValueLog(l *ValueLog) {
	if l != nil && l.Value != nil && len(l.Value) > 0 {
		e.Value = l.Value[len(l.Value)-1]
	}

	e.ValueLog = l
}

type Inputs struct {
	Entry []InputEntry
}

func NewInputs(e ...InputEntry) *Inputs {
	result := &Inputs{
		Entry: e,
	}

	return result
}
