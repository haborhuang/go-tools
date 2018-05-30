package monitor

import (
	"time"

	"github.com/haborhuang/go-tools/disk-usage/types"
	"github.com/haborhuang/go-tools/disk-usage/du"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	l, err := zap.NewDevelopment()
	if nil != err {
		panic(err)
	}

	logger = l.Sugar()
}

func StartDUMonitorOrDie(mc types.MetricsConfig, ps types.DisksPaths) {
	if len(ps) < 1 {
		panic("Missing paths of disks")
	}

	for d, p := range ps {
		if d == "" || p == "" {
			panic("Empty disk name or path")
		}

		startMonitor(d, p, newDUMetrics(mc.Namespace, mc.SubSys))
	}
}

type monitor struct {
	duMetrics *diskUsageMetrics
	name string
	diskPath string
}

func startMonitor(name, diskPath string, duMetrics *diskUsageMetrics) {
	m := monitor{
		name: name,
		diskPath: diskPath,
		duMetrics: duMetrics,
	}

	go m.loop()
}

func (m *monitor) loop() {
	logger.Infof("DU monitor for '%s' in '%s' started", m.name, m.diskPath)
	for {
		<-time.After(10 * time.Second)

		_, c, u, err := m.du()
		if nil != err {
			logger.Errorf("Get disk usage of '%s' error: %v", m.name, err)
		}

		m.duMetrics.set(m.name, float64(u), float64(c))
	}
}

// Get disk info of the mounted path
func (m *monitor) du() (int64, int64, int64, error) {
	return du.DiskUsage(m.diskPath)
}
