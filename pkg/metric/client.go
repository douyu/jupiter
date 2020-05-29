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

// ClientMetrics ...
type ClientMetrics struct {
	handledCounter   *prometheus.CounterVec
	handledHistogram *prometheus.HistogramVec
}

// NewClientMetrics ...
func NewClientMetrics() *ClientMetrics {
	return &ClientMetrics{
		handledCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_handle_total",
				Help: "",
			}, []string{"type", "name", "method", "server", "code"},
		),
		handledHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "client_handle_seconds",
				Help: "",
			}, []string{"type", "name", "method", "server"},
		),
	}
}

// GetHandlerCounter ...
func (s *ClientMetrics) GetHandlerCounter() *prometheus.CounterVec {
	return s.handledCounter
}

// GetHandlerHistogram ...
func (s *ClientMetrics) GetHandlerHistogram() *prometheus.HistogramVec {
	return s.handledHistogram
}
