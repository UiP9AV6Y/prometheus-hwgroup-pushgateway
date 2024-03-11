package dao

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
)

type DAO struct {
	logger log.Logger
	store  map[uint64]model.Device
	mut    *sync.RWMutex
}

func New(logger log.Logger) *DAO {
	mut := &sync.RWMutex{}
	store := make(map[uint64]model.Device)
	result := &DAO{
		logger: logger,
		store:  store,
		mut:    mut,
	}

	return result
}

func (d *DAO) ImportFile(f string) error {
	fd, err := os.Open(f)
	if err == nil {
		return d.ImportReader(fd)
	} else if errors.Is(err, fs.ErrNotExist) {
		// no file = no data
		return nil
	}

	return err
}

func (d *DAO) ImportReader(r io.Reader) error {
	importer := json.NewDecoder(r)
	payload := []model.Device{}

	if err := importer.Decode(&payload); err != nil {
		return err
	}

	return d.ImportDevices(payload)
}

func (d *DAO) ImportDevices(ms []model.Device) error {
	store := make(map[uint64]model.Device, len(ms))
	for _, m := range ms {
		store[m.ID()] = m
	}

	d.mut.Lock()
	d.store = store
	d.mut.Unlock()

	return nil
}

func (d *DAO) ExportFile(f string) error {
	fd, err := os.Create(f)
	if err != nil {
		return err
	}

	return d.ExportWriter(fd)
}

func (d *DAO) ExportWriter(w io.Writer) error {
	exporter := json.NewEncoder(w)
	payload, err := d.GetAllDevices()
	if err != nil {
		return err
	}

	return exporter.Encode(payload)
}

func (d *DAO) GetAllDevices() ([]model.Device, error) {
	d.mut.RLock()
	defer d.mut.RUnlock()

	result := make([]model.Device, 0, len(d.store))
	for _, m := range d.store {
		result = append(result, m)
	}

	level.Debug(d.logger).Log("msg", "SELECT * FROM devices", "result_set", len(result))

	return result, nil
}

func (d *DAO) GetDeviceByID(id uint64) (model.Device, bool, error) {
	d.mut.RLock()
	defer d.mut.RUnlock()

	m, ok := d.store[id]
	if ok {
		level.Debug(d.logger).Log("msg", "SELECT * FROM devices WHERE id = $ID", "query", m.ID(), "result_set", 1)
		return m, true, nil
	}

	n := model.NewDevice(id)
	level.Debug(d.logger).Log("msg", "SELECT * FROM devices WHERE id = $ID", "query", m.ID(), "result_set", 0)

	return *n, false, nil
}

func (d *DAO) UpsertDevice(m model.Device) error {
	d.mut.Lock()
	d.store[m.ID()] = m
	d.mut.Unlock()

	level.Debug(d.logger).Log("msg", "UPSERT INTO devices", "device", m.ID())

	return nil
}
