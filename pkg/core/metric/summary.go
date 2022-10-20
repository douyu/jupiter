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

// SummaryVecOpts ...
type SummaryVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type summaryVec struct {
	*prometheus.SummaryVec
}

// Build ...
func (opts SummaryVecOpts) Build() *summaryVec {
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &summaryVec{
		SummaryVec: vec,
	}
}

// Observe ...
func (summary *summaryVec) Observe(v float64, labels ...string) {
	summary.WithLabelValues(labels...).Observe(v)
}
