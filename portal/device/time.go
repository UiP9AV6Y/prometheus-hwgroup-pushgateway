package device

import (
	"time"

	"github.com/UiP9AV6Y/prometheus-hwgroup-pushgateway/portal/dto/common"
)

func (f *Factory) BuildTime() (*common.Time, error) {
	result := common.NewTime()
	result.Now()

	return result, nil
}

func (f *Factory) UpdateTime(m *common.Time) error {
	if m == nil {
		return nil
	}

	now := time.Now()
	uptime := now.Sub(f.start)

	m.UpTime = int(uptime.Seconds())
	m.Timestamp = common.Timestamp(now)

	return nil
}
