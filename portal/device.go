package portal

import (
	"net/netip"
	"time"

	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dao/model"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/common"
	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/request"
)

func (p *Portal) processAgent(addr *netip.Addr, a *request.Agent) (*model.Device, string, error) {
	if a == nil || a.SerNum == 0 {
		level.Info(p.logger).Log("msg", "request agent information is missing", "client", addr)
		return nil, "Insufficient push data provided", nil
	}

	device, _, err := p.db.GetDeviceByID(a.SerNum)
	if err != nil {
		return nil, "", err
	}

	device.LastSeen = uint64(time.Now().Unix())
	device.Model = uint8(a.Model)
	device.Firmware = a.Version
	device.Address = addr.String()

	level.Debug(p.logger).Log("msg", "processing agent information", "device", device.ID())

	return &device, "", nil
}

func (p *Portal) processTime(device *model.Device, t *common.Time) (string, error) {
	if t == nil {
		level.Debug(p.logger).Log("msg", "client request contains no time information", "device", device.ID())
		return "", nil
	}

	device.Uptime = uint64(t.UpTime)

	level.Debug(p.logger).Log("msg", "processing time information", "device", device.ID())

	return "", nil
}

func (p *Portal) processPortal(device *model.Device, q *request.Portal) (string, error) {
	if q == nil {
		level.Debug(p.logger).Log("msg", "client request contains no portal information", "device", device.ID())
		return "", nil
	}

	device.Name = q.DeviceName

	level.Debug(p.logger).Log("msg", "processing portal information", "device", device.ID())

	return "", nil
}
