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

import "github.com/prometheus/client_golang/prometheus"

// WorkerMetrics ...
type WorkerMetrics struct {
	handledCounter   *prometheus.CounterVec
	handledHistogram *prometheus.HistogramVec
}

// NewWorkerMetrics ...
func NewWorkerMetrics() *WorkerMetrics {
	return &WorkerMetrics{
		handledCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "job_handle_total",
				Help: "",
			}, []string{"type", "name", "code"},
		),
		handledHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "job_handle_seconds",
				Help: "",
			}, []string{"type", "name"},
		),
	}
}

// GetHandlerCounter ...
func (s *WorkerMetrics) GetHandlerCounter() *prometheus.CounterVec {
	return s.handledCounter
}

// GetHandlerHistogram ...
func (s *WorkerMetrics) GetHandlerHistogram() *prometheus.HistogramVec {
	return s.handledHistogram
}
