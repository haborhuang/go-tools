package monitor

import "github.com/prometheus/client_golang/prometheus"

const (
	nameLabel = "disk_name"
)

type diskUsageMetrics struct {
	duGaugeVec *prometheus.GaugeVec
	dcGaugeVec *prometheus.GaugeVec
}

func newDUMetrics(namespace, subsystem string) *diskUsageMetrics {
	duGaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts {
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "disk_used",
			Help:      "bytes of used space of disk",
		},
		[]string{
			// disk name
			nameLabel,
		},
	)
	dcGaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts {
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "disk_capacity",
			Help:      "bytes of space capacity of disk",
		},
		[]string{
			// disk name
			nameLabel,
		},
	)

	prometheus.MustRegister(duGaugeVec)
	prometheus.MustRegister(dcGaugeVec)

	return &diskUsageMetrics{
		duGaugeVec: duGaugeVec,
		dcGaugeVec: dcGaugeVec,
	}
}

func (m *diskUsageMetrics) set(name string, used, capacity float64) {
	m.duGaugeVec.With(prometheus.Labels{
		nameLabel:   name,
	}).Set(used)

	m.dcGaugeVec.With(prometheus.Labels{
		nameLabel:   name,
	}).Set(capacity)
}