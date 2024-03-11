package device

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func (f *Factory) BuildMeters() (*request.Meters, error) {
	power := request.NewMeterValueEntry(1001, "Power")
	power.Unit = "kWh"
	power.Value = 84
	power.Exp = -3
	power.Status = request.ValueOKSensorStatus

	values := request.NewMeterValues(*power)

	landis := request.NewMeterEntry(1, "Landis")
	landis.Address = 1
	landis.SecAddress = 94653673
	landis.Values = *values

	result := request.NewMeters(*landis)

	return result, nil
}

func (f *Factory) UpdateMeters(m *request.Meters) error {
	entries := make([]request.MeterEntry, 0, len(m.Entry))
	for _, e := range m.Entry {
		values := make([]request.MeterValueEntry, 0, len(e.Values.Entry))
		for _, v := range e.Values.Entry {
			value := request.NewMeterValueEntry(v.ID, v.Name)
			value.Value = 10 + f.gen.Intn(540)
			value.Unit = v.Unit
			value.Exp = v.Exp
			value.Status = v.Status

			values = append(values, *value)
		}

		entry := request.NewMeterEntry(e.ID, e.Name)
		entry.Address = e.Address
		entry.SecAddress = e.SecAddress
		entry.Values.Entry = values

		entries = append(entries, *entry)
	}

	m.Entry = entries

	return nil
}
