package device

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func (f *Factory) BuildOutputs() (*request.Outputs, error) {
	lightState := request.NewValueLog(2, 0, 0, 0, 0)

	light := request.NewOutputEntry(151, "Light")
	light.SetValueLog(lightState)

	result := request.NewOutputs(*light)

	return result, nil
}

func (f *Factory) UpdateOutputs(m *request.Outputs) error {
	entries := make([]request.OutputEntry, 0, len(m.Entry))
	for _, e := range m.Entry {
		if e.ValueLog != nil {
			f.UpdateValueLog(e.ValueLog, 0, 1)
		}

		entry := request.NewOutputEntry(e.ID, e.Name)
		entry.SetValueLog(e.ValueLog)

		entries = append(entries, *entry)
	}

	m.Entry = entries

	return nil
}
