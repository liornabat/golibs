package instrument

type HttpInstrument struct {
	Name         string
	counterVec   *Instrument
	histogramVec *Instrument
}

type HttpCallInterface interface {
}

func (i *Instrument) SetHttpInstrument(name string) (*HttpInstrument, error) {
	hci := &HttpInstrument{
		Name: name,
	}
	counterVec, err := NewCounterVec(i, name+"_totals", []string{"type", "func", "result"}, "counters for total results of function calls")
	if err != nil {
		return nil, err
	}
	hci.counterVec = counterVec

	histogramVec, err := NewHistogramVec(i, name+"_stats", []string{"type", "func"}, []float64{0.01, 0.5, 1, 2, 5, 10, 20, 30}, "histogram for stats results of function calls")
	if err != nil {
		return nil, err
	}
	hci.histogramVec = histogramVec

	return hci, nil
}

func (hci *HttpInstrument) SetResult(src, funcName string, ok bool) *HttpInstrument {
	if ok {
		hci.counterVec.Inc(src, funcName, "ok")

	} else {
		hci.counterVec.Inc(src, funcName, "fail")
	}
	return hci
}

func (hci *HttpInstrument) SetObserve(src, funcName string, value float64) *HttpInstrument {
	hci.histogramVec.Observe(value, funcName)
	return hci
}
