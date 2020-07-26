package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// CounterVecOpts ...
type SummaryVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
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
		vec,
	}
}

type summaryVec struct {
	*prometheus.SummaryVec
}

// Observe ...
func (counter *summaryVec) XObserve(v float64, labels ...string) {
	counter.WithLabelValues(labels...).Observe(v)
}
