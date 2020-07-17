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

package leastloaded_test

import (
	"testing"

	"github.com/douyu/jupiter/pkg/util/xp2c/leastloaded"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
)

func TestLeastLoaded(t *testing.T) {
	t.Run("0 item", func(t *testing.T) {
		ll := leastloaded.New()
		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Nil(t, item)
	})

	t.Run("1 item", func(t *testing.T) {
		ll := leastloaded.New()
		ll.Add(1)
		item, done := ll.Next()
		done(balancer.DoneInfo{})
		assert.Equal(t, 1, item)
	})

	t.Run("3 items", func(t *testing.T) {
		ll := leastloaded.New()
		ll.Add(1)
		ll.Add(2)
		ll.Add(3)

		countMap := make(map[interface{}]int)

		totalCount := 10000
		for i := 0; i < totalCount; i++ {
			item, done := ll.Next()
			done(balancer.DoneInfo{})

			countMap[item]++
		}

		total := 0
		for _, count := range countMap {
			total += count
			assert.Less(t, totalCount/3-200, count)
		}

		assert.Equal(t, totalCount, total)
	})
}

func TestLeastLoadedAbnormal(t *testing.T) {
	t.Run("fixed inflight", func(t *testing.T) {
		ll := leastloaded.New()
		ll.Add(1)
		ll.Add(2)
		ll.Add(3)

		item, _ := ll.Next()

		for i := 0; i < 1000; i++ {
			next, done := ll.Next()
			done(balancer.DoneInfo{})

			assert.NotEqual(t, item, next)
		}
	})
}
