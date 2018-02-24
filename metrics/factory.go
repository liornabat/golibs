// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"github.com/liornabat/golibs/cache"
)

// Factory creates new metrics
type Factory interface {
	Counter(name string, tags map[string]string) Counter
	Timer(name string, tags map[string]string) Timer
	Gauge(name string, tags map[string]string) Gauge

	// Namespace returns a nested metrics factory.
	Namespace(name string, tags map[string]string) Factory
}

// NullFactory is a metrics factory that returns NullCounter, NullTimer, and NullGauge.
var NullFactory Factory = nullFactory{}

type nullFactory struct{}

func (nullFactory) Counter(name string, tags map[string]string) Counter   { return NullCounter }
func (nullFactory) Timer(name string, tags map[string]string) Timer       { return NullTimer }
func (nullFactory) Gauge(name string, tags map[string]string) Gauge       { return NullGauge }
func (nullFactory) Namespace(name string, tags map[string]string) Factory { return NullFactory }

type MetricsFactory struct {
	metricsCache *cache.LRU
	promFactory  *PrometheusFactory
}

func NewMetricsFactory() *MetricsFactory {
	return &MetricsFactory{
		metricsCache: cache.NewLRU(1000),
		promFactory:  New(),
	}
}

func (mf *MetricsFactory) SetNamespace(name string, tags map[string]string) *MetricsFactory {
	mf.promFactory.Namespace(name, tags)
	return mf
}

func (mf *MetricsFactory) AddCounter(key, name string, tags map[string]string) *Counter {
	c := mf.promFactory.Counter(name, tags)
	mf.metricsCache.Put(key, &c)
	return &c
}

func (mf *MetricsFactory) AddGauge(key, name string, tags map[string]string) *Gauge {
	g := mf.promFactory.Gauge(name, tags)
	mf.metricsCache.Put(key, &g)
	return &g
}

func (mf *MetricsFactory) AddTimer(key, name string, tags map[string]string) *Timer {
	t := mf.promFactory.Timer(name, tags)
	mf.metricsCache.Put(key, &t)
	return &t
}

func (mf *MetricsFactory) GetCounter(key string) *Counter {
	v, ok := mf.metricsCache.Get(key).(*Counter)
	if ok {
		return v
	}
	return nil
}

func (mf *MetricsFactory) GetGauge(key string) *Gauge {
	v, ok := mf.metricsCache.Get(key).(*Gauge)
	if ok {
		return v
	}
	return nil

}
func (mf *MetricsFactory) GetTimer(key string) *Timer {
	v, ok := mf.metricsCache.Get(key).(*Timer)
	if ok {
		return v
	}
	return nil
}
