package instrument

type InstrumentArray struct {
	Namespace  string
	Subsystem  string
	Name       string
	counters   *Instrument
	gauges     *Instrument
	histograms *Instrument
}

func NewInstrumentArray(nameSpace, subSystem, name string) *InstrumentArray {
	ia := &InstrumentArray{
		Namespace: nameSpace,
		Subsystem: subSystem,
		Name:      name,
	}
	return ia
}

func (ia *InstrumentArray) AddCounter(labels []string, help string) *InstrumentArray {
	ins, err := NewCounterMetric(ia.Namespace, ia.Subsystem, ia.Name+"_totals", help, labels)
	if err == nil {
		ia.counters = ins
	}
	return ia
}

func (ia *InstrumentArray) AddGauge(labels []string, help string) *InstrumentArray {
	ins, err := NewGaugeMetric(ia.Namespace, ia.Subsystem, ia.Name+"_metrics", help, labels)
	if err == nil {
		ia.gauges = ins
	}
	return ia
}

func (ia *InstrumentArray) AddHistogram(labels []string, buckets []float64, help string) *InstrumentArray {
	ins, err := NewHistogramMetric(ia.Namespace, ia.Subsystem, ia.Name+"_Observations", help, labels, buckets)
	if err == nil {
		ia.histograms = ins
	}
	return ia
}

func (ia *InstrumentArray) AddToCounter(value float64, lvs ...string) *InstrumentArray {
	if ia.counters != nil {
		ia.counters.Add(value, lvs...)
	}
	return ia
}

func (ia *InstrumentArray) AddToGauge(value float64, lvs ...string) *InstrumentArray {

	if ia.gauges != nil {
		ia.gauges.Add(value, lvs...)
	}
	return ia
}

func (ia *InstrumentArray) SubFromGauge(value float64, lvs ...string) *InstrumentArray {

	if ia.gauges != nil {
		ia.gauges.Sub(value, lvs...)
	}

	return ia
}

func (ia *InstrumentArray) IncToCounter(lvs ...string) *InstrumentArray {
	if ia.counters != nil {
		ia.counters.Inc(lvs...)
	}
	return ia
}

func (ia *InstrumentArray) IncToGauge(lvs ...string) *InstrumentArray {
	if ia.gauges != nil {
		ia.gauges.Inc(lvs...)
	}
	return ia
}

func (ia *InstrumentArray) DecFromGauge(lvs ...string) *InstrumentArray {

	if ia.gauges != nil {
		ia.gauges.Dec(lvs...)
	}
	return ia
}

func (ia *InstrumentArray) SetGauge(value float64, lvs ...string) *InstrumentArray {

	if ia.gauges != nil {
		ia.gauges.Set(value, lvs...)
	}
	return ia
}

func (ia *InstrumentArray) ObserveHistogram(value float64, lvs ...string) *InstrumentArray {
	if ia.histograms != nil {
		ia.histograms.Observe(value, lvs...)
	}
	return ia
}

func (ia *InstrumentArray) UnRegister() {

	if ia.histograms != nil {
		ia.histograms.SetUnregistered()
	}
	if ia.counters != nil {
		ia.counters.SetUnregistered()
	}
	if ia.gauges != nil {
		ia.gauges.SetUnregistered()
	}

}
