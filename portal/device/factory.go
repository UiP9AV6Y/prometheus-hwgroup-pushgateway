package device

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

type Factory struct {
	start time.Time
	addr  *url.URL
	gen   *rand.Rand
}

func NewFactory(u *url.URL, r *rand.Rand) *Factory {
	result := &Factory{
		start: time.Now(),
		addr:  u,
		gen:   r,
	}

	return result
}

func NewNowFactory(u *url.URL) *Factory {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	return NewFactory(u, rnd)
}

func (f *Factory) Update(r *request.Root) error {
	if err := f.UpdateAgent(r.Agent); err != nil {
		return err
	}

	if err := f.UpdateNetwork(r.Network); err != nil {
		return err
	}

	if err := f.UpdateTime(r.Time); err != nil {
		return err
	}

	if err := f.UpdatePortal(r.Portal); err != nil {
		return err
	}

	if err := f.UpdateSensors(r.Sensors); err != nil {
		return err
	}

	if err := f.UpdateMeters(r.Meters); err != nil {
		return err
	}

	if err := f.UpdateInputs(r.Inputs); err != nil {
		return err
	}

	if err := f.UpdateOutputs(r.Outputs); err != nil {
		return err
	}

	return nil
}

func (f *Factory) UpdateValueLog(m *request.ValueLog, min, max int) error {
	if m.Value == nil {
		return nil
	}

	limit := max - min

	for i, _ := range m.Value {
		m.Value[i] = min + f.gen.Intn(limit)
	}

	return nil
}

func (f *Factory) Build() (*request.Root, error) {
	a, err := f.BuildAgent()
	if err != nil {
		return nil, err
	}

	n, err := f.BuildNetwork()
	if err != nil {
		return nil, err
	}

	t, err := f.BuildTime()
	if err != nil {
		return nil, err
	}

	p, err := f.BuildPortal()
	if err != nil {
		return nil, err
	}

	s, err := f.BuildSensors()
	if err != nil {
		return nil, err
	}

	m, err := f.BuildMeters()
	if err != nil {
		return nil, err
	}

	i, err := f.BuildInputs()
	if err != nil {
		return nil, err
	}

	o, err := f.BuildOutputs()
	if err != nil {
		return nil, err
	}

	result := request.NewRoot()
	result.Agent = a
	result.Network = n
	result.Time = t
	result.Portal = p
	result.Sensors = s
	result.Meters = m
	result.Inputs = i
	result.Outputs = o

	return result, nil
}
