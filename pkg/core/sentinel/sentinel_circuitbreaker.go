// Copyright 2022 Douyu
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

package sentinel

import (
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
)

type CircuitBreakerRule struct {
	Enable bool `json:"enable"`
	// resource name
	Resource string                  `json:"resource"`
	Strategy circuitbreaker.Strategy `json:"strategy"`
	// RetryTimeoutMs represents recovery timeout (in milliseconds) before the circuit breaker opens.
	// During the open period, no requests are permitted until the timeout has elapsed.
	// After that, the circuit breaker will transform to half-open state for trying a few "trial" requests.
	RetryTimeoutMs uint32 `json:"retryTimeoutMs"`
	// MinRequestAmount represents the minimum number of requests (in an active statistic time span)
	// that can trigger circuit breaking.
	MinRequestAmount uint64 `json:"minRequestAmount"`
	// StatIntervalMs represents statistic time interval of the internal circuit breaker (in ms).
	// Currently the statistic interval is collected by sliding window.
	StatIntervalMs uint32 `json:"statIntervalMs"`
	// StatSlidingWindowBucketCount represents the bucket count of statistic sliding window.
	// The statistic will be more precise as the bucket count increases, but the memory cost increases too.
	// The following must be true — “StatIntervalMs % StatSlidingWindowBucketCount == 0”,
	// otherwise StatSlidingWindowBucketCount will be replaced by 1.
	// If it is not set, default value 1 will be used.
	StatSlidingWindowBucketCount uint32 `json:"statSlidingWindowBucketCount"`
	// MaxAllowedRtMs indicates that any invocation whose response time exceeds this value (in ms)
	// will be recorded as a slow request.
	// MaxAllowedRtMs only takes effect for SlowRequestRatio strategy
	MaxAllowedRtMs uint64 `json:"maxAllowedRtMs"`
	// Threshold represents the threshold of circuit breaker.
	// for SlowRequestRatio, it represents the max slow request ratio
	// for ErrorRatio, it represents the max error request ratio
	// for ErrorCount, it represents the max error request count
	Threshold float64 `json:"threshold"`
}

func convertCbRules(rules []*CircuitBreakerRule) []*circuitbreaker.Rule {
	cb := make([]*circuitbreaker.Rule, 0, len(rules))

	for _, rule := range rules {
		if rule.Enable {
			cb = append(cb, &circuitbreaker.Rule{
				Resource:                     rule.Resource,
				Strategy:                     rule.Strategy,
				RetryTimeoutMs:               rule.RetryTimeoutMs,
				MinRequestAmount:             rule.MinRequestAmount,
				StatIntervalMs:               rule.StatIntervalMs,
				StatSlidingWindowBucketCount: rule.StatSlidingWindowBucketCount,
				MaxAllowedRtMs:               rule.MaxAllowedRtMs,
				Threshold:                    rule.Threshold,
			})
		}
	}

	return cb
}

func labels(resource string) []string {
	return []string{resource, language, pkg.Name(), pkg.AppID(),
		pkg.AppRegion(), pkg.AppZone(), pkg.AppInstance(), conf.GetString("app.mode"),
	}
}

type stateChangeTestListener struct{}

// OnTransformToClosed ...
func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	sentinelState.WithLabelValues(labels(rule.ResourceName())...).Set(float64(circuitbreaker.Closed))
}

// OnTransformToOpen ...
func (s *stateChangeTestListener) OnTransformToOpen(
	prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	sentinelState.WithLabelValues(labels(rule.ResourceName())...).Set(float64(circuitbreaker.Open))
}

// OnTransformToHalfOpen ...
func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	sentinelState.WithLabelValues(labels(rule.ResourceName())...).Set(float64(circuitbreaker.HalfOpen))
}
