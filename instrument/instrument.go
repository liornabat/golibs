package instrument

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Instrument struct {
	nameSpace string
	subSystem string
	kind      metricType
	metric    interface{}
}
type metricType int

const (
	undefined    metricType = 0
	counterVec   metricType = 1
	gaugeVec     metricType = 2
	histogramVec metricType = 3
)

func NewCounterMetric(nameSpace, subSystem, name, help string, labels []string) (*Instrument, error) {
	i := &Instrument{
		kind:      counterVec,
		nameSpace: nameSpace,
		subSystem: subSystem,
	}
	i.metric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: i.nameSpace,
			Subsystem: i.subSystem,
			Name:      name,
			Help:      help,
		},
		labels)

	err := prometheus.Register(i.metric.(*prometheus.CounterVec))
	if err != nil {
		return nil, err
	}
	return i, nil
}

func NewGaugeMetric(nameSpace, subSystem, name, help string, labels []string) (*Instrument, error) {
	i := &Instrument{
		kind:      gaugeVec,
		nameSpace: nameSpace,
		subSystem: subSystem,
	}
	i.metric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: i.nameSpace,
			Subsystem: i.subSystem,
			Name:      name,
			Help:      help,
		},
		labels)

	err := prometheus.Register(i.metric.(*prometheus.GaugeVec))
	if err != nil {
		return nil, err
	}
	return i, nil
}

func NewHistogramMetric(nameSpace, subSystem, name, help string, labels []string, buckets []float64) (*Instrument, error) {

	i := &Instrument{
		kind:      histogramVec,
		nameSpace: nameSpace,
		subSystem: subSystem,
	}
	i.metric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: i.nameSpace,
			Subsystem: i.subSystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labels)

	err := prometheus.Register(i.metric.(*prometheus.HistogramVec))
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *Instrument) Add(value float64, lvs ...string) *Instrument {
	switch i.kind {
	case counterVec:
		metric, err := i.metric.(*prometheus.CounterVec).GetMetricWithLabelValues(lvs...)
		if err == nil {
			metric.Add(value)
		}
	case gaugeVec:
		metric, err := i.metric.(*prometheus.GaugeVec).GetMetricWithLabelValues(lvs...)
		if err == nil {
			metric.Add(value)
		}
	}
	return i
}

func (i *Instrument) Sub(value float64, lvs ...string) *Instrument {
	switch i.kind {

	case gaugeVec:
		metric, err := i.metric.(*prometheus.GaugeVec).GetMetricWithLabelValues(lvs...)
		if err == nil {
			metric.Sub(value)
		}
	}
	return i
}

func (i *Instrument) Inc(lvs ...string) *Instrument {
	switch i.kind {
	case counterVec:
		metric, err := i.metric.(*prometheus.CounterVec).GetMetricWithLabelValues(lvs...)
		if err == nil {
			metric.Inc()
		}
	case gaugeVec:
		metric, err := i.metric.(*prometheus.GaugeVec).GetMetricWithLabelValues(lvs...)
		if err == nil {
			metric.Inc()
		}
	}
	return i
}
func (i *Instrument) Dec(lvs ...string) *Instrument {
	switch i.kind {

	case gaugeVec:
		metric, err := i.metric.(*prometheus.GaugeVec).GetMetricWithLabelValues(lvs...)
		if err == nil {
			metric.Dec()
		}
	}
	return i
}

func (i *Instrument) Set(value float64, lvs ...string) *Instrument {
	switch i.kind {

	case gaugeVec:
		metric, err := i.metric.(*prometheus.GaugeVec).GetMetricWithLabelValues(lvs...)
		if err == nil {
			metric.Set(value)
		}
	}
	return i
}

func (i *Instrument) Observe(value float64, lvs ...string) *Instrument {
	switch i.kind {

	case histogramVec:
		metric, err := i.metric.(*prometheus.HistogramVec).GetMetricWithLabelValues(lvs...)
		if err == nil {
			metric.Observe(value)
		}
	}
	return i
}

func (i *Instrument) SetUnregistered() {
	switch i.kind {
	case counterVec:
		metric := i.metric.(*prometheus.CounterVec)
		prometheus.Unregister(metric)
	case gaugeVec:
		metric := i.metric.(*prometheus.GaugeVec)
		prometheus.Unregister(metric)
	case histogramVec:
		metric := i.metric.(*prometheus.HistogramVec)
		prometheus.Unregister(metric)
	}

}
