package device

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func (f *Factory) BuildInputs() (*request.Inputs, error) {
	doorState := request.NewValueLog(2, 0, 0, 0, 0)
	windowState := request.NewValueLog(2, 0, 0, 0, 0)

	door := request.NewInputEntry(1, "Door")
	door.SetValueLog(doorState)

	window := request.NewInputEntry(2, "Window")
	window.SetValueLog(windowState)

	result := request.NewInputs(*door, *window)

	return result, nil
}

func (f *Factory) UpdateInputs(m *request.Inputs) error {
	entries := make([]request.InputEntry, 0, len(m.Entry))
	for _, e := range m.Entry {
		if e.ValueLog != nil {
			f.UpdateValueLog(e.ValueLog, 0, 1)
		}

		entry := request.NewInputEntry(e.ID, e.Name)
		entry.SetValueLog(e.ValueLog)

		entries = append(entries, *entry)
	}

	m.Entry = entries

	return nil
}
