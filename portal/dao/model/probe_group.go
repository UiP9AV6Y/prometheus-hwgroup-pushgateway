package model

type ProbeGroup string

func (t ProbeGroup) String() string {
	return string(t)
}

const (
	MeterProbeGroup  ProbeGroup = "meter"
	InputProbeGroup  ProbeGroup = "input"
	OutputProbeGroup ProbeGroup = "output"
	SensorProbeGroup ProbeGroup = "sensor"
)
