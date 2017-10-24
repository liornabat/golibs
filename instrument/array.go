package instrument

import (
	"golang.org/x/sync/syncmap"
)

type InstrumentArray struct {
	counters   syncmap.Map
	gauges     syncmap.Map
	histograms syncmap.Map
}

func NewInstrumentArray() *InstrumentArray {
	ia := &InstrumentArray{}
	return ia
}

func (ia *InstrumentArray) AddCounter(nameSpace, subSystem, name string, labels []string, help ...string) *InstrumentArray {
	ins, err := NewInstrument().
		SetNameSpace(nameSpace).
		SetSubSystem(subSystem).
		NewCounterVec(name, labels, help...)
	if err == nil {
		ia.counters.LoadOrStore(name, ins)
	}
	return ia
}

func (ia *InstrumentArray) AddGauge(nameSpace, subSystem, name string, labels []string, help ...string) *InstrumentArray {
	ins, err := NewInstrument().
		SetNameSpace(nameSpace).
		SetSubSystem(subSystem).
		NewGaugeVec(name, labels, help...)
	if err == nil {
		ia.gauges.LoadOrStore(name, ins)
	}
	return ia
}

func (ia *InstrumentArray) AddHistogram(nameSpace, subSystem, name string, labels []string, buckets []float64, help ...string) *InstrumentArray {
	ins, err := NewInstrument().
		SetNameSpace(nameSpace).
		SetSubSystem(subSystem).
		NewHistogramVec(name, labels, buckets, help...)
	if err == nil {
		ia.gauges.LoadOrStore(name, ins)
	}
	return ia
}

func (ia *InstrumentArray) AddToCounter(name string, value float64, lvs ...string) *InstrumentArray {
	v, ok := ia.counters.Load(name)
	if ok {
		ins := v.(*Instrument)
		ins.Add(value, lvs...)
	}
	return ia
}

func (ia *InstrumentArray) AddToGauge(name string, value float64, lvs ...string) *InstrumentArray {
	v, ok := ia.gauges.Load(name)
	if ok {
		ins := v.(*Instrument)
		ins.Add(value, lvs...)
	}
	return ia
}

func (ia *InstrumentArray) SubFromGauge(name string, value float64, lvs ...string) *InstrumentArray {
	v, ok := ia.gauges.Load(name)
	if ok {
		ins := v.(*Instrument)
		ins.Sub(value, lvs...)
	}
	return ia
}

func (ia *InstrumentArray) IncToCounter(name string, value float64, lvs ...string) *InstrumentArray {
	v, ok := ia.counters.Load(name)
	if ok {
		ins := v.(*Instrument)
		ins.Inc(lvs...)
	}
	return ia
}

func (ia *InstrumentArray) IncToGauge(name string, lvs ...string) *InstrumentArray {
	v, ok := ia.gauges.Load(name)
	if ok {
		ins := v.(*Instrument)
		ins.Inc(lvs...)
	}
	return ia
}

func (ia *InstrumentArray) DecFromGauge(name string, lvs ...string) *InstrumentArray {
	v, ok := ia.gauges.Load(name)
	if ok {
		ins := v.(*Instrument)
		ins.Dec(lvs...)
	}
	return ia
}

func (ia *InstrumentArray) SetGauge(name string, value float64, lvs ...string) *InstrumentArray {
	v, ok := ia.gauges.Load(name)
	if ok {
		ins := v.(*Instrument)
		ins.Set(value, lvs...)
	}
	return ia
}

func (ia *InstrumentArray) ObserveHistogram(name string, value float64, lvs ...string) *InstrumentArray {
	v, ok := ia.histograms.Load(name)
	if ok {
		ins := v.(*Instrument)
		ins.Observe(value, lvs...)
	}
	return ia
}
