package device

import (
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func (f *Factory) BuildSensors() (*request.Sensors, error) {
	temps := request.NewValueLog(2, 245, 245, 245, 245)
	hums := request.NewValueLog(2, 265, 265, 265, 265)

	temp := request.NewSensorEntry(215, "Sensor 215")
	temp.Status = request.ValueOKSensorStatus
	temp.Exp = -1
	temp.Units = "C"
	temp.SetValueLog(temps)

	hum := request.NewSensorEntry(217, "Sensor 217")
	hum.Status = request.ValueOKSensorStatus
	hum.Exp = -1
	hum.Units = "%RH"
	hum.SetValueLog(hums)

	result := request.NewSensors(*temp, *hum)

	return result, nil
}

func (f *Factory) UpdateSensors(m *request.Sensors) error {
	entries := make([]request.SensorEntry, 0, len(m.Entry))
	for _, e := range m.Entry {
		if e.ValueLog != nil {
			f.UpdateValueLog(e.ValueLog, 100, 340)
		}

		entry := request.NewSensorEntry(e.ID, e.Name)
		entry.Exp = e.Exp
		entry.Units = e.Units
		entry.Status = e.Status
		entry.SetValueLog(e.ValueLog)

		entries = append(entries, *entry)
	}

	m.Entry = entries

	return nil
}
