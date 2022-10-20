// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/prometheus/client_golang/prometheus"
)

// HistogramVecOpts ...
type HistogramVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

type histogramVec struct {
	*prometheus.HistogramVec
}

// Build ...
func (opts HistogramVecOpts) Build() *histogramVec {
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
			Buckets:   opts.Buckets,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &histogramVec{
		HistogramVec: vec,
	}
}

// Observe ...
func (histogram *histogramVec) Observe(v float64, labels ...string) {
	histogram.WithLabelValues(labels...).Observe(v)
}

// NewHistogramVec ...
func NewHistogramVec(name string, labels []string) *histogramVec {
	return HistogramVecOpts{
		Namespace: constant.DefaultNamespace,
		Name:      name,
		Help:      name,
		Labels:    labels,
	}.Build()
}
